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
	"io"
	"mime/multipart"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/yazl-tech/beauty-rating-server/pkg/exception"
	"github.com/yazl-tech/beauty-rating-server/service/dto"
)

func (bs *BeautyRatingService) UploadImage(ctx context.Context, userId int, fh *multipart.FileHeader) (string, error) {
	imageId, err := bs.analysisSrv.UploadAnalysisImage(ctx, fh)
	if err != nil {
		plog.Errorc(ctx, "upload analysis image failed: %v", err)
		return "", exception.ErrUploadImage
	}

	return imageId, nil
}

func (bs *BeautyRatingService) GetImage(ctx context.Context, imageId string, writer io.Writer) error {
	err := bs.analysisSrv.GetAnalysisImage(ctx, imageId, writer)
	if err != nil {
		plog.Errorc(ctx, "get analysis image failed: %v", err)
		return exception.ErrGetImage
	}

	return nil
}

func (bs *BeautyRatingService) DoAnalysis(ctx context.Context, userId int, imageId string) (*dto.DoAnalysisResponse, error) {
	result, err := bs.analysisSrv.DoAnalysis(ctx, userId, imageId)
	if err != nil {
		plog.Errorc(ctx, "do analysis failed: %v", err)
		return nil, exception.ErrDoAnalysis
	}

	return &dto.DoAnalysisResponse{Detail: result}, nil
}

func (bs *BeautyRatingService) DoFavorite(ctx context.Context, userId int, recordId int) error {
	err := bs.analysisSrv.Favorite(ctx, userId, recordId)
	if err != nil {
		plog.Errorc(ctx, "do favorite failed: %v", err)
		return exception.ErrDoFavorite
	}

	return nil
}

func (bs *BeautyRatingService) DoUnfavorite(ctx context.Context, userId int, recordId int) error {
	err := bs.analysisSrv.UnFavorite(ctx, userId, recordId)
	if err != nil {
		plog.Errorc(ctx, "do unfavorite failed: %v", err)
		return exception.ErrDoUnFavorite
	}

	return nil
}

func (bs *BeautyRatingService) DeleteAnalysis(ctx context.Context, userId int, recordId int) error {
	err := bs.analysisSrv.DeleteAnalysis(ctx, userId, recordId)
	if err != nil {
		return exception.ErrDeleteAnalysis
	}

	return nil
}
