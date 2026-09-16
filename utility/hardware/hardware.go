package hardware

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"
)

// 系统平台标识
const (
	PlatformLinux   = "linux"
	PlatformWindows = "windows"
	PlatformMacOS   = "macos"

	// cmdTimeout 外部命令执行的超时时间
	cmdTimeout = 10 * time.Second
)

// GetHardwareInfo 返回当前运行系统所属平台：linux / windows / macos。
// 无法识别的系统平台返回错误。
func GetHardwareInfo() (string, error) {
	switch runtime.GOOS {
	case "linux":
		return PlatformLinux, nil
	case "windows":
		return PlatformWindows, nil
	case "darwin":
		return PlatformMacOS, nil
	default:
		return "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// GetCPUTemp 返回当前 CPU 温度（摄氏度）。
func GetCPUTemp() (float64, error) {
	sys := runtime.GOOS
	switch sys {
	case "linux":
		return getCPUTempLinux()
	case "windows":
		return getCPUTempWindows()
	case "darwin":
		return getCPUTempDarwin()
	default:
		return 0, fmt.Errorf("unsupported platform: %s", sys)
	}
}

// getCPUTempLinux 通过 sysfs 读取 CPU 温度：
// 优先 CPU 相关的 thermal zone，其次 CPU 相关的 hwmon 采样点。
func getCPUTempLinux() (float64, error) {
	if temp, err := getThermalZoneTemp(); err == nil {
		return temp, nil
	}
	if temp, err := getHwmonTemp(); err == nil {
		return temp, nil
	}
	return 0, errors.New("no cpu temperature sensor found under /sys/class/thermal or /sys/class/hwmon")
}

// getThermalZoneTemp 读取 thermal zone 温度（毫摄氏度），优先选择 CPU 相关的 zone。
func getThermalZoneTemp() (float64, error) {
	zones, err := filepath.Glob("/sys/class/thermal/thermal_zone*")
	if err != nil || len(zones) == 0 {
		return 0, errors.New("no thermal zone found")
	}
	var fallback []float64
	for _, zone := range zones {
		typ, _ := os.ReadFile(filepath.Join(zone, "type"))
		zoneType := strings.ToLower(strings.TrimSpace(string(typ)))
		temp, err := readMilliCelsius(filepath.Join(zone, "temp"))
		if err != nil {
			continue
		}
		if strings.Contains(zoneType, "cpu") || strings.Contains(zoneType, "pkg_temp") || strings.Contains(zoneType, "soc") {
			return temp, nil
		}
		fallback = append(fallback, temp)
	}
	if len(fallback) > 0 {
		return fallback[0], nil
	}
	return 0, errors.New("no valid thermal zone temperature")
}

// getHwmonTemp 读取 hwmon 中 CPU 相关驱动（coretemp/k10temp 等）的采样点温度。
func getHwmonTemp() (float64, error) {
	hwmons, err := filepath.Glob("/sys/class/hwmon/hwmon*")
	if err != nil || len(hwmons) == 0 {
		return 0, errors.New("no hwmon device found")
	}
	cpuDrivers := []string{"coretemp", "k10temp", "zenpower", "k8temp", "cpu_thermal", "soc_thermal"}
	for _, hwmon := range hwmons {
		name, err := os.ReadFile(filepath.Join(hwmon, "name"))
		if err != nil {
			continue
		}
		if !slices.Contains(cpuDrivers, strings.ToLower(strings.TrimSpace(string(name)))) {
			continue
		}
		if temp, err := readMilliCelsius(filepath.Join(hwmon, "temp1_input")); err == nil {
			return temp, nil
		}
	}
	return 0, errors.New("no cpu hwmon temperature found")
}

// winTempQuery WMI 热区温度查询脚本，输出单位为 0.1 开尔文的数值（每个热区一行）。
const winTempQuery = `$ErrorActionPreference='SilentlyContinue'; ` +
	`$t=(Get-CimInstance -Namespace root/wmi -ClassName MSAcpi_ThermalZoneTemperature).CurrentTemperature; ` +
	`if(-not $t){$t=(Get-CimInstance -ClassName Win32_PerfFormattedData_Counters_ThermalZoneInformation).HighPrecisionTemperature}; ` +
	`$t`

// getCPUTempWindows 通过 WMI 查询 ACPI 热区温度（单位为 0.1 开尔文），取最高值作为 CPU 温度。
func getCPUTempWindows() (float64, error) {
	out, err := runCommand("powershell", "-NoProfile", "-NonInteractive", "-Command", winTempQuery)
	if err != nil {
		return 0, fmt.Errorf("query windows thermal zone failed: %w", err)
	}
	var best float64
	found := false
	for _, line := range strings.Split(out, "\n") {
		v, err := strconv.ParseFloat(strings.TrimSpace(line), 64)
		if err != nil {
			continue
		}
		celsius := v/10 - 273.15
		if !validTempC(celsius) {
			continue
		}
		if !found || celsius > best {
			best, found = celsius, true
		}
	}
	if !found {
		return 0, errors.New("no cpu temperature sensor found on windows")
	}
	return best, nil
}

// getCPUTempDarwin 读取 macOS CPU 温度：
// 优先使用无需 root 的 osx-cpu-temp，其次使用需要 root 的 powermetrics。
func getCPUTempDarwin() (float64, error) {
	if path, err := exec.LookPath("osx-cpu-temp"); err == nil {
		if out, err := runCommand(path); err == nil {
			if temp, err := parseFirstTemp(out); err == nil {
				return temp, nil
			}
		}
	}
	if os.Geteuid() == 0 {
		if path, err := exec.LookPath("powermetrics"); err == nil {
			if out, err := runCommand(path, "--samplers", "smc", "-n", "1", "-i", "1000"); err == nil {
				for _, line := range strings.Split(out, "\n") {
					if !strings.Contains(line, "CPU die temperature") {
						continue
					}
					if temp, err := parseFirstTemp(line); err == nil {
						return temp, nil
					}
				}
			}
		}
	}
	return 0, errors.New("reading cpu temperature on macOS requires root (powermetrics) or the osx-cpu-temp utility")
}

// readMilliCelsius 读取 sysfs 中单位为毫摄氏度的温度文件并转换为摄氏度。
func readMilliCelsius(path string) (float64, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
	if err != nil {
		return 0, err
	}
	celsius := v / 1000
	if !validTempC(celsius) {
		return 0, fmt.Errorf("invalid temperature %.2f from %s", celsius, path)
	}
	return celsius, nil
}

// tempNumberRe 用于从命令输出中提取第一个数字（温度值）。
var tempNumberRe = regexp.MustCompile(`-?\d+(?:\.\d+)?`)

// parseFirstTemp 从输出文本中提取第一个温度数值（摄氏度）。
func parseFirstTemp(s string) (float64, error) {
	match := tempNumberRe.FindString(s)
	if match == "" {
		return 0, errors.New("no temperature number found in output")
	}
	temp, err := strconv.ParseFloat(match, 64)
	if err != nil {
		return 0, err
	}
	if !validTempC(temp) {
		return 0, fmt.Errorf("invalid temperature %.2f", temp)
	}
	return temp, nil
}

// validTempC 判断温度是否为合理的摄氏度范围。
func validTempC(t float64) bool {
	return t > -50 && t < 150
}

// runCommand 在超时限制内执行外部命令并返回标准输出。
func runCommand(name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
	defer cancel()
	out, err := exec.CommandContext(ctx, name, args...).Output()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("%s timed out after %s", name, cmdTimeout)
		}
		return "", err
	}
	return string(out), nil
}
