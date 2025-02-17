// File:		dice.go
// Created by:	Hoven
// Created on:	2025-02-18
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package dice

import (
	"math/rand"
)

type Dice struct {
	total    int
	weights  []int
	original []int
}

func NewDice(weights []int) *Dice {
	total := 0

	for _, w := range weights {
		total += w
	}
	return &Dice{
		total:    total,
		weights:  append([]int(nil), weights...),
		original: append([]int(nil), weights...),
	}
}

func (d *Dice) Next() int {
	if d.total == 0 {
		return -1
	}

	v := rand.Intn(d.total)
	for i, w := range d.weights {
		if v < w {
			d.total -= w
			d.weights[i] = 0
			return i
		}
		v -= w
	}
	return -1
}

func (d *Dice) Reset() {
	d.total = 0
	copy(d.weights, d.original)
	for _, w := range d.weights {
		d.total += w
	}
}
