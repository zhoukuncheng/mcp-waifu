package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
	"resty.dev/v3"
)

type SearchCharacterResponse struct {
	Data []struct {
		ID        int     `json:"id"`
		Name      string  `json:"name"`
		Summary   string  `json:"summary"`
		Gender    string  `json:"gender"`
		BirthYear *int    `json:"birth_year"`
		BirthMon  *int    `json:"birth_mon"`
		BirthDay  *int    `json:"birth_day"`
		BloodType *string `json:"blood_type"`
		Images    struct {
			Small  string `json:"small"`
			Grid   string `json:"grid"`
			Large  string `json:"large"`
			Medium string `json:"medium"`
		} `json:"images"`
		Stat    Stat          `json:"stat"`
		Locked  bool          `json:"locked"`
		Type    int           `json:"type"`
		Infobox []InfoboxItem `json:"infobox"`
		NSFW    bool          `json:"nsfw"`
	} `json:"data"`
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type Character struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	// 其他字段可根据需要添加
}

func handleBGMCharacterSearch(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	arguments := request.Params.Arguments
	name, ok := arguments["name"].(string)
	if !ok || name == "" {
		return nil, errors.New("Invalid name")
	}
	slog.Info("search name: ", "name", name)

	character, err := searchBangumiCharacter(ctx, name)
	if err != nil {
		slog.Error("Error searching bangumi character: ", "error", err)
		return nil, fmt.Errorf("failed to search character: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Character: %s", character.Name),
			},
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Description: %s", character.Description),
			},
			mcp.TextContent{
				Type: "text",
				Text: character.ImageURL,
			},
		},
	}, nil
}

type InfoboxItem struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"` // Can be string or array of objects
}

type Stat struct {
	Comments int `json:"comments"`
	Collects int `json:"collects"`
}

func searchBangumiCharacter(ctx context.Context, name string) (*Character, error) {
	client := resty.New()

	var response SearchCharacterResponse

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "ozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36").
		SetBody(map[string]string{
			"keyword": name,
		}).
		SetResult(&response).
		Post("https://api.bgm.tv/v0/search/characters?limit=1")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode())
	}

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("no character found with name: %s", name)
	}

	character := &Character{
		Name:        response.Data[0].Name,
		Description: response.Data[0].Summary,
		ImageURL:    response.Data[0].Images.Large, // Using the large image URL
	}
	slog.Info("search character: ", "character", character)
	return character, nil
}
