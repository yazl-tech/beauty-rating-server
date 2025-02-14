// File:		analysis.go
// Created by:	Hoven
// Created on:	2025-02-10
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package service

import (
	"context"
	"mime/multipart"
	"net/http"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/yazl-tech/beauty-rating-server/domain/analysis"
	"github.com/yazl-tech/beauty-rating-server/pkg/exception"
	"github.com/yazl-tech/beauty-rating-server/service/dto"
)

func (bs *BeautyRatingService) GetShareDetail(ctx context.Context, shareToken *dto.GetShareDetailRequest) (*dto.GetDetailResponse, error) {
	detail, err := bs.analysisSrv.GetShareDetail(ctx, &analysis.ShareDetailToken{
		Sig:      shareToken.Sig,
		Expires:  shareToken.Expires,
		DetailId: shareToken.DetailId,
	})
	if err != nil {
		plog.Errorc(ctx, "get share detail failed: %v", err)
		return nil, exception.ParseError(err, exception.ErrGetShareDetail)
	}

	return &dto.GetDetailResponse{
		Detail: detail,
	}, nil
}

func (bs *BeautyRatingService) ShareAnalysisDetail(ctx context.Context, userId, reportId int) (*dto.ShareDetailResponse, error) {
	shareToken, err := bs.analysisSrv.ShareAnalysisDetail(ctx, userId, reportId)
	if err != nil {
		plog.Errorc(ctx, "share analysis detail failed: %v", err)
		return nil, exception.ParseError(err, exception.ErrShareAnalysisDetail)
	}

	return &dto.ShareDetailResponse{
		UrlQuery: shareToken.String(),
	}, nil
}

func (bs *BeautyRatingService) GetImage(ctx context.Context, imageId string, rw http.ResponseWriter, req *http.Request) {
	bs.analysisSrv.GetAnalysisImage(ctx, imageId, rw, req)
}

func (bs *BeautyRatingService) GetFavoriteDetails(ctx context.Context, userId int) (*dto.GetDetailsResponse, error) {
	resp, err := bs.analysisSrv.GetFavoriteDetails(ctx, userId)
	if err != nil {
		plog.Errorc(ctx, "get favorite details failed: %v", err)
		return nil, exception.ParseError(err, exception.ErrGetFavoriteDetails)
	}

	return &dto.GetDetailsResponse{
		Details: resp,
	}, nil
}

func (bs *BeautyRatingService) GetAnalysisDetails(ctx context.Context, userId int) (*dto.GetDetailsResponse, error) {
	resp, err := bs.analysisSrv.GetAnalysisDetials(ctx, userId)
	if err != nil {
		plog.Errorc(ctx, "get analysis details failed: %v", err)
		return nil, exception.ParseError(err, exception.ErrGetAnalysisDetails)
	}

	return &dto.GetDetailsResponse{
		Details: resp,
	}, nil
}

func (bs *BeautyRatingService) DoAnalysis(ctx context.Context, userId int, fh *multipart.FileHeader) (*dto.DoAnalysisResponse, error) {
	imageId, b, err := bs.analysisSrv.UploadAnalysisImage(ctx, fh)
	if err != nil {
		plog.Errorc(ctx, "upload analysis image failed: %v", err)
		return nil, exception.ParseError(err, exception.ErrUploadImage)
	}

	result, err := bs.analysisSrv.DoAnalysis(ctx, userId, imageId, b)
	if err != nil {
		plog.Errorc(ctx, "do analysis failed: %v", err)
		return nil, exception.ParseError(err, exception.ErrDoAnalysis)
	}

	return &dto.DoAnalysisResponse{Detail: result}, nil
}

func (bs *BeautyRatingService) DoFavorite(ctx context.Context, userId int, recordId int) error {
	err := bs.analysisSrv.Favorite(ctx, userId, recordId)
	if err != nil {
		plog.Errorc(ctx, "do favorite failed: %v", err)
		return exception.ParseError(err, exception.ErrDoFavorite)
	}

	return nil
}

func (bs *BeautyRatingService) DoUnfavorite(ctx context.Context, userId int, recordId int) error {
	err := bs.analysisSrv.UnFavorite(ctx, userId, recordId)
	if err != nil {
		plog.Errorc(ctx, "do unfavorite failed: %v", err)
		return exception.ParseError(err, exception.ErrDoUnFavorite)
	}

	return nil
}

func (bs *BeautyRatingService) DeleteAnalysis(ctx context.Context, userId int, recordId int) error {
	err := bs.analysisSrv.DeleteAnalysis(ctx, userId, recordId)
	if err != nil {
		plog.Errorc(ctx, "delete analysis failed: %v", err)
		return exception.ParseError(err, exception.ErrDeleteAnalysis)
	}

	return nil
}
