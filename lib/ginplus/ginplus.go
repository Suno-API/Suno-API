package ginplus

const (
	SUCCESS = "success"
)

type DataResult struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func BuildApiReturn(code string, msg string, data interface{}) DataResult {
	return DataResult{Code: code, Message: msg, Data: data}
}

func ApiRetSucc(data interface{}) *DataResult {
	return &DataResult{
		Code:    SUCCESS,
		Message: "",
		Data:    data,
	}
}
