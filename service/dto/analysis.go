// File:		analysis.go
// Created by:	Hoven
// Created on:	2025-02-10
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package dto

import "github.com/yazl-tech/beauty-rating-server/domain/analysis"

type GetImageRequest struct {
	ImageId string `uri:"imageId" binding:"required"`
}

type DoAnalysisRequest struct {
	ImageId string `json:"image_id" binding:"required"`
}

type DoAnalysisResponse struct {
	Detail *analysis.AnalysisDetail `json:"detail"`
}

type DoFavoriteRequest struct {
	ReportId int `uri:"reportId" binding:"required"`
}

type DeleteAnalysisRequest struct {
	ReportId int `uri:"reportId" binding:"required"`
}

type GetDetailsResponse struct {
	Details []*analysis.AnalysisDetail `json:"details"`
}

type ShareDetailRequest struct {
	ReportId int `uri:"reportId" binding:"required"`
}

type ShareDetailResponse struct {
	UrlQuery string `json:"url_query"`
}

type GetShareDetailRequest struct {
	DetailId int    `form:"detailId" binding:"required"`
	Expires  int64  `form:"expires" binding:"required"`
	Sig      string `form:"sig" binding:"required"`
}

type GetDetailResponse struct {
	Detail *analysis.AnalysisDetail `json:"detail"`
}
