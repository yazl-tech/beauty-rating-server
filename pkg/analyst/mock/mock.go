// File:		mock.go
// Created by:	Hoven
// Created on:	2025-02-09
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package mock

import (
	"context"
	"math/rand"
	"time"

	"github.com/yazl-tech/beauty-rating-server/pkg/analyst"
)

var (
	rander *rand.Rand
)

func init() {
	rander = rand.New(rand.NewSource(time.Now().UnixNano()))
}

var _ analyst.Analyst = (*MockAnalyst)(nil)

type MockAnalyst struct {
}

func NewMockAnalyst() *MockAnalyst {
	return &MockAnalyst{}
}

func randomNumber(min, max int) int {
	return rander.Intn(max-min+1) + min
}

func randomPick(arr []string) string {
	return arr[rander.Intn(len(arr))]
}

func (m *MockAnalyst) generateScoreDetails() []analyst.Detail {
	var details []analyst.Detail
	for _, item := range scoreLabels {
		details = append(details, analyst.Detail{
			Label: item.Label,
			Score: randomNumber(85, 98),
			Desc:  item.Descs[rander.Intn(len(item.Descs))],
		})
	}
	return details
}

func (m *MockAnalyst) DoAnalysis(_ context.Context, _ []byte) (*analyst.Result, error) {
	rander.Shuffle(len(allTags), func(i, j int) {
		allTags[i], allTags[j] = allTags[j], allTags[i]
	})

	return &analyst.Result{
		Score:       randomNumber(85, 99),
		Description: descriptions[rander.Intn(len(descriptions))],
		Tags:        allTags[:randomNumber(3, 6)],
		Details:     m.generateScoreDetails(),
	}, nil
}
