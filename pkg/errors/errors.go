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