// File:		service.go
// Created by:	Hoven
// Created on:	2025-02-09
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package analysis

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/go-puzzles/puzzles/putils"
	"github.com/pkg/errors"
	"github.com/yazl-tech/beauty-rating-server/config"
	"github.com/yazl-tech/beauty-rating-server/pkg/analyst"
	"github.com/yazl-tech/beauty-rating-server/pkg/exception"
	"github.com/yazl-tech/beauty-rating-server/pkg/oss"
	"gorm.io/gorm"
)

type Service interface {
	UploadAnalysisImage(ctx context.Context, avatarFile *multipart.FileHeader) (string, []byte, error)
	GetAnalysisImage(ctx context.Context, imageId string, rw http.ResponseWriter, req *http.Request)
	DoAnalysis(ctx context.Context, userId int, imageId string, b []byte) (*AnalysisDetail, error)
	GetFavoriteDetails(ctx context.Context, userId int) ([]*AnalysisDetail, error)
	GetAnalysisDetials(ctx context.Context, userId int) ([]*AnalysisDetail, error)
	ShareAnalysisDetail(ctx context.Context, userId, reportId int) (*ShareDetailToken, error)
	GetShareDetail(ctx context.Context, token *ShareDetailToken) (*AnalysisDetail, error)
	Favorite(ctx context.Context, userId int, detailId int) error
	UnFavorite(ctx context.Context, userId int, detailId int) error
	DeleteAnalysis(ctx context.Context, userId int, detailId int) error
}

var _ Service = (*DefaultAnalysisService)(nil)

type DefaultAnalysisService struct {
	beautyConf     *config.BeautyConfig
	analyst        analyst.Analyst
	repo           Repo
	oss            oss.IOSS
	analysisImgDir string
}

func NewAnalysisService(
	beautyConf *config.BeautyConfig,
	analyst analyst.Analyst,
	repo Repo,
	oss oss.IOSS,
) *DefaultAnalysisService {
	return &DefaultAnalysisService{
		beautyConf:     beautyConf,
		analyst:        analyst,
		repo:           repo,
		oss:            oss,
		analysisImgDir: "analysis",
	}
}

func (as *DefaultAnalysisService) verifyShareToken(token *ShareDetailToken) (err error) {
	if time.Now().Unix() > token.Expires {
		return exception.ErrShareExpires
	}

	data := fmt.Sprintf("%d/%d", token.DetailId, token.Expires)
	h := hmac.New(sha256.New, []byte(as.beautyConf.ShareSecretKey))
	h.Write([]byte(data))
	expectedSig := hex.EncodeToString(h.Sum(nil))

	if !hmac.Equal([]byte(expectedSig), []byte(token.Sig)) {
		return exception.ErrShareTokenInvalidates
	}

	return nil
}

func (as *DefaultAnalysisService) generateShareToken(detailId int, expiresDura time.Duration) *ShareDetailToken {
	expires := time.Now().Add(expiresDura).Unix()
	data := fmt.Sprintf("%d/%d", detailId, expires)

	h := hmac.New(sha256.New, []byte(as.beautyConf.ShareSecretKey))
	h.Write([]byte(data))
	signature := hex.EncodeToString(h.Sum(nil))

	return &ShareDetailToken{
		DetailId: detailId,
		Expires:  int64(expires),
		Sig:      signature,
	}
}

func (as *DefaultAnalysisService) ShareAnalysisDetail(ctx context.Context, userId, reportId int) (*ShareDetailToken, error) {
	exists := as.repo.CheckDetailExists(ctx, userId, reportId)
	if !exists {
		return nil, exception.ErrDetailNotFound
	}

	expires := time.Hour * 24
	shareToken := as.generateShareToken(reportId, expires)

	return shareToken, nil
}

func (as *DefaultAnalysisService) GetShareDetail(ctx context.Context, token *ShareDetailToken) (*AnalysisDetail, error) {
	err := as.verifyShareToken(token)
	if err != nil {
		return nil, errors.Wrap(err, "parse share token failed")
	}

	detail, err := as.repo.GetDetail(ctx, token.DetailId)
	if err != nil {
		return nil, err
	}

	return as.convertImage(ctx, detail), nil
}

func (as *DefaultAnalysisService) GetFavoriteDetails(ctx context.Context, userId int) ([]*AnalysisDetail, error) {
	resp, err := as.repo.GetUserFavoriteDetails(ctx, userId)
	if err != nil {
		return nil, err
	}

	return as.convertImages(ctx, resp), nil
}

