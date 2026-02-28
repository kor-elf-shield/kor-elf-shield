package entity

type Blocking struct {
	IP           string `json:"IP"`
	ExpireAtUnix int64  `json:"ExpireAtUnix"`
	Reason       string `json:"Reason"`
}
