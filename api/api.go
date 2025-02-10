// File:		api.go
// Created by:	Hoven
// Created on:	2025-02-09
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package api

import (
	"net/http"
	"time"

	"gitea.hoven.com/core/auth-core/pkg/sdk/middleware"
	"github.com/go-puzzles/puzzles/pgin"
	"github.com/yazl-tech/beauty-rating-server/api/handler"
	"github.com/yazl-tech/beauty-rating-server/service"
)

type BeautyRatingApi struct {
	handler http.Handler
}

func SetupRouter(
	beautyService *service.BeautyRatingService,
	storage middleware.Storage,
) *BeautyRatingApi {

	tokenManager := middleware.NewManager(
		storage,
		middleware.WithCacheTTL(time.Hour*2),
		middleware.WithCachePrefix("beauty-rating"),
	)
	middleware := middleware.NewAuthCoreHttpMiddleware[*middleware.UserToken](tokenManager)

	router := pgin.NewServerHandlerWithOptions(
		pgin.WithMiddlewares(
			middleware.InjectTokenToGrpcContext(),
			middleware.UserLoginStatMiddleware(),
		),
		pgin.WithRouters(
			"/v1",
			handler.NewUserHandler(beautyService, middleware),
		),
	)

	return &BeautyRatingApi{
		handler: router,
	}
}

func (a *BeautyRatingApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler.ServeHTTP(w, r)
}
