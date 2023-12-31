package gocd

import (
	"encoding/json"
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"
)

// GetPermissions fetches all pipelines configured in GoCD server.
func (conf *client) GetPermissions(query map[string]string) (Permission, error) {
	var permissions struct {
		Permission Permission `json:"permissions"`
	}

	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return Permission{}, err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		SetQueryParams(query).
		Get(PermissionsEndpoint)
	if err != nil {
		return Permission{}, &errors.APIError{Err: err, Message: "get permissions"}
	}

	if resp.StatusCode() != http.StatusOK {
		return Permission{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &permissions); err != nil {
		return Permission{}, &errors.MarshalError{Err: err}
	}

	return permissions.Permission, nil
}
