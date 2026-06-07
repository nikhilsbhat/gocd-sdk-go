package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"
)

func (conf *client) GetTemplates() (Templates, error) {
	newClient := conf.clone()

	var templatesCfg TemplatesConfig

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionSeven,
		}).
		Get(TemplateConfigEndpoint)
	if err != nil {
		return Templates{}, &errors.APIError{Err: err, Message: "get templates"}
	}

	if resp.StatusCode() != http.StatusOK {
		return Templates{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &templatesCfg); err != nil {
		return Templates{}, &errors.MarshalError{Err: err}
	}

	templatesCfg.Templates.ETAG = resp.Header().Get("ETag")

	return templatesCfg.Templates, nil
}

func (conf *client) GetTemplate(name string) (Template, error) {
	var template Template

	newClient := conf.clone()

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionSeven,
		}).
		Get(filepath.Join(TemplateConfigEndpoint, name))
	if err != nil {
		return template, &errors.APIError{Err: err, Message: fmt.Sprintf("get template '%s'", name)}
	}

	if resp.StatusCode() != http.StatusOK {
		return template, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &template); err != nil {
		return template, &errors.MarshalError{Err: err}
	}

	template.ETAG = resp.Header().Get("ETag")

	return template, nil
}

func (conf *client) CreateTemplate(config Template) (Template, error) {
	var template Template

	newClient := conf.clone()

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionSeven,
			"Content-Type": ContentJSON,
		}).
		SetBody(config).
		Post(TemplateConfigEndpoint)
	if err != nil {
		return template, &errors.APIError{Err: err, Message: fmt.Sprintf("create template '%s'", config.Name)}
	}

	if resp.StatusCode() != http.StatusOK {
		return template, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &template); err != nil {
		return template, &errors.MarshalError{Err: err}
	}

	template.ETAG = resp.Header().Get("ETag")

	return template, nil
}

func (conf *client) UpdateTemplate(config Template) (Template, error) {
	var template Template

	newClient := conf.clone()

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionSeven,
			"Content-Type": ContentJSON,
			"If-Match":     config.ETAG,
		}).
		SetBody(config).
		Put(filepath.Join(TemplateConfigEndpoint, config.Name))
	if err != nil {
		return template, &errors.APIError{Err: err, Message: fmt.Sprintf("update template '%s'", config.Name)}
	}

	if resp.StatusCode() != http.StatusOK {
		return template, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &template); err != nil {
		return template, &errors.MarshalError{Err: err}
	}

	template.ETAG = resp.Header().Get("ETag")

	return template, nil
}

func (conf *client) DeleteTemplate(name string) error {
	newClient := conf.clone()

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionSeven,
		}).
		Delete(filepath.Join(TemplateConfigEndpoint, name))
	if err != nil {
		return &errors.APIError{Err: err, Message: fmt.Sprintf("delete template '%s'", name)}
	}

	if resp.StatusCode() != http.StatusOK {
		return &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	return nil
}

func (conf *client) GetTemplateParameters(name string) (TemplateParameters, error) {
	var parameters TemplateParameters

	newClient := conf.clone()

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionSeven,
		}).
		Get(filepath.Join(TemplateConfigEndpoint, name, "parameters"))
	if err != nil {
		return parameters, &errors.APIError{Err: err, Message: fmt.Sprintf("get template parameters '%s'", name)}
	}

	if resp.StatusCode() != http.StatusOK {
		return parameters, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &parameters); err != nil {
		return parameters, &errors.MarshalError{Err: err}
	}

	parameters.ETAG = resp.Header().Get("ETag")

	return parameters, nil
}

func (conf *client) GetTemplateAuthorization(name string) (TemplateAuthorization, error) {
	var authorization TemplateAuthorization

	newClient := conf.clone()

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		Get(filepath.Join(TemplateConfigEndpoint, name, "authorization"))
	if err != nil {
		return authorization, &errors.APIError{Err: err, Message: fmt.Sprintf("get template authorization '%s'", name)}
	}

	if resp.StatusCode() != http.StatusOK {
		return authorization, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &authorization); err != nil {
		return authorization, &errors.MarshalError{Err: err}
	}

	authorization.ETAG = resp.Header().Get("ETag")

	return authorization, nil
}

func (conf *client) UpdateTemplateAuthorization(name string, authorization TemplateAuthorization) (TemplateAuthorization, error) {
	var templateAuthorization TemplateAuthorization

	newClient := conf.clone()

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionOne,
			"Content-Type": ContentJSON,
			"If-Match":     authorization.ETAG,
		}).
		SetBody(authorization).
		Put(filepath.Join(TemplateConfigEndpoint, name, "authorization"))
	if err != nil {
		return templateAuthorization, &errors.APIError{Err: err, Message: fmt.Sprintf("update template authorization '%s'", name)}
	}

	if resp.StatusCode() != http.StatusOK {
		return templateAuthorization, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &templateAuthorization); err != nil {
		return templateAuthorization, &errors.MarshalError{Err: err}
	}

	templateAuthorization.ETAG = resp.Header().Get("ETag")

	return templateAuthorization, nil
}
