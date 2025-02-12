// File:		common.go
// Created by:	Hoven
// Created on:	2025-02-10
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package handler

import "github.com/gin-gonic/gin"

type UserMiddleware interface {
	UserLoginRequired() gin.HandlerFunc
	GrpcTokenRequired() gin.HandlerFunc
	GetCurrentUserId(c *gin.Context) (int, error)
}
