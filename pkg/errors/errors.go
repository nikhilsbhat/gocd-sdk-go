package errors

import (
	"fmt"
)

func (err NonOkError) Error() string {
	return fmt.Sprintf(`got %d from GoCD while making %s call for %s
with BODY:%s`,
		err.Code,
		err.Response.Request.Method,
		err.Response.Request.URL,
		err.Response.Body())
}

func (err NonFoundError) Error() string {
	return fmt.Sprintf(`looks like the object your are looking in GoCD is not found, the response from GoCD we got: '%s'`,
		err.Response.String())
}

func (err CipherError) Error() string {
	return "value or cipher key cannot be empty"
}

func (err APIError) Error() string {
	return fmt.Sprintf("call made to %s errored with: %v", err.Message, err.Err)
}

func (err MarshalError) Error() string {
	return fmt.Sprintf("reading response body errored with: %v", err.Err)
}

func (err GoCDError) Error() string {
	return fmt.Sprintf("%s %v", err.Message, err.Err)
}

func (err NilHeaderError) Error() string {
	return fmt.Sprintf("header %s not set, this will impact while %s", err.Header, err.Message)
}

func (err GoCDSDKError) Error() string {
	return err.Message
}

func (err PipelineValidationError) Error() string {
	return err.Message
}
