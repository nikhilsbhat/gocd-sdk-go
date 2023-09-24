package gocd

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jinzhu/copier"
	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"
)

func (conf *client) GetDefaultJobTimeout() (map[string]string, error) {
	var timeout map[string]string

	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return timeout, err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		Get(DefaultTimeoutEndpoint)
	if err != nil {
		return timeout, &errors.APIError{Err: err, Message: "get default job timeout"}
	}

	if resp.StatusCode() != http.StatusOK {
		return timeout, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &timeout); err != nil {
		return timeout, &errors.MarshalError{Err: err}
	}

	return timeout, nil
}

func (conf *client) UpdateDefaultJobTimeout(timeoutMinutes int) error {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		SetBody(
			map[string]string{
				"default_job_timeout": strconv.Itoa(timeoutMinutes),
			},
		).
		Post(DefaultTimeoutEndpoint)
	if err != nil {
		return &errors.APIError{Err: err, Message: "update default job timeout"}
	}

	if resp.StatusCode() != http.StatusOK {
		return &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	return nil
}
