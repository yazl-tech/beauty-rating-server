// File:		generate.go
// Created by:	Hoven
// Created on:	2025-02-09
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package main

import (
	"github.com/yazl-tech/beauty-rating-server/pkg/dal/model"
	"gorm.io/gen"
)

//go:generate go run generate.go
func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:       "../../pkg/dal/base",
		WithUnitTest:  false,
		FieldNullable: true,
		Mode:          gen.WithQueryInterface,
	})

	// 直接使用模型
	g.ApplyBasic(
		&model.Analysis{},
	)

	g.Execute()
}
