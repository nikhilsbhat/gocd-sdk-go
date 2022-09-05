package gocd

import (
	"fmt"
)

func APIWithCodeError(code int) error {
	return fmt.Errorf("goCd server returned code %d with message", code) //nolint:goerr113
}

func APIErrorWithBody(body string, code int) error {
	return fmt.Errorf("body: %s httpcode: %d", body, code) //nolint:goerr113
}

func ResponseReadError(msg string) error {
	return fmt.Errorf("reading response body errored with: %s", msg) //nolint:goerr113
}
