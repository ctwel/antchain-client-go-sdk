package utils

import (
	"github.com/stretchr/testify/require"
	"gitlab.alipay-inc.com/antchain/restclient-go-sdk/model"
	"testing"
)

func TestCheckCallRestBizParams_NoAccessId(t *testing.T) {
	callRestBizParam := model.CallRestBizParam{
		OrderId: "",
	}
	resp := CheckCallRestBizParams(callRestBizParam)
	require.Truef(t, !resp.Success, "cannot check without access id")
}

func TestCheckCallRestBizParams_NoToken(t *testing.T) {
	callRestBizParam := model.CallRestBizParam{
		OrderId: "",
	}
	callRestBizParam.AccessId = "accessId"
	resp := CheckCallRestBizParams(callRestBizParam)
	require.Truef(t, !resp.Success, "cannot check without token")
}

func TestCheckCallRestBizParams_NoBizId(t *testing.T) {
	callRestBizParam := model.CallRestBizParam{
		OrderId: "",
	}
	callRestBizParam.AccessId = "accessId"
	callRestBizParam.Token = "token"
	resp := CheckCallRestBizParams(callRestBizParam)
	require.Truef(t, !resp.Success, "cannot check without bizid")
}

func TestCheckCallRestBizParams_TransactionRequestWithoutSignId(t *testing.T) {
	callRestBizParam := model.CallRestBizParam{
		BaseParam: model.BaseParam{
			AccessId: "accessId",
			BizId:    "bizid",
			Token:    "token",
			Method:   model.DEPOSIT,
		},
		OrderId: "orderId",
	}
	resp := CheckCallRestBizParams(callRestBizParam)
	require.Truef(t, !resp.Success, "cannot check without signId")
}

func TestCheckCallRestBizParams_NoOrderId(t *testing.T) {
	callRestBizParam := model.CallRestBizParam{
		BaseParam: model.BaseParam{
			AccessId: "accessId",
			BizId:    "bizid",
			Token:    "token",
			Method:   model.DEPOSIT,
		},
		MykmsKeyId: "mykmsId",
	}
	resp := CheckCallRestBizParams(callRestBizParam)
	require.Truef(t, !resp.Success, "cannot check without orderId")
}

func TestCheckCallRestBizParams_DepositWithoutAccount(t *testing.T) {
	callRestBizParam := model.CallRestBizParam{
		BaseParam: model.BaseParam{
			AccessId: "accessId",
			BizId:    "bizid",
			Token:    "token",
			Method:   model.DEPOSIT,
		},
		OrderId:    "orderId",
		MykmsKeyId: "mykmsId",
	}
	resp := CheckCallRestBizParams(callRestBizParam)
	require.Truef(t, !resp.Success, "cannot check deposit without account")
}

func TestCheckCallRestBizParams_DepositWithoutContent(t *testing.T) {
	callRestBizParam := model.CallRestBizParam{
		BaseParam: model.BaseParam{
			AccessId: "accessId",
			BizId:    "bizid",
			Token:    "token",
			Method:   model.DEPOSIT,
		},
		OrderId:    "orderId",
		Account:    "account",
		MykmsKeyId: "mykmsId",
	}
	resp := CheckCallRestBizParams(callRestBizParam)
	require.Truef(t, !resp.Success, "cannot check deposit without content")
}

func TestCheckCallRestBizParams_CallContractWithoutAccount(t *testing.T) {
	callRestBizParam := model.CallRestBizParam{
		BaseParam: model.BaseParam{
			AccessId: "accessId",
			BizId:    "bizid",
			Token:    "token",
			Method:   model.CALLCONTRACTBIZASYNC,
		},
		OrderId:    "orderId",
		MykmsKeyId: "mykmsId",
	}
	resp := CheckCallRestBizParams(callRestBizParam)
	require.Truef(t, !resp.Success, "cannot check call contract without account")
}

func TestCheckCallRestBizParams_CallContractWithoutContractName(t *testing.T) {
	callRestBizParam := model.CallRestBizParam{
		BaseParam: model.BaseParam{
			AccessId: "accessId",
			BizId:    "bizid",
			Token:    "token",
			Method:   model.CALLCONTRACTBIZASYNC,
		},
		OrderId:    "orderId",
		Account:    "account",
		MykmsKeyId: "mykmsId",
	}
	resp := CheckCallRestBizParams(callRestBizParam)
	require.Truef(t, !resp.Success, "cannot check call contract without contract name")
}

