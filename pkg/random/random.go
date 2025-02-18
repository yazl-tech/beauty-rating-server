// File:		random.go
// Created by:	Hoven
// Created on:	2025-02-18
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package random

import (
	"math/rand"
	"time"
)

var (
	rander *rand.Rand
)

func init() {
	rander = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomNumber(min, max int) int {
	return rander.Intn(max-min+1) + min
}

func RandomPick(arr []string) string {
	return arr[rander.Intn(len(arr))]
}

func RandomInt(n int) int {
	return rander.Intn(n)
}

func RandomShuffle(n int, swap func(i, j int)) {
	rander.Shuffle(n, swap)
}
