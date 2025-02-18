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
	"time"

	"github.com/yazl-tech/beauty-rating-server/pkg/analyst"
	"github.com/yazl-tech/beauty-rating-server/pkg/random"
)

var _ analyst.Analyst = (*MockAnalyst)(nil)

type MockAnalyst struct {
}

func NewMockAnalyst() *MockAnalyst {
	return &MockAnalyst{}
}

func (m *MockAnalyst) generateScoreDetails() []analyst.Detail {
	var details []analyst.Detail
	for _, item := range scoreLabels {
		details = append(details, analyst.Detail{
			Label: item.Label,
			Score: random.RandomNumber(85, 98),
			Desc:  item.Descs[random.RandomInt(len(item.Descs))],
		})
	}
	return details
}

func (m *MockAnalyst) Name() string {
	return "MockAnalyst"
}

func (m *MockAnalyst) Typ() analyst.AnalystType {
	return analyst.TypeMock
}

func (m *MockAnalyst) DoAnalysis(_ context.Context, _, _ string, _ []byte) (*analyst.Result, error) {
	time.Sleep(time.Duration(random.RandomNumber(4, 9)))

	random.RandomShuffle(len(allTags), func(i, j int) {
		allTags[i], allTags[j] = allTags[j], allTags[i]
	})

	return &analyst.Result{
		AnalystType: analyst.TypeMock,
		Score:       random.RandomNumber(85, 99),
		Description: descriptions[random.RandomInt(len(descriptions))],
		Tags:        allTags[:random.RandomNumber(3, 6)],
		Details:     m.generateScoreDetails(),
	}, nil
}
