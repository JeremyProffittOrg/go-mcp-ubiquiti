package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirkirby/go-mcp-ubiquiti/managers"
)

// RegisterQoSTools registers all QoS-related tools.
func RegisterQoSTools(s *server.MCPServer, qosMgr *managers.QoSManager) {
	// unifi_list_qos_rules
	s.AddTool(
		mcp.NewTool("unifi_list_qos_rules",
			mcp.WithDescription("List all QoS (traffic shaping) rules"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			rules, err := qosMgr.ListQoSRules(ctx)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to list QoS rules: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"count":   len(rules),
				"rules":   rules,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_get_qos_rule_details
	s.AddTool(
		mcp.NewTool("unifi_get_qos_rule_details",
			mcp.WithDescription("Get detailed information about a specific QoS rule"),
			mcp.WithString("rule_id",
				mcp.Required(),
				mcp.Description("ID of the QoS rule"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			ruleID := req.GetString("rule_id", "")
			if ruleID == "" {
				return mcp.NewToolResultError("rule_id is required"), nil
			}

			rule, err := qosMgr.GetQoSRuleDetails(ctx, ruleID)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to get QoS rule: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"rule":    rule,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_toggle_qos_rule
	s.AddTool(
		mcp.NewTool("unifi_toggle_qos_rule",
			mcp.WithDescription("Enable or disable a QoS rule"),
			mcp.WithString("rule_id",
				mcp.Required(),
				mcp.Description("ID of the QoS rule"),
			),
			mcp.WithBoolean("enabled",
				mcp.Required(),
				mcp.Description("Whether to enable (true) or disable (false) the rule"),
			),
			mcp.WithBoolean("confirm",
				mcp.Required(),
				mcp.Description("Must be true to execute the toggle"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			ruleID := req.GetString("rule_id", "")
			enabled := req.GetBool("enabled", false)
			confirm := req.GetBool("confirm", false)

			if ruleID == "" {
				return mcp.NewToolResultError("rule_id is required"), nil
			}

			action := "disable"
			if enabled {
				action = "enable"
			}

			if !confirm {
				return mcp.NewToolResultText(fmt.Sprintf("Preview: Would %s QoS rule %s. Set confirm=true to execute.", action, ruleID)), nil
			}

			if err := qosMgr.ToggleQoSRule(ctx, ruleID, enabled); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to toggle QoS rule: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("QoS rule %s has been %sd", ruleID, action),
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_update_qos_rule
	s.AddTool(
		mcp.NewTool("unifi_update_qos_rule",
			mcp.WithDescription("Update an existing QoS rule"),
			mcp.WithString("rule_id",
				mcp.Required(),
				mcp.Description("ID of the QoS rule to update"),
			),
			mcp.WithString("name",
				mcp.Description("New name for the rule"),
			),
			mcp.WithNumber("download_kbps",
				mcp.Description("Download bandwidth limit in Kbps"),
			),
			mcp.WithNumber("upload_kbps",
				mcp.Description("Upload bandwidth limit in Kbps"),
			),
			mcp.WithBoolean("confirm",
				mcp.Required(),
				mcp.Description("Must be true to execute the update"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			ruleID := req.GetString("rule_id", "")
			confirm := req.GetBool("confirm", false)

			if ruleID == "" {
				return mcp.NewToolResultError("rule_id is required"), nil
			}
			if !confirm {
				return mcp.NewToolResultText(fmt.Sprintf("Preview: Would update QoS rule %s. Set confirm=true to execute.", ruleID)), nil
			}

			updates := make(map[string]interface{})
			if name := req.GetString("name", ""); name != "" {
				updates["name"] = name
			}
			// Check for numeric params - need to handle float64 for numbers
			args := req.GetArguments()
			if downVal, ok := args["download_kbps"]; ok {
				if down, ok := downVal.(float64); ok {
					updates["download_kbps"] = int(down)
				}
			}
			if upVal, ok := args["upload_kbps"]; ok {
				if up, ok := upVal.(float64); ok {
					updates["upload_kbps"] = int(up)
				}
			}

			if len(updates) == 0 {
				return mcp.NewToolResultError("no updates specified"), nil
			}

			if err := qosMgr.UpdateQoSRule(ctx, ruleID, updates); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to update QoS rule: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("QoS rule %s has been updated", ruleID),
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_create_qos_rule
	s.AddTool(
		mcp.NewTool("unifi_create_qos_rule",
			mcp.WithDescription("Create a new QoS rule"),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Name for the QoS rule"),
			),
			mcp.WithString("matching_target",
				mcp.Description("What to match: INTERNET, APP, DOMAIN, IP"),
			),
			mcp.WithNumber("download_kbps",
				mcp.Description("Download bandwidth limit in Kbps"),
			),
			mcp.WithNumber("upload_kbps",
				mcp.Description("Upload bandwidth limit in Kbps"),
			),
			mcp.WithBoolean("confirm",
				mcp.Required(),
				mcp.Description("Must be true to execute creation"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			name := req.GetString("name", "")
			confirm := req.GetBool("confirm", false)

			if name == "" {
				return mcp.NewToolResultError("name is required"), nil
			}
			if !confirm {
				return mcp.NewToolResultText(fmt.Sprintf("Preview: Would create QoS rule '%s'. Set confirm=true to execute.", name)), nil
			}

			rule := &managers.QoSRule{
				Name:    name,
				Enabled: true,
			}

			if target := req.GetString("matching_target", ""); target != "" {
				rule.MatchingTarget = target
			}
			rule.DownloadKbps = req.GetInt("download_kbps", 0)
			rule.UploadKbps = req.GetInt("upload_kbps", 0)

			created, err := qosMgr.CreateQoSRule(ctx, rule)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to create QoS rule: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("QoS rule '%s' has been created", name),
				"rule":    created,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)
}
