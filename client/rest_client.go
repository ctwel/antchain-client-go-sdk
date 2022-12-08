package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gitlab.alipay-inc.com/antchain/restclient-go-sdk/client/config"
	"gitlab.alipay-inc.com/antchain/restclient-go-sdk/model"
	"gitlab.alipay-inc.com/antchain/restclient-go-sdk/response"
	"gitlab.alipay-inc.com/antchain/restclient-go-sdk/utils"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	DefaultMaxIdleConns     = 10
	DefaultIdleConnTimeout  = 30
	DefaultRetryMaxAttempts = 5
	DefaultBackOffPeriod    = 500
)

const (
	ShakeHandPath       = "/api/contract/shakeHand"
	ChainCallPath       = "/api/contract/chainCall"
	ChainCallForBizPath = "/api/contract/chainCallForBiz"
)

const (
	ChainCallForBiz = "chainCallForBiz"
	ChainCall       = "chainCall"
)

type RestClient struct {
	RestClientProperties config.RestClientProperties
	RestToken            string
	httpClient           *http.Client
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

func NewRestClient(restClientPropertiesPath string) (*RestClient, error) {
	data, err := ioutil.ReadFile(restClientPropertiesPath)
	if err != nil {
		log.WithFields(log.Fields{
			"restClientPropertiesPath": restClientPropertiesPath,
			"err":                      err.Error(),
		}).Error("fail to read restClientProperties")
		return nil, err
	}
	restClientProperties := config.RestClientProperties{}
	err = json.Unmarshal(data, &restClientProperties)
	if err != nil {
		log.WithFields(log.Fields{
			"restClientPropertiesPath": restClientPropertiesPath,
			"err":                      err.Error(),
		}).Error("fail to parse restClientProperties")
		return nil, err
	}

	maxIdleConns := DefaultMaxIdleConns
	if restClientProperties.MaxIdleConns != 0 {
		maxIdleConns = restClientProperties.MaxIdleConns
	}
	idleConnTimeout := DefaultIdleConnTimeout
	if restClientProperties.IdleConnTimeout != 0 {
		idleConnTimeout = restClientProperties.IdleConnTimeout
	}
	tr := &http.Transport{
		MaxIdleConns:    maxIdleConns,
		IdleConnTimeout: time.Duration(idleConnTimeout) * time.Second,
	}
	client := &http.Client{Transport: tr}

	restClient := &RestClient{
		RestClientProperties: restClientProperties,
		httpClient:           client,
	}

	err = restClient.shake()
	if err != nil {
		return nil, err
	}
	return restClient, nil
}

func (client *RestClient) CreateQueryAccountParam(queryAccount string) (model.ClientParam, error) {
	queryAccountRequest := model.AccountRequest{QueryAccount: queryAccount}

	jsonStr, err := json.Marshal(&queryAccountRequest)
	if err != nil {
		return model.ClientParam{}, err
	}
	queryAccountParam := model.ClientParam{
		SignData: string(jsonStr),
	}
	return queryAccountParam, nil
}

func (client *RestClient) shake() error {
	log.Info("start shake hand")
	nowMill := time.Now().UnixNano() / 1e6
	secret, err := utils.Sign(fmt.Sprintf("%v%v", client.RestClientProperties.AccessId, nowMill), client.RestClientProperties.AccessSecret)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("fail to sign secret")
		return err
	}
	shakeRequest := &model.ShakeRequest{
		AccessId: client.RestClientProperties.AccessId,
		Time:     fmt.Sprintf("%v", nowMill),
		Secret:   secret,
	}
	jsonStr, err := json.Marshal(shakeRequest)
	//if err != nil {
	//	log.WithFields(log.Fields{
	//		"shakeRequest": shakeRequest,
	//		"err":          err.Error(),
	//	}).Error("fail to marshal shakeRequest")
	//	return err
	//}
	req, err := http.NewRequest(http.MethodPost, client.RestClientProperties.RestUrl+ShakeHandPath, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.WithFields(log.Fields{
			"shakeRequest": shakeRequest,
			"err":          err.Error(),
		}).Error("fail to new shakeRequest")
		return err
	}
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	resp, err := client.httpClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"req": req,
			"err": err.Error(),
		}).Error("fail to get shakeResponse")
		return err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	baseResp := response.BaseResp{}
	err = json.Unmarshal(body, &baseResp)
	//if err != nil {
	//	log.WithFields(log.Fields{
	//		"body": string(body),
	//		"err":  err.Error(),
	//	}).Error("fail to unmarshal shakeResponse")
	//	return err
	//}
	client.RestToken = baseResp.Data
	log.Info("new rest token:" + client.RestToken)
	return nil
}

