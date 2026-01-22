package managers

import (
	"context"
	"encoding/json"
	"fmt"
)

// StaticRoute represents a static route configuration.
type StaticRoute struct {
	ID          string `json:"_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Enabled     bool   `json:"enabled"`
	SiteID      string `json:"site_id,omitempty"`
	Type        string `json:"type,omitempty"` // nexthop-route, interface-route, blackhole
	Network     string `json:"static-route_network,omitempty"`
	Distance    int    `json:"static-route_distance,omitempty"`
	NextHop     string `json:"static-route_nexthop,omitempty"`
	Interface   string `json:"static-route_interface,omitempty"`
	Gateway     string `json:"gateway_device,omitempty"`
	GatewayType string `json:"gateway_type,omitempty"`
}

// ActiveRoute represents an active route in the routing table.
type ActiveRoute struct {
	Protocol    string `json:"protocol,omitempty"`
	Network     string `json:"network,omitempty"`
	NextHop     string `json:"nexthop,omitempty"`
	Interface   string `json:"interface,omitempty"`
	Metric      int    `json:"metric,omitempty"`
	Distance    int    `json:"distance,omitempty"`
	Uptime      int64  `json:"uptime,omitempty"`
	Selected    bool   `json:"selected"`
	FIB         bool   `json:"fib"`
}

// RoutingManager handles routing operations.
type RoutingManager struct {
	conn *ConnectionManager
}

// NewRoutingManager creates a new routing manager.
func NewRoutingManager(conn *ConnectionManager) *RoutingManager {
	return &RoutingManager{conn: conn}
}

// ListStaticRoutes returns all static routes.
func (rm *RoutingManager) ListStaticRoutes(ctx context.Context) ([]StaticRoute, error) {
	endpoint := rm.conn.GetSitePath("/rest/routing")

	data, err := rm.conn.Request(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get static routes: %w", err)
	}

	var routes []StaticRoute
	if err := json.Unmarshal(data, &routes); err != nil {
		return nil, fmt.Errorf("failed to parse static routes: %w", err)
	}

	return routes, nil
}

// GetStaticRouteDetails returns details for a specific static route.
func (rm *RoutingManager) GetStaticRouteDetails(ctx context.Context, routeID string) (*StaticRoute, error) {
	routes, err := rm.ListStaticRoutes(ctx)
	if err != nil {
		return nil, err
	}

	for _, r := range routes {
		if r.ID == routeID {
			return &r, nil
		}
	}

	return nil, fmt.Errorf("static route not found: %s", routeID)
}

// CreateStaticRoute creates a new static route.
func (rm *RoutingManager) CreateStaticRoute(ctx context.Context, route *StaticRoute) (*StaticRoute, error) {
	endpoint := rm.conn.GetSitePath("/rest/routing")

	data, err := rm.conn.Request(ctx, "POST", endpoint, route)
	if err != nil {
		return nil, fmt.Errorf("failed to create static route: %w", err)
	}

	var routes []StaticRoute
	if err := json.Unmarshal(data, &routes); err != nil {
		return nil, fmt.Errorf("failed to parse created route: %w", err)
	}

	if len(routes) == 0 {
		return nil, fmt.Errorf("no route returned after creation")
	}

	return &routes[0], nil
}

// UpdateStaticRoute updates an existing static route.
func (rm *RoutingManager) UpdateStaticRoute(ctx context.Context, routeID string, updates map[string]interface{}) error {
	endpoint := rm.conn.GetSitePath("/rest/routing/" + routeID)

	_, err := rm.conn.Request(ctx, "PUT", endpoint, updates)
	if err != nil {
		return fmt.Errorf("failed to update static route: %w", err)
	}

	return nil
}

// DeleteStaticRoute deletes a static route.
func (rm *RoutingManager) DeleteStaticRoute(ctx context.Context, routeID string) error {
	endpoint := rm.conn.GetSitePath("/rest/routing/" + routeID)

	_, err := rm.conn.Request(ctx, "DELETE", endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to delete static route: %w", err)
	}

	return nil
}

// ListActiveRoutes returns the active routing table.
func (rm *RoutingManager) ListActiveRoutes(ctx context.Context) ([]ActiveRoute, error) {
	// This requires a command to get the active routes
	cmd := map[string]interface{}{
		"cmd": "show-route",
	}

	data, err := rm.conn.Request(ctx, "POST", rm.conn.GetSitePath("/cmd/stat"), cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to get active routes: %w", err)
	}

	var routes []ActiveRoute
	if err := json.Unmarshal(data, &routes); err != nil {
		// The response format varies by controller version
		return nil, fmt.Errorf("failed to parse active routes: %w", err)
	}

	return routes, nil
}
