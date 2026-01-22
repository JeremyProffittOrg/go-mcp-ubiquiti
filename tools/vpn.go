package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirkirby/go-mcp-ubiquiti/managers"
)

// RegisterVPNTools registers all VPN-related tools.
func RegisterVPNTools(s *server.MCPServer, vpnMgr *managers.VPNManager) {
	// unifi_list_vpn_clients
	s.AddTool(
		mcp.NewTool("unifi_list_vpn_clients",
			mcp.WithDescription("List all VPN client configurations"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			clients, err := vpnMgr.ListVPNClients(ctx)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to list VPN clients: %v", err)), nil
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

	// unifi_get_vpn_client_details
	s.AddTool(
		mcp.NewTool("unifi_get_vpn_client_details",
			mcp.WithDescription("Get detailed information about a specific VPN client"),
			mcp.WithString("client_id",
				mcp.Required(),
				mcp.Description("ID of the VPN client"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			clientID := req.GetString("client_id", "")
			if clientID == "" {
				return mcp.NewToolResultError("client_id is required"), nil
			}

			client, err := vpnMgr.GetVPNClientDetails(ctx, clientID)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to get VPN client details: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"client":  client,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_update_vpn_client_state
	s.AddTool(
		mcp.NewTool("unifi_update_vpn_client_state",
			mcp.WithDescription("Enable or disable a VPN client"),
			mcp.WithString("client_id",
				mcp.Required(),
				mcp.Description("ID of the VPN client"),
			),
			mcp.WithBoolean("enabled",
				mcp.Required(),
				mcp.Description("Whether to enable (true) or disable (false) the client"),
			),
			mcp.WithBoolean("confirm",
				mcp.Required(),
				mcp.Description("Must be true to execute the update"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			clientID := req.GetString("client_id", "")
			enabled := req.GetBool("enabled", false)
			confirm := req.GetBool("confirm", false)

			if clientID == "" {
				return mcp.NewToolResultError("client_id is required"), nil
			}

			action := "disable"
			if enabled {
				action = "enable"
			}

			if !confirm {
				return mcp.NewToolResultText(fmt.Sprintf("Preview: Would %s VPN client %s. Set confirm=true to execute.", action, clientID)), nil
			}

			if err := vpnMgr.UpdateVPNClientState(ctx, clientID, enabled); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to update VPN client state: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("VPN client %s has been %sd", clientID, action),
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_list_vpn_servers
	s.AddTool(
		mcp.NewTool("unifi_list_vpn_servers",
			mcp.WithDescription("List all VPN server configurations"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			servers, err := vpnMgr.ListVPNServers(ctx)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to list VPN servers: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"count":   len(servers),
				"servers": servers,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_get_vpn_server_details
	s.AddTool(
		mcp.NewTool("unifi_get_vpn_server_details",
			mcp.WithDescription("Get detailed information about a specific VPN server"),
			mcp.WithString("server_id",
				mcp.Required(),
				mcp.Description("ID of the VPN server"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			serverID := req.GetString("server_id", "")
			if serverID == "" {
				return mcp.NewToolResultError("server_id is required"), nil
			}

			server, err := vpnMgr.GetVPNServerDetails(ctx, serverID)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to get VPN server details: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"server":  server,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_update_vpn_server_state
	s.AddTool(
		mcp.NewTool("unifi_update_vpn_server_state",
			mcp.WithDescription("Enable or disable a VPN server"),
			mcp.WithString("server_id",
				mcp.Required(),
				mcp.Description("ID of the VPN server"),
			),
			mcp.WithBoolean("enabled",
				mcp.Required(),
				mcp.Description("Whether to enable (true) or disable (false) the server"),
			),
			mcp.WithBoolean("confirm",
				mcp.Required(),
				mcp.Description("Must be true to execute the update"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			serverID := req.GetString("server_id", "")
			enabled := req.GetBool("enabled", false)
			confirm := req.GetBool("confirm", false)

			if serverID == "" {
				return mcp.NewToolResultError("server_id is required"), nil
			}

			action := "disable"
			if enabled {
				action = "enable"
			}

			if !confirm {
				return mcp.NewToolResultText(fmt.Sprintf("Preview: Would %s VPN server %s. Set confirm=true to execute.", action, serverID)), nil
			}

			if err := vpnMgr.UpdateVPNServerState(ctx, serverID, enabled); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to update VPN server state: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("VPN server %s has been %sd", serverID, action),
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)
}
