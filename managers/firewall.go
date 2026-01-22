package managers

import (
	"context"
	"encoding/json"
	"fmt"
)

// FirewallPolicy represents a firewall policy (v2 API).
type FirewallPolicy struct {
	ID          string   `json:"_id,omitempty"`
	Name        string   `json:"name"`
	Enabled     bool     `json:"enabled"`
	Index       int      `json:"index,omitempty"`
	RuleIndex   int      `json:"rule_index,omitempty"`
	Action      string   `json:"action,omitempty"` // ACCEPT, DROP, REJECT
	Protocol    string   `json:"protocol,omitempty"` // all, tcp, udp, icmp, tcp_udp
	Logging     bool     `json:"logging"`
	Description string   `json:"description,omitempty"`
	SrcZone     string   `json:"source_zone,omitempty"`
	DstZone     string   `json:"destination_zone,omitempty"`
	SrcAddress  string   `json:"source_address,omitempty"`
	DstAddress  string   `json:"destination_address,omitempty"`
	SrcNetworkType string `json:"source_network_type,omitempty"`
	DstNetworkType string `json:"destination_network_type,omitempty"`
	SrcNetworkID string  `json:"source_network_id,omitempty"`
	DstNetworkID string  `json:"destination_network_id,omitempty"`
	SrcIPGroupID string  `json:"source_ip_group_id,omitempty"`
	DstIPGroupID string  `json:"destination_ip_group_id,omitempty"`
	SrcPort     string   `json:"source_port,omitempty"`
	DstPort     string   `json:"destination_port,omitempty"`
	SrcMAC      string   `json:"source_mac,omitempty"`
	Schedule    string   `json:"schedule,omitempty"`
	Established bool     `json:"match_state_established"`
	Related     bool     `json:"match_state_related"`
	New         bool     `json:"match_state_new"`
	Invalid     bool     `json:"match_state_invalid"`
	IPSec       string   `json:"match_ipsec,omitempty"` // match_ipsec, match_none, NONE
	ICMPType    string   `json:"icmp_typename,omitempty"`
}

