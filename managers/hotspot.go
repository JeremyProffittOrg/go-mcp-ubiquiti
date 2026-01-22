package managers

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

// Voucher represents a hotspot voucher.
type Voucher struct {
	ID             string `json:"_id,omitempty"`
	SiteID         string `json:"site_id,omitempty"`
	Code           string `json:"code,omitempty"`
	CreateTime     int64  `json:"create_time,omitempty"`
	Duration       int    `json:"duration,omitempty"` // minutes
	Quota          int    `json:"quota,omitempty"`    // 0 = unlimited, 1 = single use, n = multi-use
	Used           int    `json:"used,omitempty"`
	Note           string `json:"note,omitempty"`
	Status         string `json:"status,omitempty"`
	StatusExpires  int64  `json:"status_expires,omitempty"`
	AdminName      string `json:"admin_name,omitempty"`
	ForHotspot     bool   `json:"for_hotspot"`
	QOSOverwrite   bool   `json:"qos_overwrite"`
	QOSRateMaxUp   int    `json:"qos_rate_max_up,omitempty"`
	QOSRateMaxDown int    `json:"qos_rate_max_down,omitempty"`
	QOSUsageQuota  int    `json:"qos_usage_quota,omitempty"` // MB
	StartTime      int64  `json:"start_time,omitempty"`
	EndTime        int64  `json:"end_time,omitempty"`
}

// HotspotManager handles hotspot and voucher operations.
type HotspotManager struct {
	conn *ConnectionManager
}

// NewHotspotManager creates a new hotspot manager.
func NewHotspotManager(conn *ConnectionManager) *HotspotManager {
	return &HotspotManager{conn: conn}
}

// ListVouchers returns all vouchers.
func (hm *HotspotManager) ListVouchers(ctx context.Context) ([]Voucher, error) {
	endpoint := hm.conn.GetSitePath("/stat/voucher")

	data, err := hm.conn.Request(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get vouchers: %w", err)
	}

	var vouchers []Voucher
	if err := json.Unmarshal(data, &vouchers); err != nil {
		return nil, fmt.Errorf("failed to parse vouchers: %w", err)
	}

	return vouchers, nil
}

// GetVoucherDetails returns details for a specific voucher.
func (hm *HotspotManager) GetVoucherDetails(ctx context.Context, voucherID string) (*Voucher, error) {
	vouchers, err := hm.ListVouchers(ctx)
	if err != nil {
		return nil, err
	}

	for _, v := range vouchers {
		if v.ID == voucherID {
			return &v, nil
		}
	}

	return nil, fmt.Errorf("voucher not found: %s", voucherID)
}

// GetVoucherByCode finds a voucher by its code.
func (hm *HotspotManager) GetVoucherByCode(ctx context.Context, code string) (*Voucher, error) {
	vouchers, err := hm.ListVouchers(ctx)
	if err != nil {
		return nil, err
	}

	for _, v := range vouchers {
		if v.Code == code {
			return &v, nil
		}
	}

	return nil, fmt.Errorf("voucher not found with code: %s", code)
}

// CreateVoucher creates a new voucher.
func (hm *HotspotManager) CreateVoucher(ctx context.Context, count int, quota int, duration int, note string, upKbps, downKbps, mbytes int) ([]Voucher, error) {
	if count <= 0 {
		count = 1
	}
	if duration <= 0 {
		duration = 60 // 1 hour default
	}

	cmd := map[string]interface{}{
		"cmd":    "create-voucher",
		"n":      count,
		"quota":  quota,
		"expire": duration, // minutes
	}

	if note != "" {
		cmd["note"] = note
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

	data, err := hm.conn.Request(ctx, "POST", hm.conn.GetSitePath("/cmd/hotspot"), cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to create voucher: %w", err)
	}

	// The response contains the created voucher(s)
	var result struct {
		CreateTime int64    `json:"create_time"`
		Vouchers   []string `json:"voucher"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		// Try to parse as voucher list directly
		var vouchers []Voucher
		if err := json.Unmarshal(data, &vouchers); err != nil {
			return nil, fmt.Errorf("failed to parse created vouchers: %w", err)
		}
		return vouchers, nil
	}

	// Build voucher objects from codes
	var vouchers []Voucher
	for _, code := range result.Vouchers {
		vouchers = append(vouchers, Voucher{
			Code:       code,
			CreateTime: result.CreateTime,
			Duration:   duration,
			Quota:      quota,
			Note:       note,
		})
	}

	return vouchers, nil
}

// RevokeVoucher revokes/deletes a voucher.
func (hm *HotspotManager) RevokeVoucher(ctx context.Context, voucherID string) error {
	cmd := map[string]interface{}{
		"cmd": "delete-voucher",
		"_id": voucherID,
	}

	_, err := hm.conn.Request(ctx, "POST", hm.conn.GetSitePath("/cmd/hotspot"), cmd)
	if err != nil {
		return fmt.Errorf("failed to revoke voucher: %w", err)
	}

	return nil
}

// GenerateVoucherCode generates a random voucher code.
func GenerateVoucherCode(length int) string {
	if length <= 0 {
		length = 10
	}

	const charset = "0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	code := make([]byte, length)
	for i := range code {
		code[i] = charset[r.Intn(len(charset))]
	}

	return string(code)
}
