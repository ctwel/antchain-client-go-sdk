package model

type VMTypeEnum string

const(
	NATIVE VMTypeEnum = "NATIVE"
	EVM = "EVM"
	WASM = "WASM"
	NATIVE_PRECOMPILE = "NATIVE_PRECOMPILE"
)
