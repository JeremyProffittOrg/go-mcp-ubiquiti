package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirkirby/go-mcp-ubiquiti/managers"
)

// RegisterHotspotTools registers all hotspot and voucher tools.
func RegisterHotspotTools(s *server.MCPServer, hotspotMgr *managers.HotspotManager) {
	// unifi_list_vouchers
	s.AddTool(
		mcp.NewTool("unifi_list_vouchers",
			mcp.WithDescription("List all hotspot vouchers"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			vouchers, err := hotspotMgr.ListVouchers(ctx)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to list vouchers: %v", err)), nil
			}

			result := map[string]interface{}{
				"success":  true,
				"count":    len(vouchers),
				"vouchers": vouchers,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_get_voucher_details
	s.AddTool(
		mcp.NewTool("unifi_get_voucher_details",
			mcp.WithDescription("Get detailed information about a specific voucher"),
			mcp.WithString("voucher_id",
				mcp.Required(),
				mcp.Description("ID of the voucher"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			voucherID := req.GetString("voucher_id", "")
			if voucherID == "" {
				return mcp.NewToolResultError("voucher_id is required"), nil
			}

			voucher, err := hotspotMgr.GetVoucherDetails(ctx, voucherID)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to get voucher: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"voucher": voucher,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_create_voucher
	s.AddTool(
		mcp.NewTool("unifi_create_voucher",
			mcp.WithDescription("Create new hotspot voucher(s)"),
			mcp.WithNumber("count",
				mcp.Description("Number of vouchers to create (default: 1)"),
			),
			mcp.WithNumber("duration",
				mcp.Description("Duration in minutes (default: 60)"),
			),
			mcp.WithNumber("quota",
				mcp.Description("Usage quota: 0=unlimited, 1=single-use, n=multi-use"),
			),
			mcp.WithString("note",
				mcp.Description("Note/description for the voucher"),
			),
			mcp.WithNumber("up_kbps",
				mcp.Description("Upload bandwidth limit in Kbps"),
			),
			mcp.WithNumber("down_kbps",
				mcp.Description("Download bandwidth limit in Kbps"),
			),
			mcp.WithNumber("megabytes",
				mcp.Description("Data transfer limit in MB"),
			),
			mcp.WithBoolean("confirm",
				mcp.Required(),
				mcp.Description("Must be true to execute creation"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			confirm := req.GetBool("confirm", false)

			count := req.GetInt("count", 1)
			if count < 1 {
				count = 1
			}
			duration := req.GetInt("duration", 60)
			if duration < 1 {
				duration = 60
			}
			quota := req.GetInt("quota", 0)
			note := req.GetString("note", "")
			upKbps := req.GetInt("up_kbps", 0)
			downKbps := req.GetInt("down_kbps", 0)
			mbytes := req.GetInt("megabytes", 0)

			if !confirm {
				return mcp.NewToolResultText(fmt.Sprintf("Preview: Would create %d voucher(s) with %d minute duration. Set confirm=true to execute.", count, duration)), nil
			}

			vouchers, err := hotspotMgr.CreateVoucher(ctx, count, quota, duration, note, upKbps, downKbps, mbytes)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to create voucher: %v", err)), nil
			}

			result := map[string]interface{}{
				"success":  true,
				"message":  fmt.Sprintf("Created %d voucher(s)", len(vouchers)),
				"vouchers": vouchers,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_revoke_voucher
	s.AddTool(
		mcp.NewTool("unifi_revoke_voucher",
			mcp.WithDescription("Revoke/delete a voucher"),
			mcp.WithString("voucher_id",
				mcp.Required(),
				mcp.Description("ID of the voucher to revoke"),
			),
			mcp.WithBoolean("confirm",
				mcp.Required(),
				mcp.Description("Must be true to execute the revocation"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			voucherID := req.GetString("voucher_id", "")
			confirm := req.GetBool("confirm", false)

			if voucherID == "" {
				return mcp.NewToolResultError("voucher_id is required"), nil
			}
			if !confirm {
				return mcp.NewToolResultText(fmt.Sprintf("Preview: Would revoke voucher %s. Set confirm=true to execute.", voucherID)), nil
			}

			if err := hotspotMgr.RevokeVoucher(ctx, voucherID); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to revoke voucher: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("Voucher %s has been revoked", voucherID),
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)
}
