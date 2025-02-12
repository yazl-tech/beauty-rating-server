// File:		user.go
// Created by:	Hoven
// Created on:	2025-02-10
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package service

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/yazl-tech/beauty-rating-server/domain/user"
	"github.com/yazl-tech/beauty-rating-server/pkg/exception"
	"github.com/yazl-tech/beauty-rating-server/service/dto"
)

func (bs *BeautyRatingService) WechatLogin(ctx context.Context, code, deviceId string) (*dto.WechatLoginResponse, error) {
	token, err := bs.userSrv.WxLogin(ctx, deviceId, code)
	if err != nil {
		plog.Errorc(ctx, "wechat login failed: %v", err)
		return nil, exception.ErrWechatLogin
	}

	return &dto.WechatLoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		User: &dto.User{
			UserId: token.UserID,
		},
	}, nil
}

func (bs *BeautyRatingService) GetUserProfile(ctx context.Context) (*user.User, error) {
	u, err := bs.userSrv.GetUserInfo(ctx)
	if err != nil {
		plog.Errorc(ctx, "get user profile failed: %v", err)
		return nil, exception.ErrGetUserInfo
	}

	return u, nil
}

func (bs *BeautyRatingService) UpdateUserName(ctx context.Context, username string) error {
	err := bs.userSrv.UpdateUsername(ctx, username)
	if err != nil {
		plog.Errorc(ctx, "update user name failed: %v", err)
		return exception.ErrUpdateUsername
	}

	return nil
}

func (bs *BeautyRatingService) UpdateUserGender(ctx context.Context, gender int) error {
	err := bs.userSrv.UpdateGender(ctx, gender)
	if err != nil {
		plog.Errorc(ctx, "update user gender failed: %v", err)
		return exception.ErrUpdateGender
	}

	return nil
}

func (bs *BeautyRatingService) UploadAvatar(ctx context.Context, fh *multipart.FileHeader) (string, error) {
	avatarUrl, err := bs.userSrv.UploadAvatar(ctx, fh)
	if err != nil {
		plog.Errorc(ctx, "upload avatar failed: %v", err)
		return "", exception.ErrUploadAvatar
	}

	return avatarUrl, nil
}

func (bs *BeautyRatingService) GetAvatar(ctx context.Context, avatarId string, writer io.Writer) error {
	err := bs.userSrv.GetAvatar(ctx, avatarId, writer)
	if err != nil {
		plog.Errorc(ctx, "get avatar failed: %v", err)
		return exception.ErrGetAvatar
	}

	return nil
}
