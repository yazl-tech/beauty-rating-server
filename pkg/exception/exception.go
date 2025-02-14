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
	"net/http"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BeautyException struct {
	code  int
	cause error
}

func New(code int, cause string) *BeautyException {
	return &BeautyException{code: code, cause: errors.New(cause)}
}

func Wrap(code int, cause error) *BeautyException {
	return &BeautyException{code, errors.WithStack(cause)}
}

func (be *BeautyException) Error() string {
	return fmt.Sprintf("beautyException -> code: %d, cause: %v", be.code, be.cause)
}

func (be *BeautyException) Code() int {
	return be.code
}

func (be *BeautyException) Cause() error {
	return be.cause
}

func (be *BeautyException) Message() string {
	return be.cause.Error()
}

var (
	ErrUnauthorized          = New(http.StatusUnauthorized, "登录过期或未登录")
	ErrFileTooLarge          = New(http.StatusRequestEntityTooLarge, "文件大小超出预期")
	ErrDetailNotFound        = New(http.StatusNotFound, "分析报告不存在")
	ErrNotSpecifyDetail      = New(http.StatusBadRequest, "没有指定报告")
	ErrUploadAvatar          = New(http.StatusBadRequest, "上传头像失败")
	ErrGetAvatar             = New(http.StatusBadRequest, "获取头像失败")
	ErrUploadImage           = New(http.StatusBadRequest, "上传图片失败")
	ErrGetImage              = New(http.StatusBadRequest, "获取照片失败")
	ErrWechatLogin           = New(http.StatusBadRequest, "微信登录失败")
	ErrGetUserInfo           = New(http.StatusBadRequest, "获取用户信息失败")
	ErrUpdateUsername        = New(http.StatusBadRequest, "更新用户姓名失败")
	ErrUpdateGender          = New(http.StatusBadRequest, "更新性别失败")
	ErrDoAnalysis            = New(http.StatusBadRequest, "分析图片失败")
	ErrDoFavorite            = New(http.StatusBadRequest, "收藏失败")
	ErrDoUnFavorite          = New(http.StatusBadRequest, "取消收藏失败")
	ErrDeleteAnalysis        = New(http.StatusBadRequest, "删除分析报告失败")
	ErrGetAnalysisDetails    = New(http.StatusBadRequest, "获取分析报告列表失败")
	ErrGetFavoriteDetails    = New(http.StatusBadRequest, "获取收藏报告列表失败")
	ErrShareExpires          = New(http.StatusBadRequest, "分享已过期")
	ErrShareTokenInvalidates = New(http.StatusBadRequest, "分享链接异常")
	ErrShareAnalysisDetail   = New(http.StatusBadRequest, "分享报告失败")
	ErrGetShareDetail        = New(http.StatusBadRequest, "获取分享报告失败")
)

func CheckException(err error) bool {
	se := new(BeautyException)
	return errors.As(err, &se)
}

func ParseError(err error, defaultErr error) error {
	if CheckException(err) {
		return err
	}

	return defaultErr
}

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
