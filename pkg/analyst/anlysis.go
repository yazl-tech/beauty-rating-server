// File:		anlysis.go
// Created by:	Hoven
// Created on:	2025-02-09
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package analyst

import "context"

type Result struct {
	Score       int
	Description string
	Tags        []string
	Details     []Detail
}

type Detail struct {
	Label string
	Score int
	Desc  string
}

type Analyst interface {
	DoAnalysis(ctx context.Context, image []byte) (*Result, error)
}
