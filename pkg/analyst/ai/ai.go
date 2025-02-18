// File:		ai.go
// Created by:	Hoven
// Created on:	2025-02-17
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package ai

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/yazl-tech/beauty-rating-server/pkg/analyst"
	"github.com/yazl-tech/beauty-rating-server/pkg/random"

	botpb "github.com/yazl-tech/ai-bot/pkg/proto/bot"
	doubaopb "github.com/yazl-tech/ai-bot/pkg/proto/doubao"
)

var _ analyst.Analyst = (*AiAnalyst)(nil)

var extMimeTypeMap = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".apng": "image/png",
	".png":  "image/png",
	".gif":  "image/gif",
	".webp": "image/webp",
	".bmp":  "image/bmp",
	".tiff": "image/tiff",
	".tif":  "image/tiff",
	".ico":  "image/x-icon",
	".dib":  "image/bmp",
	".icns": "image/icns",
	".sgi":  "image/sgi",
	".j2c":  "image/jp2",
	".j2k":  "image/jp2",
	".jp2":  "image/jp2",
	".jpc":  "image/jp2",
	".jpf":  "image/jp2",
	".jpx":  "image/jp2",
}

type AiAnalyst struct {
	model        string
	doubaoClient doubaopb.DoubaoHandlerClient
}

func NewAiAnalyst(model string, doubaoClient doubaopb.DoubaoHandlerClient) *AiAnalyst {
	return &AiAnalyst{model: model, doubaoClient: doubaoClient}
}

func (a *AiAnalyst) Name() string {
	return "AiAnalyst"
}

func (a *AiAnalyst) Typ() analyst.AnalystType {
	return analyst.TypeAi
}

func (a *AiAnalyst) calcBase64(b []byte) string {
	encoder := base64.StdEncoding
	buf := make([]byte, encoder.EncodedLen(len(b)))

	encoder.Encode(buf, b)

	return string(buf)
}

func (a *AiAnalyst) generateImageUrl(imageName string, image []byte) string {
	base64Encode := a.calcBase64(image)

	ext := strings.ToLower(filepath.Ext(imageName))
	mimeType := "image/png"
	if mt, ok := extMimeTypeMap[ext]; ok {
		mimeType = mt
	}
	return fmt.Sprintf("data:%s;base64,%s", mimeType, base64Encode)
}

func (a *AiAnalyst) packRequest(imageUrl string) *botpb.ChatRequest {
	return &botpb.ChatRequest{
		Messages: []*botpb.Message{
			{
				Role: botpb.Message_system,
				Content: &botpb.Message_StringContent{
					StringContent: systemPrompt,
				},
			},
			{
				Role: botpb.Message_user,
				Content: &botpb.Message_TypeContent{
					TypeContent: &botpb.TypeMessage{
						Type: botpb.TypeMessage_image,
						ImageUrl: &botpb.TypeMessage_ImageUrl{
							Url: imageUrl,
						},
					},
				},
			},
		},
		Options: &botpb.ChatOptions{
			Model:       a.model,
			Temperature: 1,
		},
	}
}

func (a *AiAnalyst) choiceTags(tags []string) []string {
	tl := len(tags)

	if tl < 4 {
		return tags
	}

	if tl >= 4 && tl <= 6 {
		return tags
	}

	random.RandomShuffle(len(tags), func(i, j int) {
		tags[i], tags[j] = tags[j], tags[i]
	})

	return tags[:random.RandomNumber(3, 6)]
}

func (a *AiAnalyst) parseAiResp(choices []*botpb.Choice) (*analyst.Result, error) {
	if len(choices) == 0 {
		return nil, fmt.Errorf("empty choices")
	}

	choice := choices[0]
	strContent := choice.GetMessage().GetStringContent()

	re := regexp.MustCompile(`(?s)\{.*\}`)
	jsonMatch := re.FindString(strContent)
	if jsonMatch == "" {
		return nil, fmt.Errorf("no valid JSON found in response")
	}

	ret := &analyst.Result{}
	err := json.Unmarshal([]byte(jsonMatch), ret)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalChoice")
	}
	ret.Tags = a.choiceTags(ret.Tags)

	return ret, nil
}

func (a *AiAnalyst) DoAnalysis(ctx context.Context, imageName, imageUrl string, image []byte) (*analyst.Result, error) {
	imageUrl = a.generateImageUrl(imageName, image)

	resp, err := a.doubaoClient.ChatCompletions(ctx, a.packRequest(imageUrl))
	if err != nil {
		return nil, err
	}

	ret, err := a.parseAiResp(resp.GetChoices())
	if err != nil {
		return nil, err
	}

	ret.AnalystType = analyst.TypeAi
	return ret, nil
}
