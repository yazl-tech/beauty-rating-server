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

	"gitea.hoven.com/core/auth-core/pkg/sdk/middleware"
	"github.com/go-puzzles/puzzles/pgin"
	"github.com/yazl-tech/beauty-rating-server/api/handler"
	"github.com/yazl-tech/beauty-rating-server/domain/user"
	"github.com/yazl-tech/beauty-rating-server/service"
	"google.golang.org/grpc"

	sdkHttpHandler "gitea.hoven.com/core/auth-core/pkg/sdk/handler"
)

type BeautyRatingApi struct {
	handler http.Handler
}

func SetupRouter(
	tokenKey string,
	wechatConf *user.WechatConfig,
	authCoreConn grpc.ClientConnInterface,
	beautyService *service.BeautyRatingService,
) *BeautyRatingApi {
	authCoreMiddleware := middleware.NewAuthCoreHttpMiddleware()

	authCoreHandler := sdkHttpHandler.NewAuthCoreSdkHttpHandler(
		authCoreConn,
		authCoreMiddleware,
		sdkHttpHandler.WithAuthBaseRoutes(),
		sdkHttpHandler.WithAccountRoutes(),
		sdkHttpHandler.WithWechatRoutes(wechatConf.AppId, wechatConf.SecretId),
		sdkHttpHandler.WithUserRoutes(),
	)

	router := pgin.NewServerHandlerWithOptions(
		pgin.WithMiddlewares(
			authCoreMiddleware.InjectTokenToGrpcContext(),
			authCoreMiddleware.UserLoginStatMiddleware(tokenKey),
		),
		pgin.WithRouters(
			"/v1",
			authCoreHandler,
			handler.NewAnalysisHandler(beautyService, authCoreMiddleware),
		),
	)

	return &BeautyRatingApi{
		handler: router,
	}
}

func (a *BeautyRatingApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler.ServeHTTP(w, r)
}
