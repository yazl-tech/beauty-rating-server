// File:		config.go
// Created by:	Hoven
// Created on:	2025-02-13
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package config

import "errors"

type BeautyConfig struct {
	ApiHost     string
	ApiPrefix   string
	ApiVersion  string
	ApiPort     int
	AuthCoreSrv string
	TokenKey    string
}

func (bc *BeautyConfig) SetDefault() {
	bc.ApiHost = "localhost:28084"
	bc.ApiPrefix = "/api"
	bc.ApiVersion = "/v1"
	bc.ApiPort = 28084
	bc.AuthCoreSrv = "auth-core"
}

func (bc *BeautyConfig) Validate() error {
	if bc.TokenKey == "" {
		return errors.New("missing tokenKey")
	}

	return nil
}
