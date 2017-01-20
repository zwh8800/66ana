package model

type SpiderClosedPayload struct {
	WorkerId string
	RoomId   int64

	*ReportPayload
}
