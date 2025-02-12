// File:		user.go
// Created by:	Hoven
// Created on:	2025-02-10
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package dto

import "github.com/yazl-tech/beauty-rating-server/domain/user"

type User struct {
	UserId  int    `json:"user_id"`
	Name    string `json:"name"`
	Email   string `json:"email,omitempty"`
	Avatar  string `json:"avatar,omitempty"`
	Gender  string `json:"gender,omitempty"`
	Status  int    `json:"status,omitempty"`
	Role    int    `json:"role,omitempty"`
	IsAdmin bool   `json:"is_admin,omitempty"`
}

func UserEntityToDto(u *user.User) *User {
	if u == nil {
		return nil
	}

	dtoUser := &User{
		UserId:  u.ID,
		Name:    u.Name,
		Email:   u.Email,
		Avatar:  u.Avatar,
		Status:  int(u.Status),
		Role:    int(u.Role),
		Gender:  u.Gender.String(),
		IsAdmin: u.Role.IsAdmin(),
	}

	return dtoUser
}

type WechatLoginRequest struct {
	Code     string `json:"code"`
	DeviceId string `header:"X-Device-Id"`
}

type WechatLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         *User  `json:"user"`
}

type UpdateUserNameRequest struct {
	UserName string `json:"username"`
}

type UpdateUserGenderRequest struct {
	Gender int `json:"gender" binding:"oneof=1 2"`
}

type UploadAvatarResponse struct {
	AvatarUrl string `json:"avatar_url"`
}

type GetUserAvatarRequest struct {
	AvatarId string `uri:"avatar_id"`
}
