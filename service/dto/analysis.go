// File:		analysis.go
// Created by:	Hoven
// Created on:	2025-02-10
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package dto

import "github.com/yazl-tech/beauty-rating-server/domain/analysis"

type UploadImageResponse struct {
	ImageId string `json:"image_id"`
}

type GetImageRequest struct {
	ImageId string `uri:"image_id"`
}

type DoAnalysisRequest struct {
	ImageId string `json:"image_id"`
}

type DoAnalysisResponse struct {
	Detail *analysis.AnalysisDetail
}

type DoFavoriteRequest struct {
	ReportId int `uri:"report_id"`
}

type DeleteAnalysisRequest struct {
	ReportId int `uri:"report_id"`
}
