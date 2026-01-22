package managers

import (
	"context"
	"encoding/json"
	"fmt"
)

// Client represents a UniFi network client.
type Client struct {
	ID          string `json:"_id,omitempty"`
	MAC         string `json:"mac"`
	IP          string `json:"ip,omitempty"`
	Hostname    string `json:"hostname,omitempty"`
	Name        string `json:"name,omitempty"`
	OUI         string `json:"oui,omitempty"`
	IsWired     bool   `json:"is_wired"`
	IsGuest     bool   `json:"is_guest"`
	Blocked     bool   `json:"blocked"`
	Noted       bool   `json:"noted"`
	Network     string `json:"network,omitempty"`
	NetworkID   string `json:"network_id,omitempty"`
	ESSID       string `json:"essid,omitempty"`
	BSSID       string `json:"bssid,omitempty"`
	Channel     int    `json:"channel,omitempty"`
	Radio       string `json:"radio,omitempty"`
	RadioProto  string `json:"radio_proto,omitempty"`
	Signal      int    `json:"signal,omitempty"`
	RSSI        int    `json:"rssi,omitempty"`
	Noise       int    `json:"noise,omitempty"`
	Uptime      int64  `json:"uptime,omitempty"`
	TxBytes     int64  `json:"tx_bytes,omitempty"`
	RxBytes     int64  `json:"rx_bytes,omitempty"`
	TxPackets   int64  `json:"tx_packets,omitempty"`
	RxPackets   int64  `json:"rx_packets,omitempty"`
	TxRate      int64  `json:"tx_rate,omitempty"`
	RxRate      int64  `json:"rx_rate,omitempty"`
	Satisfaction int   `json:"satisfaction,omitempty"`
	FirstSeen   int64  `json:"first_seen,omitempty"`
	LastSeen    int64  `json:"last_seen,omitempty"`
	DeviceType  int    `json:"dev_id_override,omitempty"`
	FixedIP     string `json:"fixed_ip,omitempty"`
	UseFixedIP  bool   `json:"use_fixedip"`
	LocalDNS    string `json:"local_dns_record,omitempty"`
	SwitchMAC   string `json:"sw_mac,omitempty"`
	SwitchPort  int    `json:"sw_port,omitempty"`
	APMAC       string `json:"ap_mac,omitempty"`
	UserGroupID string `json:"usergroup_id,omitempty"`
}

// DisplayName returns the best available name for the client.
func (c *Client) DisplayName() string {
	if c.Name != "" {
		return c.Name
	}
	if c.Hostname != "" {
		return c.Hostname
	}
	return c.MAC
}

// ClientManager handles client operations.
type ClientManager struct {
	conn *ConnectionManager
}

// NewClientManager creates a new client manager.
func NewClientManager(conn *ConnectionManager) *ClientManager {
	return &ClientManager{conn: conn}
}

// ListClients returns all online clients.
func (cm *ClientManager) ListClients(ctx context.Context, includeOffline bool, limit int) ([]Client, error) {
	endpoint := cm.conn.GetSitePath("/stat/sta")
	if includeOffline {
		endpoint = cm.conn.GetSitePath("/rest/user")
	}

	data, err := cm.conn.Request(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get clients: %w", err)
	}

	var clients []Client
	if err := json.Unmarshal(data, &clients); err != nil {
		return nil, fmt.Errorf("failed to parse clients: %w", err)
	}

	if limit > 0 && len(clients) > limit {
		clients = clients[:limit]
	}

	return clients, nil
}

// GetClientDetails returns detailed information for a specific client.
func (cm *ClientManager) GetClientDetails(ctx context.Context, mac string) (*Client, error) {
	// First try to get from online clients
	data, err := cm.conn.Request(ctx, "GET", cm.conn.GetSitePath("/stat/sta"), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get clients: %w", err)
	}

	var clients []Client
	if err := json.Unmarshal(data, &clients); err != nil {
		return nil, fmt.Errorf("failed to parse clients: %w", err)
	}

	for _, c := range clients {
		if c.MAC == mac {
			return &c, nil
		}
	}

	// Try historical clients
	data, err = cm.conn.Request(ctx, "GET", cm.conn.GetSitePath("/rest/user"), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get all clients: %w", err)
	}

	if err := json.Unmarshal(data, &clients); err != nil {
		return nil, fmt.Errorf("failed to parse clients: %w", err)
	}

	for _, c := range clients {
		if c.MAC == mac {
			return &c, nil
		}
	}

	return nil, fmt.Errorf("client not found: %s", mac)
}

