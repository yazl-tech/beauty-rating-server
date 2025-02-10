// File:		analysis.go
// Created by:	Hoven
// Created on:	2025-02-09
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package analysis

import (
	"time"

	"github.com/yazl-tech/beauty-rating-server/pkg/analyst"
)

type AnalysisDetail struct {
	ID           int           `json:"id,omitempty"`
	UserID       int           `json:"userId,omitempty"`
	ImageUrl     string        `json:"imageUrl,omitempty"`
	Score        int           `json:"score,omitempty"`
	Percentile   int           `json:"percentile,omitempty"`
	Date         time.Time     `json:"date,omitempty"`
	Description  string        `json:"description,omitempty"`
	Tags         []string      `json:"tags,omitempty"`
	ScoreDetails []ScoreDetail `json:"scoreDetails,omitempty"`
	IsFavorite   bool          `json:"isFavorite,omitempty"`
}

type ScoreDetail struct {
	Label string `json:"label,omitempty"`
	Score int    `json:"score,omitempty"`
	Desc  string `json:"desc,omitempty"`
}

func parseAnalystDetails(d []analyst.Detail) []ScoreDetail {
	var result []ScoreDetail
	for _, detail := range d {
		result = append(result, ScoreDetail{
			Label: detail.Label,
			Score: detail.Score,
			Desc:  detail.Desc,
		})
	}
	return result
}
