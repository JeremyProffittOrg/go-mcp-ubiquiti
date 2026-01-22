package managers

import (
	"context"
	"encoding/json"
	"fmt"
)

// SystemInfo holds UniFi controller system information.
type SystemInfo struct {
	Version           string `json:"version,omitempty"`
	Build             string `json:"build,omitempty"`
	LocalHostname     string `json:"hostname,omitempty"`
	Name              string `json:"name,omitempty"`
	IPAddress         string `json:"ip_address,omitempty"`
	Uptime            int64  `json:"uptime,omitempty"`
	IsCloudConsole    bool   `json:"is_cloud_console"`
	IsCloudEnabled    bool   `json:"cloud_enabled"`
	Timezone          string `json:"timezone,omitempty"`
	Autobackup        bool   `json:"autobackup"`
	UpdateAvailable   bool   `json:"update_available"`
	PreviousVersion   string `json:"previous_version,omitempty"`
	LiveChat          string `json:"live_chat,omitempty"`
	StoreEnabled      string `json:"store_enabled,omitempty"`
	DataRetentionDays int    `json:"data_retention_time_in_hours_for_5minutes_scale,omitempty"`
	DeviceCount       int    `json:"device_count,omitempty"`
}

// SiteSettings holds site configuration.
type SiteSettings struct {
	ID                          string `json:"_id,omitempty"`
	SiteID                      string `json:"site_id,omitempty"`
	Key                         string `json:"key,omitempty"`
	CountryCode                 int    `json:"country_code,omitempty"`
	Timezone                    string `json:"timezone,omitempty"`
	AutoUpgrade                 bool   `json:"auto_upgrade"`
	SpeedTestEnabled            bool   `json:"speed_test_enabled"`
	LEDEnabled                  bool   `json:"led_enabled"`
	AlertEnabled                bool   `json:"alert_enabled"`
	UplinkType                  string `json:"uplink_type,omitempty"`
	ConnectivityMonitorEnabled  bool   `json:"connectivity_monitor_enabled"`
	ConnectivityMonitorInterval int    `json:"connectivity_monitor_interval,omitempty"`
}

// NetworkHealth holds network health information.
type NetworkHealth struct {
	Subsystem string `json:"subsystem"`
	Status    string `json:"status"`
	NumAP     int    `json:"num_ap,omitempty"`
	NumAdopted int   `json:"num_adopted,omitempty"`
	NumDisabled int  `json:"num_disabled,omitempty"`
	NumDisconnected int `json:"num_disconnected,omitempty"`
	NumPending int   `json:"num_pending,omitempty"`
	NumUser    int   `json:"num_user,omitempty"`
	NumGuest   int   `json:"num_guest,omitempty"`
	NumIOT     int   `json:"num_iot,omitempty"`
	TxBytesR   int64 `json:"tx_bytes-r,omitempty"`
	RxBytesR   int64 `json:"rx_bytes-r,omitempty"`
	NumSW      int   `json:"num_sw,omitempty"`
	NumGW      int   `json:"num_gw,omitempty"`
	NumSTA     int   `json:"num_sta,omitempty"`
	WANIP      string `json:"wan_ip,omitempty"`
	GatewaysMAC []string `json:"gateways,omitempty"`
	NameServers []string `json:"nameservers,omitempty"`
	NETmask    string `json:"netmask,omitempty"`
	ISPName    string `json:"isp_name,omitempty"`
	ISPOrg     string `json:"isp_organization,omitempty"`
	UptimeStats map[string]interface{} `json:"uptime_stats,omitempty"`
	Latency     int   `json:"latency,omitempty"`
	SpeedtestLastRun int64 `json:"speedtest_lastrun,omitempty"`
	SpeedtestStatus string `json:"speedtest_status,omitempty"`
	XputDown    float64 `json:"xput_down,omitempty"`
	XputUp      float64 `json:"xput_up,omitempty"`
}

// SystemManager handles system operations.
type SystemManager struct {
	conn *ConnectionManager
}

// NewSystemManager creates a new system manager.
func NewSystemManager(conn *ConnectionManager) *SystemManager {
	return &SystemManager{conn: conn}
}

// GetSystemInfo returns controller system information.
func (sm *SystemManager) GetSystemInfo(ctx context.Context) (*SystemInfo, error) {
	data, err := sm.conn.Request(ctx, "GET", sm.conn.GetSitePath("/stat/sysinfo"), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get system info: %w", err)
	}

	var infos []SystemInfo
	if err := json.Unmarshal(data, &infos); err != nil {
		// Try parsing as a single object
		var info SystemInfo
		if err := json.Unmarshal(data, &info); err != nil {
			return nil, fmt.Errorf("failed to parse system info: %w", err)
		}
		return &info, nil
	}

	if len(infos) == 0 {
		return nil, fmt.Errorf("no system info returned")
	}

	return &infos[0], nil
}

// GetSiteSettings returns site settings.
func (sm *SystemManager) GetSiteSettings(ctx context.Context) ([]SiteSettings, error) {
	data, err := sm.conn.Request(ctx, "GET", sm.conn.GetSitePath("/rest/setting"), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get site settings: %w", err)
	}

	var settings []SiteSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, fmt.Errorf("failed to parse site settings: %w", err)
	}

	return settings, nil
}

// GetNetworkHealth returns network health information.
func (sm *SystemManager) GetNetworkHealth(ctx context.Context) ([]NetworkHealth, error) {
	data, err := sm.conn.Request(ctx, "GET", sm.conn.GetSitePath("/stat/health"), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get network health: %w", err)
	}

	var health []NetworkHealth
	if err := json.Unmarshal(data, &health); err != nil {
		return nil, fmt.Errorf("failed to parse network health: %w", err)
	}

	return health, nil
}
