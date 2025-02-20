// File:		user.go
// Created by:	Hoven
// Created on:	2025-02-09
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package handler

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/go-puzzles/puzzles/pgin"
	"github.com/yazl-tech/beauty-rating-server/domain/user"
	"github.com/yazl-tech/beauty-rating-server/pkg/exception"
	"github.com/yazl-tech/beauty-rating-server/service/dto"
)

type UserHandlerApp interface {
	WechatLogin(ctx context.Context, code, deviceId, appName string) (*dto.WechatLoginResponse, error)
	GetUserProfile(ctx context.Context) (*user.User, error)
	UpdateUserName(ctx context.Context, username string) error
	UpdateUserGender(ctx context.Context, gender int) error
	UploadAvatar(ctx context.Context, fh *multipart.FileHeader) (string, error)
	GetAvatar(ctx context.Context, avatarId string, writer io.Writer) error
}

type UserHandler struct {
	userApp    UserHandlerApp
	middleware UserMiddleware
}

func NewUserHandler(userApp UserHandlerApp, middleware UserMiddleware) *UserHandler {
	return &UserHandler{
		userApp:    userApp,
		middleware: middleware,
	}
}

func (u *UserHandler) Init(router gin.IRouter) {
	userGroup := router.Group("user")
	userGroup.POST("login/wx", pgin.RequestResponseHandler(u.wechatLoginHandler))

	userNeedLogin := router.Group("user")
	userNeedLogin.Use(u.middleware.UserLoginRequired(), u.middleware.GrpcTokenRequired())
	userNeedLogin.GET("info", pgin.ResponseHandler(u.getUserInfoHandler))
	userNeedLogin.PUT("nickname/update", pgin.RequestWithErrorHandler(u.updateUserNameHandler))
	userNeedLogin.PUT("gender/update", pgin.RequestWithErrorHandler(u.updateUserGenderHander))
	userNeedLogin.POST("avatar/upload", pgin.ResponseHandler(u.uploadAvatarHandler))
	userNeedLogin.GET("avatar/:avatar_id", pgin.RequestWithErrorHandler(u.getAvatarHandler))
}

func (u *UserHandler) wechatLoginHandler(ctx *gin.Context, req *dto.WechatLoginRequest) (*dto.WechatLoginResponse, error) {
	return u.userApp.WechatLogin(ctx.Request.Context(), req.Code, req.DeviceId, req.AppName)
}

func (u *UserHandler) getUserInfoHandler(ctx *gin.Context) (*dto.User, error) {
	user, err := u.userApp.GetUserProfile(ctx.Request.Context())
	if err != nil {
		return nil, err
	}

	return dto.UserEntityToDto(user), nil
}

func (u *UserHandler) updateUserNameHandler(ctx *gin.Context, req *dto.UpdateUserNameRequest) error {
	err := u.userApp.UpdateUserName(ctx.Request.Context(), req.UserName)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserHandler) updateUserGenderHander(ctx *gin.Context, req *dto.UpdateUserGenderRequest) error {
	err := u.userApp.UpdateUserGender(ctx.Request.Context(), req.Gender)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserHandler) uploadAvatarHandler(ctx *gin.Context) (*dto.UploadAvatarResponse, error) {
	fh, err := ctx.FormFile("avatar")
	if err != nil {
		return nil, exception.ErrUploadAvatar
	}

	avatarUrl, err := u.userApp.UploadAvatar(ctx.Request.Context(), fh)
	if err != nil {
		return nil, err
	}

	return &dto.UploadAvatarResponse{AvatarUrl: avatarUrl}, nil
}

func (u *UserHandler) getAvatarHandler(ctx *gin.Context, req *dto.GetUserAvatarRequest) error {
	return u.userApp.GetAvatar(ctx.Request.Context(), req.AvatarId, ctx.Writer)
}
