package hardware

import (
	"errors"
	"fmt"
	"math"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// UsageStat 系统资源使用量统计（内存 / 磁盘共用），大小单位均为字节。
type UsageStat struct {
	Total       uint64  // 总大小（字节）
	Used        uint64  // 已使用大小（字节）
	UsedPercent float64 // 使用率（百分比，保留两位小数）
}

// GetUptime 返回当前系统的开机时长。
func GetUptime() (time.Duration, error) {
	sys := runtime.GOOS
	switch sys {
	case "linux":
		return getUptimeLinux()
	case "windows":
		return getUptimeWindows()
	case "darwin":
		return getUptimeDarwin()
	default:
		return 0, fmt.Errorf("unsupported platform: %s", sys)
	}
}

// GetMemoryUsage 返回物理内存使用情况（已使用大小与使用率）。
func GetMemoryUsage() (UsageStat, error) {
	sys := runtime.GOOS
	switch sys {
	case "linux":
		return getMemoryUsageLinux()
	case "windows":
		return getMemoryUsageWindows()
	case "darwin":
		return getMemoryUsageDarwin()
	default:
		return UsageStat{}, fmt.Errorf("unsupported platform: %s", sys)
	}
}

// GetDiskUsage 返回系统盘使用情况（已使用大小与使用率）：
// Windows 为系统盘（SystemDrive），Linux 为根分区，macOS 为根分区（含数据卷）。
func GetDiskUsage() (UsageStat, error) {
	sys := runtime.GOOS
	switch sys {
	case "linux":
		return getDiskUsageUnix("/")
	case "windows":
		return getDiskUsageWindows()
	case "darwin":
		// APFS 下系统卷与数据卷共享容器容量，只有两者一并统计才能反映真实占用；
		// 旧版本 macOS 不存在数据卷，df 报错时回退为仅统计根分区。
		if stat, err := getDiskUsageUnix("/", "/System/Volumes/Data"); err == nil {
			return stat, nil
		}
		return getDiskUsageUnix("/")
	default:
		return UsageStat{}, fmt.Errorf("unsupported platform: %s", sys)
	}
}

// getUptimeLinux 读取 /proc/uptime 的第一列（开机以来的秒数）得到开机时长。
func getUptimeLinux() (time.Duration, error) {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0, err
	}
	fields := strings.Fields(string(data))
	if len(fields) == 0 {
		return 0, errors.New("unexpected /proc/uptime format")
	}
	sec, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, fmt.Errorf("parse /proc/uptime failed: %w", err)
	}
	return time.Duration(sec * float64(time.Second)), nil
}

// winUptimeQuery WMI 开机时长查询脚本，输出开机以来的整数秒数。
const winUptimeQuery = `$ErrorActionPreference='SilentlyContinue'; ` +
	`$b=(Get-CimInstance Win32_OperatingSystem).LastBootUpTime; ` +
	`if($b){[long]((Get-Date)-$b).TotalSeconds}`

// getUptimeWindows 通过 WMI 的最近启动时间（LastBootUpTime）计算开机时长。
func getUptimeWindows() (time.Duration, error) {
	out, err := runCommand("powershell", "-NoProfile", "-NonInteractive", "-Command", winUptimeQuery)
	if err != nil {
		return 0, fmt.Errorf("query windows uptime failed: %w", err)
	}
	sec, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse windows uptime failed: %w", err)
	}
	return time.Duration(sec) * time.Second, nil
}

// bootTimeRe 从 sysctl kern.boottime 输出（如 "{ sec = 1757836800, usec = 123456 } ..."）中提取开机时间戳。
var bootTimeRe = regexp.MustCompile(`sec = (\d+)`)

