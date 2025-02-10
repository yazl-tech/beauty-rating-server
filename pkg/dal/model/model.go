// File:		model.go
// Created by:	Hoven
// Created on:	2025-02-09
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package model

import "github.com/go-puzzles/puzzles/pgorm"

func AllTables() []pgorm.SqlModel {
	return []pgorm.SqlModel{
		new(Analysis),
	}
}
