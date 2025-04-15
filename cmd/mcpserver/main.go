package main

import (
	"fmt"

	"github.com/Napat/mcpserver-demo/internal/mcpserver"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// สร้าง MCP server จาก package mcpserver
	s := mcpserver.CreateServer()

	// เริ่มการทำงานของ server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
