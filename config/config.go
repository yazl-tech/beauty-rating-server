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
	"github.com/yazl-tech/beauty-rating-server/pkg/analyst"
)

type BeautyConfig struct {
	ApiHost        string
	ApiPrefix      string
	ApiVersion     string
	ApiPort        int
	AuthCoreSrv    string
	TokenKey       string
	ShareSecretKey string
	AiModel        string
	AiBotSrv       string
	AnalystWeights map[analyst.AnalystType]int
}

func (bc *BeautyConfig) AnalystWeight(at analyst.AnalystType) int {
	return bc.AnalystWeights[at]
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

	if bc.AiBotSrv == "" {
		bc.AiBotSrv = "ai-bot"
	}

	if bc.ShareSecretKey == "" {
		bc.ShareSecretKey = putils.RandString(7)
	}

	if bc.AnalystWeights == nil {
		bc.AnalystWeights = map[analyst.AnalystType]int{
			analyst.TypeMock: 80,
			analyst.TypeAi:   20,
		}
	}
}

func (bc *BeautyConfig) Validate() error {
	if bc.TokenKey == "" {
		return errors.New("missing tokenKey")
	}

	if bc.AiModel == "" {
		return errors.New("missing aiModel")
	}

	return nil
}
