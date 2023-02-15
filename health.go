package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"

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
		return nil, fmt.Errorf("call made to get health info errored with %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, APIWithCodeError(resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &health); err != nil {
		return nil, ResponseReadError(err.Error())
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
		return health, fmt.Errorf("call made to get server health errored with %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return health, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &health); err != nil {
		return health, ResponseReadError(err.Error())
	}

	return health, nil
}
