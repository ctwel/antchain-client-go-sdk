package utils

import (
	"fmt"
	"github.com/ctwel/antchain-client-go-sdk/model"
	"github.com/ctwel/antchain-client-go-sdk/response"
)

func CheckCallRestBizParams(callRestBizParam model.CallRestBizParam) response.BaseResp {
	passChecked := true
	data := ""

	if callRestBizParam.AccessId == "" {
		return response.BaseResp{Success: false, Data: "no access id"}
	}
	if callRestBizParam.Token == "" {
		return response.BaseResp{Success: false, Data: "no token"}
	}
	if callRestBizParam.BizId == "" {
		return response.BaseResp{Success: false, Data: "no bizid"}
	}
	method := callRestBizParam.Method
	if method != model.QUERYTENANTKMSLIST && method != model.DEPOSITWITHADMIN && method != model.APPLYKEY &&
		method != model.QUERYACCESSLIST && method != model.RESETAPPLYKEY && method != model.CREATEACCOUNT &&
		method != model.DEPLOYNATIVECONTRACT && method != model.QUERYACCOUNT && method != model.QUERYRECEIPT &&
		method != model.QUERYTRANSACTION && method != model.QUERYRECEIPTBIZ && method != model.QUERYTRANSACTIONBIZ &&
		method != model.FROZENTENANT && method != model.UNFROZENTENANT {
		if callRestBizParam.Uid == "" && callRestBizParam.MykmsKeyId == "" {
			return response.BaseResp{Success: false, Data: "uid or mykmsKeyId must be not null"}
		}
	}
	if method != model.APPLYKEY && method != model.QUERYACCESSLIST && method != model.RESETAPPLYKEY &&
		method != model.QUERYRECEIPT && method != model.QUERYTRANSACTION && method != model.FROZENTENANT &&
		method != model.UNFROZENTENANT &&
		callRestBizParam.OrderId == "" {
		passChecked = false
		data = fmt.Sprintf("%v method must has orderId", callRestBizParam.Method)
	}

	switch method {
	case model.DEPOSIT:
		if callRestBizParam.Account == "" {
			passChecked = false
			data = fmt.Sprintf("%v method must has account", callRestBizParam.Method)
		}
		if callRestBizParam.Content == "" {
			passChecked = false
			data = fmt.Sprintf("%v method must has content", callRestBizParam.Method)
		}
	case model.CALLCONTRACTBIZ:
		fallthrough
	case model.CALLCONTRACTBIZASYNC:
		if callRestBizParam.Account == "" {
			passChecked = false
			data = fmt.Sprintf("%v method must has account", callRestBizParam.Method)
		}
		if callRestBizParam.ContractName == "" {
			passChecked = false
			data = fmt.Sprintf("%v method must has contract name", callRestBizParam.Method)
		}
		if callRestBizParam.OutTypes == "" {
			passChecked = false
			data = fmt.Sprintf("%v method must has outTypes", callRestBizParam.Method)
		}
		if callRestBizParam.MethodSignature == "" {
			passChecked = false
			data = fmt.Sprintf("%v method must has methodSignature", callRestBizParam.Method)
		}
		if callRestBizParam.InputParamListStr == "" {
			passChecked = false
			data = fmt.Sprintf("%v method must has inputParamListStr", callRestBizParam.Method)
		}
	//case model.CALLNATIVECONTRACTFORBIZASYNC:
	//case model.CALLNATIVECONTRACTFORBIZ:
	//	if callRestBizParam.Account == "" {
	//		passChecked = false
	//		data = fmt.Sprintf("%v method must has account", callRestBizParam.Method)
	//	}
	//	if callRestBizParam.ContractName == "" {
	//		passChecked = false
	//		data = fmt.Sprintf("%v method must has contract name", callRestBizParam.Method)
	//	}
	//	if callRestBizParam.MethodSignature == "" {
	//		passChecked = false
	//		data = fmt.Sprintf("%v method must has methodSignature", callRestBizParam.Method)
	//	}
	//	if callRestBizParam.NativeContractData == "" {
	//		passChecked = false
	//		data = fmt.Sprintf("%v method must has nativeContractData", callRestBizParam.Method)
	//	}
	case model.QUERYRECEIPT:
		fallthrough
	case model.QUERYTRANSACTION:
		if callRestBizParam.Hash == "" {
			passChecked = false
			data = fmt.Sprintf("%v method must has hash", callRestBizParam.Method)
		}
	case model.DEPLOYCONTRACTFORBIZ:
		if callRestBizParam.Account == "" {
			passChecked = false
			data = fmt.Sprintf("%v method must has account", callRestBizParam.Method)
		}
		if callRestBizParam.ContractName == "" {
			passChecked = false
			data = fmt.Sprintf("%v method must has contract name", callRestBizParam.Method)
		}
		if callRestBizParam.ContractCode == "" {
			passChecked = false
			data = fmt.Sprintf("%v method must has contractCode", callRestBizParam.Method)
		}
	//case model.DEPLOYWASMCONTRACT:
	//	if callRestBizParam.Account == "" {
	//		passChecked = false
	//		data = fmt.Sprintf("%v method must has account", callRestBizParam.Method)
	//	}
	//	if callRestBizParam.ContractName == "" {
	//		passChecked = false
	//		data = fmt.Sprintf("%v method must has contract name", callRestBizParam.Method)
	//	}
	//	if callRestBizParam.ContractCode == "" {
	//		passChecked = false
	//		data = fmt.Sprintf("%v method must has contractCode", callRestBizParam.Method)
	//	}
	//	break
	//case model.CALLWASMCONTRACT:
	//case model.CALLWASMCONTRACTASYNC:
	//	if callRestBizParam.Account == "" {
	//		passChecked = false
	//		data = fmt.Sprintf("%v method must has account", callRestBizParam.Method)
	//	}
	//	if callRestBizParam.ContractName == "" {
	//		passChecked = false
	//		data = fmt.Sprintf("%v method must has contract name", callRestBizParam.Method)
	//	}
	//	if callRestBizParam.OutTypes == "" {
	//		passChecked = false
	//		data = fmt.Sprintf("%v method must has outTypes", callRestBizParam.Method)
	//	}
	//	if callRestBizParam.MethodSignature == "" {
	//		passChecked = false
	//		data = fmt.Sprintf("%v method must has methodSignature", callRestBizParam.Method)
	//	}
	//	if callRestBizParam.InputParamListStr == "" {
	//		passChecked = false
	//		data = fmt.Sprintf("%v method must has inputParamListStr", callRestBizParam.Method)
	//	}
	//case model.UPDATECONTRACTFORBIZ:
	//	if callRestBizParam.Account == "" {
	//		passChecked = false
	//		data = fmt.Sprintf("%v method must has account", callRestBizParam.Method)
	//	}
	//	if callRestBizParam.ContractName == "" {
	//		passChecked = false
	//		data = fmt.Sprintf("%v method must has contract name", callRestBizParam.Method)
	//	}
	//	if callRestBizParam.ContractCode == "" {
	//		passChecked = false
	//		data = fmt.Sprintf("%v method must has contractCode", callRestBizParam.Method)
	//	}
	case model.CREATEACCOUNT:
		if callRestBizParam.Account == "" {
			passChecked = false
			data = fmt.Sprintf("%v method must has account", callRestBizParam.Method)
		}
		if callRestBizParam.MykmsKeyId == "" {
			passChecked = false
			data = fmt.Sprintf("%v method must has mykmsKeyId", callRestBizParam.Method)
		}
	}

	resp := response.BaseResp{Success: passChecked, Data: data}
	return resp
}
