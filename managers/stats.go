package managers

import (
	"context"
	"encoding/json"
	"fmt"
)

// NetworkStats holds network statistics.
type NetworkStats struct {
	Site             string  `json:"site,omitempty"`
	OUI              string  `json:"o,omitempty"`
	Time             int64   `json:"time,omitempty"`
	DateTime         string  `json:"datetime,omitempty"`
	NumSta           int     `json:"num_sta,omitempty"`
	UserNumSta       int     `json:"user-num_sta,omitempty"`
	GuestNumSta      int     `json:"guest-num_sta,omitempty"`
	NumAP            int     `json:"num_ap,omitempty"`
	NumAdoptedAP     int     `json:"num_adopted,omitempty"`
	NumDisabledAP    int     `json:"num_disabled,omitempty"`
	NumDisconnectedAP int    `json:"num_disconnected,omitempty"`
	NumPendingAP     int     `json:"num_pending,omitempty"`
	WLANRX           int64   `json:"wlan-rx_bytes,omitempty"`
	WLANTX           int64   `json:"wlan-tx_bytes,omitempty"`
	WLANNumSta       int     `json:"wlan-num_sta,omitempty"`
	LANRX            int64   `json:"lan-rx_bytes,omitempty"`
	LANTX            int64   `json:"lan-tx_bytes,omitempty"`
	LANNumSta        int     `json:"lan-num_sta,omitempty"`
	WANRX            int64   `json:"wan-rx_bytes,omitempty"`
	WANTX            int64   `json:"wan-tx_bytes,omitempty"`
	Latency          float64 `json:"latency,omitempty"`
	XputUp           float64 `json:"xput_up,omitempty"`
	XputDown         float64 `json:"xput_down,omitempty"`
}

// ClientStats holds client statistics.
type ClientStats struct {
	MAC        string `json:"mac,omitempty"`
	SiteID     string `json:"site_id,omitempty"`
	OUI        string `json:"o,omitempty"`
	Time       int64  `json:"time,omitempty"`
	RXBytes    int64  `json:"rx_bytes,omitempty"`
	TXBytes    int64  `json:"tx_bytes,omitempty"`
	RXPackets  int64  `json:"rx_packets,omitempty"`
	TXPackets  int64  `json:"tx_packets,omitempty"`
	TXRetries  int64  `json:"tx_retries,omitempty"`
	WifiTXAttempts int64 `json:"wifi_tx_attempts,omitempty"`
	RXRate     int64  `json:"rx_rate,omitempty"`
	TXRate     int64  `json:"tx_rate,omitempty"`
	Signal     int    `json:"signal,omitempty"`
	RSSI       int    `json:"rssi,omitempty"`
	Noise      int    `json:"noise,omitempty"`
}

// DeviceStats holds device statistics.
type DeviceStats struct {
	OUI       string `json:"o,omitempty"`
	SiteID    string `json:"site_id,omitempty"`
	Time      int64  `json:"time,omitempty"`
	DateTime  string `json:"datetime,omitempty"`
	MAC       string `json:"mac,omitempty"`
	RXBytes   int64  `json:"rx_bytes,omitempty"`
	TXBytes   int64  `json:"tx_bytes,omitempty"`
	RXPackets int64  `json:"rx_packets,omitempty"`
	TXPackets int64  `json:"tx_packets,omitempty"`
	RXDropped int64  `json:"rx_dropped,omitempty"`
	TXDropped int64  `json:"tx_dropped,omitempty"`
	RXErrors  int64  `json:"rx_errors,omitempty"`
	TXErrors  int64  `json:"tx_errors,omitempty"`
	NumSta    int    `json:"num_sta,omitempty"`
	UserNumSta int   `json:"user-num_sta,omitempty"`
	GuestNumSta int  `json:"guest-num_sta,omitempty"`
	CPUUsage  float64 `json:"cpu,omitempty"`
	MemUsage  float64 `json:"mem,omitempty"`
	Loadavg1  float64 `json:"loadavg_1,omitempty"`
	Loadavg5  float64 `json:"loadavg_5,omitempty"`
	Loadavg15 float64 `json:"loadavg_15,omitempty"`
}

// DPIStats holds Deep Packet Inspection statistics.
type DPIStats struct {
	Cat       int    `json:"cat,omitempty"`
	App       int    `json:"app,omitempty"`
	CatName   string `json:"catname,omitempty"`
	AppName   string `json:"appname,omitempty"`
	RXBytes   int64  `json:"rx_bytes,omitempty"`
	TXBytes   int64  `json:"tx_bytes,omitempty"`
	RXPackets int64  `json:"rx_packets,omitempty"`
	TXPackets int64  `json:"tx_packets,omitempty"`
}

