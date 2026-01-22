package managers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Event represents a UniFi event.
type Event struct {
	ID           string `json:"_id,omitempty"`
	Key          string `json:"key,omitempty"`
	Subsystem    string `json:"subsystem,omitempty"`
	SiteID       string `json:"site_id,omitempty"`
	Time         int64  `json:"time,omitempty"`
	DateTime     string `json:"datetime,omitempty"`
	Message      string `json:"msg,omitempty"`
	IsAdmin      bool   `json:"is_admin"`
	Admin        string `json:"admin,omitempty"`
	User         string `json:"user,omitempty"`
	Hostname     string `json:"hostname,omitempty"`
	SSID         string `json:"ssid,omitempty"`
	AP           string `json:"ap,omitempty"`
	APName       string `json:"ap_name,omitempty"`
	APFrom       string `json:"ap_from,omitempty"`
	APTo         string `json:"ap_to,omitempty"`
	Channel      int    `json:"channel,omitempty"`
	ChannelFrom  int    `json:"channel_from,omitempty"`
	ChannelTo    int    `json:"channel_to,omitempty"`
	Radio        string `json:"radio,omitempty"`
	RadioFrom    string `json:"radio_from,omitempty"`
	RadioTo      string `json:"radio_to,omitempty"`
	Duration     int64  `json:"duration,omitempty"`
	Bytes        int64  `json:"bytes,omitempty"`
	Guest        string `json:"guest,omitempty"`
	GuestMAC     string `json:"guest_mac,omitempty"`
	ClientMAC    string `json:"client_mac,omitempty"`
	Network      string `json:"network,omitempty"`
	SW           string `json:"sw,omitempty"`
	SWName       string `json:"sw_name,omitempty"`
	Port         int    `json:"port,omitempty"`
	GW           string `json:"gw,omitempty"`
	GWName       string `json:"gw_name,omitempty"`
	InIface      string `json:"in_iface,omitempty"`
	OutIface     string `json:"out_iface,omitempty"`
	SrcIP        string `json:"src_ip,omitempty"`
	DstIP        string `json:"dest_ip,omitempty"`
	SrcPort      int    `json:"src_port,omitempty"`
	DstPort      int    `json:"dest_port,omitempty"`
	Protocol     string `json:"proto,omitempty"`
	UniqueAlertID string `json:"unique_alertid,omitempty"`
}

// Alarm represents a UniFi alarm.
type Alarm struct {
	ID              string `json:"_id,omitempty"`
	Key             string `json:"key,omitempty"`
	Subsystem       string `json:"subsystem,omitempty"`
	SiteID          string `json:"site_id,omitempty"`
	Time            int64  `json:"time,omitempty"`
	DateTime        string `json:"datetime,omitempty"`
	Message         string `json:"msg,omitempty"`
	Archived        bool   `json:"archived"`
	HandledAdminID  string `json:"handled_admin_id,omitempty"`
	HandledTime     int64  `json:"handled_time,omitempty"`
	AP              string `json:"ap,omitempty"`
	APName          string `json:"ap_name,omitempty"`
	DeviceMAC       string `json:"device_mac,omitempty"`
	DeviceName      string `json:"device_name,omitempty"`
	GWMAC           string `json:"gw_mac,omitempty"`
	GWName          string `json:"gw_name,omitempty"`
	SWMAC           string `json:"sw_mac,omitempty"`
	SWName          string `json:"sw_name,omitempty"`
	ClientMAC       string `json:"client_mac,omitempty"`
	DestPort        int    `json:"dest_port,omitempty"`
	UniqueAlertID   string `json:"unique_alertid,omitempty"`
	CatName         string `json:"catname,omitempty"`
	InnerAlertType  string `json:"inner_alert_type,omitempty"`
	ThreatManaged   bool   `json:"threat_managed"`
}

// EventManager handles event and alarm operations.
type EventManager struct {
	conn *ConnectionManager
}

// NewEventManager creates a new event manager.
func NewEventManager(conn *ConnectionManager) *EventManager {
	return &EventManager{conn: conn}
}

// ListEvents returns recent events.
func (em *EventManager) ListEvents(ctx context.Context, hours int, limit int) ([]Event, error) {
	if hours <= 0 {
		hours = 24
	}
	if limit <= 0 {
		limit = 100
	}

	// Calculate start time
	start := time.Now().Add(-time.Duration(hours) * time.Hour).UnixMilli()

	params := map[string]interface{}{
		"_start": start,
		"_limit": limit,
	}

	data, err := em.conn.Request(ctx, "GET", em.conn.GetSitePath("/stat/event"), params)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	var events []Event
	if err := json.Unmarshal(data, &events); err != nil {
		return nil, fmt.Errorf("failed to parse events: %w", err)
	}

	return events, nil
}

// ListAlarms returns active alarms.
func (em *EventManager) ListAlarms(ctx context.Context, archived bool) ([]Alarm, error) {
	endpoint := em.conn.GetSitePath("/stat/alarm")
	if !archived {
		endpoint = em.conn.GetSitePath("/stat/alarm?archived=false")
	}

	data, err := em.conn.Request(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get alarms: %w", err)
	}

	var alarms []Alarm
	if err := json.Unmarshal(data, &alarms); err != nil {
		return nil, fmt.Errorf("failed to parse alarms: %w", err)
	}

	return alarms, nil
}

// ArchiveAlarm archives a specific alarm.
func (em *EventManager) ArchiveAlarm(ctx context.Context, alarmID string) error {
	cmd := map[string]interface{}{
		"cmd":  "archive-alarm",
		"_id":  alarmID,
	}

	_, err := em.conn.Request(ctx, "POST", em.conn.GetSitePath("/cmd/evtmgr"), cmd)
	if err != nil {
		return fmt.Errorf("failed to archive alarm: %w", err)
	}

	return nil
}

// ArchiveAllAlarms archives all alarms.
func (em *EventManager) ArchiveAllAlarms(ctx context.Context) error {
	cmd := map[string]interface{}{
		"cmd": "archive-all-alarms",
	}

	_, err := em.conn.Request(ctx, "POST", em.conn.GetSitePath("/cmd/evtmgr"), cmd)
	if err != nil {
		return fmt.Errorf("failed to archive all alarms: %w", err)
	}

	return nil
}

// GetEventTypes returns available event types.
func (em *EventManager) GetEventTypes() []string {
	return []string{
		"EVT_AP_Connected",
		"EVT_AP_Disconnected",
		"EVT_AP_Restarted",
		"EVT_AP_Upgraded",
		"EVT_AP_ChannelChanged",
		"EVT_AP_Lost_Contact",
		"EVT_AP_Possible_Interference",
		"EVT_GW_Connected",
		"EVT_GW_Disconnected",
		"EVT_GW_Restarted",
		"EVT_GW_Upgraded",
		"EVT_GW_Lost_Contact",
		"EVT_GW_WANTransition",
		"EVT_SW_Connected",
		"EVT_SW_Disconnected",
		"EVT_SW_Restarted",
		"EVT_SW_Upgraded",
		"EVT_SW_Lost_Contact",
		"EVT_SW_StpPortBlocking",
		"EVT_WU_Connected",
		"EVT_WU_Disconnected",
		"EVT_WU_Roam",
		"EVT_WU_RoamRadio",
		"EVT_LU_Connected",
		"EVT_LU_Disconnected",
		"EVT_GU_Connected",
		"EVT_GU_Disconnected",
		"EVT_AD_Login",
		"EVT_AD_Logout",
		"EVT_IDS_IPS",
	}
}
