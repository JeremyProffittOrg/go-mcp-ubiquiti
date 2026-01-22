package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirkirby/go-mcp-ubiquiti/managers"
)

// RegisterSystemTools registers all system management tools.
func RegisterSystemTools(s *server.MCPServer, systemMgr *managers.SystemManager, site string) {
	// unifi_get_system_info
	s.AddTool(
		mcp.NewTool("unifi_get_system_info",
			mcp.WithDescription("Get general system information from the UniFi Network controller (version, uptime, etc)"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			info, err := systemMgr.GetSystemInfo(ctx)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to get system info: %v", err)), nil
			}

			result := map[string]interface{}{
				"success":     true,
				"site":        site,
				"system_info": info,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_get_network_health
	s.AddTool(
		mcp.NewTool("unifi_get_network_health",
			mcp.WithDescription("Get network health status including WAN status and device counts"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			health, err := systemMgr.GetNetworkHealth(ctx)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to get network health: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"site":    site,
				"health":  health,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_get_site_settings
	s.AddTool(
		mcp.NewTool("unifi_get_site_settings",
			mcp.WithDescription("Get current site settings (country code, timezone, connectivity monitoring, etc)"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			settings, err := systemMgr.GetSiteSettings(ctx)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to get site settings: %v", err)), nil
			}

			result := map[string]interface{}{
				"success":  true,
				"site":     site,
				"settings": settings,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)
}
