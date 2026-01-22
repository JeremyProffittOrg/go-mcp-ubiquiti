package managers

import (
	"context"
	"encoding/json"
	"fmt"
)

// Device represents a UniFi device (AP, switch, gateway, etc.).
type Device struct {
	ID              string   `json:"_id,omitempty"`
	MAC             string   `json:"mac"`
	IP              string   `json:"ip,omitempty"`
	Name            string   `json:"name,omitempty"`
	Model           string   `json:"model,omitempty"`
	ModelName       string   `json:"model_name,omitempty"`
	Type            string   `json:"type"`
	Version         string   `json:"version,omitempty"`
	Serial          string   `json:"serial,omitempty"`
	Adopted         bool     `json:"adopted"`
	State           int      `json:"state"`
	Uptime          int64    `json:"uptime,omitempty"`
	LastSeen        int64    `json:"last_seen,omitempty"`
	Upgradable      bool     `json:"upgradable"`
	UpgradeToFW     string   `json:"upgrade_to_firmware,omitempty"`
	ConfigVersion   string   `json:"cfgversion,omitempty"`
	SiteID          string   `json:"site_id,omitempty"`
	LicenseState    string   `json:"license_state,omitempty"`
	InformURL       string   `json:"inform_url,omitempty"`
	InformIP        string   `json:"inform_ip,omitempty"`
	ConnectedAt     int64    `json:"connected_at,omitempty"`
	ProvisionedAt   int64    `json:"provisioned_at,omitempty"`
	Satisfaction    int      `json:"satisfaction,omitempty"`
	SystemStats     SysStats `json:"system-stats,omitempty"`
	TxBytes         int64    `json:"tx_bytes,omitempty"`
	RxBytes         int64    `json:"rx_bytes,omitempty"`
	NumSta          int      `json:"num_sta,omitempty"`
	UserNumSta      int      `json:"user-num_sta,omitempty"`
	GuestNumSta     int      `json:"guest-num_sta,omitempty"`
	RadioTable      []Radio  `json:"radio_table,omitempty"`
	VAPTable        []VAP    `json:"vap_table,omitempty"`
	PortTable       []Port   `json:"port_table,omitempty"`
	EthernetTable   []Port   `json:"ethernet_table,omitempty"`
	Temperatures    []Temp   `json:"temperatures,omitempty"`
}

