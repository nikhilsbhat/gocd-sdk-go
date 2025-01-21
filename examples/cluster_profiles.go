package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nikhilsbhat/gocd-sdk-go"
)

func clusterProfiles() {
	auth := gocd.Auth{
		UserName: "admin",
		Password: "admin",
	}
	client := gocd.NewClient("http://localhost:8153/go", auth, "info", nil)
	clusterProfileConfig := gocd.CommonConfig{
		ID:       "sample-cluster-profile",
		PluginID: "cd.go.contrib.elasticagent.kubernetes",
		Properties: []gocd.PluginConfiguration{
			{
				Key:   "go_server_url",
				Value: "https://gocd.sample.com/go",
			},
			{
				Key:   "auto_register_timeout",
				Value: "15",
			},
			{
				Key:   "pending_pods_count",
				Value: "2",
			},
			{
				Key:   "kubernetes_cluster_url",
				Value: "https://0.0.0.0:64527",
			},
			{
				Key:   "namespace",
				Value: "default",
			},
			{
				Key:   "security_token",
				Value: "dGVoZHhuZGhueGl3dWRoZnhua2R3aGZuaXdkdWZobnh3",
			},
			{
				Key:   "kubernetes_cluster_ca_cert",
				Value: "dGVoZHhuZGhueGl3dWRoZnhua2R3aGZuaXdkdWZobnh3",
			},
		},
		ETAG: "",
	}

	createClusterProfiles(client, clusterProfileConfig)
	jsonOut, err := json.MarshalIndent(getClusterProfiles(client, clusterProfileConfig.ID), " ", " ")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(string(jsonOut))
	}
}

func getClusterProfiles(client gocd.GoCd, id string) gocd.CommonConfig {
	clusterProfile, err := client.GetClusterProfile(id)
	if err != nil {
		log.Fatal(err)
	}
	return clusterProfile
}

func createClusterProfiles(client gocd.GoCd, config gocd.CommonConfig) gocd.CommonConfig {
	out, err := json.MarshalIndent(config, " ", " ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
	clusterProfile, err := client.CreateClusterProfile(config)
	if err != nil {
		log.Fatal(err)
	}
	return clusterProfile
}
