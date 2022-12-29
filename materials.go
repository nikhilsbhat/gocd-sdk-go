package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"

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
		return nil, fmt.Errorf("call made to get all available materials errored with: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &materials); err != nil {
		return nil, ResponseReadError(err.Error())
	}

	return materials.Materials.Materials, nil
}
