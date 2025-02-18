// File:		service.go
// Created by:	Hoven
// Created on:	2025-02-09
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package service

import (
	doubaopb "github.com/yazl-tech/ai-bot/pkg/proto/doubao"
	"github.com/yazl-tech/beauty-rating-server/config"
	"github.com/yazl-tech/beauty-rating-server/domain/analysis"
	"github.com/yazl-tech/beauty-rating-server/domain/user"
	"github.com/yazl-tech/beauty-rating-server/pkg/analyst"
	"github.com/yazl-tech/beauty-rating-server/pkg/analyst/ai"
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
	authCoreConn grpc.ClientConnInterface,
	aiBotConn grpc.ClientConnInterface,
	beautyConf *config.BeautyConfig,
	wechatConfig *user.WechatConfig,
) *BeautyRatingService {
	mockAnalyst := mock.NewMockAnalyst()
	doubaoClient := doubaopb.NewDoubaoHandlerClient(aiBotConn)
	aiAnalyst := ai.NewAiAnalyst(beautyConf.AiModel, doubaoClient)

	analystSelector := analyst.NewAnalystSelector(
		analyst.WithAnalysts(mockAnalyst, beautyConf.AnalystWeight(mockAnalyst.Typ())),
		analyst.WithAnalysts(aiAnalyst, beautyConf.AnalystWeight(aiAnalyst.Typ())),
	)

	analysisRepo := analysisRepo.NewAnalysisRepo(db)
	analysisSrv := analysis.NewAnalysisService(
		beautyConf,
		analystSelector,
		analysisRepo,
		oss,
	)

	userSrv := user.NewUserService(wechatConfig, authCoreConn)

	return &BeautyRatingService{
		analysisSrv: analysisSrv,
		userSrv:     userSrv,
	}
}
