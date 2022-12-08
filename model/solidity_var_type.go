package model

type SolidityVarType string

const (
	Unsupported   SolidityVarType = "unsupported"
	Int                           = "int"
	Int64                         = "int64"
	IntArray                      = "int[]"
	Int64Array                    = "int64[]"
	Uint                          = "uint"
	UintArray                     = "uint[]"
	Bool                          = "bool"
	BoolArray                     = "bool[]"
	Bytes                         = "bytes"
	BytesArray                    = "bytes[]"
	Identity                      = "identity"
	IdentityArray                 = "identity[]"
	String                        = "string"
	EncodedBytes                  = "encodedbytes"
	ListBytes                     = "list(bytes)"
)
