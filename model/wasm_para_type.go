package model

type WasmParaType string

const (
	IDENTITY       WasmParaType = "Identity"
	VECTORIDENTITY              = "Identity[]"
	INT8                        = "int8"
	INT16                       = "int16"
	INT32                       = "int32"
	INT64                       = "int64"
	VECTORINT8                  = "int8[]"
	VECTORINT16                 = "int16[]"
	VECTORINT32                 = "int32[]"
	VECTORINT64                 = "int64[]"
	STRING                      = "string"
	VECTORSTRING                = "string[]"
	UINT8                       = "uint8"
	UINT16                      = "uint16"
	UINT32                      = "uint32"
	UINT64                      = "uint64"
	VECTORUINT8                 = "uint8[]"
	VECTORUINT16                = "uint16[]"
	VECTORUINT32                = "uint32[]"
	VECTORUINT64                = "uint64[]"
	BOOL                        = "bool"
	VECTORBOOL                  = "bool[]"
	VOID                        = "void"
)
