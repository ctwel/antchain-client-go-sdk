package model

type ErrorCode string

const (
	ServiceQueryNoResult    = "404"
	ServiceTxWaitingVerify  = "413"
	ServiceTxWaitingExecute = "414"
)
