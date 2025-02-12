// File:		service.go
// Created by:	Hoven
// Created on:	2025-02-09
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package service

import (
	"github.com/yazl-tech/beauty-rating-server/domain/analysis"
	"github.com/yazl-tech/beauty-rating-server/domain/user"
	"github.com/yazl-tech/beauty-rating-server/pkg/analyst/mock"
	"github.com/yazl-tech/beauty-rating-server/pkg/oss"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	analysisRepo "github.com/yazl-tech/beauty-rating-server/pkg/dal/analysis"
)

type BeautyRatingService struct {
	analysisSrv analysis.Service
	userSrv     user.Service
}

func NewBeautyRatingService(
	db *gorm.DB,
	oss oss.IOSS,
	wechatConfig *user.WechatConfig,
	authCoreConn grpc.ClientConnInterface,
) *BeautyRatingService {
	mockAnlyst := mock.NewMockAnalyst()

	analysisRepo := analysisRepo.NewAnalysisRepo(db)
	analysisSrv := analysis.NewAnalysisService(mockAnlyst, analysisRepo, oss)

	userSrv := user.NewUserService(wechatConfig, authCoreConn)

	return &BeautyRatingService{
		analysisSrv: analysisSrv,
		userSrv:     userSrv,
	}
}