// getUptimeDarwin 读取 sysctl kern.boottime 的开机时间戳计算开机时长。
func getUptimeDarwin() (time.Duration, error) {
	out, err := runCommand("sysctl", "-n", "kern.boottime")
	if err != nil {
		return 0, fmt.Errorf("query macos boot time failed: %w", err)
	}
	m := bootTimeRe.FindStringSubmatch(out)
	if len(m) < 2 {
		return 0, errors.New("no boot time found in sysctl kern.boottime output")
	}
	sec, err := strconv.ParseInt(m[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse macos boot time failed: %w", err)
	}
	uptime := time.Since(time.Unix(sec, 0))
	if uptime < 0 {
		return 0, errors.New("invalid boot time in the future")
	}
	return uptime, nil
}

// getMemoryUsageLinux 读取 /proc/meminfo 计算内存使用情况：
// 可用内存优先取 MemAvailable，缺失（老内核）时回退 MemFree。
func getMemoryUsageLinux() (UsageStat, error) {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return UsageStat{}, err
	}
	var totalKB, freeKB, availableKB uint64
	for _, line := range strings.Split(string(data), "\n") {
		switch {
		case strings.HasPrefix(line, "MemTotal:"):
			totalKB = meminfoKB(line)
		case strings.HasPrefix(line, "MemFree:"):
			freeKB = meminfoKB(line)
		case strings.HasPrefix(line, "MemAvailable:"):
			availableKB = meminfoKB(line)
		}
	}
	if totalKB == 0 {
		return UsageStat{}, errors.New("no valid MemTotal found in /proc/meminfo")
	}
	free := availableKB
	if free == 0 {
		free = freeKB
	}
	return newUsageStat(totalKB*1024, free*1024), nil
}

// meminfoKB 从 /proc/meminfo 数据行（如 "MemTotal:   16316484 kB"）中解析第二列的数值（kB）。
func meminfoKB(line string) uint64 {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return 0
	}
	v, err := strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		return 0
	}
	return v
}

// winMemQuery WMI 内存查询脚本，输出 "总内存 空闲内存"（单位 kB）。
const winMemQuery = `$ErrorActionPreference='SilentlyContinue'; ` +
	`$os=Get-CimInstance Win32_OperatingSystem; ` +
	`if($os){"$($os.TotalVisibleMemorySize) $($os.FreePhysicalMemory)"}`

// getMemoryUsageWindows 通过 WMI 查询物理内存使用情况（WMI 数值单位为 kB）。
func getMemoryUsageWindows() (UsageStat, error) {
	out, err := runCommand("powershell", "-NoProfile", "-NonInteractive", "-Command", winMemQuery)
	if err != nil {
		return UsageStat{}, fmt.Errorf("query windows memory failed: %w", err)
	}
	fields := strings.Fields(out)
	if len(fields) < 2 {
		return UsageStat{}, errors.New("no memory info found in wmi output")
	}
	totalKB, err := strconv.ParseUint(fields[0], 10, 64)
	if err != nil {
		return UsageStat{}, fmt.Errorf("parse total memory failed: %w", err)
	}
	freeKB, err := strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		return UsageStat{}, fmt.Errorf("parse free memory failed: %w", err)
	}
	return newUsageStat(totalKB*1024, freeKB*1024), nil
}

// vmStatPageSizeRe 从 vm_stat 输出首行提取页大小（字节）。
var vmStatPageSizeRe = regexp.MustCompile(`page size of (\d+) bytes`)

// getMemoryUsageDarwin 读取 hw.memsize 与 vm_stat 计算内存使用情况：
// 已使用大小与活动监视器口径一致（active + wired + compressor 页），
// 统计字段缺失时回退为 total - free。
func getMemoryUsageDarwin() (UsageStat, error) {
	totalOut, err := runCommand("sysctl", "-n", "hw.memsize")
	if err != nil {
		return UsageStat{}, fmt.Errorf("query macos memory size failed: %w", err)
	}
	total, err := strconv.ParseUint(strings.TrimSpace(totalOut), 10, 64)
	if err != nil {
		return UsageStat{}, fmt.Errorf("parse macos memory size failed: %w", err)
	}

	vmOut, err := runCommand("vm_stat")
	if err != nil {
		return UsageStat{}, fmt.Errorf("query macos vm_stat failed: %w", err)
	}
	pageSize := uint64(4096) // 解析不到页大小时按 4K 页回退
	if m := vmStatPageSizeRe.FindStringSubmatch(vmOut); len(m) == 2 {
		if v, err := strconv.ParseUint(m[1], 10, 64); err == nil {
			pageSize = v
		}
	}
	var active, wired, compressed, free uint64
	for _, line := range strings.Split(vmOut, "\n") {
		switch {
		case strings.HasPrefix(line, "Pages active:"):
			active = vmStatPages(line)
		case strings.HasPrefix(line, "Pages wired down:"):
			wired = vmStatPages(line)
		case strings.HasPrefix(line, "Pages occupied by compressor:"):
			compressed = vmStatPages(line)
		case strings.HasPrefix(line, "Pages free:"):
			free = vmStatPages(line)
		}
	}

	usedPages := active + wired + compressed
	var used uint64
	if usedPages > 0 {
		used = usedPages * pageSize
	} else {
		// 旧版本 vm_stat 缺少统计字段时回退为 total - free
		freeBytes := free * pageSize
		if freeBytes > total {
			freeBytes = total
		}
		used = total - freeBytes
	}
	if used > total {
		used = total
	}
	return UsageStat{Total: total, Used: used, UsedPercent: usedPercent(used, total)}, nil
}

