package gocd

import (
	"encoding/json"
	"net/http"

	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"

	"github.com/jinzhu/copier"
)

func (conf *client) GetMaterials() ([]Material, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}

	var materials Materials
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionTwo,
		}).
		Get(MaterialEndpoint)
	if err != nil {
		return nil, &errors.APIError{Err: err, Message: "get all available materials"}
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &materials); err != nil {
		return nil, &errors.MarshalError{Err: err}
	}

	return materials.Materials.Materials, nil
}
