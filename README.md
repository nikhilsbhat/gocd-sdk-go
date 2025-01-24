# GoCD Golang SDK

[![Go Report Card](https://goreportcard.com/badge/github.com/nikhilsbhat/gocd-sdk-go)](https://goreportcard.com/report/github.com/nikhilsbhat/gocd-sdk-go)
[![shields](https://img.shields.io/badge/license-MIT-blue)](https://github.com/nikhilsbhat/gocd-sdk-go/blob/master/LICENSE)
[![shields](https://godoc.org/github.com/nikhilsbhat/gocd-sdk-go?status.svg)](https://godoc.org/github.com/nikhilsbhat/gocd-sdk-go)
[![shields](https://img.shields.io/github/v/tag/nikhilsbhat/gocd-sdk-go.svg)](https://github.com/nikhilsbhat/gocd-sdk-go/tags)

Golang client library for [GoCD API](https://api.gocd.org/current/) (Supports Most of the APIs).

## Introduction

This Library could be helpful while building any tools around GoCD or while interacting with GoCD to perform certain
daily activities.

This could include checking the health of all agents connected to GoCD or status of a job and many more.

## Installation

Get the latest version of GoCD sdk using `go get` command. Example:

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

- [x] [Agents](https://api.gocd.org/current/#agents)
    - [x] Get All Agents
    - [x] Get Specific Agent
    - [x] Update Agent
    - [x] Update Agents bulk
    - [x] Delete Agent
    - [x] Delete Agents bulk
    - [x] Kill running tasks iin Agent
    - [x] Agent job run history
- [x] [ConfigRepo](https://api.gocd.org/current/#config-repo)
    - [x] Get All Config repo
    - [x] Get Specific Config repo
    - [x] Create Config repo
    - [x] Update Config repo
    - [x] Delete Config repo
    - [x] Get Config repo status
    - [x] Trigger config repo update
    - [x] Preflight check of config repo configurations
    - [x] Export pipeline config to config repo format
    - [x] Definitions defined in config repo
- [x] [Maintenance Mode](https://api.gocd.org/current/#maintenance-mode)
    - [x] Enable Maintenance Mode
    - [x] Disable Maintenance Mode
    - [x] Get Maintenance Mode info
- [x] [PipelineGroup](https://api.gocd.org/current/#pipeline-group-config)
    - [x] Get All pipeline groups
    - [x] Get specific pipeline group
    - [x] Update pipeline group
    - [x] Create pipeline group
    - [x] Delete pipeline group
- [x] [Environment Config](https://api.gocd.org/current/#environment-config)
    - [x] Get All Environments
    - [x] Get specific Environment
    - [x] Create Environment
    - [x] Update Environment
    - [x] Patch Environment
    - [x] Delete Environment
- [x] [Backup-config](https://api.gocd.org/current/#backup-config)
    - [x] Get Backup Info
    - [x] Create or Update Backup
    - [x] Delete Backup Info
- [x] [Backup](https://api.gocd.org/current/#backups)
    - [x] Schedule Backup
    - [x] Get Backup
- [ ] [Pipeline](https://api.gocd.org/current/#pipelines)
    - [x] Get pipeline status
    - [x] Pause Pipeline
    - [x] UnPause Pipeline
    - [x] UnLock Pipeline
    - [x] Schedule Pipeline
    - [x] Get Pipeline Schedules
    - [ ] Compare pipeline instances
- [x] [Pipeline Instances](https://api.gocd.org/current/#pipeline-instances)
    - [x] Get Pipeline Instance
    - [x] Get Pipeline History
    - [x] Comment on Pipeline
- [x] [Pipeline Config](https://api.gocd.org/current/#pipeline-config)
    - [x] Get pipeline config
    - [x] Edit pipeline config
    - [x] Create a pipeline
    - [x] Delete a pipeline
    - [x] Extract template from pipeline
    - [x] Validate pipeline config syntax
- [ ] [Stage Instances](https://api.gocd.org/current/#stage-instances)
    - [x] Cancel stage
    - [ ] Get stage instance
    - [ ] Get stage history
    - [x] Run failed jobs
    - [x] Run selected jobs
- [x] [Stages](https://api.gocd.org/current/#stages)
    - [x] Run stage
- [ ] [Jobs](https://api.gocd.org/current/#jobs)
    - [ ] Get job instance
    - [ ] Get job history
- [ ] [Feeds](https://api.gocd.org/current/#feeds)
    - [x] Get All pipelines
    - [ ] Get Pipeline
    - [ ] Get Stage
    - [ ] Get Job
    - [ ] Get Material
    - [x] Scheduled Jobs
- [x] [Artifact Config](https://api.gocd.org/current/#artifacts-config)
    - [x] Get Artifact Config
    - [x] Update Artifact Config
- [x] [Artifact Store](https://api.gocd.org/current/#artifact-store)
    - [x] Get Artifact Stores
    - [x] Get Artifact Store
    - [x] Create Artifact Stores
    - [x] Update Artifact Stores
    - [x] Delete Artifact Stores
- [x] [Cluster Profiles](https://api.gocd.org/current/#cluster-profiles)
    - [x] Get Cluster Profiles
    - [x] Get Cluster Profile
    - [x] Create Cluster Profile
    - [x] Update Cluster Profile
    - [x] Delete Cluster Profile
- [x] [Elastic Agent Profiles](https://api.gocd.org/current/#elastic-agent-profiles)
    - [x] Get Elastic Agent Profiles
    - [x] Get Elastic Agent Profile
    - [x] Create Elastic Agent Profile
    - [x] Update Elastic Agent Profile
    - [x] Delete Elastic Agent Profile
    - [x] Get Elastic Agent Profile Usage
- [x] [Secret Configs](https://api.gocd.org/current/#secret-configs)
    - [x] Get Secret Configs
    - [x] Get Secret Config
    - [x] Create Secret Config
    - [x] Update Secret Config
    - [x] Delete Secret Config
- [x] [Package Repositories](https://api.gocd.org/current/#package-repositories)
    - [x] Get Package repositories
    - [x] Get Package Repository
    - [x] Create Package Repository
    - [x] Update Package Repository
    - [x] Delete Package Repository
- [x] [Package](https://api.gocd.org/current/#packages)
    - [x] Get Package Materials
    - [x] Get Package Material
    - [x] Create Package Material
    - [x] Update Package Material
    - [x] Delete Package Material
- [x] [Materials](hhttps://api.gocd.org/current/#materials)
    - [x] Get All Materials
    - [x] Get Materials Usage
    - [ ] Get material modifications
- [x] [Site URL](https://api.gocd.org/current/#siteurls-config)
    - [x] Get Site URL
    - [x] Create or Update Site URL
- [x] [Mail server config](https://api.gocd.org/current/#mailserver-config)
    - [x] Get Mail server config
    - [x] Create or Update Mail server config
    - [x] Update Mail server config
- [x] [Default Job timeout](https://api.gocd.org/current/#default-job-timeout)
    - [x] Get Default Job timeout
    - [x] Update Default Job timeout
- [x] [Plugin settings](https://api.gocd.org/current/#plugin-settings)
    - [x] Get Plugin settings
    - [x] Create Plugin settings
    - [x] Update Plugin settings
- [x] [Plugin Info](https://api.gocd.org/current/#plugin-info)
    - [x] Get all plugin info
    - [x] Get plugin info
- [x] [Auth Configs](https://api.gocd.org/current/#authorization-configuration)
    - [x] Get All Auth configs
    - [x] Get Specific Auth config
    - [x] Create Auth config
    - [x] Update Auth config
    - [x] Delete Auth config
- [x] [System Admin](https://api.gocd.org/current/#system-admins)
    - [x] Get All system admins
    - [x] Update system Admin
    - [x] Bulk update system admins
- [x] [Role](https://api.gocd.org/current/#roles)
    - [x] Get all roles
    - [X] Get all roles by type
    - [x] Get Specific role
    - [X] Create a GoCD role
    - [X] Create a plugin role
    - [X] Update a role
    - [x] Delete a role
    - [ ] <del>Bulk update roles<del>
- [ ] [access-tokens](https://api.gocd.org/current/#access-tokens)
    - [ ] Get all tokens for current user
    - [ ] Get one token for current user
    - [ ] Create token for current user
    - [ ] Revoke token for current user
    - [ ] Get all tokens for all users
    - [ ] Get one token for any user
    - [ ] Revoke token for any user
- [x] [current-user](https://api.gocd.org/current/#current-user)
    - [x] Get current user
    - [x] Update current user info
- [x] [Local Users](https://api.gocd.org/current/#users)
    - [x] Get all users
    - [x] Get a user
    - [x] Create a user
    - [x] Update a user
    - [x] Delete a user
    - [x] Bulk delete users
    - [x] Bulk enable/disable users
- [ ] [Notification Filter](https://api.gocd.org/current/#notification-filters)
    - [ ] Get all notification filters
    - [ ] Get a notification filter
    - [ ] Create a notification filter
    - [ ] Update a notification filter
    - [ ] Delete a notification filter
- [x] [Server Health Messages](https://api.gocd.org/current/#server-health-messages)
    - [x] Get Server Health messages
- [x] [Version](https://api.gocd.org/current/#version)
    - [x] Get Version
- [x] [Encryption](https://api.gocd.org/current/#encryption)
    - [x] Encrypt plain text value
    - [x] Decrypt encrypted text value
- [x] [Permission](https://api.gocd.org/current/#permissions)
    - [x] Show permissions one has

## Enhancements

If any of the APIs are missed, feel free to raise the PR or create issues for the same.