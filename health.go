package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jinzhu/copier"
)

// GetHealthMessages implements method that fetches the details of all warning and errors.
func (conf *client) GetHealthMessages() ([]ServerHealth, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}

	newClient.httpClient.SetHeaders(map[string]string{
		"Accept": GoCdHeaderVersionOne,
	})

	var health []ServerHealth
	resp, err := newClient.httpClient.R().SetResult(&health).Get(GoCdServerHealthEndpoint)
	if err != nil {
		return nil, fmt.Errorf("call made to get health info errored with %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, APIWithCodeError(resp.StatusCode())
	}

	if err := json.Unmarshal(resp.Body(), &health); err != nil {
		return nil, ResponseReadError(err.Error())
	}

	return health, nil
}
