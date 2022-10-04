package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jinzhu/copier"
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
		return timeout, fmt.Errorf("call made to get default job timeout errored with: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return timeout, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &timeout); err != nil {
		return timeout, ResponseReadError(err.Error())
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
		return fmt.Errorf("call made to update default job timeout errored with: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	return nil
}
