# GoCD Golang SDK

[![Go Report Card](https://goreportcard.com/badge/github.com/nikhilsbhat/gocd-sdk-go)](https://goreportcard.com/report/github.com/nikhilsbhat/gocd-sdk-go)
[![shields](https://img.shields.io/badge/license-MIT-blue)](https://github.com/nikhilsbhat/gocd-sdk-go/blob/master/LICENSE)
[![shields](https://godoc.org/github.com/nikhilsbhat/gocd-sdk-go?status.svg)](https://godoc.org/github.com/nikhilsbhat/gocd-sdk-go)

Golang client library for [GoCD API](https://api.gocd.org/current/) (Not all the API is supported).

## Introduction

This Library could be helpful while building any tools around GoCD or while interacting with GoCD to perform certain
daily activities.

This could include checking the health of all agents connected to GoCD or status of a job and many more.

## Installation

Get latest version of GoCD sdk using `go get` command. Example:

```shell
go get github.com/nikhilsbhat/gocd-sdk-go@latest
```

Get specific version of the same. Example:

```shell
go get github.com/nikhilsbhat/gocd-sdk-go@v0.0.2
```

## Usage

```shell
package main

import (
	"fmt"
	"log"

	"github.com/nikhilsbhat/gocd-sdk-go"
)

func main() {
	client := gocd.NewClient("http://localhost:8153/go", "admin", "admin", "info", nil)
	env, err := client.GetEnvironments()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(env)
}
```
More examples can be found [here](https://github.com/nikhilsbhat/gocd-sdk-go/tree/master/examples).

## Supported APIs
Below are the list of supported APIs:

- [x] Agents
    - [x] Get All Agents
    - [x] Get Specific Agent
    - [x] Update Agent
    - [x] Update Agents bulk
    - [x] Delete Agent
    - [x] Delete Agents bulk
    - [x] Kill running tasks iin Agent
    - [x] Agent job run history
- [ ] ConfigRepo
    - [x] Get All Config repo
    - [x] Get Specific Config repo
    - [x] Create Config repo
    - [x] Update Config repo
    - [x] Delete Config repo
    - [x] Get Config repo status
    - [x] Trigger config repo update
    - [ ] Export pipeline config to config repo format
    - [ ] Preflight check of config repo configurations
    - [ ] Definitions defined in config repo
- [x] Maintenance Mode
    - [x] Enable Maintenance Mode
    - [x] Disable Maintenance Mode
    - [x] Get Maintenance Mode info
- [x] PipelineGroup
    - [x] Get All pipeline groups
    - [x] Get specific pipeline group
    - [x] Update pipeline group
    - [x] Create pipeline group
    - [x] Delete pipeline group
- [x] Environment Config
    - [x] Get All Environments
    - [x] Get specific Environment
    - [x] Create Environment
    - [x] Update Environment
    - [x] Patch Environment
    - [x] Delete Environment
- [x] Backup
    - [x] Get Backup Info
    - [x] Create or Update Backup
    - [x] Delete Backup Info
- [ ] Pipeline
    - [x] Get pipeline status
    - [x] Pause Pipeline
    - [x] UnPause Pipeline
    - [x] UnLock Pipeline
    - [x] Schedule Pipeline
    - [ ] Compare pipeline instances
- [x] Pipeline Instances
    - [x] Get Pipeline Instance
    - [x] Get Pipeline History
    - [x] Comment on Pipeline
- [ ] Pipeline Config
    - [ ] Get pipeline config
    - [ ] Edit pipeline config
    - [ ] Create a pipeline
    - [ ] Delete a pipeline
    - [ ] Extract template from pipeline
- [ ] Stage Instances
    - [ ] Cancel stage
    - [ ] Get stage instance
    - [ ] Get stage history
    - [ ] Run failed jobs
    - [ ] Run selected jobs
- [ ] Stages
    - [ ] Run stage
- [ ] Jobs
    - [ ] Get job instance
    - [ ] Get job history
- [x] Feeds
    - [x] Get All pipelines
    - [ ]
- [x] Artifact Config
    - [x] Get Artifact Config
    - [x] Update Artifact Config
- [x] Site URL
    - [x] Get Site URL
    - [x] Create or Update Site URL
- [x] Mail server config
    - [x] Get Mail server config
    - [x] Create or Update Mail server config
    - [x] Update Mail server config
- [x] Default Job timeout
    - [x] Get Default Job timeout
    - [x] Update Default Job timeout
- [x] Plugin settings
    - [x] Get Plugin settings
    - [x] Create Plugin settings
    - [x] Update Plugin settings
- [ ] Plugin Info
    - [ ] Get all plugin info
    - [ ] Get plugin info
- [x] Auth Configs
  - [x] Get All Auth configs
  - [x] Get Specific Auth config
  - [x] Create Auth config
  - [x] Update Auth config
  - [x] Delete Auth config
- [ ] System Admin
    - [x] Get All system admins
    - [x] Update system Admin
    - [ ] Bulk update system admins
- [ ] Roles
    - [ ] Get all roles
    - [ ] Get all roles by type
    - [ ] Get Specific role
    - [ ] Create a GoCD role
    - [ ] Create a plugin role
    - [ ] Update a role
    - [ ] Delete a role
    - [ ] Bulk update roles
- [x] Server Health Messages
    - [x] Get Server Health messages
- [x] Version
    - [x] Get Version
- [x] Encryption
    - [x] Encrypt plain text value
