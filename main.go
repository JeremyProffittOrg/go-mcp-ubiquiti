// UniFi Network MCP Server
// A Go implementation of the Model Context Protocol server for UniFi Network Controller.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirkirby/go-mcp-ubiquiti/config"
	"github.com/sirkirby/go-mcp-ubiquiti/managers"
	"github.com/sirkirby/go-mcp-ubiquiti/tools"
)

const (
	serverName    = "unifi-network-mcp"
	serverVersion = "1.0.0"
)

func main() {
	// Configure logging to stderr (stdout is reserved for JSON-RPC)
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// Not an error if .env doesn't exist
		log.Printf("Note: .env file not loaded: %v", err)
	}

	// Load configuration
	cfg := config.Load()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Set log level
	if strings.ToUpper(cfg.LogLevel) == "DEBUG" {
		log.Printf("Debug logging enabled")
	}

	log.Printf("Starting %s v%s", serverName, serverVersion)
	log.Printf("Connecting to UniFi controller at %s (site: %s)", cfg.Host, cfg.Site)

	// Initialize connection manager
	connMgr := managers.NewConnectionManager(cfg)

	// Login to the controller
	ctx := context.Background()
	if err := connMgr.Login(ctx); err != nil {
		log.Fatalf("Failed to login to UniFi controller: %v", err)
	}

	// Initialize all managers
	clientMgr := managers.NewClientManager(connMgr)
	deviceMgr := managers.NewDeviceManager(connMgr)
	networkMgr := managers.NewNetworkManager(connMgr)
	systemMgr := managers.NewSystemManager(connMgr)
	eventMgr := managers.NewEventManager(connMgr)
	statsMgr := managers.NewStatsManager(connMgr)
	firewallMgr := managers.NewFirewallManager(connMgr)
	routingMgr := managers.NewRoutingManager(connMgr)
	vpnMgr := managers.NewVPNManager(connMgr)
	qosMgr := managers.NewQoSManager(connMgr)
	hotspotMgr := managers.NewHotspotManager(connMgr)
	userGroupMgr := managers.NewUserGroupManager(connMgr)

	// Create MCP server
	s := server.NewMCPServer(
		serverName,
		serverVersion,
		server.WithToolCapabilities(true),
	)

	// Register all tools
	log.Printf("Registering MCP tools...")

	tools.RegisterClientTools(s, clientMgr)
	tools.RegisterDeviceTools(s, deviceMgr)
	tools.RegisterNetworkTools(s, networkMgr)
	tools.RegisterSystemTools(s, systemMgr, cfg.Site)
	tools.RegisterEventTools(s, eventMgr)
	tools.RegisterStatsTools(s, statsMgr)
	tools.RegisterFirewallTools(s, firewallMgr)
	tools.RegisterRoutingTools(s, routingMgr)
	tools.RegisterVPNTools(s, vpnMgr)
	tools.RegisterQoSTools(s, qosMgr)
	tools.RegisterHotspotTools(s, hotspotMgr)
	tools.RegisterUserGroupTools(s, userGroupMgr)

	log.Printf("MCP server ready, starting stdio transport...")

	// Start stdio transport
	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
