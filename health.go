package gocd

import (
	"encoding/json"
	"net/http"

	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"

	"github.com/jinzhu/copier"
)

// GetServerHealthMessages implements method that fetches the details of all warning and errors.
func (conf *client) GetServerHealthMessages() ([]ServerHealth, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}

	var health []ServerHealth
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		SetResult(&health).Get(ServerHealthEndpoint)
	if err != nil {
		return nil, &errors.APIError{Err: err, Message: "get health info"}
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &health); err != nil {
		return nil, &errors.MarshalError{Err: err}
	}

	return health, nil
}

func (conf *client) GetServerHealth() (map[string]string, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}

	var health map[string]string
	resp, err := newClient.httpClient.R().
		SetResult(&health).Get(HealthEndpoint)
	if err != nil {
		return health, &errors.APIError{Err: err, Message: "get server health"}
	}

	if resp.StatusCode() != http.StatusOK {
		return health, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &health); err != nil {
		return health, &errors.MarshalError{Err: err}
	}

	return health, nil
}
