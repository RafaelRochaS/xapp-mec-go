package models

type PolicyConfig struct {
	OffloadThreshold int    `json:"offloadThreshold"`
	Algorithm        string `json:"algorithm"`
}

type AppConfig struct {
	Policy PolicyConfig `json:"policy"`
}