func (as *DefaultAnalysisService) convertImage(ctx context.Context, detail *AnalysisDetail) *AnalysisDetail {
	objName := fmt.Sprintf("%s/%s", as.analysisImgDir, detail.ImageUrl)
	presignedUrl, err := as.presignImageUrl(ctx, objName, 5*time.Minute)
	if err != nil {
		plog.Warnc(ctx, "presignedUrl: %v failed: %v", detail.ImageUrl, err)
		return nil
	}

	presignedUrl.Host = as.beautyConf.ApiHost
	// /api/v1/analysis/image/:imageId
	presignedUrl.Path = fmt.Sprintf("%s%s/analysis/image/%s", as.beautyConf.ApiPrefix, as.beautyConf.ApiVersion, detail.ImageUrl)

	detail.ImageUrl = presignedUrl.String()
	return detail
}

func (as *DefaultAnalysisService) convertImages(ctx context.Context, details []*AnalysisDetail) []*AnalysisDetail {
	return putils.Convert(details, func(d *AnalysisDetail) *AnalysisDetail {
		return as.convertImage(ctx, d)
	})
}

func (as *DefaultAnalysisService) presignImageUrl(ctx context.Context, objName string, expires time.Duration) (*url.URL, error) {
	u, err := as.oss.PresignedGetObject(ctx, objName, expires)
	if err != nil {
		return nil, errors.Wrap(err, "presignedImage")
	}

	return u, nil
}

func (as *DefaultAnalysisService) GetAnalysisDetials(ctx context.Context, userId int) ([]*AnalysisDetail, error) {
	resp, err := as.repo.GetUserDetails(ctx, userId)
	if err != nil {
		return nil, err
	}

	return as.convertImages(ctx, resp), nil
}

func (as *DefaultAnalysisService) UploadAnalysisImage(ctx context.Context, file *multipart.FileHeader) (string, []byte, error) {
	src, err := file.Open()
	if err != nil {
		return "", nil, err
	}
	defer src.Close()

	var fileBytes bytes.Buffer
	teeReader := io.TeeReader(src, &fileBytes)

	imageUrl, err := as.oss.UploadFile(ctx, file.Size, as.analysisImgDir, file.Filename, teeReader)
	if err != nil {
		return "", nil, errors.Wrap(err, "uploadAnalysisImage")
	}

	return imageUrl, fileBytes.Bytes(), nil
}

func (as *DefaultAnalysisService) GetAnalysisImage(ctx context.Context, imageId string, rw http.ResponseWriter, req *http.Request) {
	objName := fmt.Sprintf("%s/%s", as.analysisImgDir, imageId)
	as.oss.ProxyPresignedGetObject(objName, rw, req)
}

func (as *DefaultAnalysisService) DoAnalysis(ctx context.Context, userId int, imageId string, b []byte) (*AnalysisDetail, error) {
	d, err := as.analyst.DoAnalysis(ctx, imageId, imageId, b)
	if err != nil {
		return nil, err
	}

	detail := &AnalysisDetail{
		UserID:       userId,
		ImageUrl:     imageId,
		Score:        d.Score,
		Percentile:   rand.Intn(20) + 80,
		Description:  d.Description,
		Tags:         d.Tags,
		Date:         time.Now(),
		ScoreDetails: parseAnalystDetails(d.Details),
		AnalyisType:  int(d.AnalystType),
	}

	err = as.repo.CreateAnalysisDetail(ctx, detail)
	if err != nil {
		return nil, err
	}

	return as.convertImage(ctx, detail), nil
}

func (as *DefaultAnalysisService) Favorite(ctx context.Context, userId int, detailId int) error {
	detail, err := as.repo.GetUserDetail(ctx, userId, detailId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return exception.ErrDetailNotFound
	} else if err != nil {
		return errors.Wrapf(err, "getUserDetail. userId=%v, detailId=%v", userId, detail)
	}

	detail.IsFavorite = true
	return as.repo.UpdateAnalysisDetail(ctx, detail)
}

func (as *DefaultAnalysisService) UnFavorite(ctx context.Context, userId int, detailId int) error {
	detail, err := as.repo.GetUserDetail(ctx, userId, detailId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return exception.ErrDetailNotFound
	} else if err != nil {
		return errors.Wrapf(err, "getUserDetail. userId=%v, detailId=%v", userId, detail)
	}

	detail.IsFavorite = false
	return as.repo.UpdateAnalysisDetail(ctx, detail)
}

func (as *DefaultAnalysisService) DeleteAnalysis(ctx context.Context, userId int, detailId int) error {
	return as.repo.DeleteAnalysisDetail(ctx, userId, detailId)
}
