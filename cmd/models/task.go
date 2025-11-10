package models

type Status int

const (
	Idle Status = iota
	InProgress
	Failed
	Complete
)

func (s Status) String() string {
	switch s {
	case Idle:
		return "idle"
	case InProgress:
		return "in_progress"
	case Failed:
		return "failed"
	case Complete:
		return "complete"
	default:
		return "unknown"
	}
}

type Task struct {
	Id     string `json:"id"`
	Status Status `json:"status"`
	RegisterRequest
}

type StartTaskRequest struct {
	Id       string `json:"id"`
	DeviceId string `json:"deviceId"`
}
