package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/jinzhu/copier"
	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"
)

func (conf *client) GetPipelineVSM(pipeline, instance string) (VSM, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return VSM{}, err
	}

	var vsmObj VSM

	resp, err := newClient.httpClient.R().
		Get(filepath.Join(VSMEndpoint, pipeline, fmt.Sprintf("%s.json", instance))) //nolint:perfsprint
	if err != nil {
		return vsmObj, &errors.APIError{
			Err:     err,
			Message: fmt.Sprintf("get vsm information for pipeline '%s' of instance '%s'", pipeline, instance),
		}
	}

	if resp.StatusCode() != http.StatusOK {
		return vsmObj, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &vsmObj); err != nil {
		return vsmObj, &errors.MarshalError{Err: err}
	}

	return vsmObj, nil
}
