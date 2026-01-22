package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirkirby/go-mcp-ubiquiti/managers"
)

// RegisterNetworkTools registers all network management tools.
func RegisterNetworkTools(s *server.MCPServer, networkMgr *managers.NetworkManager) {
	// unifi_list_networks
	s.AddTool(
		mcp.NewTool("unifi_list_networks",
			mcp.WithDescription("List all networks configured on the UniFi controller"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			networks, err := networkMgr.ListNetworks(ctx)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to list networks: %v", err)), nil
			}

			result := map[string]interface{}{
				"success":  true,
				"count":    len(networks),
				"networks": networks,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_get_network_details
	s.AddTool(
		mcp.NewTool("unifi_get_network_details",
			mcp.WithDescription("Get detailed configuration for a specific network"),
			mcp.WithString("network_id",
				mcp.Required(),
				mcp.Description("ID of the network"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			networkID := req.GetString("network_id", "")
			if networkID == "" {
				return mcp.NewToolResultError("network_id is required"), nil
			}

			network, err := networkMgr.GetNetworkDetails(ctx, networkID)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to get network details: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"network": network,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_update_network
	s.AddTool(
		mcp.NewTool("unifi_update_network",
			mcp.WithDescription("Update network settings"),
			mcp.WithString("network_id",
				mcp.Required(),
				mcp.Description("ID of the network to update"),
			),
			mcp.WithString("name",
				mcp.Description("New name for the network"),
			),
			mcp.WithBoolean("enabled",
				mcp.Description("Enable or disable the network"),
			),
			mcp.WithBoolean("confirm",
				mcp.Required(),
				mcp.Description("Must be true to execute the update"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			networkID := req.GetString("network_id", "")
			confirm := req.GetBool("confirm", false)

			if networkID == "" {
				return mcp.NewToolResultError("network_id is required"), nil
			}
			if !confirm {
				return mcp.NewToolResultText(fmt.Sprintf("Preview: Would update network %s. Set confirm=true to execute.", networkID)), nil
			}

			updates := make(map[string]interface{})
			if name := req.GetString("name", ""); name != "" {
				updates["name"] = name
			}
			// For enabled, we need to check if it was explicitly provided
			// Since GetBool returns false by default, we check the raw args
			args := req.GetArguments()
			if enabledVal, ok := args["enabled"]; ok {
				if enabled, ok := enabledVal.(bool); ok {
					updates["enabled"] = enabled
				}
			}

			if len(updates) == 0 {
				return mcp.NewToolResultError("no updates specified"), nil
			}

			if err := networkMgr.UpdateNetwork(ctx, networkID, updates); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to update network: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("Network %s has been updated", networkID),
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_create_network
	s.AddTool(
		mcp.NewTool("unifi_create_network",
			mcp.WithDescription("Create a new network"),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Name for the new network"),
			),
			mcp.WithString("purpose",
				mcp.Description("Network purpose: corporate, guest, wan, vlan-only"),
			),
			mcp.WithNumber("vlan",
				mcp.Description("VLAN ID (1-4095)"),
			),
			mcp.WithString("subnet",
				mcp.Description("IP subnet (e.g., 192.168.1.0/24)"),
			),
			mcp.WithBoolean("dhcp_enabled",
				mcp.Description("Enable DHCP server"),
			),
			mcp.WithBoolean("confirm",
				mcp.Required(),
				mcp.Description("Must be true to execute creation"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			name := req.GetString("name", "")
			confirm := req.GetBool("confirm", false)

			if name == "" {
				return mcp.NewToolResultError("name is required"), nil
			}
			if !confirm {
				return mcp.NewToolResultText(fmt.Sprintf("Preview: Would create network '%s'. Set confirm=true to execute.", name)), nil
			}

			network := &managers.Network{
				Name:    name,
				Enabled: true,
			}

			if purpose := req.GetString("purpose", ""); purpose != "" {
				network.Purpose = purpose
			}
			vlan := req.GetInt("vlan", 0)
			if vlan > 0 {
				network.VLAN = vlan
				network.VLANEnabled = true
			}
			if subnet := req.GetString("subnet", ""); subnet != "" {
				network.IPSubnet = subnet
			}
			network.DHCPDEnabled = req.GetBool("dhcp_enabled", false)

			created, err := networkMgr.CreateNetwork(ctx, network)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to create network: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("Network '%s' has been created", name),
				"network": created,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_list_wlans
	s.AddTool(
		mcp.NewTool("unifi_list_wlans",
			mcp.WithDescription("List all wireless networks (WLANs)"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			wlans, err := networkMgr.ListWLANs(ctx)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to list WLANs: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"count":   len(wlans),
				"wlans":   wlans,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_get_wlan_details
	s.AddTool(
		mcp.NewTool("unifi_get_wlan_details",
			mcp.WithDescription("Get detailed configuration for a specific WLAN"),
			mcp.WithString("wlan_id",
				mcp.Required(),
				mcp.Description("ID of the WLAN"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			wlanID := req.GetString("wlan_id", "")
			if wlanID == "" {
				return mcp.NewToolResultError("wlan_id is required"), nil
			}

			wlan, err := networkMgr.GetWLANDetails(ctx, wlanID)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to get WLAN details: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"wlan":    wlan,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)
}
