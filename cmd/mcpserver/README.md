# MCPServer Sample

Model Context Protocol (MCP) เป็นโปรโตคอลเปิดที่ช่วยเชื่อมต่อโมเดล AI (LLM) กับแหล่งข้อมูลและเครื่องมือภายนอก โดยเฉพาะเมื่อใช้ไลบรารี [Go mcp-go จาก mark3labs](https://github.com/mark3labs/mcp-go)

ในปัจจุบันที่กำลังเขียน MCP สามารถแยกประเภทย่อยๆได้เป็น Tool และ Resource เป็นหลัก (อาจจะมีเพิ่มในภายหลังได้ เช่น ประเภท Prompt หรืออื่นๆ)

- Tool: เป็นฟังก์ชันหรือการดำเนินการที่โมเดล AI สามารถเรียกใช้เพื่อประมวลผลหรือทำงานบางอย่าง มีผลข้างเคียง เช่น การคำนวณ, การเขียนไฟล์, หรือการส่งคำขอ API ที่ต้องการการประมวลผล ตัวอย่างเช่น การคำนวณทางคณิตศาสตร์, การสร้างไฟล์ใหม่, หรือการเรียกใช้ API ที่ต้องการการดำเนินการ

- Resource: เป็นแหล่งข้อมูลที่ใช้สำหรับการจัดหาข้อมูลให้กับโมเดล AI หรือผู้ใช้ มีลักษณะคล้ายกับ "GET endpoint" ใน web service คือใช้สำหรับการดึงข้อมูลโดยไม่มีผลข้างเคียง (side effects) ตัวอย่างเช่น การอ่านไฟล์, การดึงข้อมูลจาก API, หรือการเข้าถึงฐานข้อมูล

วิธีการแยกแยะ:

- ถ้าเป็นการดำเนินการที่ต้องการการประมวลผลหรือมีผลข้างเคียง(Create, Update, Delete): ถือเป็น tool
- ถ้าเป็นการจัดหาข้อมูลโดยไม่มีการประมวลผลหรือผลข้างเคียง(Read): ถือเป็น resource

Resource สามารถแบ่งประเภทได้เป็น 2 ประเภทหลักๆ คือ:

- `Static Fixed URI`: เป็น URI ที่คงที่ ไม่มีส่วนที่เปลี่ยนแปลงได้ ใช้สำหรับการเข้าถึงข้อมูลที่กำหนดไว้ เช่น `docs://readme` ซึ่งหมายถึงการเข้าถึงไฟล์ README โดยตรง
- `Dynamic URI Templates`: เป็น URI ที่มีส่วนที่สามารถเปลี่ยนแปลงได้โดยใช้ template เช่น `users://{id}/profile` ซึ่ง {id} เป็นพารามิเตอร์ที่สามารถถูกแทนที่ด้วยค่าจริง เช่น `users://123/profile` สำหรับการเข้าถึงข้อมูลโปรไฟล์ของผู้ใช้

**CAUTION** อย่างไรก็ตาม MCP Client ส่วนใหญ่ในปัจจุบันยังรองรับแค่แบบ Tool เท่านั้น  
สามารถติดตามการอัปเดทได้ที่ <https://modelcontextprotocol.io/clients#feature-support-matrix>  
ถ้า client ส่วนใหญ่ในตลาดยังรองรับแค่แบบ Tool เท่านั้น ก็เขียนทุกอย่างไว้ใน tool กันไปก่อนนะ 😅 (ตอนแรกเขียนแยกเป็น resources ออกไปเรียบร้อยแล้วเลยต้องมานั่งแก้ให้เป็น tool ใหม่เพราะ client หลักๆยังใช้ไม่ได้ซักกะตัว 💀💀💀)

## การติดตั้ง

```bash
# Clone repository
git clone https://github.com/Napat/mcpserver-demo.git
cd mcpserver-sample

# Build MCP Server
cd cmd/mcpserver
make mcp-docker-build

# Optionally, push to Docker Hub
make mcp-docker-push

```

## ทดสอบการใช้งานด้วย mcphost

**NOTE** ผลลัพธ์ขึ้นอยู่กับเอา model ไหนมาใช้งานนะ ขึ้นกับงบประมาณของแต่ละคนเลย

```bash
# Run MCP Server
make mcp-docker-run

-----------------
Enter your prompt: /help
> ...

Enter your prompt: /servers
> ...

Enter your prompt: /tools
> ...

Enter your prompt: Hi
> ...

Enter your prompt: สวัสดี ทดสอบ mcp server login ด้วย base url http://host.docker.internal:8080 ด้วยอีเมล "user@example.com" และรหัสผ่าน "user123" แล้วบอก response เช่น token ที่ได้รับมาให้ผม
> ได้รับ token เรียบร้อยแล้วค่ะ: eyJhbGciOiJIUzI1NiIsInR5c...

Enter your prompt: ช่วยหาจำนวณ visitor ที่เข้ามาใช้งานระบบให้หน่อย base url คือ http://host.docker.internal:8080 
> จำนวนผู้เข้าชมระบบตอนนี้คือ 19 คนค่ะ

Enter your prompt: ช่วยอ่าน note id=1 ให้หน่อยสิ
> ...

```

## References

- [Model Context Protocol (MCP)](https://github.com/mark3labs/mcp-go)
- [mcphost](https://github.com/mark3labs/mcphost)
- [Github ใช้ Golang ทำ MCP Server](https://github.blog/changelog/2025-04-04-github-mcp-server-public-preview/)
- [github-mcp-server](https://github.com/github/github-mcp-server)
