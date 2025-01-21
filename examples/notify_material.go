package main

import (
	"fmt"
	"log"

	"github.com/nikhilsbhat/gocd-sdk-go"
)

func notifyMaterial() {
	auth := gocd.Auth{
		UserName: "admin",
		Password: "admin",
	}
	client := gocd.NewClient("http://localhost:8153/go", auth, "info", nil)

	material := gocd.Material{
		Type:    "git",
		RepoURL: "https://github.com/nikhilsbhat/helm-drift.git",
	}

	resp, err := client.NotifyMaterial(material)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(resp)
}
