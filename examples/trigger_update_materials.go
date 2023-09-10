package main

import (
	"fmt"
	"log"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"gopkg.in/yaml.v3"
)

func main() {
	auth := gocd.Auth{
		UserName: "admin",
		Password: "admin",
	}
	client := gocd.NewClient("http://localhost:8153/go", auth, "info", nil)

	resp, err := client.MaterialTriggerUpdate("3d00c7a0bbe3e425c2ecfac072f1c3ffc14580908c9da0d84a3ec4e5283fca14")
	if err != nil {
		log.Fatalln(err)
	}

	out, err := yaml.Marshal(resp)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(out))
}
