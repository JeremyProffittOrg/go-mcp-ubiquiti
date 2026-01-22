package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirkirby/go-mcp-ubiquiti/managers"
)

// RegisterStatsTools registers all statistics tools.
func RegisterStatsTools(s *server.MCPServer, statsMgr *managers.StatsManager) {
	// unifi_get_network_stats
	s.AddTool(
		mcp.NewTool("unifi_get_network_stats",
			mcp.WithDescription("Get network-wide statistics"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			stats, err := statsMgr.GetNetworkStats(ctx)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to get network stats: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"stats":   stats,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_get_client_stats
	s.AddTool(
		mcp.NewTool("unifi_get_client_stats",
			mcp.WithDescription("Get statistics for a specific client"),
			mcp.WithString("mac",
				mcp.Required(),
				mcp.Description("MAC address of the client"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			mac := req.GetString("mac", "")
			if mac == "" {
				return mcp.NewToolResultError("mac address is required"), nil
			}

			stats, err := statsMgr.GetClientStats(ctx, mac)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to get client stats: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"mac":     mac,
				"stats":   stats,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_get_device_stats
	s.AddTool(
		mcp.NewTool("unifi_get_device_stats",
			mcp.WithDescription("Get statistics for a specific device"),
			mcp.WithString("mac",
				mcp.Required(),
				mcp.Description("MAC address of the device"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			mac := req.GetString("mac", "")
			if mac == "" {
				return mcp.NewToolResultError("mac address is required"), nil
			}

			stats, err := statsMgr.GetDeviceStats(ctx, mac)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to get device stats: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"mac":     mac,
				"stats":   stats,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_get_top_clients
	s.AddTool(
		mcp.NewTool("unifi_get_top_clients",
			mcp.WithDescription("Get top clients by bandwidth usage"),
			mcp.WithNumber("limit",
				mcp.Description("Number of top clients to return (default: 10)"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			limit := req.GetInt("limit", 10)

			clients, err := statsMgr.GetTopClients(ctx, limit)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to get top clients: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"count":   len(clients),
				"clients": clients,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_get_dpi_stats
	s.AddTool(
		mcp.NewTool("unifi_get_dpi_stats",
			mcp.WithDescription("Get Deep Packet Inspection (DPI) statistics"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			stats, err := statsMgr.GetDPIStats(ctx)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to get DPI stats: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"count":   len(stats),
				"dpi":     stats,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)
}
