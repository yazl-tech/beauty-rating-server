// File:		analysis.go
// Created by:	Hoven
// Created on:	2025-02-09
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package model

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
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

func (a *Analysis) FromEntity(entity *analysis.AnalysisDetail) (err error) {
	if entity == nil {
		return nil
	}

	a.UserId = entity.UserID
	a.ImageUrl = entity.ImageUrl
	a.Score = entity.Score
	a.Description = entity.Description
	a.Tags, err = a.convertDBJson(entity.Tags)
	a.ScoreDetails, err = a.convertDBJson(entity.ScoreDetails)
	a.IsFavorite = entity.IsFavorite
	a.CreatedAt = entity.Date

	return nil
}

func (a *Analysis) convertDBJson(v any) (datatypes.JSON, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, errors.Wrap(err, "marshal")
	}

	return datatypes.JSON(b), nil
}

func (a *Analysis) ToEntity() (ad *analysis.AnalysisDetail, err error) {
	if a == nil {
		return nil, nil
	}

	ad = &analysis.AnalysisDetail{
		ID:           a.ID,
		UserID:       a.UserId,
		ImageUrl:     a.ImageUrl,
		Score:        a.Score,
		Description:  a.Description,
		IsFavorite:   a.IsFavorite,
		Date:         a.CreatedAt,
		Tags:         make([]string, 0),
		ScoreDetails: make([]analysis.ScoreDetail, 0),
	}

	err = a.parseDBJson(a.Tags, &ad.Tags)
	if err != nil {
		return nil, err
	}

	err = a.parseDBJson(a.ScoreDetails, &ad.ScoreDetails)
	if err != nil {
		return nil, err
	}

	return ad, nil
}

func (a *Analysis) parseDBJson(j datatypes.JSON, v any) error {
	return json.Unmarshal(j, v)
}
