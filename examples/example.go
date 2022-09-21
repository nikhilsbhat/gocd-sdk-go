package main

import (
	"fmt"
	"log"

	"github.com/nikhilsbhat/gocd-sdk-go"
)

func main() {
	client := gocd.NewClient("http://localhost:8153/go", "admin", "admin", "info", nil)
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
