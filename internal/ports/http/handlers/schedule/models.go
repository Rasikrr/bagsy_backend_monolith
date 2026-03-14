package schedule

//go:generate easyjson -all models.go

type setScheduleRequest struct {
	Start string        `json:"start"`
	End   string        `json:"end"`
	Slots []slotRequest `json:"slots"`
}

type slotRequest struct {
	Date      string `json:"date"`
	Type      string `json:"type"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type slotResponse struct {
	ID        string `json:"id"`
	Date      string `json:"date"`
	Type      string `json:"type"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type getScheduleResponse struct {
	Slots []slotResponse `json:"slots"`
}
