// File:		analysis.go
// Created by:	Hoven
// Created on:	2025-02-09
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package analysisRepo

import (
	"context"
	"errors"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/go-puzzles/puzzles/putils"
	"github.com/yazl-tech/beauty-rating-server/domain/analysis"
	"github.com/yazl-tech/beauty-rating-server/pkg/dal/base"
	"github.com/yazl-tech/beauty-rating-server/pkg/dal/model"
	"github.com/yazl-tech/beauty-rating-server/pkg/exception"
	"gorm.io/gorm"
)

var _ analysis.Repo = (*AnalysisRepo)(nil)

type AnalysisRepo struct {
	db *base.Query
}

func NewAnalysisRepo(db *gorm.DB) *AnalysisRepo {
	return &AnalysisRepo{db: base.Use(db)}
}

func (ar *AnalysisRepo) CreateAnalysisDetail(ctx context.Context, detail *analysis.AnalysisDetail) error {
	detailDal := new(model.Analysis)
	err := detailDal.FromEntity(detail)
	if err != nil {
		return err
	}

	err = ar.db.Analysis.WithContext(ctx).Create(detailDal)
	if err != nil {
		return err
	}

	detail.ID = detailDal.ID
	return nil
}

func (ar *AnalysisRepo) GetUserDetails(ctx context.Context, userId int) ([]*analysis.AnalysisDetail, error) {
	db := ar.db.Analysis

	details, err := db.WithContext(ctx).Where(db.UserId.Eq(userId)).Find()
	if err != nil {
		return nil, err
	}

	detailEntyties := putils.Convert(details, func(detail *model.Analysis) *analysis.AnalysisDetail {
		de, err := detail.ToEntity()
		if err != nil {
			plog.Errorc(ctx, "convert detail: %v to entity error: %v", detail.ID, err)
			return nil
		}

		return de
	})

	return detailEntyties, nil
}

func (ar *AnalysisRepo) GetUserDetail(ctx context.Context, userId int, detailId int) (*analysis.AnalysisDetail, error) {
	db := ar.db.Analysis

	detail, err := db.WithContext(ctx).Where(db.ID.Eq(detailId), db.UserId.Eq(userId)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrDetailNotFound
	} else if err != nil {
		return nil, err
	}

	return detail.ToEntity()
}

func (ar *AnalysisRepo) GetDetail(ctx context.Context, detailId int) (*analysis.AnalysisDetail, error) {
	db := ar.db.Analysis

	detail, err := db.WithContext(ctx).Where(db.ID.Eq(detailId)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrDetailNotFound
	} else if err != nil {
		return nil, err
	}

	return detail.ToEntity()
}

func (ar *AnalysisRepo) UpdateAnalysisDetail(ctx context.Context, detail *analysis.AnalysisDetail) error {
	if detail.ID == 0 {
		return exception.ErrNotSpecifyDetail
	}

	detailDal := new(model.Analysis)
	err := detailDal.FromEntity(detail)
	if err != nil {
		return err
	}

	db := ar.db.Analysis
	err = db.WithContext(ctx).Where(db.ID.Eq(detail.ID)).Save(detailDal)
	if err != nil {
		return err
	}

	return nil
}

func (ar *AnalysisRepo) GetUserFavoriteDetails(ctx context.Context, userId int) ([]*analysis.AnalysisDetail, error) {
	db := ar.db.Analysis

	details, err := db.WithContext(ctx).Where(db.UserId.Eq(userId), db.IsFavorite.Is(true)).Find()
	if err != nil {
		return nil, err
	}

	detailEntyties := putils.Convert(details, func(detail *model.Analysis) *analysis.AnalysisDetail {
		de, err := detail.ToEntity()
		if err != nil {
			plog.Errorc(ctx, "convert detail: %v to entity error: %v", detail.ID, err)
			return nil
		}

		return de
	})

	return detailEntyties, nil
}

func (ar *AnalysisRepo) DeleteAnalysisDetail(ctx context.Context, userId int, detailId int) error {
	db := ar.db.Analysis

	info, err := db.WithContext(ctx).Where(db.ID.Eq(detailId), db.UserId.Eq(userId)).Delete()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return exception.ErrDetailNotFound
	} else if err != nil {
		return err
	}

	if info.RowsAffected == 0 {
		return exception.ErrDetailNotFound
	}

	return nil
}

func (ar *AnalysisRepo) CheckDetailExists(ctx context.Context, userId, detailId int) bool {
	db := ar.db.Analysis

	count, err := db.WithContext(ctx).Where(db.ID.Eq(detailId), db.UserId.Eq(userId)).Count()
	if err != nil {
		return false
	}

	return count > 0
}
