package hardware

import (
	"fmt"
	"runtime"
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
