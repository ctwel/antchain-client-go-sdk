package response

type BaseResp struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Data    string `json:"data"`
}
