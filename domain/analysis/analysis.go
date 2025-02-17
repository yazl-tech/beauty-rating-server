// File:		analysis.go
// Created by:	Hoven
// Created on:	2025-02-09
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package analysis

import (
	"fmt"
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
	AnalyisType  int           `json:"analyisType"`
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

type ShareDetailToken struct {
	DetailId int    `json:"detailId"`
	Expires  int64  `json:"expires"`
	Sig      string `json:"sig"`
}

func (st *ShareDetailToken) String() string {
	return fmt.Sprintf("detailId=%d&expires=%d&sig=%s", st.DetailId, st.Expires, st.Sig)
}
