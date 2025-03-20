package main

import (
	"log"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"Waifu MCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	// Add tool
	searchWaifuTool := mcp.NewTool("search_anime_character",
		mcp.WithDescription("Search anime character by name"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the character to search"),
		),
	)

	// Add tool handler
	s.AddTool(searchWaifuTool, handleBGMCharacterSearch)

	// Start the server
	sseServer := server.NewSSEServer(s,
		server.WithBaseURL("http://localhost:8080"),
	)
	slog.Info("SSE server listening on :8080")
	if err := sseServer.Start(":8080"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
