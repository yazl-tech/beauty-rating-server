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
	GetImage(ctx context.Context, imageId string, rw http.ResponseWriter, req *http.Request)
	GetAnalysisDetails(ctx context.Context, userId int) (*dto.GetDetailsResponse, error)
	DoAnalysis(ctx context.Context, userId int, fh *multipart.FileHeader) (*dto.DoAnalysisResponse, error)
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

	needLoginGrp := router.Group("analysis", ah.middleware.UserLoginRequired())
	needLoginGrp.GET("", pgin.ResponseHandler(ah.getAnalysisDetails))
	needLoginGrp.POST("", pgin.ResponseHandler(ah.doAnalysisHandler))
	needLoginGrp.GET("favorite", pgin.ResponseHandler(ah.getFavoriteDetails))
	needLoginGrp.POST("favorite/:report_id", pgin.RequestWithErrorHandler(ah.doFavoriteHandler))
	needLoginGrp.POST("unfavorite/:report_id", pgin.RequestWithErrorHandler(ah.doUnFavoriteHandler))
	needLoginGrp.DELETE(":report_id", pgin.RequestWithErrorHandler(ah.deleteAnalysisHandler))
}

func (ah *AnalysisHandler) getFavoriteDetails(ctx *gin.Context) (*dto.GetDetailsResponse, error) {
	userId, err := ah.middleware.GetCurrentUserId(ctx)
	if err != nil {
		return nil, exception.ErrUnauthorized
	}

	resp, err := ah.analysisApp.GetFavoriteDetails(ctx.Request.Context(), userId)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (ah *AnalysisHandler) getAnalysisDetails(ctx *gin.Context) (*dto.GetDetailsResponse, error) {
	userId, err := ah.middleware.GetCurrentUserId(ctx)
	if err != nil {
		return nil, exception.ErrUnauthorized
	}

	resp, err := ah.analysisApp.GetAnalysisDetails(ctx.Request.Context(), userId)
	if err != nil {
		return nil, err
	}

	return resp, nil
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

	result, err := ah.analysisApp.DoAnalysis(ctx.Request.Context(), userId, fh)
	if err != nil {
		return nil, err
	}

	return result, nil
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
