package main

import (
	"github.com/WQGroup/gonvtop/pkg/info_hub"
	"github.com/WQGroup/logger"
	"time"
)

func main() {

	infoHub := info_hub.NewInfoHub("")
	defer infoHub.Close()

	println("Driver version:", infoHub.GPUDriverInfos.GetDriverVersion())
	println("NVML version:", infoHub.GPUDriverInfos.GetNVMLVersion())
	println("CUDA version:", infoHub.GPUDriverInfos.GetCUDADriverVersion())

	for true {
		err := infoHub.Refresh()
		if err != nil {
			logger.Panicln(err)
		}
		println("GPU count:", len(infoHub.GPUInfos))

		for _, gpuInfos := range infoHub.GPUInfos {
			println("------------------------")
			println("Utilization: GPU - Memory:", gpuInfos.UtilizationRates.GPU, gpuInfos.UtilizationRates.Memory)
			for _, process := range gpuInfos.Processes {
				println("\tName:", process.GetName())
				println("\tPID:", process.GetPID())
				println("\tGPU Util:", process.USample.SmUtil)
				println("\tGPU Memory:", process.USample.MemUtil)
			}
		}

		time.Sleep(1 * time.Second)
	}
}
