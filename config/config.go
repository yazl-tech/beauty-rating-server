// File:		config.go
// Created by:	Hoven
// Created on:	2025-02-13
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package config

import (
	"errors"

	"github.com/go-puzzles/puzzles/putils"
)

type BeautyConfig struct {
	ApiHost        string
	ApiPrefix      string
	ApiVersion     string
	ApiPort        int
	AuthCoreSrv    string
	TokenKey       string
	ShareSecretKey string
}

func (bc *BeautyConfig) SetDefault() {
	if bc.ApiHost == "" {
		bc.ApiHost = "localhost:28084"
	}

	if bc.ApiPrefix == "" {
		bc.ApiPrefix = "/api"
	}

	if bc.ApiVersion == "" {
		bc.ApiVersion = "/v1"
	}

	if bc.ApiPort == 0 {
		bc.ApiPort = 8080
	}

	if bc.AuthCoreSrv == "" {
		bc.AuthCoreSrv = "auth-core"
	}

	if bc.ShareSecretKey == "" {
		bc.ShareSecretKey = putils.RandString(7)
	}
}

func (bc *BeautyConfig) Validate() error {
	if bc.TokenKey == "" {
		return errors.New("missing tokenKey")
	}

	return nil
}
