package main

import (
	"fmt"
	"github.com/WQGroup/gonvtop/pkg/nvml"
	"time"
)

func main() {

	api, err := nvml.New("")
	if err != nil {
		panic(err)
	}

	defer api.Shutdown()

	err = api.Init()
	if err != nil {
		panic(err)
	}

	driverVersion, err := api.SystemGetDriverVersion()
	if err != nil {
		panic(err)
	}

	fmt.Println("Driver version:", driverVersion)

	nvmlVersion, err := api.SystemGetNVMLVersion()
	if err != nil {
		panic(err)
	}

	fmt.Println("NVML version:", nvmlVersion)

	deviceCount, err := api.DeviceGetCount()
	if err != nil {
		panic(err)
	}

	for i := uint32(0); i < deviceCount; i++ {
		handle, err := api.DeviceGetHandleByIndex(i)
		if err != nil {
			panic(err)
		}

		name, err := api.DeviceGetName(handle)
		fmt.Println("Product name:", name)

		brand, err := api.DeviceGetBrand(handle)
		if err != nil {
			panic(err)
		}

		fmt.Println("Product Brand:", brand)

		major, minor, err := api.DeviceGetCudaComputeCapability(handle)
		if err != nil {
			panic(err)
		}
		fmt.Println("CUDA Compute Capability:", major, minor)

		uuid, err := api.DeviceGetUUID(handle)
		if err != nil {
			panic(err)
		}

		fmt.Println("GPU UUID:", uuid)

		fan, err := api.DeviceGetFanSpeed(handle)
		if err != nil {
			panic(err)
		}

		fmt.Println("Fan Speed:", fan)

		powerLimit, err := api.DeviceGetPowerManagementLimit(handle)
		if err != nil {
			panic(err)
		}
		powerUsage, err := api.DeviceGetPowerUsage(handle)
		if err != nil {
			panic(err)
		}

		fmt.Println("Power Management:", powerUsage, powerLimit)

		// 获取当前时间的纳秒级别时间戳
		nanoTime := time.Now().UnixNano()
		// 将纳秒时间戳转换为微秒，并减去1000000微秒
		microTime := (nanoTime / 1000) - 1000000
		lastSeenTimeStamp := uint64(microTime)
		processUtilization, err := api.DeviceGetProcessUtilization(handle, lastSeenTimeStamp)
		if err != nil {
			panic(err)
		}
		for i2, utilizationSample := range processUtilization {

			if utilizationSample.Pid == 0 {
				continue
			}
			fmt.Println("Process:", i2)
			fmt.Println("\tPID:", utilizationSample.Pid)
			fmt.Println("\tSM:", utilizationSample.SmUtil)
			fmt.Println("\tMem:", utilizationSample.MemUtil)
			fmt.Println("\tEnc:", utilizationSample.EncUtil)
			fmt.Println("\tDec:", utilizationSample.DecUtil)
			fmt.Println("\tTimeStamp:", utilizationSample.TimeStamp)
		}

		rates, err := api.DeviceGetUtilizationRates(handle)
		if err != nil {
			panic(err)
		}
		fmt.Println("Utilization GPU Rates:\t", rates.GPU)
		fmt.Println("Utilization Memory Rates:\t", rates.Memory)

		memoryInfo, err := api.DeviceGetMemoryInfo(handle)
		if err != nil {
			panic(err)
		}
		fmt.Println("Memory Info Free:", memoryInfo.Free)
		fmt.Println("Memory Info Total:", memoryInfo.Total)
		fmt.Println("Memory Info Used:", memoryInfo.Used)

		computeRunningProcesses, err := api.DeviceGetComputeRunningProcesses(handle)
		if err != nil {
			panic(err)
		}
		for i2, process := range computeRunningProcesses {
			fmt.Println("Process:", i2)
			fmt.Println("\tPID:", process.PID)
			fmt.Println("\tUsed GPU Memory:", process.UsedGPUMemory)
		}
		//processes, err := api.DeviceGetGraphicsRunningProcesses(handle)
		//if err != nil {
		//	panic(err)
		//}
		//for i2, process := range processes {
		//	fmt.Println("Process :", i2)
		//	fmt.Println("\tPID:", process.PID)
		//	fmt.Println("\tUsed GPU Memory:", process.UsedGPUMemory)
		//}
	}

}
