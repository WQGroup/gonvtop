package info_hub

import (
	"testing"
	"time"
)

func TestInfoHub_Refresh(t *testing.T) {

	infoHub := NewInfoHub("")
	defer infoHub.Close()

	println("Driver version:", infoHub.GPUDriverInfos.GetDriverVersion())
	println("NVML version:", infoHub.GPUDriverInfos.GetNVMLVersion())
	println("CUDA version:", infoHub.GPUDriverInfos.GetCUDADriverVersion())

	for true {
		err := infoHub.Refresh()
		if err != nil {
			t.Fatal(err)
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

func BenchmarkInfoHub_Refresh(b *testing.B) {

	infoHub := NewInfoHub("")
	defer infoHub.Close()
	err := infoHub.Refresh()
	if err != nil {
		b.Fatal(err)
	}

	for _, gpuInfos := range infoHub.GPUInfos {
		for _, process := range gpuInfos.Processes {
			println("------------------------")
			println("Name:", process.GetName())
			println("PID:", process.GetPID())
			println("GPU Util:", process.GetSmUtil())
			println("GPU Memory:", process.GetMemoryUtil())
		}
	}
}