func (client *RestClient) ChainCall(hash, bizid, requestStr string, method model.Method) (response.BaseResp, error) {
	if bizid == "" {
		return response.BaseResp{}, fmt.Errorf("bizid is empty")
	}
	if method == "" {
		return response.BaseResp{}, fmt.Errorf("method is empty")
	}
	param := &model.CallRestParam{}
	param.AccessId = client.RestClientProperties.AccessId
	param.Token = client.RestToken
	param.Hash = hash
	param.BizId = bizid
	param.RequestStr = requestStr
	param.Method = method
	return client.retryableSendRequest(param, client.RestClientProperties.RestUrl+ChainCallPath, ChainCall)
}

func (client *RestClient) ChainCallForBiz(param model.CallRestBizParam) (response.BaseResp, error) {
	param.Token = client.RestToken
	baseResp := utils.CheckCallRestBizParams(param)
	if !baseResp.Success {
		return response.BaseResp{}, fmt.Errorf("%v", baseResp.Data)
	}
	if param.Method == model.CREATEACCOUNT || param.Method == model.DEPLOYNATIVECONTRACT || param.Method == model.QUERYACCOUNT {
		if param.MykmsKeyId == "" {
			return client.ChainCall("", param.BizId, param.RequestStr, param.Method)
		}
	}

	return client.retryableSendRequest(param, client.RestClientProperties.RestUrl+ChainCallForBizPath, ChainCallForBiz)
}

func (client *RestClient) retryableSendRequest(param interface{}, url string, chainCallType string) (response.BaseResp, error) {
	retryMaxAttempts := DefaultRetryMaxAttempts
	if client.RestClientProperties.RetryMaxAttempts != 0 {
		retryMaxAttempts = client.RestClientProperties.RetryMaxAttempts
	}
	backoffPeriod := DefaultBackOffPeriod
	if client.RestClientProperties.BackOffPeriod != 0 {
		backoffPeriod = client.RestClientProperties.BackOffPeriod
	}

	tick := time.Tick(time.Duration(backoffPeriod) * time.Millisecond)
	for i := 0; i < retryMaxAttempts; i++ {
		jsonStr, err := json.Marshal(&param)
		if err != nil {
			return response.BaseResp{}, err
		}
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr))
		if err != nil {
			log.WithFields(log.Fields{
				"req": req,
				"err": err.Error(),
			}).Errorf("fail to new %v request", chainCallType)
			return response.BaseResp{}, err
		}
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
		resp, err := client.httpClient.Do(req)
		if err != nil {
			log.WithFields(log.Fields{
				"req": req,
				"err": err.Error(),
			}).Errorf("fail to get %v response", chainCallType)
			// retry later An error is returned if caused by client policy (such as
			// CheckRedirect), or failure to speak HTTP (such as a network
			// connectivity problem). A non-2xx status code doesn't cause an
			// error.
			select {
			case <-tick:
				log.WithFields(log.Fields{
					"req": req,
				}).Infof("retry %v request", chainCallType)
			}
			resp.Body.Close()
		} else {
			if resp.StatusCode >= 300 && resp.StatusCode < 600 {
				log.WithFields(log.Fields{
					"req":        req,
					"statusCode": resp.StatusCode,
				}).Warnf("%v return non 2xx code", chainCallType)
				resp.Body.Close()
				return response.BaseResp{}, fmt.Errorf("%v return non 2xx code,statusCode:%v", chainCallType, resp.StatusCode)
			} else {
				body, _ := ioutil.ReadAll(resp.Body)
				baseResp := response.BaseResp{}
				err = json.Unmarshal(body, &baseResp)
				//if err != nil {
				//	log.WithFields(log.Fields{
				//		"body": string(body),
				//		"err":  err.Error(),
				//	}).Errorf("fail to unmarshal %v", chainCallType)
				//	resp.Body.Close()
				//	return response.BaseResp{}, fmt.Errorf("fail to unmarshal %v,err:%+v", chainCallType, err)
				//}
				log.WithFields(log.Fields{
					"param": param,
					"resp":  baseResp,
				}).Info("request and resp")
				if !baseResp.Success {
					if baseResp.Code == "202" {
						client.shake()
						switch param.(type) {
						case model.CallRestParam:
							newParam := param.(model.CallRestParam)
							newParam.Token = client.RestToken
							param = newParam
						case model.CallRestBizParam:
							newParam := param.(model.CallRestBizParam)
							newParam.Token = client.RestToken
							param = newParam
						}
					}
					if baseResp.Code == "202" || strings.HasPrefix(baseResp.Code, "5") {
						log.WithFields(log.Fields{
							"restCode": baseResp.Code,
						}).Warnf("fail to get %v successfully", chainCallType)
						resp.Body.Close()
						continue // retry next time
					}
				}
				resp.Body.Close()
				return baseResp, nil
			}
		}
	}
	return response.BaseResp{}, fmt.Errorf("fail to get %v response", chainCallType)
}

