package managers

import (
	"context"
	"encoding/json"
	"fmt"
)

// Network represents a UniFi network configuration.
type Network struct {
	ID                    string `json:"_id,omitempty"`
	Name                  string `json:"name"`
	Purpose               string `json:"purpose,omitempty"`
	Enabled               bool   `json:"enabled"`
	SiteID                string `json:"site_id,omitempty"`
	VLANEnabled           bool   `json:"vlan_enabled"`
	VLAN                  int    `json:"vlan,omitempty"`
	IPSubnet              string `json:"ip_subnet,omitempty"`
	NetworkGroup          string `json:"networkgroup,omitempty"`
	DHCPDEnabled          bool   `json:"dhcpd_enabled"`
	DHCPDStart            string `json:"dhcpd_start,omitempty"`
	DHCPDStop             string `json:"dhcpd_stop,omitempty"`
	DHCPDLeasetime        int    `json:"dhcpd_leasetime,omitempty"`
	DHCPDGatewayEnabled   bool   `json:"dhcpd_gateway_enabled"`
	DHCPDGateway          string `json:"dhcpd_gateway,omitempty"`
	DHCPDDNSEnabled       bool   `json:"dhcpd_dns_enabled"`
	DHCPDDNS1             string `json:"dhcpd_dns_1,omitempty"`
	DHCPDDNS2             string `json:"dhcpd_dns_2,omitempty"`
	DHCPDUnifiController  string `json:"dhcpd_unifi_controller,omitempty"`
	DomainName            string `json:"domain_name,omitempty"`
	IGMPSnooping          bool   `json:"igmp_snooping"`
	DHCPGuardEnabled      bool   `json:"dhcp_guard_enabled"`
	DHCPRelayEnabled      bool   `json:"dhcpd_relay_enabled"`
	InternetAccessEnabled bool   `json:"internet_access_enabled"`
	IntraNetworkEnabled   bool   `json:"intra_network_enabled"`
	NATOutboundEnabled    bool   `json:"nat_outbound_enabled"`
	AutoScaleEnabled      bool   `json:"auto_scale_enabled"`
	DHCPDBootEnabled      bool   `json:"dhcpd_boot_enabled"`
	DHCPDTFTPServer       string `json:"dhcpd_tftp_server,omitempty"`
	DHCPDBootFilename     string `json:"dhcpd_boot_filename,omitempty"`
	DHCPDWPADServer       string `json:"dhcpd_wpad_server,omitempty"`
	LTELANEnabled         bool   `json:"lte_lan_enabled"`
	SettingPreference     string `json:"setting_preference,omitempty"`
}

// WLAN represents a UniFi wireless network.
type WLAN struct {
	ID                   string   `json:"_id,omitempty"`
	Name                 string   `json:"name"`
	Enabled              bool     `json:"enabled"`
	SiteID               string   `json:"site_id,omitempty"`
	SSID                 string   `json:"x_passphrase,omitempty"` // Note: actual SSID is in 'name'
	Security             string   `json:"security,omitempty"`
	WPAMode              string   `json:"wpa_mode,omitempty"`
	WPAEnc               string   `json:"wpa_enc,omitempty"`
	Passphrase           string   `json:"x_passphrase,omitempty"`
	IsGuest              bool     `json:"is_guest"`
	NetworkID            string   `json:"networkconf_id,omitempty"`
	UserGroupID          string   `json:"usergroup_id,omitempty"`
	APGroupIDs           []string `json:"ap_group_ids,omitempty"`
	WLANGroupID          string   `json:"wlangroup_id,omitempty"`
	HiddenSSID           bool     `json:"hide_ssid"`
	NoMDNS               bool     `json:"no_mdns"`
	L2Isolation          bool     `json:"l2_isolation"`
	GroupRekey           int      `json:"group_rekey,omitempty"`
	DTIMMode             string   `json:"dtim_mode,omitempty"`
	DTIMNG               int      `json:"dtim_ng,omitempty"`
	DTIMNA               int      `json:"dtim_na,omitempty"`
	MinrateNGEnabled     bool     `json:"minrate_ng_enabled"`
	MinrateNGDataRateKbps int     `json:"minrate_ng_data_rate_kbps,omitempty"`
	MinrateNAEnabled     bool     `json:"minrate_na_enabled"`
	MinrateNADataRateKbps int     `json:"minrate_na_data_rate_kbps,omitempty"`
	MACFilterEnabled     bool     `json:"mac_filter_enabled"`
	MACFilterPolicy      string   `json:"mac_filter_policy,omitempty"`
	MACFilterList        []string `json:"mac_filter_list,omitempty"`
	RadiusMACAuthEnabled bool     `json:"radius_mac_auth_enabled"`
	ScheduleEnabled      bool     `json:"schedule_enabled"`
	Schedule             []string `json:"schedule,omitempty"`
	RADIUSProfiles       []string `json:"radiusprofile_id,omitempty"`
	FastRoamingEnabled   bool     `json:"fast_roaming_enabled"`
	PMFMode              string   `json:"pmf_mode,omitempty"`
	SaeAntiBruteForce    bool     `json:"sae_anti_brute_force_enabled"`
	SaeSync              int      `json:"sae_sync,omitempty"`
	BSSTrans             bool     `json:"bss_transition"`
	UAPSDEnabled         bool     `json:"uapsd_enabled"`
	ProxyARP             bool     `json:"proxy_arp"`
	OptimizeIoT          bool     `json:"optimize_iot_wifi_connectivity"`
}

