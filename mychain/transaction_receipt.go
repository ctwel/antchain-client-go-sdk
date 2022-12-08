package mychain

type TransactionReceipt struct {
	Result  int64  `json:"result,omitempty"`
	GasUsed int64  `json:"gasUsed,omitempty"`
	Output  string `json"output,omitempty"`
}
