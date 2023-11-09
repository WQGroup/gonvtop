package info_hub

import (
	"github.com/WQGroup/gonvtop/pkg/nvml"
	"github.com/WQGroup/logger"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"sync"
	"time"
)

type InfoHub struct {
	api             *nvml.API
	HostSystemInfos *HostSystemInfos
	GPUDriverInfos  *GPUDriverInfos
	GPUInfos        []*GPUInfos // UUID -> GPUInfos
	locker4Refresh  sync.Mutex

	cacheInfo    CacheInfo
	locker4Cache sync.Mutex
}

func NewInfoHub(nvmlDllPath string) *InfoHub {

	api, err := nvml.New(nvmlDllPath)
	if err != nil {
		logger.Panicln(err)
	}
	err = api.Init()
	if err != nil {
		logger.Panicln(err)
	}

	driverVersion, err := api.SystemGetDriverVersion()
	if err != nil {
		logger.Panicln(err)
	}

	nvmlVersion, err := api.SystemGetNVMLVersion()
	if err != nil {
		logger.Panicln(err)
	}

	cudaDriverVersion, err := api.SystemGetCudaDriverVersion()
	if err != nil {
		logger.Panicln(err)
	}

	systemInfos := NewGPUDriverInfos(driverVersion, nvmlVersion, cudaDriverVersion)

	deviceCount, err := api.DeviceGetCount()
	if err != nil {
		logger.Panicln(err)
	}

	return &InfoHub{
		api:            api,
		GPUDriverInfos: systemInfos,
		GPUInfos:       make([]*GPUInfos, deviceCount),
	}
}

func (i *InfoHub) Close() {
	if i.api != nil {
		_ = i.api.Shutdown()
	}
}

func (i *InfoHub) Monitor(sleepTime time.Duration) {

	for {
		err := i.refreshBaseInfo()
		if err != nil {
			logger.Panicln(err)
		}
		time.Sleep(sleepTime)
	}
}

func (i *InfoHub) Refresh() error {
	return i.refreshBaseInfo()
}

func (i *InfoHub) GetCacheInfo() CacheInfo {
	i.locker4Cache.Lock()
	defer i.locker4Cache.Unlock()
	return i.cacheInfo
}

func (i *InfoHub) getHostSystemInfo() error {

	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return err
	}
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	nowMem := NewMemory(memInfo.Total, memInfo.Available, memInfo.Used, memInfo.UsedPercent, memInfo.Free)

	i.HostSystemInfos = NewHostSystemInfos(percent[0], nowMem)
	return nil
}

func (i *InfoHub) refreshBaseInfo() error {

	i.locker4Refresh.Lock()
	defer func() {

		i.locker4Cache.Lock()
		i.cacheInfo = *NewCacheInfo(*i.HostSystemInfos, *i.GPUDriverInfos, i.GPUInfos)
		i.locker4Cache.Unlock()

		i.locker4Refresh.Unlock()
	}()

	err := i.getHostSystemInfo()
	if err != nil {
		return err
	}

	deviceCount, err := i.api.DeviceGetCount()
	if err != nil {
		return err
	}

	for j := uint32(0); j < deviceCount; j++ {

		var handle nvml.Device
		var name, uuid string
		var major, minor int32
		var brand nvml.BrandType
		var fan, powerLimit, powerUsage, temperature uint32
		var memoryInfo nvml.Memory
		var utilizationRates nvml.Utilization
		var processUtilizations []nvml.ProcessUtilizationSample

		handle, err = i.api.DeviceGetHandleByIndex(j)
		if err != nil {
			return err
		}

		name, err = i.api.DeviceGetName(handle)
		if err != nil {
			return err
		}

		brand, err = i.api.DeviceGetBrand(handle)
		if err != nil {
			return err
		}

		major, minor, err = i.api.DeviceGetCudaComputeCapability(handle)
		if err != nil {
			return err
		}

		uuid, err = i.api.DeviceGetUUID(handle)
		if err != nil {
			return err
		}

		fan, err = i.api.DeviceGetFanSpeed(handle)
		if err != nil {
			return err
		}
		powerLimit, err = i.api.DeviceGetPowerManagementLimit(handle)
		if err != nil {
			return err
		}
		powerUsage, err = i.api.DeviceGetPowerUsage(handle)
		if err != nil {
			return err
		}

		memoryInfo, err = i.api.DeviceGetMemoryInfo(handle)
		if err != nil {
			return err
		}

		temperature, err = i.api.DeviceGetTemperature(handle, nvml.TemperatureGPU)
		if err != nil {
			return err
		}

		utilizationRates, err = i.api.DeviceGetUtilizationRates(handle)
		if err != nil {
			return err
		}

		// 获取当前时间的纳秒级别时间戳
		nanoTime := time.Now().UnixNano()
		// 将纳秒时间戳转换为微秒，并减去1000000微秒
		microTime := (nanoTime / 1000) - 1000000
		lastSeenTimeStamp := uint64(microTime)
		processUtilizations, err = i.api.DeviceGetProcessUtilization(handle, lastSeenTimeStamp)
		if err != nil {
			return err
		}

		nowGPUInfo := i.GPUInfos[j]
		if nowGPUInfo == nil {
			// 新建
			nowGPUInfo = &GPUInfos{
				Index:             j,
				Name:              name,
				BrandType:         brand,
				UUID:              uuid,
				Fan:               fan,
				Temperature:       temperature,
				UtilizationRates:  &utilizationRates,
				Memory:            &memoryInfo,
				Power:             NewPowerInfo(powerLimit, powerUsage),
				ComputeCapability: NewComputeCapabilityInfo(uint32(major), uint32(minor)),
				Processes:         make(map[uint32]*ProcessInfo),
			}
			i.GPUInfos[j] = nowGPUInfo

		} else {
			// 更新已有的
			nowGPUInfo.Fan = fan
			nowGPUInfo.Temperature = temperature
			nowGPUInfo.UtilizationRates = &utilizationRates
			nowGPUInfo.Memory = &memoryInfo
			nowGPUInfo.Power = NewPowerInfo(powerLimit, powerUsage)
			nowGPUInfo.Processes = make(map[uint32]*ProcessInfo)
		}
		// 更新进程信息
		for _, processUtilization := range processUtilizations {

			if processUtilization.Pid == 0 {
				continue
			}
			pName, err := i.api.SystemGetProcessName(uint(processUtilization.Pid))
			if err != nil {
				continue
			}
			// 不管如何都要更新
			nowGPUInfo.Processes[processUtilization.Pid] = NewProcessInfo(pName, processUtilization)
		}
	}

	return err
}

type CacheInfo struct {
	HostSystemInfos HostSystemInfos `json:"host_system_infos"`
	GPUDriverInfos  GPUDriverInfos  `json:"gpu_driver_infos"`
	GPUInfos        []*GPUInfos     `json:"gpu_infos"` // UUID -> GPUInfos
}

func NewCacheInfo(hostSystemInfos HostSystemInfos, GPUDriverInfos GPUDriverInfos, GPUInfos []*GPUInfos) *CacheInfo {
	return &CacheInfo{HostSystemInfos: hostSystemInfos, GPUDriverInfos: GPUDriverInfos, GPUInfos: GPUInfos}
}