func TestCheckCallRestBizParams_CallContractWithoutOutTypes(t *testing.T) {
	callRestBizParam := model.CallRestBizParam{
		BaseParam: model.BaseParam{
			AccessId: "accessId",
			BizId:    "bizid",
			Token:    "token",
			Method:   model.CALLCONTRACTBIZASYNC,
		},
		OrderId:      "orderId",
		Account:      "account",
		ContractName: "contractName",
		MykmsKeyId:   "mykmsId",
	}
	resp := CheckCallRestBizParams(callRestBizParam)
	require.Truef(t, !resp.Success, "cannot check call contract without out types")
}
func TestCheckCallRestBizParams_CallContractWithoutMethodSignature(t *testing.T) {
	callRestBizParam := model.CallRestBizParam{
		BaseParam: model.BaseParam{
			AccessId: "accessId",
			BizId:    "bizid",
			Token:    "token",
			Method:   model.CALLCONTRACTBIZASYNC,
		},
		OrderId:      "orderId",
		MykmsKeyId:   "mykmsId",
		Account:      "account",
		ContractName: "contractName",
		OutTypes:     "outtype",
	}
	resp := CheckCallRestBizParams(callRestBizParam)
	require.Truef(t, !resp.Success, "cannot check call contract without out signature")
}

func TestCheckCallRestBizParams_CallContractWithoutInputParamList(t *testing.T) {
	callRestBizParam := model.CallRestBizParam{
		BaseParam: model.BaseParam{
			AccessId: "accessId",
			BizId:    "bizid",
			Token:    "token",
			Method:   model.CALLCONTRACTBIZASYNC,
		},
		OrderId:         "orderId",
		MykmsKeyId:      "mykmsId",
		Account:         "account",
		ContractName:    "contractName",
		OutTypes:        "outtype",
		MethodSignature: "foo()",
	}
	resp := CheckCallRestBizParams(callRestBizParam)
	require.Truef(t, !resp.Success, "cannot check call contract without out input param list")
}

func TestCheckCallRestBizParams_QueryTransactionWithoutHash(t *testing.T) {
	callRestBizParam := model.CallRestBizParam{
		BaseParam: model.BaseParam{
			AccessId: "accessId",
			BizId:    "bizid",
			Token:    "token",
			Method:   model.QUERYTRANSACTION,
		},
		OrderId: "orderId",
	}
	resp := CheckCallRestBizParams(callRestBizParam)
	require.Truef(t, !resp.Success, "cannot check query transaction without out hash")
}

func TestCheckCallRestBizParams_DeployContractWithoutAccount(t *testing.T) {
	callRestBizParam := model.CallRestBizParam{
		BaseParam: model.BaseParam{
			AccessId: "accessId",
			BizId:    "bizid",
			Token:    "token",
			Method:   model.DEPLOYCONTRACTFORBIZ,
		},
		OrderId: "orderId",
		MykmsKeyId: "kmsId",
	}
	resp := CheckCallRestBizParams(callRestBizParam)
	require.Truef(t, !resp.Success, "cannot check deploy contract without out account")
}

func TestCheckCallRestBizParams_DeployContractWithoutContractName(t *testing.T) {
	callRestBizParam := model.CallRestBizParam{
		BaseParam: model.BaseParam{
			AccessId: "accessId",
			BizId:    "bizid",
			Token:    "token",
			Method:   model.DEPLOYCONTRACTFORBIZ,
		},
		OrderId: "orderId",
		MykmsKeyId: "kmsId",
		Account: "account",
	}
	resp := CheckCallRestBizParams(callRestBizParam)
	require.Truef(t, !resp.Success, "cannot check deploy contract without out contract name")
}

func TestCheckCallRestBizParams_DeployContractWithoutContractCode(t *testing.T) {
	callRestBizParam := model.CallRestBizParam{
		BaseParam: model.BaseParam{
			AccessId: "accessId",
			BizId:    "bizid",
			Token:    "token",
			Method:   model.DEPLOYCONTRACTFORBIZ,
		},
		OrderId:      "orderId",
		MykmsKeyId: "kmsId",
		Account:      "account",
		ContractName: "contractName",
	}
	resp := CheckCallRestBizParams(callRestBizParam)
	require.Truef(t, !resp.Success, "cannot check deploy contract without out contract code")
}

func TestCheckCallRestBizParams_CreateAccountWithoutAccount(t *testing.T) {
	callRestBizParam := model.CallRestBizParam{
		BaseParam: model.BaseParam{
			AccessId: "accessId",
			BizId:    "bizid",
			Token:    "token",
			Method:   model.CREATEACCOUNT,
		},
		OrderId: "orderId",
	}
	resp := CheckCallRestBizParams(callRestBizParam)
	require.Truef(t, !resp.Success, "cannot check create account without account")
}

func TestCheckCallRestBizParams_CreateAccountWithoutMykmsId(t *testing.T) {
	callRestBizParam := model.CallRestBizParam{
		BaseParam: model.BaseParam{
			AccessId: "accessId",
			BizId:    "bizid",
			Token:    "token",
			Method:   model.CREATEACCOUNT,
		},
		OrderId: "orderId",
		Account: "account",
	}
	resp := CheckCallRestBizParams(callRestBizParam)
	require.Truef(t, !resp.Success, "cannot check create account without kmsid")
}
