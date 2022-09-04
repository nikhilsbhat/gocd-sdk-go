package gocd

import (
	"fmt"
)

func apiWithCodeError(code int) error {
	return fmt.Errorf("goCd server returned code %d with message", code) //nolint:goerr113
}

func responseReadError(msg string) error {
	return fmt.Errorf("reading resposne body errored with: %s", msg) //nolint:goerr113
}
