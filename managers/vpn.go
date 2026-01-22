package managers

import (
	"context"
	"encoding/json"
	"fmt"
)

// VPNClient represents a VPN client configuration.
type VPNClient struct {
	ID              string `json:"_id,omitempty"`
	Name            string `json:"name"`
	Enabled         bool   `json:"enabled"`
	SiteID          string `json:"site_id,omitempty"`
	VPNType         string `json:"vpn_type,omitempty"` // wireguard, openvpn, ipsec
	Interface       string `json:"interface,omitempty"`
	ServerAddress   string `json:"server_address,omitempty"`
	Username        string `json:"username,omitempty"`
	Password        string `json:"x_password,omitempty"`
	PreSharedKey    string `json:"x_psk,omitempty"`
	PublicKey       string `json:"public_key,omitempty"`
	PrivateKey      string `json:"x_private_key,omitempty"`
	LocalAddress    string `json:"local_address,omitempty"`
	RemoteAddress   string `json:"remote_address,omitempty"`
	AllowedIPs      string `json:"allowed_ips,omitempty"`
	PersistentKeepalive int `json:"persistent_keepalive,omitempty"`
	DNSServers      string `json:"dns_servers,omitempty"`
	RouteDistance   int    `json:"route_distance,omitempty"`
	MTU             int    `json:"mtu,omitempty"`
	Status          string `json:"status,omitempty"`
}

// VPNServer represents a VPN server configuration.
type VPNServer struct {
	ID               string   `json:"_id,omitempty"`
	Name             string   `json:"name"`
	Enabled          bool     `json:"enabled"`
	SiteID           string   `json:"site_id,omitempty"`
	VPNType          string   `json:"vpn_type,omitempty"` // wireguard, openvpn, l2tp, pptp
	Interface        string   `json:"interface,omitempty"`
	Port             int      `json:"port,omitempty"`
	Protocol         string   `json:"protocol,omitempty"` // udp, tcp
	NetworkID        string   `json:"network_id,omitempty"`
	Subnet           string   `json:"subnet,omitempty"`
	PublicKey        string   `json:"public_key,omitempty"`
	PrivateKey       string   `json:"x_private_key,omitempty"`
	PreSharedKey     string   `json:"x_psk,omitempty"`
	DNSServers       []string `json:"dns_servers,omitempty"`
	AuthType         string   `json:"auth_type,omitempty"`
	RADIUSProfileID  string   `json:"radiusprofile_id,omitempty"`
	RequireMSCHAPv2  bool     `json:"require_mschapv2"`
	RequireStrongCrypto bool  `json:"require_strong_crypto"`
	MTU              int      `json:"mtu,omitempty"`
	ConnectedClients int      `json:"connected_clients,omitempty"`
	Status           string   `json:"status,omitempty"`
}

// VPNManager handles VPN operations.
type VPNManager struct {
	conn *ConnectionManager
}

// NewVPNManager creates a new VPN manager.
func NewVPNManager(conn *ConnectionManager) *VPNManager {
	return &VPNManager{conn: conn}
}

// ListVPNClients returns all VPN client configurations.
func (vm *VPNManager) ListVPNClients(ctx context.Context) ([]VPNClient, error) {
	var endpoint string
	if vm.conn.IsUnifiOS() {
		endpoint = "/v2/api/site/" + vm.conn.Site() + "/vpn/clients"
	} else {
		endpoint = vm.conn.GetSitePath("/rest/vpnclient")
	}

	data, err := vm.conn.Request(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get VPN clients: %w", err)
	}

	var clients []VPNClient
	if err := json.Unmarshal(data, &clients); err != nil {
		return nil, fmt.Errorf("failed to parse VPN clients: %w", err)
	}

	return clients, nil
}

// GetVPNClientDetails returns details for a specific VPN client.
func (vm *VPNManager) GetVPNClientDetails(ctx context.Context, clientID string) (*VPNClient, error) {
	clients, err := vm.ListVPNClients(ctx)
	if err != nil {
		return nil, err
	}

	for _, c := range clients {
		if c.ID == clientID {
			return &c, nil
		}
	}

	return nil, fmt.Errorf("VPN client not found: %s", clientID)
}

// UpdateVPNClientState enables or disables a VPN client.
func (vm *VPNManager) UpdateVPNClientState(ctx context.Context, clientID string, enabled bool) error {
	var endpoint string
	if vm.conn.IsUnifiOS() {
		endpoint = "/v2/api/site/" + vm.conn.Site() + "/vpn/clients/" + clientID
	} else {
		endpoint = vm.conn.GetSitePath("/rest/vpnclient/" + clientID)
	}

	update := map[string]interface{}{
		"enabled": enabled,
	}

	_, err := vm.conn.Request(ctx, "PUT", endpoint, update)
	if err != nil {
		return fmt.Errorf("failed to update VPN client state: %w", err)
	}

	return nil
}

// ListVPNServers returns all VPN server configurations.
func (vm *VPNManager) ListVPNServers(ctx context.Context) ([]VPNServer, error) {
	var endpoint string
	if vm.conn.IsUnifiOS() {
		endpoint = "/v2/api/site/" + vm.conn.Site() + "/vpn/servers"
	} else {
		endpoint = vm.conn.GetSitePath("/rest/vpnserver")
	}

	data, err := vm.conn.Request(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get VPN servers: %w", err)
	}

	var servers []VPNServer
	if err := json.Unmarshal(data, &servers); err != nil {
		return nil, fmt.Errorf("failed to parse VPN servers: %w", err)
	}

	return servers, nil
}

// GetVPNServerDetails returns details for a specific VPN server.
func (vm *VPNManager) GetVPNServerDetails(ctx context.Context, serverID string) (*VPNServer, error) {
	servers, err := vm.ListVPNServers(ctx)
	if err != nil {
		return nil, err
	}

	for _, s := range servers {
		if s.ID == serverID {
			return &s, nil
		}
	}

	return nil, fmt.Errorf("VPN server not found: %s", serverID)
}

// UpdateVPNServerState enables or disables a VPN server.
func (vm *VPNManager) UpdateVPNServerState(ctx context.Context, serverID string, enabled bool) error {
	var endpoint string
	if vm.conn.IsUnifiOS() {
		endpoint = "/v2/api/site/" + vm.conn.Site() + "/vpn/servers/" + serverID
	} else {
		endpoint = vm.conn.GetSitePath("/rest/vpnserver/" + serverID)
	}

	update := map[string]interface{}{
		"enabled": enabled,
	}

	_, err := vm.conn.Request(ctx, "PUT", endpoint, update)
	if err != nil {
		return fmt.Errorf("failed to update VPN server state: %w", err)
	}

	return nil
}
