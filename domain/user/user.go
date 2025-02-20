// File:		user.go
// Created by:	Hoven
// Created on:	2025-02-09
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package user

import sdkHttpHandler "gitea.hoven.com/core/auth-core/pkg/sdk/handler"

type WechatConfig struct {
	AppMap   map[string]*sdkHttpHandler.WechatAppSecret
	AppId    string
	SecretId string
}

func (wc *WechatConfig) GetWechatAppConfig(appName string) *sdkHttpHandler.WechatAppSecret {
	if app, ok := wc.AppMap[appName]; ok {
		return app
	}

	return &sdkHttpHandler.WechatAppSecret{
		AppId:    wc.AppId,
		SecretId: wc.SecretId,
	}
}

type BaseInfo struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	Password string `json:"password"`
}

type User struct {
	ID     int
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Gender Gender `json:"gender"`
	Email  string `json:"email"`
	Status Status `json:"status"`
	Role   Role   `json:"role"`
}

type Token struct {
	UserID       int
	AccessToken  string
	RefreshToken string
}

type Gender int

const (
	GenderUnknown Gender = iota
	GenderMale
	GenderFemale
)

func (g Gender) String() string {
	switch g {
	case GenderUnknown:
		return "未知"
	case GenderMale:
		return "男"
	case GenderFemale:
		return "女"
	default:
		return "未知"
	}
}

type Status int

const (
	StatusUnknown Status = iota
	StatusActive
	StatusInactive
)

func (s Status) ID() int {
	return int(s)
}

type Role int

const (
	RoleUnknown Role = iota
	RoleUser
	RoleAdmin
	RolePro
)

func (r Role) ID() int {
	return int(r)
}

func (r Role) IsAdmin() bool {
	return r == RoleAdmin
}

func (r Role) IsPro() bool {
	return r == RolePro
}
