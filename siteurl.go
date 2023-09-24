package gocd

import (
	"encoding/json"
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"
)

// GetSiteURL fetches the site url config configured from GoCD.
func (conf *client) GetSiteURL() (SiteURLConfig, error) {
	var site SiteURLConfig

	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return site, err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		Get(SiteURLEndpoint)
	if err != nil {
		return site, &errors.APIError{Err: err, Message: "get site url"}
	}

	if resp.StatusCode() != http.StatusOK {
		return site, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &site); err != nil {
		return site, err
	}

	return site, nil
}

// CreateOrUpdateSiteURL creates/updates the site url configured in GoCD.
func (conf *client) CreateOrUpdateSiteURL(config SiteURLConfig) (SiteURLConfig, error) {
	var site SiteURLConfig

	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return site, err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionOne,
			"Content-Type": ContentJSON,
		}).
		SetBody(config).
		Post(SiteURLEndpoint)
	if err != nil {
		return site, &errors.APIError{Err: err, Message: "create/update site url"}
	}

	if resp.StatusCode() != http.StatusOK {
		return site, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &site); err != nil {
		return site, err
	}

	return site, nil
}
