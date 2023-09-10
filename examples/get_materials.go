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

	resp, err := client.GetMaterials()
	if err != nil {
		log.Fatalln(err)
	}

	out, err := yaml.Marshal(resp)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(out))
}
