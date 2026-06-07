package gocd

import (
	"encoding/xml"
	"net/http"

	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"
)

func (conf *client) GetCCTray() ([]Project, error) {
	newClient := conf.clone()

	var projectsConf Projects

	resp, err := newClient.httpClient.R().
		Get("cctray.xml")
	if err != nil {
		return nil, &errors.APIError{Err: err, Message: "get cctray"}
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	//nolint:musttag
	if err = xml.Unmarshal(resp.Body(), &projectsConf); err != nil {
		return nil, &errors.MarshalError{Err: err}
	}

	return projectsConf.Project, nil
}
