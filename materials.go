package gocd

import (
	"encoding/json"
	"fmt"
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

func (conf *client) GetMaterialUsage(materialID string) ([]string, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}

	var materialUsage MaterialUsage
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionZero,
		}).
		Get(fmt.Sprintf(MaterialUsageEndpoint, materialID))
	if err != nil {
		return nil, &errors.APIError{Err: err, Message: fmt.Sprintf("get material usage '%s'", materialID)}
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &materialUsage); err != nil {
		return nil, &errors.MarshalError{Err: err}
	}

	return materialUsage.Usages, nil
}
