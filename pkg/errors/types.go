package errors

import (
	"github.com/go-resty/resty/v2"
)

type NonOkError struct {
	Code     int
	Response *resty.Response
}

type CipherError struct{}

type APIError struct {
	Err     error
	Message string
}

type MarshalError struct {
	Err error
}

type GoCDError struct {
	Message string
	Err     error
}

type NilHeaderError struct {
	Header  string
	Message string
}

type GoCDSDKError struct {
	Message string
}

type PipelineValidationError struct {
	Message string
}

type NonFoundError struct {
	Code     int
	Response *resty.Response
}
