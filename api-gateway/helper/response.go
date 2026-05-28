package helper

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   ErrorDetail `json:"error"`
}

func Success(data interface{}) SuccessResponse {
	return SuccessResponse{Success: true, Data: data}
}

func Fail(code, message string) ErrorResponse {
	return ErrorResponse{Success: false, Error: ErrorDetail{Code: code, Message: message}}
}