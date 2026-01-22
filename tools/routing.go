package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirkirby/go-mcp-ubiquiti/managers"
)

// RegisterRoutingTools registers all routing-related tools.
func RegisterRoutingTools(s *server.MCPServer, routingMgr *managers.RoutingManager) {
	// unifi_list_static_routes
	s.AddTool(
		mcp.NewTool("unifi_list_static_routes",
			mcp.WithDescription("List all static routes"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			routes, err := routingMgr.ListStaticRoutes(ctx)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to list static routes: %v", err)), nil
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

	// unifi_get_static_route_details
	s.AddTool(
		mcp.NewTool("unifi_get_static_route_details",
			mcp.WithDescription("Get detailed information about a specific static route"),
			mcp.WithString("route_id",
				mcp.Required(),
				mcp.Description("ID of the static route"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			routeID := req.GetString("route_id", "")
			if routeID == "" {
				return mcp.NewToolResultError("route_id is required"), nil
			}

			route, err := routingMgr.GetStaticRouteDetails(ctx, routeID)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to get static route: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"route":   route,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_create_static_route
	s.AddTool(
		mcp.NewTool("unifi_create_static_route",
			mcp.WithDescription("Create a new static route"),
			mcp.WithString("name",
				mcp.Description("Name for the static route"),
			),
			mcp.WithString("network",
				mcp.Required(),
				mcp.Description("Destination network (e.g., 10.0.0.0/24)"),
			),
			mcp.WithString("next_hop",
				mcp.Description("Next hop IP address"),
			),
			mcp.WithString("interface",
				mcp.Description("Interface name for interface routes"),
			),
			mcp.WithString("type",
				mcp.Description("Route type: nexthop-route, interface-route, blackhole"),
			),
			mcp.WithNumber("distance",
				mcp.Description("Administrative distance (default: 1)"),
			),
			mcp.WithBoolean("confirm",
				mcp.Required(),
				mcp.Description("Must be true to execute creation"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			network := req.GetString("network", "")
			confirm := req.GetBool("confirm", false)

			if network == "" {
				return mcp.NewToolResultError("network is required"), nil
			}
			if !confirm {
				return mcp.NewToolResultText(fmt.Sprintf("Preview: Would create static route to %s. Set confirm=true to execute.", network)), nil
			}

			route := &managers.StaticRoute{
				Network: network,
				Enabled: true,
			}

			if name := req.GetString("name", ""); name != "" {
				route.Name = name
			}
			if nextHop := req.GetString("next_hop", ""); nextHop != "" {
				route.NextHop = nextHop
			}
			if iface := req.GetString("interface", ""); iface != "" {
				route.Interface = iface
			}
			if routeType := req.GetString("type", ""); routeType != "" {
				route.Type = routeType
			}
			distance := req.GetInt("distance", 0)
			if distance > 0 {
				route.Distance = distance
			}

			created, err := routingMgr.CreateStaticRoute(ctx, route)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to create static route: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("Static route to %s has been created", network),
				"route":   created,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_update_static_route
	s.AddTool(
		mcp.NewTool("unifi_update_static_route",
			mcp.WithDescription("Update an existing static route"),
			mcp.WithString("route_id",
				mcp.Required(),
				mcp.Description("ID of the static route to update"),
			),
			mcp.WithString("name",
				mcp.Description("New name for the route"),
			),
			mcp.WithBoolean("enabled",
				mcp.Description("Enable or disable the route"),
			),
			mcp.WithBoolean("confirm",
				mcp.Required(),
				mcp.Description("Must be true to execute the update"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			routeID := req.GetString("route_id", "")
			confirm := req.GetBool("confirm", false)

			if routeID == "" {
				return mcp.NewToolResultError("route_id is required"), nil
			}
			if !confirm {
				return mcp.NewToolResultText(fmt.Sprintf("Preview: Would update static route %s. Set confirm=true to execute.", routeID)), nil
			}

			updates := make(map[string]interface{})
			if name := req.GetString("name", ""); name != "" {
				updates["name"] = name
			}
			// Check for enabled param explicitly
			args := req.GetArguments()
			if enabledVal, ok := args["enabled"]; ok {
				if enabled, ok := enabledVal.(bool); ok {
					updates["enabled"] = enabled
				}
			}

			if len(updates) == 0 {
				return mcp.NewToolResultError("no updates specified"), nil
			}

			if err := routingMgr.UpdateStaticRoute(ctx, routeID, updates); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to update static route: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("Static route %s has been updated", routeID),
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_list_active_routes
	s.AddTool(
		mcp.NewTool("unifi_list_active_routes",
			mcp.WithDescription("List active routes from the routing table"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			routes, err := routingMgr.ListActiveRoutes(ctx)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to list active routes: %v", err)), nil
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
}
