package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nikhilsbhat/gocd-sdk-go"
)

func main() {
	client := gocd.NewClient("http://localhost:8153/go", "admin", "admin", "info", nil)
	secretsConfig := gocd.CommonConfig{
		ID:          "sample-kube-secret-config",
		PluginID:    "cd.go.contrib.secrets.kubernetes",
		Description: "sample secret config",
		Properties: []gocd.PluginConfiguration{
			{
				Key:   "kubernetes_secret_name",
				Value: "ci_secret",
			},
			{
				Key:   "kubernetes_cluster_url",
				Value: "https://0.0.0.0:64527",
			},
			{
				Key:   "security_token",
				Value: "dGVoZHhuZGhueGl3dWRoZnhua2R3aGZuaXdkdWZobnh3",
			},
			{
				Key:   "kubernetes_cluster_ca_cert",
				Value: "dGVoZHhuZGhueGl3dWRoZnhua2R3aGZuaXdkdWZobnh3",
			},
			{
				Key:   "namespace",
				Value: "default",
			},
		},
		Rules: []map[string]string{
			{
				"directive": "allow",
				"action":    "refer",
				"type":      "*",
				"resource":  "*",
			},
		},
	}

	fmt.Println(createSecretsConfig(client, secretsConfig))
	jsonOut, err := json.MarshalIndent(getSecretsConfig(client, secretsConfig.ID), " ", " ")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(string(jsonOut))
	}
}

func getSecretsConfig(client gocd.GoCd, id string) gocd.CommonConfig {
	clusterProfile, err := client.GetSecretConfig(id)
	if err != nil {
		log.Fatal(err)
	}
	return clusterProfile
}

func createSecretsConfig(client gocd.GoCd, config gocd.CommonConfig) gocd.CommonConfig {
	out, err := json.MarshalIndent(config, " ", " ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
	secretConfig, err := client.CreateSecretConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	return secretConfig
}
