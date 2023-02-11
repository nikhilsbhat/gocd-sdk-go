package main

import (
	"fmt"
	"log"

	"github.com/nikhilsbhat/gocd-sdk-go"
)

func exampleMain() {
	auth := gocd.Auth{
		UserName: "admin",
		Password: "admin",
	}
	client := gocd.NewClient("http://localhost:8153/go", auth, "info", nil)
	fmt.Println(environments(client))
	fmt.Println(configRepos(client))
}

func environments(client gocd.GoCd) []gocd.Environment {
	env, err := client.GetEnvironments()
	if err != nil {
		log.Fatal(err)
	}
	return env
}

func configRepos(client gocd.GoCd) []gocd.ConfigRepo {
	repos, err := client.GetConfigRepos()
	if err != nil {
		log.Fatal(err)
	}
	return repos
}
