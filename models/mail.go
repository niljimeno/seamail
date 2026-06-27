package models

type Mail struct {
	Id   int64 `json:"id"`
	Read bool  `json:"read"`
}

type MailDetailed struct {
	Id      int64  `json:"id"`
	Read    bool   `json:"read"`
	Address string `json:"address"`
	Subject string `json:"subject"`
}