func (client *RestClient) DepositSyncWithTransaction(bizid, orderId, account, tenantId, content, mykmsKeyId string, gas int64) (response.BaseResp, error) {
	baseResp, err := client.Deposit(bizid, orderId, account, tenantId, content, mykmsKeyId, gas)
	if err != nil {
		return response.BaseResp{}, err
	}
	if !baseResp.Success || baseResp.Code != "200" {
		return response.BaseResp{}, fmt.Errorf("deposit failed,code:%+v err msg:%+v", baseResp.Code, baseResp.Data)
	}
	return client.MultipleQueryTransaction(bizid, baseResp.Data)
}

//func (client *RestClient) CallSolcContractSyncWithReceipt(abi abi.ABI, bizid, orderId, account, tenantId, kmsId, contractName, methodSignature, inputParamListStr, outTypes string, gas int64, respStruct interface{}) (response.BaseResp, error) {
//	callRestBizParam := model.CallRestBizParam{
//		BaseParam: model.BaseParam{
//			AccessId: client.RestClientProperties.AccessId,
//			BizId:    bizid,
//			Method:   model.CALLCONTRACTBIZASYNC,
//		},
//		OrderId:           orderId,
//		Account:           account,
//		TenantId:          tenantId,
//		ContractName:      contractName,
//		MethodSignature:   methodSignature,
//		InputParamListStr: inputParamListStr,
//		OutTypes:          outTypes,
//		MykmsKeyId:        kmsId,
//		Gas:               gas, // 0表示不受限
//	}
//	callResp, err := client.ChainCallForBiz(callRestBizParam)
//	if err != nil {
//		return response.BaseResp{}, err
//	}
//	if callResp.Success && callResp.Code == "200" {
//		baseResp, err := client.MultipleQueryReceipt(bizid, callResp.Data)
//		if err != nil {
//			return response.BaseResp{}, err
//		}
//		transactionReceipt := mychain.GetDefaultTransactionReceipt()
//		err = json.Unmarshal([]byte(baseResp.Data), &transactionReceipt)
//		if err != nil {
//			return response.BaseResp{}, err
//		}
//		output := make([]string, 0)
//		err = json.Unmarshal([]byte(outTypes), &output)
//		if err != nil {
//			return response.BaseResp{}, err
//		}
//		if len(output) > 0 {
//			if transactionReceipt.Output == "" && output[0] != model.VOID {
//				return response.BaseResp{}, fmt.Errorf("function has no any output")
//			}
//			decodedOutput, err := base64.StdEncoding.DecodeString(transactionReceipt.Output)
//			if err != nil {
//				return response.BaseResp{}, err
//			}
//			strParts := strings.Split(methodSignature, "(")
//			methodName := strParts[0]
//			err = abi.Unpack(respStruct, methodName, decodedOutput)
//			if err != nil {
//				return response.BaseResp{}, err
//			}
//			jsonStr, err := json.Marshal(respStruct)
//			if err != nil {
//				return response.BaseResp{}, err
//			}
//			return response.BaseResp{Success: true, Code: "200", Data: string(jsonStr)}, nil
//		}
//	}
//	return response.BaseResp{}, fmt.Errorf("no succ call contract resp,Success:%v Code:%v", callResp.Success, callResp.Code)
//}

func (client *RestClient) QueryAccount(bizid, account string) (response.BaseResp, error) {
	clientParam, err := client.CreateQueryAccountParam(account)
	if err != nil {
		return response.BaseResp{}, err
	}
	return client.ChainCall(clientParam.Hash, bizid, clientParam.SignData, model.QUERYACCOUNT)
}

func (client *RestClient) CreateAccountWithKmsId(bizid, orderId, account, tenantId, kmsId string) (response.BaseResp, error) {
	callRestBizParam := model.CallRestBizParam{
		BaseParam: model.BaseParam{
			AccessId: client.RestClientProperties.AccessId,
			BizId:    bizid,
			Method:   model.CREATEACCOUNT,
		},
		OrderId:    orderId,
		TenantId:   tenantId,
		Account:    account,
		MykmsKeyId: kmsId,
	}
	return client.ChainCallForBiz(callRestBizParam)
}

