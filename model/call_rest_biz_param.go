package model

type CallRestBizParam struct {
	BaseParam
	OrderId            string     `json:"orderId,omitempty"`
	Account            string     `json:"account,omitempty"`
	Content            string     `json:"content,omitempty"`
	TenantId           string     `json:"tenantid,omitempty"`
	Uid                string     `json:"uid,omitempty"`
	ContractName       string     `json:"contractName,omitempty"`
	ContractCode       string     `json:"contractCode,omitempty"`
	OutTypes           string     `json:"outTypes,omitempty"`
	MethodSignature    string     `json:"methodSignature,omitempty"`
	InputParamListStr  string     `json:"inputParamListStr,omitempty"`
	NativeContractData string     `json:"nativeContractData,omitempty"`
	MykmsKeyId         string     `json:"mykmsKeyId,omitempty"`
	BlockNumber        int64      `json:"blockNumber,omitempty"`
	IsLocalTransaction bool       `json:"isLocalTransaction,omitempty"`
	ApplyAccessKey     string     `json:"applyAccessKey,omitempty"`
	Gas                int64      `json:"gas,omitempty"`
	VmTypeEnum         VMTypeEnum `json:"vmTypeEnum,omitempty"`
	Abi                string     `json:"abi,omitempty"`
}