// NetworkManager handles network operations.
type NetworkManager struct {
	conn *ConnectionManager
}

// NewNetworkManager creates a new network manager.
func NewNetworkManager(conn *ConnectionManager) *NetworkManager {
	return &NetworkManager{conn: conn}
}

// ListNetworks returns all networks.
func (nm *NetworkManager) ListNetworks(ctx context.Context) ([]Network, error) {
	data, err := nm.conn.Request(ctx, "GET", nm.conn.GetSitePath("/rest/networkconf"), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get networks: %w", err)
	}

	var networks []Network
	if err := json.Unmarshal(data, &networks); err != nil {
		return nil, fmt.Errorf("failed to parse networks: %w", err)
	}

	return networks, nil
}

// GetNetworkDetails returns detailed information for a specific network.
func (nm *NetworkManager) GetNetworkDetails(ctx context.Context, networkID string) (*Network, error) {
	data, err := nm.conn.Request(ctx, "GET", nm.conn.GetSitePath("/rest/networkconf/"+networkID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get network: %w", err)
	}

	var networks []Network
	if err := json.Unmarshal(data, &networks); err != nil {
		return nil, fmt.Errorf("failed to parse network: %w", err)
	}

	if len(networks) == 0 {
		return nil, fmt.Errorf("network not found: %s", networkID)
	}

	return &networks[0], nil
}

// CreateNetwork creates a new network.
func (nm *NetworkManager) CreateNetwork(ctx context.Context, network *Network) (*Network, error) {
	data, err := nm.conn.Request(ctx, "POST", nm.conn.GetSitePath("/rest/networkconf"), network)
	if err != nil {
		return nil, fmt.Errorf("failed to create network: %w", err)
	}

	var networks []Network
	if err := json.Unmarshal(data, &networks); err != nil {
		return nil, fmt.Errorf("failed to parse created network: %w", err)
	}

	if len(networks) == 0 {
		return nil, fmt.Errorf("no network returned after creation")
	}

	return &networks[0], nil
}

// UpdateNetwork updates an existing network.
func (nm *NetworkManager) UpdateNetwork(ctx context.Context, networkID string, updates map[string]interface{}) error {
	_, err := nm.conn.Request(ctx, "PUT", nm.conn.GetSitePath("/rest/networkconf/"+networkID), updates)
	if err != nil {
		return fmt.Errorf("failed to update network: %w", err)
	}

	return nil
}

// ListWLANs returns all wireless networks.
func (nm *NetworkManager) ListWLANs(ctx context.Context) ([]WLAN, error) {
	data, err := nm.conn.Request(ctx, "GET", nm.conn.GetSitePath("/rest/wlanconf"), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get WLANs: %w", err)
	}

	var wlans []WLAN
	if err := json.Unmarshal(data, &wlans); err != nil {
		return nil, fmt.Errorf("failed to parse WLANs: %w", err)
	}

	return wlans, nil
}

// GetWLANDetails returns detailed information for a specific WLAN.
func (nm *NetworkManager) GetWLANDetails(ctx context.Context, wlanID string) (*WLAN, error) {
	data, err := nm.conn.Request(ctx, "GET", nm.conn.GetSitePath("/rest/wlanconf/"+wlanID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get WLAN: %w", err)
	}

	var wlans []WLAN
	if err := json.Unmarshal(data, &wlans); err != nil {
		return nil, fmt.Errorf("failed to parse WLAN: %w", err)
	}

	if len(wlans) == 0 {
		return nil, fmt.Errorf("WLAN not found: %s", wlanID)
	}

	return &wlans[0], nil
}
