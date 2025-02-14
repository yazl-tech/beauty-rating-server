// File:		main.go
// Created by:	Hoven
// Created on:	2025-02-09
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package main

import (
	"github.com/go-puzzles/puzzles/cores"
	"github.com/go-puzzles/puzzles/dialer/grpc"
	"github.com/go-puzzles/puzzles/pflags"
	"github.com/go-puzzles/puzzles/pgorm"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/yazl-tech/beauty-rating-server/api"
	"github.com/yazl-tech/beauty-rating-server/config"
	"github.com/yazl-tech/beauty-rating-server/domain/user"
	"github.com/yazl-tech/beauty-rating-server/pkg/dal/model"
	"github.com/yazl-tech/beauty-rating-server/pkg/oss/minio"
	"github.com/yazl-tech/beauty-rating-server/service"

	consulpuzzle "github.com/go-puzzles/puzzles/cores/puzzles/consul-puzzle"
	httppuzzle "github.com/go-puzzles/puzzles/cores/puzzles/http-puzzle"
)

var (
	beautyConfFlag    = pflags.Struct("beautyConf", (*config.BeautyConfig)(nil), "beauty configuration")
	mysqlConfFlag     = pflags.Struct("mysqlAuth", (*pgorm.MysqlConfig)(nil), "mysql auth config")
	minioConfFlag     = pflags.Struct("minioAuth", (*minio.MinioConfig)(nil), "minio auth config")
	wechatSdkConfFlag = pflags.Struct("wechat", (*user.WechatConfig)(nil), "wechat sdk config")
)

func main() {
	pflags.Parse()

	beautyConf := new(config.BeautyConfig)
	plog.PanicError(beautyConfFlag(beautyConf))
	minioConf := new(minio.MinioConfig)
	plog.PanicError(minioConfFlag(minioConf))
	mysqlConf := new(pgorm.MysqlConfig)
	plog.PanicError(mysqlConfFlag(mysqlConf))
	wechatConf := new(user.WechatConfig)
	plog.PanicError(wechatSdkConfFlag(wechatConf))

	plog.Debugf("beautyConf: %v", plog.Jsonify(beautyConf))

	authCoreConn, err := grpc.DialGrpc(beautyConf.AuthCoreSrv)
	plog.PanicError(err)

	minioClient := minio.NewMinioOss(minioConf)

	plog.PanicError(pgorm.RegisterSqlModelWithConf(mysqlConf, model.AllTables()...))
	plog.PanicError(pgorm.AutoMigrate(mysqlConf))
	db := pgorm.GetDbByConf(mysqlConf)

	beautyService := service.NewBeautyRatingService(db, minioClient, authCoreConn, beautyConf, wechatConf)
	router := api.SetupRouter(beautyConf, wechatConf, authCoreConn, beautyService)

	coreSrv := cores.NewPuzzleCore(
		cores.WithService(pflags.GetServiceName()),
		consulpuzzle.WithConsulRegister(),
		httppuzzle.WithCoreHttpCORS(),
		httppuzzle.WithCoreHttpPuzzle(beautyConf.ApiPrefix, router),
	)
	plog.PanicError(cores.Start(coreSrv, beautyConf.ApiPort))
}
