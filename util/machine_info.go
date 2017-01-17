package util

import "runtime"

func CpuCount() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}

func MemUsage() int64 {
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	return m.Alloc
}
