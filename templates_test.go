package gocd_test

import (
	_ "embed"
	"fmt"
	"net/http"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed internal/fixtures/templates.json
var templatesJSON string

//go:embed internal/fixtures/template.json
var templateJSON string

//go:embed internal/fixtures/template_parameters.json
var templateParametersJSON string

//go:embed internal/fixtures/template_authorization.json
var templateAuthorizationJSON string

func Test_client_GetTemplates(t *testing.T) {
	correctTemplateHeader := map[string]string{"Accept": gocd.HeaderVersionSeven}
	server := mockServer([]byte(templatesJSON), http.StatusOK,
		correctTemplateHeader, false, map[string]string{"ETag": "templates-etag"})
	client := gocd.NewClient(server.URL, auth, "info", nil)

	expected := gocd.Templates{
		Templates: []gocd.Template{
			{
				Name:          "template1",
				CanEdit:       true,
				CanAdminister: true,
				Pipelines: &gocd.TemplatePipelines{
					Pipelines: []gocd.TemplatePipeline{
						{Name: "up42", CanAdminister: true},
					},
				},
			},
		},
		ETAG: "templates-etag",
	}

	actual, err := client.GetTemplates()
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func Test_client_GetTemplate(t *testing.T) {
	correctTemplateHeader := map[string]string{"Accept": gocd.HeaderVersionSeven}
	server := mockServer([]byte(templateJSON), http.StatusOK,
		correctTemplateHeader, false, map[string]string{"ETag": "template-etag"})
	client := gocd.NewClient(server.URL, auth, "info", nil)

	actual, err := client.GetTemplate("template1")
	require.NoError(t, err)
	assert.Equal(t, "template1", actual.Name)
	assert.Equal(t, "defaultStage", actual.Stages[0].Name)
	assert.Equal(t, "defaultJob", actual.Stages[0].Jobs[0].Name)
	assert.Equal(t, "template-etag", actual.ETAG)
}

func Test_client_CreateTemplate(t *testing.T) {
	correctTemplateHeader := map[string]string{"Accept": gocd.HeaderVersionSeven, "Content-Type": gocd.ContentJSON}
	server := mockServer([]byte(templateJSON), http.StatusOK,
		correctTemplateHeader, false, map[string]string{"ETag": "created-etag"})
	client := gocd.NewClient(server.URL, auth, "info", nil)

	actual, err := client.CreateTemplate(gocd.Template{Name: "template1"})
	require.NoError(t, err)
	assert.Equal(t, "template1", actual.Name)
	assert.Equal(t, "created-etag", actual.ETAG)
}

func Test_client_UpdateTemplate(t *testing.T) {
	correctTemplateHeader := map[string]string{
		"Accept":       gocd.HeaderVersionSeven,
		"Content-Type": gocd.ContentJSON,
		"If-Match":     "old-etag",
	}
	server := mockServer([]byte(templateJSON), http.StatusOK,
		correctTemplateHeader, false, map[string]string{"ETag": "updated-etag"})
	client := gocd.NewClient(server.URL, auth, "info", nil)

	actual, err := client.UpdateTemplate(gocd.Template{Name: "template1", ETAG: "old-etag"})
	require.NoError(t, err)
	assert.Equal(t, "template1", actual.Name)
	assert.Equal(t, "updated-etag", actual.ETAG)
}

func Test_client_DeleteTemplate(t *testing.T) {
	correctTemplateHeader := map[string]string{"Accept": gocd.HeaderVersionSeven}
	server := mockServer([]byte(`{"message":"ok"}`), http.StatusOK, correctTemplateHeader, false, nil)
	client := gocd.NewClient(server.URL, auth, "info", nil)

	err := client.DeleteTemplate("template1")
	require.NoError(t, err)
}

func Test_client_GetTemplateParameters(t *testing.T) {
	correctTemplateHeader := map[string]string{"Accept": gocd.HeaderVersionSeven}
	server := mockServer([]byte(templateParametersJSON), http.StatusOK,
		correctTemplateHeader, false, map[string]string{"ETag": "parameters-etag"})
	client := gocd.NewClient(server.URL, auth, "info", nil)

	actual, err := client.GetTemplateParameters("template1")
	require.NoError(t, err)
	assert.Equal(t, gocd.TemplateParameters{
		Name:       "template1",
		Parameters: []string{"resources", "command"},
		ETAG:       "parameters-etag",
	}, actual)
}

func Test_client_GetTemplateAuthorization(t *testing.T) {
	correctTemplateHeader := map[string]string{"Accept": gocd.HeaderVersionOne}
	server := mockServer([]byte(templateAuthorizationJSON), http.StatusOK,
		correctTemplateHeader, false, map[string]string{"ETag": "authorization-etag"})
	client := gocd.NewClient(server.URL, auth, "info", nil)

	actual, err := client.GetTemplateAuthorization("template1")
	require.NoError(t, err)
	assert.Equal(t, gocd.TemplateAuthorization{
		AllGroupAdminsAreViewUsers: true,
		Admin: gocd.AuthorizationConfig{
			Roles: []string{"template-admins"},
			Users: []string{},
		},
		View: gocd.AuthorizationConfig{
			Roles: []string{},
			Users: []string{"template-viewer"},
		},
		ETAG: "authorization-etag",
	}, actual)
}

func Test_client_UpdateTemplateAuthorization(t *testing.T) {
	correctTemplateHeader := map[string]string{
		"Accept":       gocd.HeaderVersionOne,
		"Content-Type": gocd.ContentJSON,
		"If-Match":     "old-etag",
	}
	server := mockServer([]byte(templateAuthorizationJSON), http.StatusOK,
		correctTemplateHeader, false, map[string]string{"ETag": "updated-authorization-etag"})
	client := gocd.NewClient(server.URL, auth, "info", nil)

	actual, err := client.UpdateTemplateAuthorization("template1", gocd.TemplateAuthorization{ETAG: "old-etag"})
	require.NoError(t, err)
	assert.True(t, actual.AllGroupAdminsAreViewUsers)
	assert.Equal(t, "updated-authorization-etag", actual.ETAG)
}

func Test_client_TemplateConfigs_ShouldReturnNonOkErrors(t *testing.T) {
	templateV7Header := map[string]string{"Accept": gocd.HeaderVersionSeven}
	templateV7JSONHeader := map[string]string{"Accept": gocd.HeaderVersionSeven, "Content-Type": gocd.ContentJSON}
	templateV7UpdateHeader := map[string]string{"Accept": gocd.HeaderVersionSeven, "Content-Type": gocd.ContentJSON, "If-Match": "old-etag"}
	templateV1Header := map[string]string{"Accept": gocd.HeaderVersionOne}
	templateV1UpdateHeader := map[string]string{"Accept": gocd.HeaderVersionOne, "Content-Type": gocd.ContentJSON, "If-Match": "old-etag"}

	tests := []struct {
		name   string
		header map[string]string
		call   func(gocd.GoCd) error
	}{
		{
			name:   "get templates",
			header: templateV7Header,
			call: func(client gocd.GoCd) error {
				_, err := client.GetTemplates()
				return err
			},
		},
		{
			name:   "get template",
			header: templateV7Header,
			call: func(client gocd.GoCd) error {
				_, err := client.GetTemplate("template1")
				return err
			},
		},
		{
			name:   "create template",
			header: templateV7JSONHeader,
			call: func(client gocd.GoCd) error {
				_, err := client.CreateTemplate(gocd.Template{Name: "template1"})
				return err
			},
		},
		{
			name:   "update template",
			header: templateV7UpdateHeader,
			call: func(client gocd.GoCd) error {
				_, err := client.UpdateTemplate(gocd.Template{Name: "template1", ETAG: "old-etag"})
				return err
			},
		},
		{
			name:   "delete template",
			header: templateV7Header,
			call: func(client gocd.GoCd) error {
				return client.DeleteTemplate("template1")
			},
		},
		{
			name:   "get template parameters",
			header: templateV7Header,
			call: func(client gocd.GoCd) error {
				_, err := client.GetTemplateParameters("template1")
				return err
			},
		},
		{
			name:   "get template authorization",
			header: templateV1Header,
			call: func(client gocd.GoCd) error {
				_, err := client.GetTemplateAuthorization("template1")
				return err
			},
		},
		{
			name:   "update template authorization",
			header: templateV1UpdateHeader,
			call: func(client gocd.GoCd) error {
				_, err := client.UpdateTemplateAuthorization("template1", gocd.TemplateAuthorization{ETAG: "old-etag"})
				return err
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := mockServer([]byte("template error"), http.StatusBadGateway, test.header, false, nil)
			client := gocd.NewClient(server.URL, auth, "info", nil)

			err := test.call(client)
			expected := fmt.Sprintf(
				"got 502 from GoCD while making %s call for %s\nwith BODY:template error",
				methodForTemplateTest(test.name),
				server.URL+pathForTemplateTest(test.name),
			)
			require.EqualError(t, err, expected)
		})
	}
}

func Test_client_TemplateConfigs_ShouldReturnMarshalErrors(t *testing.T) {
	templateV7Header := map[string]string{"Accept": gocd.HeaderVersionSeven}
	templateV7JSONHeader := map[string]string{"Accept": gocd.HeaderVersionSeven, "Content-Type": gocd.ContentJSON}
	templateV7UpdateHeader := map[string]string{"Accept": gocd.HeaderVersionSeven, "Content-Type": gocd.ContentJSON, "If-Match": "old-etag"}
	templateV1Header := map[string]string{"Accept": gocd.HeaderVersionOne}
	templateV1UpdateHeader := map[string]string{"Accept": gocd.HeaderVersionOne, "Content-Type": gocd.ContentJSON, "If-Match": "old-etag"}

	tests := []struct {
		name   string
		header map[string]string
		call   func(gocd.GoCd) error
	}{
		{
			name:   "get templates",
			header: templateV7Header,
			call: func(client gocd.GoCd) error {
				_, err := client.GetTemplates()
				return err
			},
		},
		{
			name:   "get template",
			header: templateV7Header,
			call: func(client gocd.GoCd) error {
				_, err := client.GetTemplate("template1")
				return err
			},
		},
		{
			name:   "create template",
			header: templateV7JSONHeader,
			call: func(client gocd.GoCd) error {
				_, err := client.CreateTemplate(gocd.Template{Name: "template1"})
				return err
			},
		},
		{
			name:   "update template",
			header: templateV7UpdateHeader,
			call: func(client gocd.GoCd) error {
				_, err := client.UpdateTemplate(gocd.Template{Name: "template1", ETAG: "old-etag"})
				return err
			},
		},
		{
			name:   "get template parameters",
			header: templateV7Header,
			call: func(client gocd.GoCd) error {
				_, err := client.GetTemplateParameters("template1")
				return err
			},
		},
		{
			name:   "get template authorization",
			header: templateV1Header,
			call: func(client gocd.GoCd) error {
				_, err := client.GetTemplateAuthorization("template1")
				return err
			},
		},
		{
			name:   "update template authorization",
			header: templateV1UpdateHeader,
			call: func(client gocd.GoCd) error {
				_, err := client.UpdateTemplateAuthorization("template1", gocd.TemplateAuthorization{ETAG: "old-etag"})
				return err
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := mockServer([]byte("templateJSON"), http.StatusOK, test.header, false, nil)
			client := gocd.NewClient(server.URL, auth, "info", nil)

			err := test.call(client)
			require.EqualError(t, err, "reading response body errored with: invalid character 'e' in literal true (expecting 'r')")
		})
	}
}

func Test_client_TemplateConfigs_ShouldReturnAPIErrors(t *testing.T) {
	client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)
	client.SetRetryCount(1)
	client.SetRetryWaitTime(1)

	tests := []struct {
		name    string
		message string
		call    func(gocd.GoCd) error
	}{
		{
			name:    "get templates",
			message: "get templates",
			call: func(client gocd.GoCd) error {
				_, err := client.GetTemplates()
				return err
			},
		},
		{
			name:    "get template",
			message: "get template 'template1'",
			call: func(client gocd.GoCd) error {
				_, err := client.GetTemplate("template1")
				return err
			},
		},
		{
			name:    "create template",
			message: "create template 'template1'",
			call: func(client gocd.GoCd) error {
				_, err := client.CreateTemplate(gocd.Template{Name: "template1"})
				return err
			},
		},
		{
			name:    "update template",
			message: "update template 'template1'",
			call: func(client gocd.GoCd) error {
				_, err := client.UpdateTemplate(gocd.Template{Name: "template1", ETAG: "old-etag"})
				return err
			},
		},
		{
			name:    "delete template",
			message: "delete template 'template1'",
			call: func(client gocd.GoCd) error {
				return client.DeleteTemplate("template1")
			},
		},
		{
			name:    "get template parameters",
			message: "get template parameters 'template1'",
			call: func(client gocd.GoCd) error {
				_, err := client.GetTemplateParameters("template1")
				return err
			},
		},
		{
			name:    "get template authorization",
			message: "get template authorization 'template1'",
			call: func(client gocd.GoCd) error {
				_, err := client.GetTemplateAuthorization("template1")
				return err
			},
		},
		{
			name:    "update template authorization",
			message: "update template authorization 'template1'",
			call: func(client gocd.GoCd) error {
				_, err := client.UpdateTemplateAuthorization("template1", gocd.TemplateAuthorization{ETAG: "old-etag"})
				return err
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.call(client)
			require.Contains(t, err.Error(), fmt.Sprintf("call made to %s errored with:", test.message))
			require.Contains(t, err.Error(), "localhost:8156")
		})
	}
}

func methodForTemplateTest(name string) string {
	switch name {
	case "create template":
		return http.MethodPost
	case "update template", "update template authorization":
		return http.MethodPut
	case "delete template":
		return http.MethodDelete
	default:
		return http.MethodGet
	}
}

func pathForTemplateTest(name string) string {
	switch name {
	case "get templates", "create template":
		return "/api/admin/templates"
	case "get template", "update template", "delete template":
		return "/api/admin/templates/template1"
	case "get template parameters":
		return "/api/admin/templates/template1/parameters"
	case "get template authorization", "update template authorization":
		return "/api/admin/templates/template1/authorization"
	default:
		return ""
	}
}