// vmStatPages 从 vm_stat 数据行（如 "Pages active: 1102979."）中解析页面数。
func vmStatPages(line string) uint64 {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return 0
	}
	v, err := strconv.ParseUint(strings.TrimSuffix(fields[len(fields)-1], "."), 10, 64)
	if err != nil {
		return 0
	}
	return v
}

// getDiskUsageUnix 通过 df -Pk 统计给定挂载点的磁盘使用情况：
// 总大小取各挂载点的最大容量（APFS 共享容器下各卷容量相同），已使用大小为各挂载点之和。
func getDiskUsageUnix(targets ...string) (UsageStat, error) {
	if len(targets) == 0 {
		return UsageStat{}, errors.New("no mount point given")
	}
	args := append([]string{"-Pk"}, targets...)
	out, err := runCommand("df", args...)
	if err != nil {
		return UsageStat{}, fmt.Errorf("df failed: %w", err)
	}
	var totalBytes, usedBytes uint64
	found := false
	for i, line := range strings.Split(out, "\n") {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue // 跳过表头与空行
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		blocks, err1 := strconv.ParseUint(fields[1], 10, 64)
		used, err2 := strconv.ParseUint(fields[2], 10, 64)
		if err1 != nil || err2 != nil {
			continue
		}
		if total := blocks * 1024; total > totalBytes {
			totalBytes = total
		}
		usedBytes += used * 1024
		found = true
	}
	if !found {
		return UsageStat{}, errors.New("no valid df output")
	}
	if usedBytes > totalBytes {
		usedBytes = totalBytes
	}
	return UsageStat{
		Total:       totalBytes,
		Used:        usedBytes,
		UsedPercent: usedPercent(usedBytes, totalBytes),
	}, nil
}

// winDiskQuery WMI 系统盘查询脚本，输出 "总大小 空闲空间"（单位字节）。
const winDiskQuery = `$ErrorActionPreference='SilentlyContinue'; ` +
	`$d=Get-CimInstance Win32_LogicalDisk -Filter "DeviceID='$env:SystemDrive'"; ` +
	`if($d){"$($d.Size) $($d.FreeSpace)"}`

// getDiskUsageWindows 通过 WMI 查询系统盘（SystemDrive，通常为 C:）使用情况。
func getDiskUsageWindows() (UsageStat, error) {
	out, err := runCommand("powershell", "-NoProfile", "-NonInteractive", "-Command", winDiskQuery)
	if err != nil {
		return UsageStat{}, fmt.Errorf("query windows system drive failed: %w", err)
	}
	fields := strings.Fields(out)
	if len(fields) < 2 {
		return UsageStat{}, errors.New("no system drive info found in wmi output")
	}
	total, err := strconv.ParseUint(fields[0], 10, 64)
	if err != nil {
		return UsageStat{}, fmt.Errorf("parse disk size failed: %w", err)
	}
	free, err := strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		return UsageStat{}, fmt.Errorf("parse disk free space failed: %w", err)
	}
	return newUsageStat(total, free), nil
}

// newUsageStat 由总大小与空闲大小（字节）组装使用量统计，空闲大于总大小时已使用计为 0。
func newUsageStat(total, free uint64) UsageStat {
	var used uint64
	if total > free {
		used = total - free
	}
	return UsageStat{
		Total:       total,
		Used:        used,
		UsedPercent: usedPercent(used, total),
	}
}

// usedPercent 返回 used/total 的百分比（保留两位小数）；total 为 0 时返回 0。
func usedPercent(used, total uint64) float64 {
	if total == 0 {
		return 0
	}
	return math.Round(float64(used)/float64(total)*10000) / 100
}
