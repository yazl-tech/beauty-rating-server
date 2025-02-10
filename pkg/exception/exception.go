// File:		exception.go
// Created by:	Hoven
// Created on:	2025-02-09
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package exception

import (
	"fmt"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrDetailNotFound = errors.New("分析报告不存在")
	ErrUnauthorized   = errors.New("登录过期或未登录")
	ErrFileTooLarge   = errors.New("文件大小超出预期")

	ErrUploadAvatar   = errors.New("上传头像失败")
	ErrGetAvatar      = errors.New("获取头像失败")
	ErrUploadImage    = errors.New("上传图片失败")
	ErrGetImage       = errors.New("获取照片失败")
	ErrGetUserInfo    = errors.New("获取用户信息失败")
	ErrUpdateUsername = errors.New("更新用户姓名失败")
	ErrUpdateGender   = errors.New("更新性别失败")
	ErrDoAnalysis     = errors.New("分析图片失败")
	ErrDoFavorite     = errors.New("收藏失败")
	ErrDoUnFavorite   = errors.New("取消收藏失败")
	ErrDeleteAnalysis = errors.New("删除分析报告失败")
)

func ParseGrpcError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	if st.Code() == codes.Unauthenticated {
		return ErrUnauthorized
	}

	return fmt.Errorf("%v", st.Message())
}
