package hardware

import (
	"context"
	"errors"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

const (
	sysInfoFetchInterval = time.Second
	// diskFetchInterval 系统盘使用情况的采样间隔：磁盘变化缓慢且查询较重（尤其 Windows WMI），
	// 低频采样避免拖慢整体采样周期。
	diskFetchInterval = 30 * time.Second
)

var (
	// HS 全局唯一的硬件采样服务；最后一个调用方释放后由 release 置回 nil。
	HS *HardwareService
	mu sync.Mutex
)

// HardwareService 全局单例：所有调用方共享同一实例，只由内部协程读取系统信息；
// 引用计数归零后自动停止采样并释放全局实例。
type HardwareService struct {
	sampleMu   sync.RWMutex // 保护最近一次采样结果，避免读写数据竞争
	cpuThermal float64
	cpuErr     error         // 最近一次 CPU 温度采样错误
	hasCpuTemp bool          // CPU 温度是否已完成首次采样
	sysInfo    SysInfoSample // 最近一次系统信息采样结果（开机时长 / 内存 / 磁盘）
	hasSysInfo bool          // 系统信息是否已完成首轮采样
	lastDiskAt time.Time     // 上一次磁盘采样时间（磁盘按 diskFetchInterval 低频采样）
	ticker     *time.Ticker
	done       chan struct{}
	stopOnce   sync.Once
	count      atomic.Int64 // 当前活跃调用方数量
}

// SysInfoSample 系统信息的一次采样结果（开机时长、内存与系统盘使用情况）。
// 某项采样成功时携带数值；失败时对应 Err 非空，该项数值无意义。
type SysInfoSample struct {
	Uptime    time.Duration // 开机时长
	UptimeErr error         // 开机时长采样错误
	Memory    UsageStat     // 内存使用情况
	MemoryErr error         // 内存采样错误
	Disk      UsageStat     // 系统盘使用情况
	DiskErr   error         // 磁盘采样错误
}

// NewHardwareService 获取全局唯一的 HardwareService 并将引用计数加一；
// 传入的 ctx 结束时自动释放本次引用，当最后一个调用方结束时实例停止采样、
// HS 置为 nil，之后再次调用会创建新的实例。
func NewHardwareService(ctx context.Context) *HardwareService {
	mu.Lock()
	if HS == nil {
		HS = &HardwareService{
			done:   make(chan struct{}),
			ticker: time.NewTicker(sysInfoFetchInterval),
		}
		HS.run() // 只有一处读取系统信息，避免重复
	}
	HS.count.Add(1)
	h := HS
	mu.Unlock()

	context.AfterFunc(ctx, h.release)
	return h
}

// release 释放一次引用；计数归零时停止采样并将全局实例置为 nil。
func (h *HardwareService) release() {
	mu.Lock()
	defer mu.Unlock()
	if h.count.Add(-1) > 0 {
		return
	}
	h.stop()
	if HS == h {
		HS = nil
	}
}

// stop 通知采样协程退出，重复调用安全。
func (h *HardwareService) stop() {
	h.stopOnce.Do(func() {
		close(h.done)
	})
}

func (h *HardwareService) run() {
	go func() {
		defer h.ticker.Stop()
		h.getSystemInfo() // 建立后立即采样一次，避免调用方等待首个数据
		for {
			select {
			case <-h.ticker.C:
				h.getSystemInfo()
			case <-h.done:
				return
			}
		}
	}()
}

// GetInfo 返回最近一次采样到的 CPU 温度（摄氏度，保留两位小数）；
// 首个采样完成前或采样失败时返回错误。
func (h *HardwareService) GetInfo() (float64, error) {
	h.sampleMu.RLock()
	defer h.sampleMu.RUnlock()
	if !h.hasCpuTemp {
		return 0, errors.New("cpu temperature sample not ready")
	}
	return h.cpuThermal, h.cpuErr
}

// GetSysInfo 返回最近一次采样到的系统信息（开机时长、内存与系统盘使用情况）；
// 首轮采样完成前返回错误；就绪后各项采样错误通过 Sample 内对应的 Err 字段暴露。
func (h *HardwareService) GetSysInfo() (SysInfoSample, error) {
	h.sampleMu.RLock()
	defer h.sampleMu.RUnlock()
	if !h.hasSysInfo {
		return SysInfoSample{}, errors.New("sys info sample not ready")
	}
	return h.sysInfo, nil
}

func (h *HardwareService) getSystemInfo() {
	h.getCPUTemp()
	h.getUptime()
	h.getMemoryUsage()
	h.getDiskUsage()
	// 整轮采样结束后标记系统信息就绪，避免调用方读到尚未赋值的字段
	h.sampleMu.Lock()
	h.hasSysInfo = true
	h.sampleMu.Unlock()
}

func (h *HardwareService) getCPUTemp() {
	temp, err := GetCPUTemp()
	temp = math.Round(temp*100) / 100 // 摄氏度保留两位小数
	h.sampleMu.Lock()
	h.cpuThermal, h.cpuErr, h.hasCpuTemp = temp, err, true
	h.sampleMu.Unlock()
}

// getUptime 采样一次开机时长并写入最近一次结果。
func (h *HardwareService) getUptime() {
	uptime, err := GetUptime()
	h.sampleMu.Lock()
	h.sysInfo.Uptime, h.sysInfo.UptimeErr = uptime, err
	h.sampleMu.Unlock()
}

// getMemoryUsage 采样一次内存使用情况并写入最近一次结果。
func (h *HardwareService) getMemoryUsage() {
	memory, err := GetMemoryUsage()
	h.sampleMu.Lock()
	h.sysInfo.Memory, h.sysInfo.MemoryErr = memory, err
	h.sampleMu.Unlock()
}

// getDiskUsage 采样一次系统盘使用情况并写入最近一次结果；
// 距上次采样不足 diskFetchInterval 时跳过，避免高频磁盘查询拖慢采样周期。
func (h *HardwareService) getDiskUsage() {
	h.sampleMu.RLock()
	last := h.lastDiskAt
	h.sampleMu.RUnlock()
	if !last.IsZero() && time.Since(last) < diskFetchInterval {
		return
	}
	disk, err := GetDiskUsage()
	h.sampleMu.Lock()
	h.sysInfo.Disk, h.sysInfo.DiskErr = disk, err
	h.lastDiskAt = time.Now()
	h.sampleMu.Unlock()
}
