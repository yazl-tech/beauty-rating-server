// File:		anlysis.go
// Created by:	Hoven
// Created on:	2025-02-09
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package analyst

import (
	"context"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/yazl-tech/beauty-rating-server/pkg/dice"
)

type Result struct {
	Score       int
	Description string
	Tags        []string
	Details     []Detail `json:"scoreDetails"`
}

type Detail struct {
	Label string
	Score int
	Desc  string
}

type Analyst interface {
	Name() string
	DoAnalysis(ctx context.Context, imageName, imageUrl string, image []byte) (*Result, error)
}

type SelectorOption func(*AnalystSelector)

func WithAnalysts(analyst Analyst, weight int) SelectorOption {
	return func(s *AnalystSelector) {
		s.analysts = append(s.analysts, analyst)
		s.weights = append(s.weights, weight)
	}
}

type AnalystSelector struct {
	weights  []int
	analysts []Analyst
	dice     *dice.Dice
}

func NewAnalystSelector(opts ...SelectorOption) *AnalystSelector {
	s := &AnalystSelector{
		analysts: []Analyst{},
	}
	for _, opt := range opts {
		opt(s)
	}

	s.dice = dice.NewDice(s.weights)
	return s
}

func (s *AnalystSelector) GetAnalyst() Analyst {
	idx := s.dice.Next()
	if idx < 0 || idx >= len(s.analysts) {
		s.dice.Reset()
		idx = s.dice.Next()
	}
	return s.analysts[idx]
}

func (s *AnalystSelector) Name() string {
	return "AnalystSelector"
}

func (s *AnalystSelector) DoAnalysis(ctx context.Context, imageName, imageUrl string, image []byte) (*Result, error) {
	analyst := s.GetAnalyst()

	plog.Debugc(ctx, "GetAnalyst: %v", analyst.Name())
	return analyst.DoAnalysis(ctx, imageName, imageUrl, image)
}
