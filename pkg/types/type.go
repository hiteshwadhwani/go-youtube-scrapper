package types

type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

func NewSuccessResponse(data interface{}) *Response {
	return &Response{
		Status: "success",
		Data:   data,
	}
}

func NewErrorResponse(err string) *Response {
	return &Response{
		Status: "error",
		Error:  err,
	}
}
