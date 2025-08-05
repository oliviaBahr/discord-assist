# AI Service Tools

This directory contains the AI service with tool support for the Discord bot.

## Tools

The AI service currently supports the following tools:

### 1. `get_current_time`
- **Description**: Get the current time and date for a specified timezone
- **Parameters**:
  - `timezone` (optional): The timezone to get the time for (e.g., 'UTC', 'America/New_York')
  - Default: 'UTC'

### 2. `get_weather`
- **Description**: Get current weather information for a specific location
- **Parameters**:
  - `location` (required): The city and state, or city and country
  - `unit` (optional): The temperature unit to use ('celsius' or 'fahrenheit')
  - Default unit: 'celsius'

### 3. `search_web`
- **Description**: Search the web for current information on a specific topic
- **Parameters**:
  - `query` (required): The search query

## Adding New Tools

To add a new tool:

1. **Define the tool in `tools.go`**:
   ```go
   "your_tool_name": {
       ToolParam: anthropic.ToolParam{
           Type:        anthropic.ToolTypeCustom,
           Name:        "your_tool_name",
           Description: anthropic.String("Description of what your tool does"),
           InputSchema: anthropic.ToolInputSchemaParam{
               Type: "object",
               Properties: map[string]interface{}{
                   "parameter_name": map[string]interface{}{
                       "type":        "string",
                       "description": "Description of the parameter",
                   },
               },
               Required: []string{"parameter_name"},
           },
       },
       Execute: func(params map[string]interface{}) (string, error) {
           // Your tool logic here
           return "Tool result", nil
       },
   },
   ```

2. **Add the tool execution logic in `service.go`**:
   ```go
   case "your_tool_name":
       return s.yourToolFunction(toolUse.Input)
   ```

3. **Implement the tool function**:
   ```go
   func (s *Service) yourToolFunction(input json.RawMessage) (string, error) {
       var params map[string]interface{}
       if err := json.Unmarshal(input, &params); err != nil {
           return "", fmt.Errorf("failed to parse tool input: %w", err)
       }
       
       // Your tool logic here
       return "Tool result", nil
   }
   ```

## Tool Parameters

When defining tool parameters, you can use these JSON schema types:
- `"string"`: Text values
- `"number"`: Numeric values
- `"boolean"`: True/false values
- `"object"`: Complex objects
- `"array"`: Lists of values

You can also add constraints like:
- `"enum"`: List of allowed values
- `"default"`: Default value
- `"description"`: Parameter description

## Example Usage

Users can ask the bot to use tools like:
- "What time is it in Tokyo?"
- "What's the weather like in New York?"
- "Search for the latest news about AI" 