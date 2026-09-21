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

type MailContent struct {
	ContentType string `json:"content-type"`
	Data        string `json:"data"`
}

type MailFull struct {
	Id      int64         `json:"id"`
	Read    bool          `json:"read"`
	Address string        `json:"address"`
	Subject string        `json:"subject"`
	Date    string        `json:"date"`
	Content []MailContent `json:"content"`
}
