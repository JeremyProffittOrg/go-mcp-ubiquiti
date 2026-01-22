package managers

import (
	"context"
	"encoding/json"
	"fmt"
)

// QoSRule represents a QoS (traffic shaping) rule.
type QoSRule struct {
	ID               string   `json:"_id,omitempty"`
	Name             string   `json:"name"`
	Enabled          bool     `json:"enabled"`
	SiteID           string   `json:"site_id,omitempty"`
	Description      string   `json:"description,omitempty"`
	TargetDevices    []TargetDevice `json:"target_devices,omitempty"`
	NetworkID        string   `json:"network_id,omitempty"`
	MatchingTarget   string   `json:"matching_target,omitempty"` // INTERNET, APP, DOMAIN, IP
	AppCategoryID    string   `json:"app_category_id,omitempty"`
	AppID            string   `json:"app_id,omitempty"`
	Domains          []string `json:"domains,omitempty"`
	IPAddresses      []string `json:"ip_addresses,omitempty"`
	BandwidthProfile string   `json:"bandwidth_profile,omitempty"` // preset name or custom
	DownloadKbps     int      `json:"download_kbps,omitempty"`
	UploadKbps       int      `json:"upload_kbps,omitempty"`
	DownloadLimit    int      `json:"download_limit,omitempty"`
	UploadLimit      int      `json:"upload_limit,omitempty"`
	Priority         int      `json:"priority,omitempty"`
	Action           string   `json:"action,omitempty"` // RATE_LIMIT, BLOCK
	Schedule         string   `json:"schedule,omitempty"`
}

// QoSManager handles QoS operations.
type QoSManager struct {
	conn *ConnectionManager
}

// NewQoSManager creates a new QoS manager.
func NewQoSManager(conn *ConnectionManager) *QoSManager {
	return &QoSManager{conn: conn}
}

// ListQoSRules returns all QoS rules.
func (qm *QoSManager) ListQoSRules(ctx context.Context) ([]QoSRule, error) {
	var endpoint string
	if qm.conn.IsUnifiOS() {
		endpoint = "/v2/api/site/" + qm.conn.Site() + "/traffic-rules"
	} else {
		endpoint = qm.conn.GetSitePath("/rest/trafficrule")
	}

	data, err := qm.conn.Request(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get QoS rules: %w", err)
	}

	var rules []QoSRule
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, fmt.Errorf("failed to parse QoS rules: %w", err)
	}

	return rules, nil
}

// GetQoSRuleDetails returns details for a specific QoS rule.
func (qm *QoSManager) GetQoSRuleDetails(ctx context.Context, ruleID string) (*QoSRule, error) {
	rules, err := qm.ListQoSRules(ctx)
	if err != nil {
		return nil, err
	}

	for _, r := range rules {
		if r.ID == ruleID {
			return &r, nil
		}
	}

	return nil, fmt.Errorf("QoS rule not found: %s", ruleID)
}

// ToggleQoSRule enables or disables a QoS rule.
func (qm *QoSManager) ToggleQoSRule(ctx context.Context, ruleID string, enabled bool) error {
	var endpoint string
	if qm.conn.IsUnifiOS() {
		endpoint = "/v2/api/site/" + qm.conn.Site() + "/traffic-rules/" + ruleID
	} else {
		endpoint = qm.conn.GetSitePath("/rest/trafficrule/" + ruleID)
	}

	update := map[string]interface{}{
		"enabled": enabled,
	}

	_, err := qm.conn.Request(ctx, "PUT", endpoint, update)
	if err != nil {
		return fmt.Errorf("failed to toggle QoS rule: %w", err)
	}

	return nil
}

// CreateQoSRule creates a new QoS rule.
func (qm *QoSManager) CreateQoSRule(ctx context.Context, rule *QoSRule) (*QoSRule, error) {
	var endpoint string
	if qm.conn.IsUnifiOS() {
		endpoint = "/v2/api/site/" + qm.conn.Site() + "/traffic-rules"
	} else {
		endpoint = qm.conn.GetSitePath("/rest/trafficrule")
	}

	data, err := qm.conn.Request(ctx, "POST", endpoint, rule)
	if err != nil {
		return nil, fmt.Errorf("failed to create QoS rule: %w", err)
	}

	var rules []QoSRule
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, fmt.Errorf("failed to parse created rule: %w", err)
	}

	if len(rules) == 0 {
		return nil, fmt.Errorf("no rule returned after creation")
	}

	return &rules[0], nil
}

// UpdateQoSRule updates an existing QoS rule.
func (qm *QoSManager) UpdateQoSRule(ctx context.Context, ruleID string, updates map[string]interface{}) error {
	var endpoint string
	if qm.conn.IsUnifiOS() {
		endpoint = "/v2/api/site/" + qm.conn.Site() + "/traffic-rules/" + ruleID
	} else {
		endpoint = qm.conn.GetSitePath("/rest/trafficrule/" + ruleID)
	}

	_, err := qm.conn.Request(ctx, "PUT", endpoint, updates)
	if err != nil {
		return fmt.Errorf("failed to update QoS rule: %w", err)
	}

	return nil
}

// DeleteQoSRule deletes a QoS rule.
func (qm *QoSManager) DeleteQoSRule(ctx context.Context, ruleID string) error {
	var endpoint string
	if qm.conn.IsUnifiOS() {
		endpoint = "/v2/api/site/" + qm.conn.Site() + "/traffic-rules/" + ruleID
	} else {
		endpoint = qm.conn.GetSitePath("/rest/trafficrule/" + ruleID)
	}

	_, err := qm.conn.Request(ctx, "DELETE", endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to delete QoS rule: %w", err)
	}

	return nil
}
