package info_hub

import (
	"fmt"
	"github.com/WQGroup/gonvtop/pkg/nvml"
	"github.com/shirou/gopsutil/process"
)

type GPUDriverInfos struct {
	DriverVersion     string `json:"driver_version"`      // 显卡驱动版本
	NVMLVersion       string `json:"nvml_version"`        // NVML 版本
	CudaDriverVersion int32  `json:"cuda_driver_version"` // CUDA 版本
}

func NewGPUDriverInfos(driverVersion, nVMLVersion string, cudaDriverVersion int32) *GPUDriverInfos {
	return &GPUDriverInfos{
		DriverVersion:     driverVersion,
		NVMLVersion:       nVMLVersion,
		CudaDriverVersion: cudaDriverVersion,
	}
}

func (v GPUDriverInfos) GetDriverVersion() string {
	return v.DriverVersion
}

func (v GPUDriverInfos) GetNVMLVersion() string {
	return v.NVMLVersion
}

func (v GPUDriverInfos) GetCUDADriverVersion() string {
	return fmt.Sprintf("%d.%d", v.CudaDriverVersion/1000, (v.CudaDriverVersion%100)/10)
}

// -----------------------

type GPUInfos struct {
	Index             uint32                     `json:"index"`              // 索引
	Name              string                     `json:"name"`               // 名称
	BrandType         nvml.BrandType             `json:"brand_type"`         // 型号分类
	UUID              string                     `json:"uuid"`               // UUID
	Fan               uint32                     `json:"fan"`                // 风扇速度的百分比，满速是 100%
	Temperature       uint32                     `json:"temperature"`        // 温度,C
	UtilizationRates  *nvml.Utilization          `json:"utilization_rates"`  // 利用率（GPU and 显存）
	Memory            *nvml.Memory               `json:"memory"`             // 显存使用信息
	Power             *PowerInfo                 `json:"power"`              // 电源信息
	ComputeCapability *ComputeCapabilityInfo     `json:"compute_capability"` // CUDA 计算能力版本
	Processes         map[uint32]*GPUProcessInfo `json:"processes"`          // 进程信息
}

func (g GPUInfos) GetComputeCapability() string {
	return g.ComputeCapability.Version()
}

// -----------------------

type ComputeCapabilityInfo struct {
	Major uint32 `json:"major"` // 主版本号
	Minor uint32 `json:"minor"` // 次版本号
}

func NewComputeCapabilityInfo(major uint32, minor uint32) *ComputeCapabilityInfo {
	return &ComputeCapabilityInfo{Major: major, Minor: minor}
}

func (c ComputeCapabilityInfo) Version() string {
	return fmt.Sprintf("%d.%d", c.Major, c.Minor)
}

// -----------------------

type PowerInfo struct {
	Usage uint32 `json:"usage"` // 电源当前的功耗
	Limit uint32 `json:"limit"` // 电源的当前最大的限制功耗
}

func NewPowerInfo(usage uint32, limit uint32) *PowerInfo {
	return &PowerInfo{Usage: usage, Limit: limit}
}

// -----------------------

type GPUProcessInfo struct {
	Name    string                        `json:"name"`     // 进程名称
	Pid     uint32                        `json:"pid"`      // 进程 ID
	USample nvml.ProcessUtilizationSample `json:"u_sample"` // 进程利用率
}

func NewGPUProcessInfo(name string, Pid uint32, uSample nvml.ProcessUtilizationSample) *GPUProcessInfo {
	return &GPUProcessInfo{Name: name, Pid: Pid, USample: uSample}
}

func (p GPUProcessInfo) GetPID() uint32 {
	return p.USample.Pid
}

func (p GPUProcessInfo) GetName() string {
	return p.Name
}

func (p GPUProcessInfo) GetSmUtil() uint32 {
	return p.USample.SmUtil
}

func (p GPUProcessInfo) GetMemoryUtil() uint32 {
	return p.USample.MemUtil
}

// -----------------------

type HostSystemInfos struct {
	CPUPercent float64 `json:"cpu_percent"` // CPU 利用率
	Memory     *Memory `json:"memory"`      // 内存使用信息
}

func NewHostSystemInfos(CPUPercent float64, Memory *Memory) *HostSystemInfos {
	return &HostSystemInfos{CPUPercent: CPUPercent, Memory: Memory}
}

type Memory struct {
	// Total amount of RAM on this system
	Total uint64 `json:"total"`

	// RAM available for programs to allocate
	//
	// This value is computed from the kernel specific values.
	Available uint64 `json:"available"`

	// RAM used by programs
	//
	// This value is computed from the kernel specific values.
	Used uint64 `json:"used"`

	// Percentage of RAM used by programs
	//
	// This value is computed from the kernel specific values.
	UsedPercent float64 `json:"usedPercent"`

	// This is the kernel's notion of free memory; RAM chips whose bits nobody
	// cares about the value of right now. For a human consumable number,
	// Available is what you really want.
	Free uint64 `json:"free"`
}

func NewMemory(total uint64, available uint64, used uint64, usedPercent float64, free uint64) *Memory {
	return &Memory{Total: total, Available: available, Used: used, UsedPercent: usedPercent, Free: free}
}

// ----------------------------

type ProcessInfo struct {
	PID        uint32                          `json:"pid"` // 进程 ID
	Name       string                          `json:"name"`
	Environ    []string                        `json:"environ"`
	Cmdline    string                          `json:"cmd_line"`
	CpuPercent float64                         `json:"cpu_percent"`
	MemPercent float32                         `json:"mem_percent"`
	MemInfo    *process.MemoryInfoStat         `json:"mem_info"`
	GPUUSample []nvml.ProcessUtilizationSample `json:"gpu_u_sample"`
}

func NewProcessInfo(PID uint32, name, Cmdline string, environ []string, cpuPercent float64, memPercent float32,
	memInfo *process.MemoryInfoStat, gpuCounter int) *ProcessInfo {
	return &ProcessInfo{PID: PID, Name: name, Cmdline: Cmdline, Environ: environ, CpuPercent: cpuPercent, MemPercent: memPercent,
		MemInfo: memInfo, GPUUSample: make([]nvml.ProcessUtilizationSample, gpuCounter)}
}