// TopClient represents a top client by usage.
type TopClient struct {
	MAC       string `json:"mac,omitempty"`
	RXBytes   int64  `json:"rx_bytes,omitempty"`
	TXBytes   int64  `json:"tx_bytes,omitempty"`
	TotalBytes int64 `json:"total_bytes,omitempty"`
}

// StatsManager handles statistics operations.
type StatsManager struct {
	conn *ConnectionManager
}

// NewStatsManager creates a new stats manager.
func NewStatsManager(conn *ConnectionManager) *StatsManager {
	return &StatsManager{conn: conn}
}

// GetNetworkStats returns network statistics.
func (sm *StatsManager) GetNetworkStats(ctx context.Context) (*NetworkStats, error) {
	// Get stats from site stats endpoint
	data, err := sm.conn.Request(ctx, "GET", sm.conn.GetSitePath("/stat/sites"), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get network stats: %w", err)
	}

	var stats []NetworkStats
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, fmt.Errorf("failed to parse network stats: %w", err)
	}

	if len(stats) == 0 {
		return nil, fmt.Errorf("no stats returned")
	}

	return &stats[0], nil
}

// GetClientStats returns statistics for a specific client.
func (sm *StatsManager) GetClientStats(ctx context.Context, mac string) (*ClientStats, error) {
	endpoint := sm.conn.GetSitePath("/stat/user/" + mac)

	data, err := sm.conn.Request(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get client stats: %w", err)
	}

	var stats []ClientStats
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, fmt.Errorf("failed to parse client stats: %w", err)
	}

	if len(stats) == 0 {
		return nil, fmt.Errorf("no stats for client: %s", mac)
	}

	return &stats[0], nil
}

// GetDeviceStats returns statistics for a specific device.
func (sm *StatsManager) GetDeviceStats(ctx context.Context, mac string) (*DeviceStats, error) {
	// Get device basic stats
	data, err := sm.conn.Request(ctx, "GET", sm.conn.GetSitePath("/stat/device/"+mac), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get device stats: %w", err)
	}

	var stats []DeviceStats
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, fmt.Errorf("failed to parse device stats: %w", err)
	}

	if len(stats) == 0 {
		return nil, fmt.Errorf("no stats for device: %s", mac)
	}

	return &stats[0], nil
}

// GetTopClients returns top clients by usage.
func (sm *StatsManager) GetTopClients(ctx context.Context, limit int) ([]TopClient, error) {
	if limit <= 0 {
		limit = 10
	}

	// Get all clients with stats
	data, err := sm.conn.Request(ctx, "GET", sm.conn.GetSitePath("/stat/sta"), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get clients: %w", err)
	}

	var clients []Client
	if err := json.Unmarshal(data, &clients); err != nil {
		return nil, fmt.Errorf("failed to parse clients: %w", err)
	}

	// Convert to TopClient and calculate total
	var topClients []TopClient
	for _, c := range clients {
		topClients = append(topClients, TopClient{
			MAC:        c.MAC,
			RXBytes:    c.RxBytes,
			TXBytes:    c.TxBytes,
			TotalBytes: c.RxBytes + c.TxBytes,
		})
	}

	// Sort by total bytes (simple bubble sort for small lists)
	for i := 0; i < len(topClients)-1; i++ {
		for j := 0; j < len(topClients)-i-1; j++ {
			if topClients[j].TotalBytes < topClients[j+1].TotalBytes {
				topClients[j], topClients[j+1] = topClients[j+1], topClients[j]
			}
		}
	}

	if len(topClients) > limit {
		topClients = topClients[:limit]
	}

	return topClients, nil
}

// GetDPIStats returns Deep Packet Inspection statistics.
func (sm *StatsManager) GetDPIStats(ctx context.Context) ([]DPIStats, error) {
	data, err := sm.conn.Request(ctx, "GET", sm.conn.GetSitePath("/stat/sitedpi"), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get DPI stats: %w", err)
	}

	var stats []DPIStats
	if err := json.Unmarshal(data, &stats); err != nil {
		// Try parsing as wrapper
		var wrapper struct {
			ByApp []DPIStats `json:"by_app,omitempty"`
			ByCat []DPIStats `json:"by_cat,omitempty"`
		}
		if err := json.Unmarshal(data, &wrapper); err != nil {
			return nil, fmt.Errorf("failed to parse DPI stats: %w", err)
		}
		return wrapper.ByApp, nil
	}

	return stats, nil
}
