package main

import (
	"log"
	"os"

	"github.com/nikhilsbhat/gocd-sdk-go"
)

func exampleWithSSL() {
	auth := gocd.Auth{
		UserName: "admin",
		Password: "admin",
	}
	ca, err := os.ReadFile("path/to/ca.pem")
	if err != nil {
		log.Fatal(err)
	}

	client := gocd.NewClient("http://localhost:8153/go", auth, "info", ca)

	if err = client.CommentOnPipeline(gocd.PipelineObject{
		Name:    "sample_pipeline",
		Counter: 1,
		Message: "message to comment",
	}); err != nil {
		log.Fatal("commenting on pipeline errored with", err)
	}
}
