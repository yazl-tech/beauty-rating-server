// File:		oss.go
// Created by:	Hoven
// Created on:	2024-11-05
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package oss

import (
	"context"
	"io"
)

type IOSS interface {
	UploadFile(ctx context.Context, size int64, dir, objName string, obj io.Reader) (uri string, err error)
	GetFile(ctx context.Context, objName string, w io.Writer) error
}
