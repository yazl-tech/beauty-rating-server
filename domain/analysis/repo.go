// File:		repo.go
// Created by:	Hoven
// Created on:	2025-02-09
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package analysis

import "context"

type Repo interface {
	CreateAnalysisDetail(ctx context.Context, detail *AnalysisDetail) error
	GetUserDetails(ctx context.Context, userId int) ([]*AnalysisDetail, error)
	GetUserDetail(ctx context.Context, userId, detailId int) (*AnalysisDetail, error)
	GetUserFavoriteDetails(ctx context.Context, userId int) ([]*AnalysisDetail, error)
	UpdateAnalysisDetail(ctx context.Context, detail *AnalysisDetail) error
	DeleteAnalysisDetail(ctx context.Context, userId, detailId int) error
}
