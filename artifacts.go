package gocd

import (
	"encoding/json"
	"net/http"

	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"

	"github.com/jinzhu/copier"
)

// UpdateArtifactConfig updates the artifact config with the latest config provided.
func (conf *client) UpdateArtifactConfig(info ArtifactInfo) (ArtifactInfo, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return ArtifactInfo{}, err
	}

	var artifactInfo ArtifactInfo
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionOne,
			"Content-Type": ContentJSON,
			"If-Match":     info.ETAG,
		}).
		SetBody(info).
		Post(ArtifactInfoEndpoint)
	if err != nil {
		return ArtifactInfo{}, &errors.APIError{Err: err, Message: "update artifacts info"}
	}

	if resp.StatusCode() != http.StatusOK {
		return ArtifactInfo{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &artifactInfo); err != nil {
		return ArtifactInfo{}, &errors.MarshalError{Err: err}
	}

	return artifactInfo, nil
}

// GetArtifactConfig fetches the latest artifact config available from GoCD.
func (conf *client) GetArtifactConfig() (ArtifactInfo, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return ArtifactInfo{}, err
	}

	var artifactInfo ArtifactInfo
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		Get(ArtifactInfoEndpoint)
	if err != nil {
		return ArtifactInfo{}, &errors.APIError{Err: err, Message: "get artifacts info"}
	}

	if resp.StatusCode() != http.StatusOK {
		return ArtifactInfo{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &artifactInfo); err != nil {
		return ArtifactInfo{}, &errors.MarshalError{Err: err}
	}

	artifactInfo.ETAG = resp.Header().Get("ETag")

	return artifactInfo, nil
}
