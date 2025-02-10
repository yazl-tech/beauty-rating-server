// File:		analysis.go
// Created by:	Hoven
// Created on:	2025-02-09
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package model

import (
	"time"

	"github.com/yazl-tech/beauty-rating-server/domain/analysis"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Analysis struct {
	ID           int    `gorm:"primaryKey;autoIncrement"`
	UserId       int    `gorm:"not null"`
	ImageUrl     string `gorm:"not null;type:varchar(256)"`
	Score        int    `gorm:"not null"`
	Description  string `gorm:"type:text"`
	Tags         datatypes.JSON
	ScoreDetails datatypes.JSON
	IsFavorite   bool

	CreatedAt time.Time      `gorm:"comment:创建时间"`
	UpdatedAt time.Time      `gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:软删除时间"`
}

func (a *Analysis) TableName() string {
	return "analysises"
}

func (a *Analysis) FromEntity(entity *analysis.AnalysisDetail) *Analysis {
	if entity == nil {
		return a
	}

	return a
}

func (a *Analysis) ToEntity() *analysis.AnalysisDetail {
	return &analysis.AnalysisDetail{}
}
