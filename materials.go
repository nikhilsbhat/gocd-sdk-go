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
			"Accept": HeaderVersionZero,
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

	return materials.Materials, nil
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

func (conf *client) NotifyMaterial(material Material) (string, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return "", err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionTwo,
			"Content-Type": ContentJSON,
		}).
		SetBody(map[string]string{
			"repository_url": material.RepoURL,
		}).
		Post(fmt.Sprintf(MaterialNotifyEndpoint, material.Type))
	if err != nil {
		return "", &errors.APIError{Err: err, Message: fmt.Sprintf("notify material '%s' of type %s", material.RepoURL, material.Type)}
	}

	if resp.StatusCode() != http.StatusAccepted {
		return "", &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	notifyMessage := map[string]string{}
	if err = json.Unmarshal(resp.Body(), &notifyMessage); err != nil {
		return "", &errors.MarshalError{Err: err}
	}

	return notifyMessage["message"], nil
}

func (conf *client) MaterialTriggerUpdate(materialID string) (map[string]string, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":      HeaderVersionZero,
			HeaderConfirm: "true",
		}).
		Post(fmt.Sprintf(MaterialTriggerUpdate, materialID))
	if err != nil {
		return nil, &errors.APIError{Err: err, Message: fmt.Sprintf("trigger update '%s'", materialID)}
	}

	if resp.StatusCode() != http.StatusCreated {
		return nil, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	updateMessage := map[string]string{}
	if err = json.Unmarshal(resp.Body(), &updateMessage); err != nil {
		return nil, &errors.MarshalError{Err: err}
	}

	return updateMessage, nil
}
