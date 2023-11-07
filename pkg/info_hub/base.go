package info_hub

import (
	"fmt"
	"github.com/WQGroup/gonvtop/pkg/nvml"
)

type SystemInfos struct {
	driverVersion     string // 显卡驱动版本
	nVMLVersion       string // NVML 版本
	cudaDriverVersion int32  // CUDA 版本
}

func NewSystemInfos(driverVersion, nVMLVersion string, cudaDriverVersion int32) *SystemInfos {
	return &SystemInfos{
		driverVersion:     driverVersion,
		nVMLVersion:       nVMLVersion,
		cudaDriverVersion: cudaDriverVersion,
	}
}

func (v SystemInfos) GetDriverVersion() string {
	return v.driverVersion
}

func (v SystemInfos) GetNVMLVersion() string {
	return v.nVMLVersion
}

func (v SystemInfos) GetCUDADriverVersion() string {
	return fmt.Sprintf("%d.%d", v.cudaDriverVersion/1000, (v.cudaDriverVersion%100)/10)
}

// -----------------------

type GPUInfos struct {
	Index             uint32                  // 索引
	Name              string                  // 名称
	BrandType         nvml.BrandType          // 型号分类
	UUID              string                  // UUID
	Fan               uint32                  // 风扇速度的百分比，满速是 100%
	Temperature       uint32                  // 温度,C
	UtilizationRates  nvml.Utilization        // 利用率（GPU and 显存）
	Memory            nvml.Memory             // 显存使用信息
	Power             PowerInfo               // 电源信息
	computeCapability *ComputeCapabilityInfo  // CUDA 计算能力版本
	Processes         map[uint32]*ProcessInfo // 进程信息
}

func (g GPUInfos) GetComputeCapability() string {
	return g.computeCapability.Version()
}

// -----------------------

type ComputeCapabilityInfo struct {
	major uint32 // 主版本号
	minor uint32 // 次版本号
}

func NewComputeCapabilityInfo(major uint32, minor uint32) *ComputeCapabilityInfo {
	return &ComputeCapabilityInfo{major: major, minor: minor}
}

func (c ComputeCapabilityInfo) Version() string {
	return fmt.Sprintf("%d.%d", c.major, c.minor)
}

// -----------------------

type PowerInfo struct {
	Usage uint32 // 电源当前的功耗
	Limit uint32 // 电源的当前最大的限制功耗
}

func NewPowerInfo(usage uint32, limit uint32) *PowerInfo {
	return &PowerInfo{Usage: usage, Limit: limit}
}

// -----------------------

type ProcessInfo struct {
	name    string // 进程名称
	uSample nvml.ProcessUtilizationSample
}

func NewProcessInfo(name string, uSample nvml.ProcessUtilizationSample) *ProcessInfo {
	return &ProcessInfo{name: name, uSample: uSample}
}
