package model

type BaseParam struct {
	AccessId   string `json:"accessId,omitempty"`
	BizId      string `json:"bizid,omitempty"`
	Hash       string `json:"hash,omitempty"`
	Token      string `json:"token,omitempty"`
	RequestStr string `json:"requestStr,omitempty"`
	Method     Method `json:"method,omitempty"`
	SecretKey  string `json:"secretKey,omitempty"`
}
