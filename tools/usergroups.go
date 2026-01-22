package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirkirby/go-mcp-ubiquiti/managers"
)

// RegisterUserGroupTools registers all user group (bandwidth profile) tools.
func RegisterUserGroupTools(s *server.MCPServer, userGroupMgr *managers.UserGroupManager) {
	// unifi_list_user_groups
	s.AddTool(
		mcp.NewTool("unifi_list_user_groups",
			mcp.WithDescription("List all user groups (bandwidth profiles)"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			groups, err := userGroupMgr.ListUserGroups(ctx)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to list user groups: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"count":   len(groups),
				"groups":  groups,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_get_user_group_details
	s.AddTool(
		mcp.NewTool("unifi_get_user_group_details",
			mcp.WithDescription("Get detailed information about a specific user group"),
			mcp.WithString("group_id",
				mcp.Required(),
				mcp.Description("ID of the user group"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			groupID := req.GetString("group_id", "")
			if groupID == "" {
				return mcp.NewToolResultError("group_id is required"), nil
			}

			group, err := userGroupMgr.GetUserGroupDetails(ctx, groupID)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to get user group: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"group":   group,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_create_user_group
	s.AddTool(
		mcp.NewTool("unifi_create_user_group",
			mcp.WithDescription("Create a new user group (bandwidth profile)"),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Name for the user group"),
			),
			mcp.WithNumber("download_kbps",
				mcp.Description("Download bandwidth limit in Kbps (-1 for unlimited)"),
			),
			mcp.WithNumber("upload_kbps",
				mcp.Description("Upload bandwidth limit in Kbps (-1 for unlimited)"),
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
				return mcp.NewToolResultText(fmt.Sprintf("Preview: Would create user group '%s'. Set confirm=true to execute.", name)), nil
			}

			group := &managers.UserGroup{
				Name: name,
			}

			group.QOSRateMaxDown = req.GetInt("download_kbps", 0)
			group.QOSRateMaxUp = req.GetInt("upload_kbps", 0)

			created, err := userGroupMgr.CreateUserGroup(ctx, group)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to create user group: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("User group '%s' has been created", name),
				"group":   created,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_update_user_group
	s.AddTool(
		mcp.NewTool("unifi_update_user_group",
			mcp.WithDescription("Update an existing user group"),
			mcp.WithString("group_id",
				mcp.Required(),
				mcp.Description("ID of the user group to update"),
			),
			mcp.WithString("name",
				mcp.Description("New name for the group"),
			),
			mcp.WithNumber("download_kbps",
				mcp.Description("Download bandwidth limit in Kbps (-1 for unlimited)"),
			),
			mcp.WithNumber("upload_kbps",
				mcp.Description("Upload bandwidth limit in Kbps (-1 for unlimited)"),
			),
			mcp.WithBoolean("confirm",
				mcp.Required(),
				mcp.Description("Must be true to execute the update"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			groupID := req.GetString("group_id", "")
			confirm := req.GetBool("confirm", false)

			if groupID == "" {
				return mcp.NewToolResultError("group_id is required"), nil
			}
			if !confirm {
				return mcp.NewToolResultText(fmt.Sprintf("Preview: Would update user group %s. Set confirm=true to execute.", groupID)), nil
			}

			updates := make(map[string]interface{})
			if name := req.GetString("name", ""); name != "" {
				updates["name"] = name
			}
			// Check for numeric params explicitly
			args := req.GetArguments()
			if downVal, ok := args["download_kbps"]; ok {
				if down, ok := downVal.(float64); ok {
					updates["qos_rate_max_down"] = int(down)
				}
			}
			if upVal, ok := args["upload_kbps"]; ok {
				if up, ok := upVal.(float64); ok {
					updates["qos_rate_max_up"] = int(up)
				}
			}

			if len(updates) == 0 {
				return mcp.NewToolResultError("no updates specified"), nil
			}

			if err := userGroupMgr.UpdateUserGroup(ctx, groupID, updates); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to update user group: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("User group %s has been updated", groupID),
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)
}
