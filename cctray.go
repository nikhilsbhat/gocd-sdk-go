package gocd

import (
	"encoding/xml"
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"
)

func (conf *client) GetCCTray() ([]Project, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}

	var projectsConf Projects

	resp, err := newClient.httpClient.R().
		Get("cctray.xml")
	if err != nil {
		return nil, &errors.APIError{Err: err, Message: "get cctray"}
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = xml.Unmarshal(resp.Body(), &projectsConf); err != nil {
		return nil, &errors.MarshalError{Err: err}
	}

	return projectsConf.Project, nil
}
