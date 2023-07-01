package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/nikhilsbhat/gocd-sdk-go/pkg/plugin"
)

func main() {
	pluginCfg := plugin.NewPluginConfig("", "", "")
	auth := gocd.Auth{
		UserName: "admin",
		Password: "admin",
	}
	client := gocd.NewClient("http://localhost:8153/go", auth, "info", nil)

	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}

	pipelinePath := filepath.Join(homePath, "gocd-sdk-go/internal/fixtures/sample-pipeline.gocd.yaml")
	success, err := client.ValidatePipelineSyntax(pluginCfg, []string{pipelinePath})
	if err != nil {
		log.Fatalln(err)
	}

	if !success {
		log.Fatalln("pipeline validation errored")
	}
	log.Println("pipeline validation is success")
}
