package resp

type Result struct {
	Code    int      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}


func Success(data interface{}) Result {
	return Result{Code: 0, Message: "success", Data: data}
}

func Fail(code int, message string) Result {
	return Result{Code: code, Message: message}
}
