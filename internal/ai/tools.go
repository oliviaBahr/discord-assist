package ai

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
)

// Tool represents a unified tool definition with both schema and execution logic
type Tool struct {
	anthropic.ToolParam
	Execute func(params map[string]any) (string, error)
}

// ToolRegistry holds all available tools
type ToolRegistry struct {
	tools map[string]*Tool
}

// Global tool registry instance
var GlobalToolRegistry = &ToolRegistry{
	tools: map[string]*Tool{
		"get_current_time": {
			ToolParam: anthropic.ToolParam{
				Type:        anthropic.ToolTypeCustom,
				Name:        "get_current_time",
				Description: anthropic.String("Get the current time and date for a specified timezone"),
				InputSchema: anthropic.ToolInputSchemaParam{
					Type: "object",
					Properties: map[string]any{
						"timezone": map[string]any{
							"type":        "string",
							"description": "The timezone to get the time for (e.g., 'UTC', 'America/New_York')",
							"default":     "UTC",
						},
					},
					Required: []string{},
				},
			},
			Execute: func(params map[string]any) (string, error) {
				timezone := "UTC"
				if tz, ok := params["timezone"].(string); ok {
					timezone = tz
				}

				loc, err := time.LoadLocation(timezone)
				if err != nil {
					loc = time.UTC
				}

				now := time.Now().In(loc)
				return fmt.Sprintf("Current time in %s: %s", timezone, now.Format("2006-01-02 15:04:05 MST")), nil
			},
		},
		"get_weather": {
			ToolParam: anthropic.ToolParam{
				Type:        anthropic.ToolTypeCustom,
				Name:        "get_weather",
				Description: anthropic.String("Get current weather information for a specific location"),
				InputSchema: anthropic.ToolInputSchemaParam{
					Type: "object",
					Properties: map[string]any{
						"location": map[string]any{
							"type":        "string",
							"description": "The city and state, or city and country",
						},
						"unit": map[string]any{
							"type":        "string",
							"enum":        []string{"celsius", "fahrenheit"},
							"description": "The temperature unit to use",
							"default":     "celsius",
						},
					},
					Required: []string{"location"},
				},
			},
			Execute: func(params map[string]any) (string, error) {
				location, ok := params["location"].(string)
				if !ok {
					return "", fmt.Errorf("location parameter is required")
				}

				unit := "celsius"
				if u, ok := params["unit"].(string); ok {
					unit = u
				}

				// Mock weather data - in a real implementation, you'd call a weather API
				return fmt.Sprintf("Weather in %s: 22°%s, Partly Cloudy, Humidity: 65%%", location, unit), nil
			},
		},
		"search_web": {
			ToolParam: anthropic.ToolParam{
				Type:        anthropic.ToolTypeCustom,
				Name:        "search_web",
				Description: anthropic.String("Search the web for current information on a specific topic"),
				InputSchema: anthropic.ToolInputSchemaParam{
					Type: "object",
					Properties: map[string]any{
						"query": map[string]any{
							"type":        "string",
							"description": "The search query",
						},
					},
					Required: []string{"query"},
				},
			},
			Execute: func(params map[string]any) (string, error) {
				query, ok := params["query"].(string)
				if !ok {
					return "", fmt.Errorf("query parameter is required")
				}

				// Mock search results - in a real implementation, you'd call a search API
				return fmt.Sprintf("Search results for '%s': Found 1,234 results. Here are the top 3:\n1. Example result 1\n2. Example result 2\n3. Example result 3", query), nil
			},
		},
	},
}

// ExecuteTool executes a tool by name with the given JSON input and returns a tool result block
func (tr *ToolRegistry) ExecuteTool(name string, input json.RawMessage, toolUseID string) (anthropic.ContentBlockParamUnion, error) {
	tool, exists := tr.tools[name]
	if !exists {
		return anthropic.ContentBlockParamUnion{}, fmt.Errorf("unknown tool: %s", name)
	}

	var params map[string]any
	if err := json.Unmarshal(input, &params); err != nil {
		return anthropic.ContentBlockParamUnion{}, fmt.Errorf("failed to parse tool input: %w", err)
	}

	result, err := tool.Execute(params)
	if err != nil {
		return anthropic.ContentBlockParamUnion{}, err
	}

	return anthropic.NewToolResultBlock(toolUseID, result, false), nil
}

// GetTools returns the available tools for the AI service
func GetTools() []anthropic.ToolUnionParam {
	var unionTools []anthropic.ToolUnionParam
	for _, tool := range GlobalToolRegistry.tools {
		unionTools = append(unionTools, anthropic.ToolUnionParam{
			OfTool: &tool.ToolParam,
		})
	}

	return unionTools
}
