package enum

// SSEDataType SSE数据类型
//
//	@author centonhuang
//	@update 2025-11-08 04:20:42
type SSEDataType string

const (
	// SSEDataTypeText 文本数据
	//
	//	@author centonhuang
	//	@update 2025-11-08 04:20:42
	SSEDataTypeText = "text"

	// SSEDataTypeError 错误数据
	//	@author centonhuang
	//	@update 2025-11-08 04:39:06
	SSEDataTypeError = "error"

	// SSEDataTypePing ping数据
	//	@author centonhuang
	//	@update 2025-11-08 04:39:27
	SSEDataTypePing = "ping"
)
