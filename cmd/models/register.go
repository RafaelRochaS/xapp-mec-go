package models

type RegisterRequest struct {
	Task         string `json:"task"`
	Image        string `json:"image"`
	CPU          string `json:"cpu"`
	Mem          string `json:"mem"`
	DeadlineSecs int64  `json:"deadlineSecs"`
	CallbackUrl  string `json:"callbackUrl,omitempty"`
}
