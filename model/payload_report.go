package model

type BasicWorkerInfo struct {
	WorkerId string
	Capacity int
	Working  int
	CpuCount int
	MemUsage int64
}

type ReportPayload struct {
	*BasicWorkerInfo

	Workers []*WorkerInfo
}

type WorkerInfo struct {
	RoomId int64
	Speed  float64
}
