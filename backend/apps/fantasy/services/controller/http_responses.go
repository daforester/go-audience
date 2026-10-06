package controller

import "net/http"

type Response struct {
	ErrorMessage string
	Response     interface{}
	ResponseCode int
}

func ErrorResponse(status int, message string) *Response {
	return &Response{
		ErrorMessage: message,
		ResponseCode: status,
		Response:     nil,
	}
}

func StandardResponse(x interface{}) *Response {
	return &Response{
		ErrorMessage: "",
		ResponseCode: http.StatusOK,
		Response:     x,
	}
}
