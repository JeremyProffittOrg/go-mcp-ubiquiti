// Package tools provides MCP tool definitions for UniFi operations.
package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirkirby/go-mcp-ubiquiti/managers"
)

// RegisterClientTools registers all client management tools.
func RegisterClientTools(s *server.MCPServer, clientMgr *managers.ClientManager) {
	// unifi_list_clients
	s.AddTool(
		mcp.NewTool("unifi_list_clients",
			mcp.WithDescription("List all clients connected to the UniFi network"),
			mcp.WithBoolean("include_offline",
				mcp.Description("Include offline clients in results"),
			),
			mcp.WithNumber("limit",
				mcp.Description("Maximum number of clients to return"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			includeOffline := req.GetBool("include_offline", false)
			limit := req.GetInt("limit", 0)

			clients, err := clientMgr.ListClients(ctx, includeOffline, limit)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to list clients: %v", err)), nil
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

	// unifi_get_client_details
	s.AddTool(
		mcp.NewTool("unifi_get_client_details",
			mcp.WithDescription("Get detailed information about a specific client"),
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

			client, err := clientMgr.GetClientDetails(ctx, mac)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to get client details: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"client":  client,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_list_blocked_clients
	s.AddTool(
		mcp.NewTool("unifi_list_blocked_clients",
			mcp.WithDescription("List all blocked clients"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			clients, err := clientMgr.GetBlockedClients(ctx)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to list blocked clients: %v", err)), nil
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

	// unifi_block_client
	s.AddTool(
		mcp.NewTool("unifi_block_client",
			mcp.WithDescription("Block a client by MAC address"),
			mcp.WithString("mac",
				mcp.Required(),
				mcp.Description("MAC address of the client to block"),
			),
			mcp.WithBoolean("confirm",
				mcp.Required(),
				mcp.Description("Must be true to execute the block operation"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			mac := req.GetString("mac", "")
			confirm := req.GetBool("confirm", false)

			if mac == "" {
				return mcp.NewToolResultError("mac address is required"), nil
			}
			if !confirm {
				return mcp.NewToolResultText(fmt.Sprintf("Preview: Would block client %s. Set confirm=true to execute.", mac)), nil
			}

			if err := clientMgr.BlockClient(ctx, mac); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to block client: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("Client %s has been blocked", mac),
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_unblock_client
	s.AddTool(
		mcp.NewTool("unifi_unblock_client",
			mcp.WithDescription("Unblock a client by MAC address"),
			mcp.WithString("mac",
				mcp.Required(),
				mcp.Description("MAC address of the client to unblock"),
			),
			mcp.WithBoolean("confirm",
				mcp.Required(),
				mcp.Description("Must be true to execute the unblock operation"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			mac := req.GetString("mac", "")
			confirm := req.GetBool("confirm", false)

			if mac == "" {
				return mcp.NewToolResultError("mac address is required"), nil
			}
			if !confirm {
				return mcp.NewToolResultText(fmt.Sprintf("Preview: Would unblock client %s. Set confirm=true to execute.", mac)), nil
			}

			if err := clientMgr.UnblockClient(ctx, mac); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to unblock client: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("Client %s has been unblocked", mac),
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_rename_client
	s.AddTool(
		mcp.NewTool("unifi_rename_client",
			mcp.WithDescription("Rename/alias a client"),
			mcp.WithString("mac",
				mcp.Required(),
				mcp.Description("MAC address of the client to rename"),
			),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("New name for the client"),
			),
			mcp.WithBoolean("confirm",
				mcp.Required(),
				mcp.Description("Must be true to execute the rename operation"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			mac := req.GetString("mac", "")
			name := req.GetString("name", "")
			confirm := req.GetBool("confirm", false)

			if mac == "" || name == "" {
				return mcp.NewToolResultError("mac and name are required"), nil
			}
			if !confirm {
				return mcp.NewToolResultText(fmt.Sprintf("Preview: Would rename client %s to '%s'. Set confirm=true to execute.", mac, name)), nil
			}

			if err := clientMgr.RenameClient(ctx, mac, name); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to rename client: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("Client %s has been renamed to '%s'", mac, name),
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_force_reconnect_client
	s.AddTool(
		mcp.NewTool("unifi_force_reconnect_client",
			mcp.WithDescription("Force a client to reconnect"),
			mcp.WithString("mac",
				mcp.Required(),
				mcp.Description("MAC address of the client to reconnect"),
			),
			mcp.WithBoolean("confirm",
				mcp.Required(),
				mcp.Description("Must be true to execute the reconnect operation"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			mac := req.GetString("mac", "")
			confirm := req.GetBool("confirm", false)

			if mac == "" {
				return mcp.NewToolResultError("mac address is required"), nil
			}
			if !confirm {
				return mcp.NewToolResultText(fmt.Sprintf("Preview: Would force reconnect client %s. Set confirm=true to execute.", mac)), nil
			}

			if err := clientMgr.ForceReconnect(ctx, mac); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to force reconnect: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("Client %s has been kicked and will reconnect", mac),
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_authorize_guest
	s.AddTool(
		mcp.NewTool("unifi_authorize_guest",
			mcp.WithDescription("Authorize a guest client"),
			mcp.WithString("mac",
				mcp.Required(),
				mcp.Description("MAC address of the guest client"),
			),
			mcp.WithNumber("minutes",
				mcp.Description("Duration of authorization in minutes (default: 60)"),
			),
			mcp.WithNumber("up_kbps",
				mcp.Description("Upload bandwidth limit in Kbps"),
			),
			mcp.WithNumber("down_kbps",
				mcp.Description("Download bandwidth limit in Kbps"),
			),
			mcp.WithNumber("megabytes",
				mcp.Description("Data transfer limit in MB"),
			),
			mcp.WithBoolean("confirm",
				mcp.Required(),
				mcp.Description("Must be true to execute the authorization"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			mac := req.GetString("mac", "")
			confirm := req.GetBool("confirm", false)
			minutes := req.GetInt("minutes", 60)
			upKbps := req.GetInt("up_kbps", 0)
			downKbps := req.GetInt("down_kbps", 0)
			mbytes := req.GetInt("megabytes", 0)

			if mac == "" {
				return mcp.NewToolResultError("mac address is required"), nil
			}
			if !confirm {
				return mcp.NewToolResultText(fmt.Sprintf("Preview: Would authorize guest %s for %d minutes. Set confirm=true to execute.", mac, minutes)), nil
			}

			if err := clientMgr.AuthorizeGuest(ctx, mac, minutes, upKbps, downKbps, mbytes); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to authorize guest: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("Guest %s has been authorized for %d minutes", mac, minutes),
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_unauthorize_guest
	s.AddTool(
		mcp.NewTool("unifi_unauthorize_guest",
			mcp.WithDescription("Revoke guest authorization"),
			mcp.WithString("mac",
				mcp.Required(),
				mcp.Description("MAC address of the guest client"),
			),
			mcp.WithBoolean("confirm",
				mcp.Required(),
				mcp.Description("Must be true to execute the operation"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			mac := req.GetString("mac", "")
			confirm := req.GetBool("confirm", false)

			if mac == "" {
				return mcp.NewToolResultError("mac address is required"), nil
			}
			if !confirm {
				return mcp.NewToolResultText(fmt.Sprintf("Preview: Would unauthorize guest %s. Set confirm=true to execute.", mac)), nil
			}

			if err := clientMgr.UnauthorizeGuest(ctx, mac); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to unauthorize guest: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("Guest %s authorization has been revoked", mac),
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)
}