// FirewallZone represents a firewall zone.
type FirewallZone struct {
	ID          string   `json:"_id,omitempty"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	NetworkIDs  []string `json:"network_ids,omitempty"`
	IsDefault   bool     `json:"is_default"`
}

// IPGroup represents an IP address group.
type IPGroup struct {
	ID          string   `json:"_id,omitempty"`
	Name        string   `json:"name"`
	GroupType   string   `json:"group_type,omitempty"` // address-group, ipv6-address-group
	GroupMembers []string `json:"group_members,omitempty"`
	SiteID      string   `json:"site_id,omitempty"`
}

// PortForward represents a port forwarding rule.
type PortForward struct {
	ID          string `json:"_id,omitempty"`
	Name        string `json:"name"`
	Enabled     bool   `json:"enabled"`
	SiteID      string `json:"site_id,omitempty"`
	PFWDInterface string `json:"pfwd_interface,omitempty"` // wan, wan2
	Src         string `json:"src,omitempty"` // any, limited
	SrcNetworkID string `json:"src_network_conf_id,omitempty"`
	FwdPort     string `json:"fwd_port,omitempty"`
	FwdIP       string `json:"fwd,omitempty"`
	DstPort     string `json:"dst_port,omitempty"`
	Proto       string `json:"proto,omitempty"` // tcp, udp, tcp_udp
	Log         bool   `json:"log"`
}

// TrafficRoute represents a traffic route (policy-based routing).
type TrafficRoute struct {
	ID              string   `json:"_id,omitempty"`
	Name            string   `json:"name"`
	Enabled         bool     `json:"enabled"`
	Description     string   `json:"description,omitempty"`
	MatchingTarget  string   `json:"matching_target,omitempty"` // INTERNET, DOMAIN, IP, REGION
	Domains         []string `json:"domains,omitempty"`
	Regions         []string `json:"regions,omitempty"`
	IPAddresses     []string `json:"ip_addresses,omitempty"`
	IPRanges        []string `json:"ip_ranges,omitempty"`
	TargetDevices   []TargetDevice `json:"target_devices,omitempty"`
	NetworkID       string   `json:"network_id,omitempty"`
	KillSwitchEnabled bool   `json:"kill_switch_enabled"`
	Interface       string   `json:"interface,omitempty"` // wan, wan2, vpn-client-x
	FallbackInterface string `json:"fallback_interface,omitempty"`
}

// TargetDevice identifies a device for traffic routing.
type TargetDevice struct {
	ClientMAC string `json:"client_mac,omitempty"`
	NetworkID string `json:"network_id,omitempty"`
	Type      string `json:"type,omitempty"` // ALL_CLIENTS, CLIENT, NETWORK
}

// FirewallManager handles firewall operations.
type FirewallManager struct {
	conn *ConnectionManager
}

// NewFirewallManager creates a new firewall manager.
func NewFirewallManager(conn *ConnectionManager) *FirewallManager {
	return &FirewallManager{conn: conn}
}

// ListFirewallPolicies returns all firewall policies.
func (fm *FirewallManager) ListFirewallPolicies(ctx context.Context) ([]FirewallPolicy, error) {
	var endpoint string
	if fm.conn.IsUnifiOS() {
		endpoint = "/v2/api/site/" + fm.conn.Site() + "/firewall-policies"
	} else {
		endpoint = fm.conn.GetSitePath("/rest/firewallrule")
	}

	data, err := fm.conn.Request(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get firewall policies: %w", err)
	}

	var policies []FirewallPolicy
	if err := json.Unmarshal(data, &policies); err != nil {
		return nil, fmt.Errorf("failed to parse firewall policies: %w", err)
	}

	return policies, nil
}

// GetFirewallPolicyDetails returns details for a specific policy.
func (fm *FirewallManager) GetFirewallPolicyDetails(ctx context.Context, policyID string) (*FirewallPolicy, error) {
	policies, err := fm.ListFirewallPolicies(ctx)
	if err != nil {
		return nil, err
	}

	for _, p := range policies {
		if p.ID == policyID {
			return &p, nil
		}
	}

	return nil, fmt.Errorf("firewall policy not found: %s", policyID)
}

// ToggleFirewallPolicy enables or disables a firewall policy.
func (fm *FirewallManager) ToggleFirewallPolicy(ctx context.Context, policyID string, enabled bool) error {
	var endpoint string
	if fm.conn.IsUnifiOS() {
		endpoint = "/v2/api/site/" + fm.conn.Site() + "/firewall-policies/" + policyID
	} else {
		endpoint = fm.conn.GetSitePath("/rest/firewallrule/" + policyID)
	}

	update := map[string]interface{}{
		"enabled": enabled,
	}

	_, err := fm.conn.Request(ctx, "PUT", endpoint, update)
	if err != nil {
		return fmt.Errorf("failed to toggle firewall policy: %w", err)
	}

	return nil
}

// CreateFirewallPolicy creates a new firewall policy.
func (fm *FirewallManager) CreateFirewallPolicy(ctx context.Context, policy *FirewallPolicy) (*FirewallPolicy, error) {
	var endpoint string
	if fm.conn.IsUnifiOS() {
		endpoint = "/v2/api/site/" + fm.conn.Site() + "/firewall-policies"
	} else {
		endpoint = fm.conn.GetSitePath("/rest/firewallrule")
	}

	data, err := fm.conn.Request(ctx, "POST", endpoint, policy)
	if err != nil {
		return nil, fmt.Errorf("failed to create firewall policy: %w", err)
	}

	var policies []FirewallPolicy
	if err := json.Unmarshal(data, &policies); err != nil {
		return nil, fmt.Errorf("failed to parse created policy: %w", err)
	}

	if len(policies) == 0 {
		return nil, fmt.Errorf("no policy returned after creation")
	}

	return &policies[0], nil
}

// ListFirewallZones returns all firewall zones.
func (fm *FirewallManager) ListFirewallZones(ctx context.Context) ([]FirewallZone, error) {
	var endpoint string
	if fm.conn.IsUnifiOS() {
		endpoint = "/v2/api/site/" + fm.conn.Site() + "/firewall-zones"
	} else {
		// Standalone may not support zones
		return nil, fmt.Errorf("firewall zones not supported on standalone controller")
	}

	data, err := fm.conn.Request(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get firewall zones: %w", err)
	}

	var zones []FirewallZone
	if err := json.Unmarshal(data, &zones); err != nil {
		return nil, fmt.Errorf("failed to parse firewall zones: %w", err)
	}

	return zones, nil
}

// ListIPGroups returns all IP groups.
func (fm *FirewallManager) ListIPGroups(ctx context.Context) ([]IPGroup, error) {
	endpoint := fm.conn.GetSitePath("/rest/firewallgroup")

	data, err := fm.conn.Request(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get IP groups: %w", err)
	}

	var groups []IPGroup
	if err := json.Unmarshal(data, &groups); err != nil {
		return nil, fmt.Errorf("failed to parse IP groups: %w", err)
	}

	return groups, nil
}

// ListPortForwards returns all port forwarding rules.
func (fm *FirewallManager) ListPortForwards(ctx context.Context) ([]PortForward, error) {
	endpoint := fm.conn.GetSitePath("/rest/portforward")

	data, err := fm.conn.Request(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get port forwards: %w", err)
	}

	var forwards []PortForward
	if err := json.Unmarshal(data, &forwards); err != nil {
		return nil, fmt.Errorf("failed to parse port forwards: %w", err)
	}

	return forwards, nil
}

// GetPortForwardDetails returns details for a specific port forward.
func (fm *FirewallManager) GetPortForwardDetails(ctx context.Context, forwardID string) (*PortForward, error) {
	forwards, err := fm.ListPortForwards(ctx)
	if err != nil {
		return nil, err
	}

	for _, f := range forwards {
		if f.ID == forwardID {
			return &f, nil
		}
	}

	return nil, fmt.Errorf("port forward not found: %s", forwardID)
}

// TogglePortForward enables or disables a port forward.
func (fm *FirewallManager) TogglePortForward(ctx context.Context, forwardID string, enabled bool) error {
	endpoint := fm.conn.GetSitePath("/rest/portforward/" + forwardID)

	update := map[string]interface{}{
		"enabled": enabled,
	}

	_, err := fm.conn.Request(ctx, "PUT", endpoint, update)
	if err != nil {
		return fmt.Errorf("failed to toggle port forward: %w", err)
	}

	return nil
}

// CreatePortForward creates a new port forward.
func (fm *FirewallManager) CreatePortForward(ctx context.Context, forward *PortForward) (*PortForward, error) {
	endpoint := fm.conn.GetSitePath("/rest/portforward")

	data, err := fm.conn.Request(ctx, "POST", endpoint, forward)
	if err != nil {
		return nil, fmt.Errorf("failed to create port forward: %w", err)
	}

	var forwards []PortForward
	if err := json.Unmarshal(data, &forwards); err != nil {
		return nil, fmt.Errorf("failed to parse created port forward: %w", err)
	}

	if len(forwards) == 0 {
		return nil, fmt.Errorf("no port forward returned after creation")
	}

	return &forwards[0], nil
}

// ListTrafficRoutes returns all traffic routes.
func (fm *FirewallManager) ListTrafficRoutes(ctx context.Context) ([]TrafficRoute, error) {
	var endpoint string
	if fm.conn.IsUnifiOS() {
		endpoint = "/v2/api/site/" + fm.conn.Site() + "/traffic-rules"
	} else {
		endpoint = fm.conn.GetSitePath("/rest/trafficroute")
	}

	data, err := fm.conn.Request(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get traffic routes: %w", err)
	}

	var routes []TrafficRoute
	if err := json.Unmarshal(data, &routes); err != nil {
		return nil, fmt.Errorf("failed to parse traffic routes: %w", err)
	}

	return routes, nil
}

// GetTrafficRouteDetails returns details for a specific traffic route.
func (fm *FirewallManager) GetTrafficRouteDetails(ctx context.Context, routeID string) (*TrafficRoute, error) {
	routes, err := fm.ListTrafficRoutes(ctx)
	if err != nil {
		return nil, err
	}

	for _, r := range routes {
		if r.ID == routeID {
			return &r, nil
		}
	}

	return nil, fmt.Errorf("traffic route not found: %s", routeID)
}

// ToggleTrafficRoute enables or disables a traffic route.
func (fm *FirewallManager) ToggleTrafficRoute(ctx context.Context, routeID string, enabled bool) error {
	var endpoint string
	if fm.conn.IsUnifiOS() {
		endpoint = "/v2/api/site/" + fm.conn.Site() + "/traffic-rules/" + routeID
	} else {
		endpoint = fm.conn.GetSitePath("/rest/trafficroute/" + routeID)
	}

	update := map[string]interface{}{
		"enabled": enabled,
	}

	_, err := fm.conn.Request(ctx, "PUT", endpoint, update)
	if err != nil {
		return fmt.Errorf("failed to toggle traffic route: %w", err)
	}

	return nil
}