// GetBlockedClients returns all blocked clients.
func (cm *ClientManager) GetBlockedClients(ctx context.Context) ([]Client, error) {
	clients, err := cm.ListClients(ctx, true, 0)
	if err != nil {
		return nil, err
	}

	var blocked []Client
	for _, c := range clients {
		if c.Blocked {
			blocked = append(blocked, c)
		}
	}

	return blocked, nil
}

// BlockClient blocks a client by MAC address.
func (cm *ClientManager) BlockClient(ctx context.Context, mac string) error {
	cmd := map[string]interface{}{
		"cmd": "block-sta",
		"mac": mac,
	}

	_, err := cm.conn.Request(ctx, "POST", cm.conn.GetSitePath("/cmd/stamgr"), cmd)
	if err != nil {
		return fmt.Errorf("failed to block client: %w", err)
	}

	return nil
}

// UnblockClient unblocks a client by MAC address.
func (cm *ClientManager) UnblockClient(ctx context.Context, mac string) error {
	cmd := map[string]interface{}{
		"cmd": "unblock-sta",
		"mac": mac,
	}

	_, err := cm.conn.Request(ctx, "POST", cm.conn.GetSitePath("/cmd/stamgr"), cmd)
	if err != nil {
		return fmt.Errorf("failed to unblock client: %w", err)
	}

	return nil
}

// RenameClient updates the name of a client.
func (cm *ClientManager) RenameClient(ctx context.Context, mac, name string) error {
	// First get the client to get its ID
	client, err := cm.GetClientDetails(ctx, mac)
	if err != nil {
		return err
	}

	if client.ID == "" {
		return fmt.Errorf("client has no ID, cannot rename")
	}

	update := map[string]interface{}{
		"name": name,
	}

	_, err = cm.conn.Request(ctx, "PUT", cm.conn.GetSitePath("/rest/user/"+client.ID), update)
	if err != nil {
		return fmt.Errorf("failed to rename client: %w", err)
	}

	return nil
}

// ForceReconnect forces a client to reconnect.
func (cm *ClientManager) ForceReconnect(ctx context.Context, mac string) error {
	cmd := map[string]interface{}{
		"cmd": "kick-sta",
		"mac": mac,
	}

	_, err := cm.conn.Request(ctx, "POST", cm.conn.GetSitePath("/cmd/stamgr"), cmd)
	if err != nil {
		return fmt.Errorf("failed to kick client: %w", err)
	}

	return nil
}

// AuthorizeGuest authorizes a guest client.
func (cm *ClientManager) AuthorizeGuest(ctx context.Context, mac string, minutes int, upKbps, downKbps, mbytes int) error {
	cmd := map[string]interface{}{
		"cmd":     "authorize-guest",
		"mac":     mac,
		"minutes": minutes,
	}

	if upKbps > 0 {
		cmd["up"] = upKbps
	}
	if downKbps > 0 {
		cmd["down"] = downKbps
	}
	if mbytes > 0 {
		cmd["bytes"] = mbytes
	}

	_, err := cm.conn.Request(ctx, "POST", cm.conn.GetSitePath("/cmd/stamgr"), cmd)
	if err != nil {
		return fmt.Errorf("failed to authorize guest: %w", err)
	}

	return nil
}

// UnauthorizeGuest revokes guest authorization.
func (cm *ClientManager) UnauthorizeGuest(ctx context.Context, mac string) error {
	cmd := map[string]interface{}{
		"cmd": "unauthorize-guest",
		"mac": mac,
	}

	_, err := cm.conn.Request(ctx, "POST", cm.conn.GetSitePath("/cmd/stamgr"), cmd)
	if err != nil {
		return fmt.Errorf("failed to unauthorize guest: %w", err)
	}

	return nil
}
