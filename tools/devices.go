package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirkirby/go-mcp-ubiquiti/managers"
)

// RegisterDeviceTools registers all device management tools.
func RegisterDeviceTools(s *server.MCPServer, deviceMgr *managers.DeviceManager) {
	// unifi_list_devices
	s.AddTool(
		mcp.NewTool("unifi_list_devices",
			mcp.WithDescription("List all UniFi devices (APs, switches, gateways)"),
			mcp.WithString("type",
				mcp.Description("Filter by device type: uap (access point), usw (switch), ugw (gateway), udm (dream machine)"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			deviceType := req.GetString("type", "")

			devices, err := deviceMgr.ListDevices(ctx, deviceType)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to list devices: %v", err)), nil
			}

			// Add state string to each device
			type DeviceWithState struct {
				*managers.Device
				StateString string `json:"state_string"`
			}

			var devicesWithState []DeviceWithState
			for _, d := range devices {
				dCopy := d
				devicesWithState = append(devicesWithState, DeviceWithState{
					Device:      &dCopy,
					StateString: d.StateString(),
				})
			}

			result := map[string]interface{}{
				"success": true,
				"count":   len(devices),
				"devices": devicesWithState,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_get_device_details
	s.AddTool(
		mcp.NewTool("unifi_get_device_details",
			mcp.WithDescription("Get detailed information about a specific device"),
			mcp.WithString("mac",
				mcp.Required(),
				mcp.Description("MAC address of the device"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			mac := req.GetString("mac", "")
			if mac == "" {
				return mcp.NewToolResultError("mac address is required"), nil
			}

			device, err := deviceMgr.GetDeviceDetails(ctx, mac)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to get device details: %v", err)), nil
			}

			result := map[string]interface{}{
				"success":      true,
				"device":       device,
				"state_string": device.StateString(),
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_reboot_device
	s.AddTool(
		mcp.NewTool("unifi_reboot_device",
			mcp.WithDescription("Reboot a UniFi device"),
			mcp.WithString("mac",
				mcp.Required(),
				mcp.Description("MAC address of the device to reboot"),
			),
			mcp.WithBoolean("confirm",
				mcp.Required(),
				mcp.Description("Must be true to execute the reboot"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			mac := req.GetString("mac", "")
			confirm := req.GetBool("confirm", false)

			if mac == "" {
				return mcp.NewToolResultError("mac address is required"), nil
			}
			if !confirm {
				return mcp.NewToolResultText(fmt.Sprintf("Preview: Would reboot device %s. Set confirm=true to execute.", mac)), nil
			}

			if err := deviceMgr.RebootDevice(ctx, mac); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to reboot device: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("Device %s is rebooting", mac),
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_rename_device
	s.AddTool(
		mcp.NewTool("unifi_rename_device",
			mcp.WithDescription("Rename a UniFi device"),
			mcp.WithString("mac",
				mcp.Required(),
				mcp.Description("MAC address of the device to rename"),
			),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("New name for the device"),
			),
			mcp.WithBoolean("confirm",
				mcp.Required(),
				mcp.Description("Must be true to execute the rename"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			mac := req.GetString("mac", "")
			name := req.GetString("name", "")
			confirm := req.GetBool("confirm", false)

			if mac == "" || name == "" {
				return mcp.NewToolResultError("mac and name are required"), nil
			}
			if !confirm {
				return mcp.NewToolResultText(fmt.Sprintf("Preview: Would rename device %s to '%s'. Set confirm=true to execute.", mac, name)), nil
			}

			if err := deviceMgr.RenameDevice(ctx, mac, name); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to rename device: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("Device %s has been renamed to '%s'", mac, name),
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_adopt_device
	s.AddTool(
		mcp.NewTool("unifi_adopt_device",
			mcp.WithDescription("Adopt a pending UniFi device"),
			mcp.WithString("mac",
				mcp.Required(),
				mcp.Description("MAC address of the device to adopt"),
			),
			mcp.WithBoolean("confirm",
				mcp.Required(),
				mcp.Description("Must be true to execute the adoption"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			mac := req.GetString("mac", "")
			confirm := req.GetBool("confirm", false)

			if mac == "" {
				return mcp.NewToolResultError("mac address is required"), nil
			}
			if !confirm {
				return mcp.NewToolResultText(fmt.Sprintf("Preview: Would adopt device %s. Set confirm=true to execute.", mac)), nil
			}

			if err := deviceMgr.AdoptDevice(ctx, mac); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to adopt device: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("Device %s is being adopted", mac),
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// unifi_upgrade_device
	s.AddTool(
		mcp.NewTool("unifi_upgrade_device",
			mcp.WithDescription("Upgrade firmware on a UniFi device"),
			mcp.WithString("mac",
				mcp.Required(),
				mcp.Description("MAC address of the device to upgrade"),
			),
			mcp.WithBoolean("confirm",
				mcp.Required(),
				mcp.Description("Must be true to execute the upgrade"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			mac := req.GetString("mac", "")
			confirm := req.GetBool("confirm", false)

			if mac == "" {
				return mcp.NewToolResultError("mac address is required"), nil
			}
			if !confirm {
				return mcp.NewToolResultText(fmt.Sprintf("Preview: Would upgrade device %s. Set confirm=true to execute.", mac)), nil
			}

			if err := deviceMgr.UpgradeDevice(ctx, mac); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to upgrade device: %v", err)), nil
			}

			result := map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("Firmware upgrade started for device %s", mac),
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)
}