func (client *RestClient) CallContract(bizid, orderId, account, tenantId, contractName, methodSignature, inputParamListStr, outTypes, kmsId string, isLocal bool, gas int64) (response.BaseResp, error) {
	callRestBizParam := model.CallRestBizParam{
		BaseParam: model.BaseParam{
			AccessId: client.RestClientProperties.AccessId,
			BizId:    bizid,
			Method:   model.CALLCONTRACTBIZ,
		},
		OrderId:            orderId,
		Account:            account,
		TenantId:           tenantId,
		ContractName:       contractName,
		MethodSignature:    methodSignature,
		InputParamListStr:  inputParamListStr,
		OutTypes:           outTypes,
		MykmsKeyId:         kmsId,
		IsLocalTransaction: isLocal,
		Gas:                gas, // 0表示不受限
	}
	return client.ChainCallForBiz(callRestBizParam)
}

func (client *RestClient) DeployContract(bizid, orderId, account, tenantId, kmsId, contractName, contractCode string, gas int64) (response.BaseResp, error) {
	//deploy contract
	callRestBizParam := model.CallRestBizParam{
		BaseParam: model.BaseParam{
			AccessId: client.RestClientProperties.AccessId,
			BizId:    bizid,
			Method:   model.DEPLOYCONTRACTFORBIZ,
		},
		OrderId:      orderId,
		Account:      account,
		MykmsKeyId:   kmsId,
		TenantId:     tenantId,
		ContractName: contractName,
		ContractCode: contractCode,
		Gas:          gas, // 0表示不受限
	}
	return client.ChainCallForBiz(callRestBizParam)
}

func (client *RestClient) Deposit(bizid, orderId, account, tenantId, content, mykmsKeyId string, gas int64) (response.BaseResp, error) {
	callRestBizParam := model.CallRestBizParam{
		BaseParam: model.BaseParam{
			AccessId: client.RestClientProperties.AccessId,
			BizId:    bizid,
			Method:   model.DEPOSIT,
		},
		OrderId:    orderId,
		Account:    account,
		Content:    content,
		MykmsKeyId: mykmsKeyId,
		TenantId:   tenantId,
		Gas:        gas, // 0表示不受限
	}
	return client.ChainCallForBiz(callRestBizParam)
}

func (client *RestClient) QueryReceipt(bizid, hash string) (response.BaseResp, error) {
	callRestBizParam := model.CallRestBizParam{
		BaseParam: model.BaseParam{
			AccessId: client.RestClientProperties.AccessId,
			BizId:    bizid,
			Hash:     hash,
			Method:   model.QUERYRECEIPT,
		},
	}
	return client.ChainCallForBiz(callRestBizParam)
}

func (client *RestClient) QueryTransaction(bizid, hash string) (response.BaseResp, error) {
	callRestBizParam := model.CallRestBizParam{
		BaseParam: model.BaseParam{
			AccessId: client.RestClientProperties.AccessId,
			BizId:    bizid,
			Hash:     hash,
			Method:   model.QUERYTRANSACTION,
		},
	}
	return client.ChainCallForBiz(callRestBizParam)
}

func (client *RestClient) MultipleQueryReceipt(bizid, hash string) (response.BaseResp, error) {
	var baseResp response.BaseResp
	var err error
	for i := 0; i < client.RestClientProperties.RetryMaxAttempts; i++ {
		baseResp, err = client.ChainCall(hash, bizid, "", model.QUERYRECEIPT)
		if err != nil {
			return baseResp, err
		} else if !baseResp.Success && (baseResp.Code == model.ServiceQueryNoResult ||
			baseResp.Code == model.ServiceTxWaitingVerify ||
			baseResp.Code == model.ServiceTxWaitingExecute) {
			continue
		}
		return baseResp, err
	}
	return baseResp, err
}

func (client *RestClient) MultipleQueryTransaction(bizid, hash string) (response.BaseResp, error) {
	var baseResp response.BaseResp
	var err error
	for i := 0; i < client.RestClientProperties.RetryMaxAttempts; i++ {
		baseResp, err = client.ChainCall(hash, bizid, "", model.QUERYTRANSACTION)
		if err != nil {
			return baseResp, err
		} else if !baseResp.Success && (baseResp.Code == model.ServiceQueryNoResult ||
			baseResp.Code == model.ServiceTxWaitingVerify ||
			baseResp.Code == model.ServiceTxWaitingExecute) {
			continue
		}
		return baseResp, err
	}
	return baseResp, err
}
