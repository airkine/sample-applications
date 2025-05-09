package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/dicedb/dicedb-go"
	"github.com/dicedb/dicedb-go/wire"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Initialize the MCP server with a name and version
	s := server.NewMCPServer(
		"DiceDB MCP", // Server name
		"0.1.0",      // Server version
		// Indicates that the list of tools is static
		server.WithToolCapabilities(false),
	)

	// Define the 'ping' tool to check DiceDB connectivity
	pingTool := mcp.NewTool("ping",
		mcp.WithDescription("Ping the DiceDB server to check connectivity"),
		mcp.WithString("url",
			mcp.Description("The URL of the DiceDB server in format 'host:port'"),
			mcp.DefaultString("localhost:7379"),
		),
	)

	// Define the 'get' tool to retrieve a value by key
	getTool := mcp.NewTool("get",
		mcp.WithDescription("Retrieve the value for a given key from DiceDB"),
		mcp.WithString("url",
			mcp.Description("The URL of the DiceDB server in format 'host:port'"),
			mcp.DefaultString("localhost:7379"),
		),
		mcp.WithString("key",
			mcp.Description("The key to retrieve from DiceDB"),
		),
	)

	// Define the 'set' tool to store a key-value pair
	setTool := mcp.NewTool("set",
		mcp.WithDescription("Set a value for a given key in DiceDB"),
		mcp.WithString("url",
			mcp.Description("The URL of the DiceDB server in format 'host:port'"),
			mcp.DefaultString("localhost:7379"),
		),
		mcp.WithString("key",
			mcp.Description("The key to set in DiceDB"),
		),
		mcp.WithString("value",
			mcp.Description("The value to associate with the key"),
		),
	)

	// Register tools and their handlers with the server
	s.AddTool(pingTool, pingHandler)
	s.AddTool(getTool, getHandler)
	s.AddTool(setTool, setHandler)

	// Start the MCP server using stdio transport
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// Handler for the 'ping' tool
func pingHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url, ok := req.Params.Arguments["url"].(string)
	if !ok || url == "" {
		return nil, errors.New("invalid or missing 'url' parameter")
	}

	conn, err := connectToDiceDB(url)
	if err != nil {
		return nil, fmt.Errorf("connection error: %v", err)
	}
	defer conn.Close()

	// Send a PING command to DiceDB
	if err := conn.Write(wire.NewCommand("PING")); err != nil {
		return nil, fmt.Errorf("failed to send PING: %v", err)
	}

	// Read the response
	resp, err := conn.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	return mcp.NewToolResultText(fmt.Sprintf("DiceDB response: %s", resp.String())), nil
}

// Handler for the 'get' tool
func getHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url, ok := req.Params.Arguments["url"].(string)
	if !ok || url == "" {
		return nil, errors.New("invalid or missing 'url' parameter")
	}

	key, ok := req.Params.Arguments["key"].(string)
	if !ok || key == "" {
		return nil, errors.New("invalid or missing 'key' parameter")
	}

	conn, err := connectToDiceDB(url)
	if err != nil {
		return nil, fmt.Errorf("connection error: %v", err)
	}
	defer conn.Close()

	// Send a GET command to DiceDB
	if err := conn.Write(wire.NewCommand("GET", key)); err != nil {
		return nil, fmt.Errorf("failed to send GET: %v", err)
	}

	// Read the response
	resp, err := conn.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	return mcp.NewToolResultText(fmt.Sprintf("Value for key '%s': %s", key, resp.String())), nil
}

// Handler for the 'set' tool
func setHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url, ok := req.Params.Arguments["url"].(string)
	if !ok || url == "" {
		return nil, errors.New("invalid or missing 'url' parameter")
	}

	key, ok := req.Params.Arguments["key"].(string)
	if !ok || key == "" {
		return nil, errors.New("invalid or missing 'key' parameter")
	}

	value, ok := req.Params.Arguments["value"].(string)
	if !ok {
		return nil, errors.New("invalid or missing 'value' parameter")
	}

	conn, err := connectToDiceDB(url)
	if err != nil {
		return nil, fmt.Errorf("connection error: %v", err)
	}
	defer conn.Close()

	// Send a SET command to DiceDB
	if err := conn.Write(wire.NewCommand("SET", key, value)); err != nil {
		return nil, fmt.Errorf("failed to send SET: %v", err)
	}

	// Read the response
	resp, err := conn.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set key '%s' to value '%s': %s", key, value, resp.String())), nil
}

// Helper function to connect to DiceDB
func connectToDiceDB(url string) (*dicedb.Conn, error) {
	// Split the URL into host and port
	parts := strings.Split(url, ":")
	if len(parts) != 2 {
		return nil, errors.New("invalid URL format; expected 'host:port'")
	}

	host := parts[0]
	port := parts[1]

	// Establish a connection to DiceDB
	return dicedb.Dial(host, port)
}
