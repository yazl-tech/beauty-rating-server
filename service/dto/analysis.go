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
	ImageId string                   `json:"image_id"`
	Detail  *analysis.AnalysisDetail `json:"detail"`
}

type DoFavoriteRequest struct {
	ReportId int `uri:"report_id" binding:"required"`
}

type DeleteAnalysisRequest struct {
	ReportId int `uri:"report_id" binding:"required"`
}

type GetDetailsResponse struct {
	Details []*analysis.AnalysisDetail `json:"details"`
}
