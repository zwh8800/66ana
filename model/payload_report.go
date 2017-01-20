package model

type ReportPayload struct {
	WorkerId string
	Capacity int
	Working  int

	Workers []*WorkerInfo

	CpuCount int
	MemUsage int64
}

type WorkerInfo struct {
	RoomId int64
	Speed  float64
}
