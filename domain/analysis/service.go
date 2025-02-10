// File:		service.go
// Created by:	Hoven
// Created on:	2025-02-09
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package analysis

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/pkg/errors"
	"github.com/yazl-tech/beauty-rating-server/pkg/analyst"
	"github.com/yazl-tech/beauty-rating-server/pkg/exception"
	"github.com/yazl-tech/beauty-rating-server/pkg/oss"
	"gorm.io/gorm"
)

type Service interface {
	UploadAnalysisImage(ctx context.Context, avatarFile *multipart.FileHeader) (string, error)
	GetAnalysisImage(ctx context.Context, imagePath string, writer io.Writer) error
	DoAnalysis(ctx context.Context, userId int, imageUrl string) (*AnalysisDetail, error)
	Favorite(ctx context.Context, userId int, detailId int) error
	UnFavorite(ctx context.Context, userId int, detailId int) error
	DeleteAnalysis(ctx context.Context, userId int, detailId int) error
}

type DefaultAnalysisService struct {
	analyst analyst.Analyst
	repo    Repo
	oss     oss.IOSS
}

func NewAnalysisService(
	analyst analyst.Analyst,
	repo Repo,
	oss oss.IOSS,
) *DefaultAnalysisService {
	return &DefaultAnalysisService{analyst: analyst, repo: repo, oss: oss}
}

func (as *DefaultAnalysisService) UploadAnalysisImage(ctx context.Context, file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	imageUrl, err := as.oss.UploadFile(ctx, file.Size, "analysisImages", file.Filename, src)
	if err != nil {
		return "", errors.Wrap(err, "uploadAnalysisImage")
	}

	return imageUrl, nil
}

func (as *DefaultAnalysisService) GetAnalysisImage(ctx context.Context, imagePath string, writer io.Writer) error {
	return as.oss.GetFile(ctx, imagePath, writer)
}

func (as *DefaultAnalysisService) DoAnalysis(ctx context.Context, userId int, imageUrl string) (*AnalysisDetail, error) {
	d, err := as.analyst.DoAnalysis(ctx, []byte{})
	if err != nil {
		return nil, err
	}

	detail := &AnalysisDetail{
		UserID:       userId,
		ImageUrl:     imageUrl,
		Score:        d.Score,
		Description:  d.Description,
		Tags:         d.Tags,
		ScoreDetails: parseAnalystDetails(d.Details),
	}

	err = as.repo.CreateAnalysisDetail(ctx, detail)
	if err != nil {
		return nil, err
	}

	return detail, nil
}

func (as *DefaultAnalysisService) Favorite(ctx context.Context, userId int, detailId int) error {
	detail, err := as.repo.GetUserDetail(ctx, userId, detailId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return exception.ErrDetailNotFound
	} else if err != nil {
		return errors.Wrapf(err, "getUserDetail. userId=%v, detailId=%v", userId, detail)
	}

	detail.IsFavorite = true
	return as.repo.UpdateAnalysisDetail(ctx, detail)
}

func (as *DefaultAnalysisService) UnFavorite(ctx context.Context, userId int, detailId int) error {
	detail, err := as.repo.GetUserDetail(ctx, userId, detailId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return exception.ErrDetailNotFound
	} else if err != nil {
		return errors.Wrapf(err, "getUserDetail. userId=%v, detailId=%v", userId, detail)
	}

	detail.IsFavorite = false
	return as.repo.UpdateAnalysisDetail(ctx, detail)
}

func (as *DefaultAnalysisService) DeleteAnalysis(ctx context.Context, userId int, detailId int) error {
	return as.repo.DeleteAnalysisDetail(ctx, userId, detailId)
}
