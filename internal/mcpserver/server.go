package mcpserver

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func CreateDocTool() mcp.Tool {
	return mcp.NewTool("doc",
		mcp.WithDescription("แสดงเอกสารการใช้งาน MCP Server"),
	)
}

// DocHandler เป็นฟังก์ชันสำหรับจัดการคำขอ /doc เพื่อแสดงเอกสารการใช้งานของเซิร์ฟเวอร์
func DocHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	documentation := `MCPServer Sample Documentation

เครื่องมือที่มีให้ใช้งาน:
- login: ใช้สำหรับล็อกอินและรับ token
    พารามิเตอร์: base_url, email, password
- visitor_count: สำหรับแสดงจำนวนผู้เข้าชม
    พารามิเตอร์: base_url
- note: สำหรับดึงข้อมูลบันทึกตาม ID (dynamic resource)
    รูปแบบ: note://{id}
    พารามิเตอร์: base_url, token
- doc: แสดงเอกสารการใช้งาน MCP Server
`
	return mcp.NewToolResultText(documentation), nil
}

func CreateServer() *server.MCPServer {

	s := server.NewMCPServer(
		"MCPServerSample",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
	)

	loginTool := CreateLoginTool()
	s.AddTool(loginTool, LoginHandler)

	visitorTool := CreateVisitorCountTool()
	s.AddTool(visitorTool, VisitorCountHandler)

	noteTool := CreateGetNoteTool()
	s.AddTool(noteTool, GetNoteHandler)

	docTool := CreateDocTool()
	s.AddTool(docTool, DocHandler)

	return s
}
