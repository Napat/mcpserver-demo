package mcpserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// LoginResponse คือโครงสร้างสำหรับข้อมูล response จากการล็อกอิน
type LoginResponse struct {
	Token string `json:"token"`
}

// Note คือโครงสร้างสำหรับข้อมูลบันทึก
type Note struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// VisitorResponse คือโครงสร้างสำหรับข้อมูลจำนวนผู้เข้าชม
type VisitorResponse struct {
	Count int `json:"visitor_count"`
}

// สร้าง Tool สำหรับการล็อกอิน
func CreateLoginTool() mcp.Tool {
	return mcp.NewTool("login",
		mcp.WithDescription("Login to the API and get a token"),
		mcp.WithString("base_url",
			mcp.Required(),
			mcp.Description("Base URL of the API (e.g., http://localhost:8001)"),
		),
		mcp.WithString("email",
			mcp.Required(),
			mcp.Description("Email for login"),
		),
		mcp.WithString("password",
			mcp.Required(),
			mcp.Description("Password for login"),
		),
	)
}

// LoginHandler เป็นฟังก์ชันสำหรับจัดการคำขอล็อกอิน
func LoginHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	baseURL, ok := request.Params.Arguments["base_url"].(string)
	if !ok {
		return nil, errors.New("base_url must be a string")
	}

	email, ok := request.Params.Arguments["email"].(string)
	if !ok {
		return nil, errors.New("email must be a string")
	}

	password, ok := request.Params.Arguments["password"].(string)
	if !ok {
		return nil, errors.New("password must be a string")
	}

	// ตัดเครื่องหมาย / ถ้ามีที่ท้าย baseURL
	baseURL = strings.TrimSuffix(baseURL, "/")

	// สร้าง HTTP client และส่ง request
	loginURL := fmt.Sprintf("%s/api/auth/login", baseURL)
	payload := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)

	// สร้าง HTTP client ที่มีการตั้งค่า timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("POST", loginURL, strings.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// เพิ่ม header
	req.Header.Set("Content-Type", "application/json")

	// ส่ง request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("login request failed: %v", err)
	}
	defer resp.Body.Close()

	// อ่าน response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(body))
	}

	// แปลง response เป็น struct
	var loginResp LoginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if loginResp.Token == "" {
		return nil, errors.New("no token received in response")
	}

	// ส่งคืนผลลัพธ์เป็น text แทน JSON
	return mcp.NewToolResultText(fmt.Sprintf("Token: %s", loginResp.Token)), nil
}

// สร้าง Tool สำหรับดึงจำนวนผู้เข้าชม
func CreateVisitorCountTool() mcp.Tool {
	return mcp.NewTool("get_visitor_count",
		mcp.WithDescription("Get the current visitor count"),
		mcp.WithString("base_url",
			mcp.Required(),
			mcp.Description("Base URL of the API (e.g., http://localhost:8001)"),
		),
	)
}

// VisitorCountHandler เป็นฟังก์ชันสำหรับดึงจำนวนผู้เข้าชม
func VisitorCountHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	baseURL, ok := request.Params.Arguments["base_url"].(string)
	if !ok {
		return nil, errors.New("base_url must be a string")
	}

	// ตัดเครื่องหมาย / ถ้ามีที่ท้าย baseURL
	baseURL = strings.TrimSuffix(baseURL, "/")

	// เรียก API เพื่อดึงจำนวนผู้เข้าชม
	visitorURL := fmt.Sprintf("%s/api/visitors", baseURL)
	resp, err := http.Get(visitorURL)
	if err != nil {
		return nil, fmt.Errorf("visitor count request failed: %v", err)
	}
	defer resp.Body.Close()

	// อ่าน response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("visitor count request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// แปลง response เป็น struct
	var visitorResp VisitorResponse
	if err := json.Unmarshal(body, &visitorResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	// ส่งคืนผลลัพธ์เป็น text
	return mcp.NewToolResultText(fmt.Sprintf("Visitor count: %d", visitorResp.Count)), nil
}

// สร้าง Tool สำหรับดึงข้อมูลบันทึกตาม ID
func CreateGetNoteTool() mcp.Tool {
	return mcp.NewTool("get_note",
		mcp.WithDescription("Get a note by ID"),
		mcp.WithString("base_url",
			mcp.Required(),
			mcp.Description("Base URL of the API (e.g., http://localhost:8001)"),
		),
		mcp.WithString("token",
			mcp.Required(),
			mcp.Description("JWT token for authentication"),
		),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("ID of the note to retrieve"),
		),
	)
}

// GetNoteHandler เป็นฟังก์ชันสำหรับดึงข้อมูลบันทึกตาม ID
func GetNoteHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	baseURL, ok := request.Params.Arguments["base_url"].(string)
	if !ok {
		return nil, errors.New("base_url must be a string")
	}

	token, ok := request.Params.Arguments["token"].(string)
	if !ok {
		return nil, errors.New("token must be a string")
	}

	id, ok := request.Params.Arguments["id"].(string)
	if !ok {
		return nil, errors.New("id must be a string")
	}

	// ตัดเครื่องหมาย / ถ้ามีที่ท้าย baseURL
	baseURL = strings.TrimSuffix(baseURL, "/")

	// สร้าง HTTP request เพื่อดึงข้อมูล note
	noteURL := fmt.Sprintf("%s/api/notes/%s", baseURL, id)
	req, err := http.NewRequest("GET", noteURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// เพิ่ม Authorization header
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	// ส่ง request
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("note request failed: %v", err)
	}
	defer resp.Body.Close()

	// อ่าน response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("note request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// แปลง response เป็น struct
	var note Note
	if err := json.Unmarshal(body, &note); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	// ส่งคืนผลลัพธ์เป็น JSON
	noteJSON, err := json.MarshalIndent(note, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal note: %v", err)
	}

	return mcp.NewToolResultText(string(noteJSON)), nil
}
