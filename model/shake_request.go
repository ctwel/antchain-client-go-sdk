package model

type ShakeRequest struct {
	AccessId string `json:"accessId"`
	Time     string `json:"time"`
	Secret   string `json:"secret"`
}
