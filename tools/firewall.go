package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirkirby/go-mcp-ubiquiti/managers"
)

// RegisterFirewallTools registers all firewall-related tools.
func RegisterFirewallTools(s *server.MCPServer, firewallMgr *managers.FirewallManager) {
	// unifi_list_firewall_policies
	s.AddTool(
		mcp.NewTool("unifi_list_firewall_policies",
			mcp.WithDescription("List all firewall policies"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			policies, err := firewallMgr.ListFirewallPolicies(ctx)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to list firewall policies: %v", err)), nil
			}

			result := map[string]interface{}{
				"success":  true,
				"count":    len(policies),
				"policies": policies,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_get_firewall_policy_details
	s.AddTool(
		mcp.NewTool("unifi_get_firewall_policy_details",
			mcp.WithDescription("Get detailed information about a specific firewall policy"),
			mcp.WithString("policy_id",
				mcp.Required(),
				mcp.Description("ID of the firewall policy"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			policyID := req.GetString("policy_id", "")
			if policyID == "" {
				return mcp.NewToolResultError("policy_id is required"), nil
			}

			policy, err := firewallMgr.GetFirewallPolicyDetails(ctx, policyID)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to get firewall policy: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"policy":  policy,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_toggle_firewall_policy
	s.AddTool(
		mcp.NewTool("unifi_toggle_firewall_policy",
			mcp.WithDescription("Enable or disable a firewall policy"),
			mcp.WithString("policy_id",
				mcp.Required(),
				mcp.Description("ID of the firewall policy"),
			),
			mcp.WithBoolean("enabled",
				mcp.Required(),
				mcp.Description("Whether to enable (true) or disable (false) the policy"),
			),
			mcp.WithBoolean("confirm",
				mcp.Required(),
				mcp.Description("Must be true to execute the toggle"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			policyID := req.GetString("policy_id", "")
			enabled := req.GetBool("enabled", false)
			confirm := req.GetBool("confirm", false)

			if policyID == "" {
				return mcp.NewToolResultError("policy_id is required"), nil
			}

			action := "disable"
			if enabled {
				action = "enable"
			}

			if !confirm {
				return mcp.NewToolResultText(fmt.Sprintf("Preview: Would %s firewall policy %s. Set confirm=true to execute.", action, policyID)), nil
			}

			if err := firewallMgr.ToggleFirewallPolicy(ctx, policyID, enabled); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to toggle firewall policy: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("Firewall policy %s has been %sd", policyID, action),
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_list_firewall_zones
	s.AddTool(
		mcp.NewTool("unifi_list_firewall_zones",
			mcp.WithDescription("List all firewall zones (UniFi OS only)"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			zones, err := firewallMgr.ListFirewallZones(ctx)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to list firewall zones: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"count":   len(zones),
				"zones":   zones,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_list_ip_groups
	s.AddTool(
		mcp.NewTool("unifi_list_ip_groups",
			mcp.WithDescription("List all IP address groups"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			groups, err := firewallMgr.ListIPGroups(ctx)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to list IP groups: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"count":   len(groups),
				"groups":  groups,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_list_port_forwards
	s.AddTool(
		mcp.NewTool("unifi_list_port_forwards",
			mcp.WithDescription("List all port forwarding rules"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			forwards, err := firewallMgr.ListPortForwards(ctx)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to list port forwards: %v", err)), nil
			}

			result := map[string]interface{}{
				"success":  true,
				"count":    len(forwards),
				"forwards": forwards,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_get_port_forward_details
	s.AddTool(
		mcp.NewTool("unifi_get_port_forward_details",
			mcp.WithDescription("Get detailed information about a specific port forward"),
			mcp.WithString("forward_id",
				mcp.Required(),
				mcp.Description("ID of the port forward rule"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			forwardID := req.GetString("forward_id", "")
			if forwardID == "" {
				return mcp.NewToolResultError("forward_id is required"), nil
			}

			forward, err := firewallMgr.GetPortForwardDetails(ctx, forwardID)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to get port forward: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"forward": forward,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_toggle_port_forward
	s.AddTool(
		mcp.NewTool("unifi_toggle_port_forward",
			mcp.WithDescription("Enable or disable a port forward rule"),
			mcp.WithString("forward_id",
				mcp.Required(),
				mcp.Description("ID of the port forward rule"),
			),
			mcp.WithBoolean("enabled",
				mcp.Required(),
				mcp.Description("Whether to enable (true) or disable (false) the rule"),
			),
			mcp.WithBoolean("confirm",
				mcp.Required(),
				mcp.Description("Must be true to execute the toggle"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			forwardID := req.GetString("forward_id", "")
			enabled := req.GetBool("enabled", false)
			confirm := req.GetBool("confirm", false)

			if forwardID == "" {
				return mcp.NewToolResultError("forward_id is required"), nil
			}

			action := "disable"
			if enabled {
				action = "enable"
			}

			if !confirm {
				return mcp.NewToolResultText(fmt.Sprintf("Preview: Would %s port forward %s. Set confirm=true to execute.", action, forwardID)), nil
			}

			if err := firewallMgr.TogglePortForward(ctx, forwardID, enabled); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to toggle port forward: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("Port forward %s has been %sd", forwardID, action),
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_list_traffic_routes
	s.AddTool(
		mcp.NewTool("unifi_list_traffic_routes",
			mcp.WithDescription("List all traffic routes (policy-based routing)"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			routes, err := firewallMgr.ListTrafficRoutes(ctx)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to list traffic routes: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"count":   len(routes),
				"routes":  routes,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_get_traffic_route_details
	s.AddTool(
		mcp.NewTool("unifi_get_traffic_route_details",
			mcp.WithDescription("Get detailed information about a specific traffic route"),
			mcp.WithString("route_id",
				mcp.Required(),
				mcp.Description("ID of the traffic route"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			routeID := req.GetString("route_id", "")
			if routeID == "" {
				return mcp.NewToolResultError("route_id is required"), nil
			}

			route, err := firewallMgr.GetTrafficRouteDetails(ctx, routeID)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to get traffic route: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"route":   route,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_toggle_traffic_route
	s.AddTool(
		mcp.NewTool("unifi_toggle_traffic_route",
			mcp.WithDescription("Enable or disable a traffic route"),
			mcp.WithString("route_id",
				mcp.Required(),
				mcp.Description("ID of the traffic route"),
			),
			mcp.WithBoolean("enabled",
				mcp.Required(),
				mcp.Description("Whether to enable (true) or disable (false) the route"),
			),
			mcp.WithBoolean("confirm",
				mcp.Required(),
				mcp.Description("Must be true to execute the toggle"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			routeID := req.GetString("route_id", "")
			enabled := req.GetBool("enabled", false)
			confirm := req.GetBool("confirm", false)

			if routeID == "" {
				return mcp.NewToolResultError("route_id is required"), nil
			}

			action := "disable"
			if enabled {
				action = "enable"
			}

			if !confirm {
				return mcp.NewToolResultText(fmt.Sprintf("Preview: Would %s traffic route %s. Set confirm=true to execute.", action, routeID)), nil
			}

			if err := firewallMgr.ToggleTrafficRoute(ctx, routeID, enabled); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to toggle traffic route: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("Traffic route %s has been %sd", routeID, action),
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)
}
