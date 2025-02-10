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

	"github.com/yazl-tech/beauty-rating-server/domain/analysis"
	"github.com/yazl-tech/beauty-rating-server/pkg/dal/base"
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
	panic("not implemented") // TODO: Implement
}

func (ar *AnalysisRepo) GetUserDetails(ctx context.Context, userId int) ([]*analysis.AnalysisDetail, error) {
	panic("not implemented") // TODO: Implement
}

func (ar *AnalysisRepo) GetUserDetail(ctx context.Context, userId int, detailId int) (*analysis.AnalysisDetail, error) {
	panic("not implemented") // TODO: Implement
}

func (ar *AnalysisRepo) UpdateAnalysisDetail(ctx context.Context, detail *analysis.AnalysisDetail) error {
	panic("not implemented") // TODO: Implement
}

func (ar *AnalysisRepo) DeleteAnalysisDetail(ctx context.Context, userId int, detailId int) error {
	panic("not implemented") // TODO: Implement
}
