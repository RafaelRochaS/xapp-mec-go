package models

type RegisterRequest struct {
	DeviceId     int    `json:"deviceId"`
	Task         string `json:"task"`
	Image        string `json:"image"`
	CPU          string `json:"cpu"`
	Mem          int    `json:"mem"`
	DeadlineSecs int    `json:"deadlineSecs,omitempty"`
	Workload     int    `json:"workload"`
	CallbackUrl  string `json:"callbackUrl,omitempty"`
}
