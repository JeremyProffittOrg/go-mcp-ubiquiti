package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirkirby/go-mcp-ubiquiti/managers"
)

// RegisterEventTools registers all event and alarm tools.
func RegisterEventTools(s *server.MCPServer, eventMgr *managers.EventManager) {
	// unifi_list_events
	s.AddTool(
		mcp.NewTool("unifi_list_events",
			mcp.WithDescription("List recent events from the UniFi controller"),
			mcp.WithNumber("hours",
				mcp.Description("Number of hours to look back (default: 24)"),
			),
			mcp.WithNumber("limit",
				mcp.Description("Maximum number of events to return (default: 100)"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			hours := req.GetInt("hours", 24)
			limit := req.GetInt("limit", 100)

			events, err := eventMgr.ListEvents(ctx, hours, limit)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to list events: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"count":   len(events),
				"events":  events,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_list_alarms
	s.AddTool(
		mcp.NewTool("unifi_list_alarms",
			mcp.WithDescription("List active alarms from the UniFi controller"),
			mcp.WithBoolean("include_archived",
				mcp.Description("Include archived alarms (default: false)"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			archived := req.GetBool("include_archived", false)

			alarms, err := eventMgr.ListAlarms(ctx, archived)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to list alarms: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"count":   len(alarms),
				"alarms":  alarms,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_archive_alarm
	s.AddTool(
		mcp.NewTool("unifi_archive_alarm",
			mcp.WithDescription("Archive a specific alarm"),
			mcp.WithString("alarm_id",
				mcp.Required(),
				mcp.Description("ID of the alarm to archive"),
			),
			mcp.WithBoolean("confirm",
				mcp.Required(),
				mcp.Description("Must be true to execute the archive"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			alarmID := req.GetString("alarm_id", "")
			confirm := req.GetBool("confirm", false)

			if alarmID == "" {
				return mcp.NewToolResultError("alarm_id is required"), nil
			}
			if !confirm {
				return mcp.NewToolResultText(fmt.Sprintf("Preview: Would archive alarm %s. Set confirm=true to execute.", alarmID)), nil
			}

			if err := eventMgr.ArchiveAlarm(ctx, alarmID); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to archive alarm: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("Alarm %s has been archived", alarmID),
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_archive_all_alarms
	s.AddTool(
		mcp.NewTool("unifi_archive_all_alarms",
			mcp.WithDescription("Archive all active alarms"),
			mcp.WithBoolean("confirm",
				mcp.Required(),
				mcp.Description("Must be true to execute the archive"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			confirm := req.GetBool("confirm", false)

			if !confirm {
				return mcp.NewToolResultText("Preview: Would archive all alarms. Set confirm=true to execute."), nil
			}

			if err := eventMgr.ArchiveAllAlarms(ctx); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to archive all alarms: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"message": "All alarms have been archived",
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_get_event_types
	s.AddTool(
		mcp.NewTool("unifi_get_event_types",
			mcp.WithDescription("Get a list of available event types"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			types := eventMgr.GetEventTypes()

			result := map[string]interface{}{
				"success":     true,
				"count":       len(types),
				"event_types": types,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)
}
