// File:		analysis.go
// Created by:	Hoven
// Created on:	2025-02-10
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package handler

import (
	"context"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-puzzles/puzzles/pgin"
	"github.com/yazl-tech/beauty-rating-server/pkg/exception"
	"github.com/yazl-tech/beauty-rating-server/service/dto"
)

type AnalysisHandlerApp interface {
	DoAnalysis(ctx context.Context, userId int, fh *multipart.FileHeader) (*dto.DoAnalysisResponse, error)
	GetImage(ctx context.Context, imageId string, rw http.ResponseWriter, req *http.Request)
	GetAnalysisDetails(ctx context.Context, userId int) (*dto.GetDetailsResponse, error)
	ShareAnalysisDetail(ctx context.Context, userId, reportId int) (*dto.ShareDetailResponse, error)
	GetShareDetail(ctx context.Context, shareToken *dto.GetShareDetailRequest) (*dto.GetDetailResponse, error)
	DoFavorite(ctx context.Context, userId int, recordId int) error
	DoUnfavorite(ctx context.Context, userId int, recordId int) error
	GetFavoriteDetails(ctx context.Context, userId int) (*dto.GetDetailsResponse, error)
	DeleteAnalysis(ctx context.Context, userId int, recordId int) error
}

type AnalysisHandler struct {
	analysisApp AnalysisHandlerApp
	middleware  UserMiddleware
}

func NewAnalysisHandler(analysisApp AnalysisHandlerApp, middleware UserMiddleware) *AnalysisHandler {
	return &AnalysisHandler{
		analysisApp: analysisApp,
		middleware:  middleware,
	}
}

func (ah *AnalysisHandler) Init(router gin.IRouter) {
	analysisGrp := router.Group("analysis")
	analysisGrp.GET("image/:imageId", pgin.RequestHandler(ah.getImageHandler))
	analysisGrp.GET("share/detail", pgin.RequestResponseHandler(ah.getShareDetail))

	needLoginGrp := router.Group("analysis", ah.middleware.UserLoginRequired())
	needLoginGrp.POST("", pgin.ResponseHandler(ah.doAnalysisHandler))
	needLoginGrp.GET("", pgin.ResponseHandler(ah.getAnalysisDetails))
	needLoginGrp.POST("share/detail/:reportId", pgin.RequestResponseHandler(ah.shareAnalusysDetail))
	needLoginGrp.GET("favorite", pgin.ResponseHandler(ah.getFavoriteDetails))
	needLoginGrp.POST("favorite/:reportId", pgin.RequestWithErrorHandler(ah.doFavoriteHandler))
	needLoginGrp.POST("unfavorite/:reportId", pgin.RequestWithErrorHandler(ah.doUnFavoriteHandler))
	needLoginGrp.DELETE(":reportId", pgin.RequestWithErrorHandler(ah.deleteAnalysisHandler))
}

func (ah *AnalysisHandler) shareAnalusysDetail(ctx *gin.Context, req *dto.ShareDetailRequest) (*dto.ShareDetailResponse, error) {
	userId, err := ah.middleware.GetCurrentUserId(ctx)
	if err != nil {
		return nil, exception.ErrUnauthorized
	}

	return ah.analysisApp.ShareAnalysisDetail(ctx.Request.Context(), userId, req.ReportId)
}

func (ah *AnalysisHandler) getShareDetail(ctx *gin.Context, req *dto.GetShareDetailRequest) (*dto.GetDetailResponse, error) {
	return ah.analysisApp.GetShareDetail(ctx.Request.Context(), req)
}

func (ah *AnalysisHandler) getFavoriteDetails(ctx *gin.Context) (*dto.GetDetailsResponse, error) {
	userId, err := ah.middleware.GetCurrentUserId(ctx)
	if err != nil {
		return nil, exception.ErrUnauthorized
	}

	return ah.analysisApp.GetFavoriteDetails(ctx.Request.Context(), userId)
}

func (ah *AnalysisHandler) getAnalysisDetails(ctx *gin.Context) (*dto.GetDetailsResponse, error) {
	userId, err := ah.middleware.GetCurrentUserId(ctx)
	if err != nil {
		return nil, exception.ErrUnauthorized
	}

	return ah.analysisApp.GetAnalysisDetails(ctx.Request.Context(), userId)
}

func (ah *AnalysisHandler) getImageHandler(ctx *gin.Context, req *dto.GetImageRequest) {
	ah.analysisApp.GetImage(ctx.Request.Context(), req.ImageId, ctx.Writer, ctx.Request)
}

func (ah *AnalysisHandler) doAnalysisHandler(ctx *gin.Context) (*dto.DoAnalysisResponse, error) {
	userId, err := ah.middleware.GetCurrentUserId(ctx)
	if err != nil {
		return nil, exception.ErrUnauthorized
	}

	fh, err := ctx.FormFile("image")
	if err != nil {
		return nil, exception.ErrUploadAvatar
	}

	return ah.analysisApp.DoAnalysis(ctx.Request.Context(), userId, fh)
}

func (ah *AnalysisHandler) doFavoriteHandler(ctx *gin.Context, req *dto.DoFavoriteRequest) error {
	userId, err := ah.middleware.GetCurrentUserId(ctx)
	if err != nil {
		return exception.ErrUnauthorized
	}

	return ah.analysisApp.DoFavorite(ctx.Request.Context(), userId, req.ReportId)
}

func (ah *AnalysisHandler) doUnFavoriteHandler(ctx *gin.Context, req *dto.DoFavoriteRequest) error {
	userId, err := ah.middleware.GetCurrentUserId(ctx)
	if err != nil {
		return exception.ErrUnauthorized
	}

	return ah.analysisApp.DoUnfavorite(ctx.Request.Context(), userId, req.ReportId)
}

func (ah *AnalysisHandler) deleteAnalysisHandler(ctx *gin.Context, req *dto.DeleteAnalysisRequest) error {
	userId, err := ah.middleware.GetCurrentUserId(ctx)
	if err != nil {
		return exception.ErrUnauthorized
	}

	return ah.analysisApp.DeleteAnalysis(ctx.Request.Context(), userId, req.ReportId)
}