// SysStats holds system statistics.
type SysStats struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"mem,omitempty"`
	Uptime string `json:"uptime,omitempty"`
}

// Radio represents a wireless radio.
type Radio struct {
	Name           string `json:"name,omitempty"`
	Radio          string `json:"radio,omitempty"`
	Channel        int    `json:"channel,omitempty"`
	HT             int    `json:"ht,omitempty"`
	TxPower        int    `json:"tx_power,omitempty"`
	TxPowerMode    string `json:"tx_power_mode,omitempty"`
	MinTxPower     int    `json:"min_txpower,omitempty"`
	MaxTxPower     int    `json:"max_txpower,omitempty"`
	NumSta         int    `json:"num_sta,omitempty"`
	CurrentChannel int    `json:"current_channel,omitempty"`
}

// VAP represents a virtual access point.
type VAP struct {
	ID        string `json:"_id,omitempty"`
	Name      string `json:"name,omitempty"`
	ESSID     string `json:"essid,omitempty"`
	BSSID     string `json:"bssid,omitempty"`
	Radio     string `json:"radio,omitempty"`
	Channel   int    `json:"channel,omitempty"`
	NumSta    int    `json:"num_sta,omitempty"`
	IsGuest   bool   `json:"is_guest"`
	IsWep     bool   `json:"is_wep"`
	Up        bool   `json:"up"`
	TxBytes   int64  `json:"tx_bytes,omitempty"`
	RxBytes   int64  `json:"rx_bytes,omitempty"`
	TxPackets int64  `json:"tx_packets,omitempty"`
	RxPackets int64  `json:"rx_packets,omitempty"`
}

// Port represents a switch port.
type Port struct {
	PortIdx       int    `json:"port_idx,omitempty"`
	Name          string `json:"name,omitempty"`
	Media         string `json:"media,omitempty"`
	Speed         int    `json:"speed,omitempty"`
	FullDuplex    bool   `json:"full_duplex"`
	Enable        bool   `json:"enable"`
	Up            bool   `json:"up"`
	POEEnable     bool   `json:"poe_enable"`
	POEMode       string `json:"poe_mode,omitempty"`
	POEGood       bool   `json:"poe_good"`
	POECurrent    string `json:"poe_current,omitempty"`
	POEVoltage    string `json:"poe_voltage,omitempty"`
	POEPower      string `json:"poe_power,omitempty"`
	TxBytes       int64  `json:"tx_bytes,omitempty"`
	RxBytes       int64  `json:"rx_bytes,omitempty"`
	TxPackets     int64  `json:"tx_packets,omitempty"`
	RxPackets     int64  `json:"rx_packets,omitempty"`
	TxBroadcast   int64  `json:"tx_broadcast,omitempty"`
	RxBroadcast   int64  `json:"rx_broadcast,omitempty"`
	TxMulticast   int64  `json:"tx_multicast,omitempty"`
	RxMulticast   int64  `json:"rx_multicast,omitempty"`
	TxDropped     int64  `json:"tx_dropped,omitempty"`
	RxDropped     int64  `json:"rx_dropped,omitempty"`
	TxErrors      int64  `json:"tx_errors,omitempty"`
	RxErrors      int64  `json:"rx_errors,omitempty"`
	STPState      string `json:"stp_state,omitempty"`
	STPPathCost   int    `json:"stp_pathcost,omitempty"`
	NetworkName   string `json:"network_name,omitempty"`
}

// Temp represents a temperature sensor reading.
type Temp struct {
	Name  string  `json:"name,omitempty"`
	Type  string  `json:"type,omitempty"`
	Value float64 `json:"value,omitempty"`
}

// DisplayName returns the best available name for the device.
func (d *Device) DisplayName() string {
	if d.Name != "" {
		return d.Name
	}
	if d.ModelName != "" {
		return d.ModelName
	}
	return d.MAC
}

// StateString returns a human-readable state string.
func (d *Device) StateString() string {
	states := map[int]string{
		0: "offline",
		1: "connected",
		2: "pending",
		4: "upgrading",
		5: "provisioning",
		6: "heartbeat missed",
		7: "adopting",
		9: "adoption failed",
	}
	if s, ok := states[d.State]; ok {
		return s
	}
	return fmt.Sprintf("unknown (%d)", d.State)
}

// DeviceManager handles device operations.
type DeviceManager struct {
	conn *ConnectionManager
}

// NewDeviceManager creates a new device manager.
func NewDeviceManager(conn *ConnectionManager) *DeviceManager {
	return &DeviceManager{conn: conn}
}

// ListDevices returns all adopted devices.
func (dm *DeviceManager) ListDevices(ctx context.Context, deviceType string) ([]Device, error) {
	data, err := dm.conn.Request(ctx, "GET", dm.conn.GetSitePath("/stat/device"), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get devices: %w", err)
	}

	var devices []Device
	if err := json.Unmarshal(data, &devices); err != nil {
		return nil, fmt.Errorf("failed to parse devices: %w", err)
	}

	if deviceType != "" {
		var filtered []Device
		for _, d := range devices {
			if d.Type == deviceType {
				filtered = append(filtered, d)
			}
		}
		return filtered, nil
	}

	return devices, nil
}

// GetDeviceDetails returns detailed information for a specific device.
func (dm *DeviceManager) GetDeviceDetails(ctx context.Context, mac string) (*Device, error) {
	devices, err := dm.ListDevices(ctx, "")
	if err != nil {
		return nil, err
	}

	for _, d := range devices {
		if d.MAC == mac {
			return &d, nil
		}
	}

	return nil, fmt.Errorf("device not found: %s", mac)
}

// RebootDevice reboots a device.
func (dm *DeviceManager) RebootDevice(ctx context.Context, mac string) error {
	cmd := map[string]interface{}{
		"cmd":      "restart",
		"mac":      mac,
		"reboot_type": "soft",
	}

	_, err := dm.conn.Request(ctx, "POST", dm.conn.GetSitePath("/cmd/devmgr"), cmd)
	if err != nil {
		return fmt.Errorf("failed to reboot device: %w", err)
	}

	return nil
}

// RenameDevice updates the name of a device.
func (dm *DeviceManager) RenameDevice(ctx context.Context, mac, name string) error {
	device, err := dm.GetDeviceDetails(ctx, mac)
	if err != nil {
		return err
	}

	if device.ID == "" {
		return fmt.Errorf("device has no ID, cannot rename")
	}

	update := map[string]interface{}{
		"name": name,
	}

	_, err = dm.conn.Request(ctx, "PUT", dm.conn.GetSitePath("/rest/device/"+device.ID), update)
	if err != nil {
		return fmt.Errorf("failed to rename device: %w", err)
	}

	return nil
}

// AdoptDevice adopts a pending device.
func (dm *DeviceManager) AdoptDevice(ctx context.Context, mac string) error {
	cmd := map[string]interface{}{
		"cmd": "adopt",
		"mac": mac,
	}

	_, err := dm.conn.Request(ctx, "POST", dm.conn.GetSitePath("/cmd/devmgr"), cmd)
	if err != nil {
		return fmt.Errorf("failed to adopt device: %w", err)
	}

	return nil
}

// UpgradeDevice initiates a firmware upgrade.
func (dm *DeviceManager) UpgradeDevice(ctx context.Context, mac string) error {
	cmd := map[string]interface{}{
		"cmd": "upgrade",
		"mac": mac,
	}

	_, err := dm.conn.Request(ctx, "POST", dm.conn.GetSitePath("/cmd/devmgr"), cmd)
	if err != nil {
		return fmt.Errorf("failed to upgrade device: %w", err)
	}

	return nil
}
