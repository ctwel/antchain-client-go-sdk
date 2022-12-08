package config

type RestClientProperties struct {
	RestUrl           string `json:"RestUrl"`
	AccessId          string `json:"AccessId""`
	AccessSecret      string `json:"AccessSecret"`

	MaxIdleConns    int `json:"MaxIdleConns"`
	IdleConnTimeout int `json:"IdleConnTimeout"` // 单位为秒

	RetryMaxAttempts int `json:"RetryMaxAttempts"` // http.client 重试次数
	BackOffPeriod    int `json:"BackOffPeriod"`    // http.client重试间隔,单位为毫秒
}
