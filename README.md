# MCP Server Demo

## Project Structure

```sh
.
├── cmd/                  # Main applications
│   ├── api/              # API server entry point
│   ├── migrate/          # Database migration tool
│   └── mcpserver/        # MCP server implementation for AI interaction
├── configs/              # Configuration files (nginx, api .env)
├── frontend/             # Next.js frontend application
│   ├── app/              # Next.js App Router pages and layouts
│   ├── public/           # Static assets for frontend
│   ├── src/              # Frontend source code (components, hooks, utils, etc.)
│   ├── next.config.mjs   # Next.js configuration (configured for static export)
│   ├── package.json
│   ├── tailwind.config.ts
│   └── tsconfig.json
├── internal/
│   ├── handler/          # Request handlers (HTTP)
│   ├── mcpserver/        # MCP server business logic and tools
│   ├── migrations/       # Database migrations
│   ├── models/           # Database models (structs)
│   ├── repository/       # Data access layer (database, storage interactions)
│   ├── router/           # API route definitions
│   └── service/          # Business logic layer
├── pkg/
│   ├── cache/            # Cache management utilities
│   ├── database/         # Database connection helper
│   ├── middleware/       # Echo middleware (JWT, roles)
│   ├── storage/          # File storage helpers (e.g., MinIO)
│   └── validator/        # Request validation helper
├── .env                  # Local environment variables (copied from configs/api/.env)
├── .gitignore
├── Dockerfile            # Multi-stage Docker build for backend
├── docker-compose.yml    # Docker Compose setup (Postgres, PgAdmin, Backend App, Nginx)
├── go.mod
├── go.sum
├── Makefile              # Make commands for development and deployment
└── README.md
```

## Getting Started

### Prerequisites

- Go (latest stable version recommended)
- Docker & Docker Compose
- Node.js & npm (for frontend development)
- Make

### Setup

#### 1.  **Clone the repository:**

```bash
git clone <repository-url>
cd mcpserver-sample
```

#### 2.  **Run setup command:** This will clean previous backend builds, copy the backend `.env` file (if it doesn't exist), install frontend dependencies, and tidy Go modules

```bash
make setup
```

#### 3.  **Review Environment Variables:**

Check the `.env` file (copied from `configs/api/.env`) and adjust database credentials, JWT secret, MinIO settings, etc. if necessary.

### Running the Backend

- **Run Backend API with Docker (Recommended for Full Environment):**

```bash
make up
```

The API will be available at `http://localhost:8080` (or the port defined in `.env`). Requires a running Postgres database accessible from your local machine.

- **Stop Docker Services:**

```bash
make down
```

### Running the Frontend (Next.js)

- **Run Frontend Development Server:**

```bash
make frontend-dev
```

Access the frontend at `http://localhost:3000` (or the next available port).
The development server supports hot-reloading.

- **Build Static Frontend:**

```bash
make frontend-build
```

This generates the static export in the `frontend/out` directory. These are the files you would deploy to a static web server.

- **Lint Frontend Code:**

```bash
make frontend-lint
```

## ทดลองใช้ MCP Server

ตัวอย่างการใช้งาน MCP (Model Context Protocol) Server สำหรับ PostgreSQL โดยใช้ Go และ Echo Framework ให้สามารถสื่อสารกับ Large Language Model ได้

## Project setup

1. Clone โปรเจค

   ```sh
   git clone https://github.com/Napat/mcpserver-demo.git
   cd mcpserver-sample
   ```

2. สร้างและรันคอนเทนเนอร์

   ```sh
   make up  # deploy all services(backend + frontend) with auto migration
   ```

3. เข้าถึง PGAdmin ที่ <http://localhost:5050>
   - อีเมล: <admin@example.com>
   - รหัสผ่าน: admin123

   Connect to database with: postgres/postgres
   MCP PostgreSQL Server > Database > mcpserver > Schemas > public > Tables > users > Right click > View/Edit Data > All Rows

## การใช้งาน MCP

ในที่นี้ใช้ VSCode เป็นหลัก แต่สามารถประยุกต์ใช้ใน client อื่นๆ เช่น Cursor, etc. ได้ไม่ต่างกัน

### MCP server: time + postgres

- ref: <https://github.com/modelcontextprotocol/servers/tree/main/src/postgres>

จากเอกสาร ref จะเห็นว่า mcp server postgres จะเพิ่ม tools `#query` ให้เราใช้งาน

ติดตั้ง mcp server time + postgres ใน vscode

`cmd+shift+P` then `>Preferences: Open User Settings (JSON)`

```json
    "mcp": {
        "inputs": [],
        "servers": {
            "time": {
                "command": "docker",
                "args": [
                    "run", "-i", "--rm", "mcp/time", "--local-timezone", "Asia/Bangkok"
                ]
            },
            "postgres": {
                "command": "docker",
                "args": [
                    "run", 
                    "-i", 
                    "--rm", 
                    "mcp/postgres", 
                    "postgresql://postgres:postgres@host.docker.internal:5432/mcpserver"]
            }
        }
    },
    "chat.mcp.discovery.enabled": true 
```

ตัวอย่างการใช้งาน #query เพื่อสอบถาม mcp server postgres

```chat
#query ตรวจสอบ schema ของ postgres database และแสดงผลออกมา
```

```chat
#query ตรวจสอบ code ใน repository ของ project นี้ว่าเข้ากันได้กับ database schema ที่ทำงานอยู่จริง ปรับปรุง code และให้แนะนำเพิ่มเติม เช่นการตั้งค่า index
```

```chat
#query ปัจจุบันมีผู้ใช้อยู่กี่คน
```

```chat
user@example.com เข้าสู่ระบบครั้งสุดท้ายเมื่อไหร่ ช่วยแปลงเป็นเวลาไทยให้ด้วยนะ

-----AI-----
user@example.com เข้าสู่ระบบครั้งสุดท้ายเมื่อวันที่ 12 เมษายน 2025 เวลา 17:58:50 น. ตามเวลาประเทศไทย (หรือประมาณ 19 ชั่วโมงที่แล้ว)
```

```chat
#query ลองดูรายละเอียดของผู้ใช้แต่ละคนพร้อมบทบาท

-----AI-----
จากผลการตรวจสอบพบว่าปัจจุบันมีผู้ใช้ในระบบทั้งหมด 5 คน แบ่งตามบทบาทดังนี้:

ผู้ใช้ทั่วไป (RoleUser = 1): 1 คน
พนักงาน (RoleStaff = 2): 1 คน
ผู้จัดการ (RoleManager = 4): 1 คน
แอดมิน (RoleAdmin = 8): 1 คน
ซูเปอร์แอดมิน (RoleSuperAdmin = 16): 1 คน
ทุกบัญชีมีสถานะ active = true และยังไม่มีการตั้งค่ารูปโปรไฟล์ (profile_image_url = null)
```

### MCP server: playwright

ref: <https://github.com/microsoft/playwright-mcp>

ติดตั้ง mcp server time + postgres ใน vscode

`cmd+shift+P` then `>Preferences: Open User Settings (JSON)`

```json
    "mcp": {
        "inputs": [],
        "servers": {
            "playwright": {
                "command": "npx",
                "args": [
                    "@playwright/mcp@latest"
                ]
            }
        }
    },
    "chat.mcp.discovery.enabled": true 
```

ตัวอย่างการใช้งาน

```chat
- เปิด browser หน้า http://localhost:8001/login
- ในช่อง email กรอก login Email: "user@example.com"
- ในช่อง password กรอก Password: "user123"
- กดปุ่ม login
```

```chat
ตรวจสอบการแสดงผล http://localhost:8001/ ว่าแสดงจำนวนผู้เข้าชมถูกต้องหรือไม่ ตัวเลขเท่ากับเราเรียก GET /api/visitors เองรึเปล่า แจ้งตัวเลขที่แสดงในหน้าเว็บกับตัวเลขที่ได้จาก api มาให้ทราบหน่อย
```

ลองทำ ui automation test ด้วย playwright

FLOW: Register new user + Login + Logout

```chat
ช่วยเปิดเว็บและทดสอบลงทะเบียน user ใหม่ให้หน่อย โดยทดสอบปุ่มลงทะเบียนทั้งสองตำแหน่ง step คร่าวๆคือดังนี้

1. เปิดหน้าแรกที่ url http://localhost:8001/
2. ถ้าเห็นว่ามีปุ่ม "ออกจากระบบ" หรือ "Logout" ให้กดเพื่อออกจากระบบและรอจนกว่าจะกลับมาหน้าแรกที่ยังไม่ได้ล็อกอิน
3. ตรวจสอบว่ามีปุ่มลงทะเบียนกี่ตำแหน่งในหน้าแรก (บนเมนูและในส่วนเนื้อหา)
4. ทดสอบกดปุ่ม "ลงทะเบียน" บนเมนูนำทาง
5. กรอกข้อมูลใน register form ให้ครบถ้วน โดยสร้างข้อมูลที่สมจริง เช่น:
   - Email: ใช้รูปแบบอีเมลที่ถูกต้อง
   - Password: ใช้รหัสผ่านที่มีความปลอดภัย
   - ชื่อ-นามสกุล: ใช้ชื่อภาษาไทยให้สมจริง
   - เพศ: เลือกตัวเลือกใดตัวเลือกหนึ่ง
6. บันทึกข้อมูลที่กรอกทั้งหมดไว้และแจ้งให้ผมทราบ
7. กดปุ่มลงทะเบียนและสังเกตว่ามีการรีไดเรคไปที่หน้าใด
8. ทดสอบฟีเจอร์หลักจากหน้าที่รีไดเรคไป เช่น การดูโปรไฟล์ หรือการดูบันทึก
9. กดปุ่ม "ออกจากระบบ" หรือ "Logout"
10. ล็อกอินกลับเข้าระบบด้วยข้อมูลที่ลงทะเบียนไว้
11. ตรวจสอบว่าข้อมูลโปรไฟล์ตรงกับที่ลงทะเบียน และสามารถใช้งานฟีเจอร์ต่างๆ ได้
12. ถ้ามีเวลา ให้กลับไปหน้าแรก ล็อกเอาท์ และทดสอบซ้ำโดยใช้ปุ่ม "ลงทะเบียน" อีกตำแหน่งหนึ่ง (ในส่วนเนื้อหา)

โปรดรายงานสิ่งที่พบ ความแตกต่างระหว่างการใช้ปุ่มลงทะเบียนทั้งสองตำแหน่ง และประเด็นที่อาจเป็นปัญหาหรือจุดที่ควรปรับปรุง (ถ้ามี)
```

FLOW: Add/Edit/Delete Note

```chat
ใช้ mcp ช่วยทดสอบการทำงานของฟีเจอร์การเพิ่มบันทึก (Note) ตามขั้นตอนดังนี้:

1. เปิดเว็บไซต์ที่ http://localhost:8001/
2. ถ้ายังไม่ได้ล็อกอิน ให้ล็อกอินด้วยบัญชีที่มีอยู่แล้ว หรือลงทะเบียนบัญชีใหม่
3. เมื่อล็อกอินแล้ว ไปที่หน้า Notes โดยคลิกที่เมนู "Notes" หรือลิงก์ "บันทึกของฉัน"
4. ตรวจสอบว่ามีปุ่ม "เพิ่มบันทึกใหม่" หรือไม่ และคลิกที่ปุ่มนี้
5. ทดสอบกรอกข้อมูลบันทึกใหม่:
   - ใส่หัวข้อบันทึก (Title): ใช้ชื่อที่มีความหมาย เช่น "ทดสอบการเพิ่มบันทึกใหม่"
   - ใส่เนื้อหาบันทึก (Content): เขียนข้อความที่มีความยาวพอสมควร อย่างน้อย 2-3 ประโยค
   - เลือกประเภท/แท็ก (หากมี)
6. กดปุ่มบันทึก หรือ "สร้างบันทึก"
7. ตรวจสอบว่าระบบแสดงข้อความยืนยันการสร้างบันทึกสำเร็จหรือไม่
8. ตรวจสอบว่าระบบพาไปที่หน้ารายการบันทึกหรือไม่ และบันทึกใหม่ปรากฏในรายการหรือไม่
9. คลิกที่บันทึกที่เพิ่งสร้างเพื่อดูรายละเอียด
10. ตรวจสอบว่าข้อมูลที่แสดงตรงกับที่กรอกไปหรือไม่ (หัวข้อ, เนื้อหา, วันที่สร้าง)
11. กลับไปที่หน้ารายการบันทึก และตรวจสอบว่ามีฟีเจอร์การจัดการบันทึก (เช่น ลบ, แก้ไข) หรือไม่
12. ทดสอบการค้นหาบันทึก (ถ้ามี) โดยใช้คำที่อยู่ในหัวข้อหรือเนื้อหาของบันทึกที่เพิ่งสร้าง

โปรดแจ้งผลการทดสอบแต่ละขั้นตอน รวมถึงสิ่งที่พบระหว่างการทดสอบ เช่น ข้อผิดพลาด หรือปัญหาในการใช้งาน (ถ้ามี) และข้อเสนอแนะในการปรับปรุง UI/UX
```

Senario: แก้ bug ที่ AI report bug ใน prompt ก่อนหน้า

```chat
แก้ไขปัญหา 
- การแสดงผลประวัติการล็อกอิน: เวลาที่แสดงในประวัติล็อกอินไม่ถูกต้อง แสดงเป็น "1/1/1, 6:42:04 AM" 
ก่อนแก้ไขให้ตรวจสอบก่อนว่าปัญหาเกิดจาก data ที่รับผ่าน api ของ backend หรือ frontend รับค่ามาแสดงผลไม่ถูกต้อง จะได้แก้ได้ถูกจุด
```

## MCP Server Development using Golang

การพัฒนา MCP Server โดยใช้ Golang สามารถดูตัวอย่างได้ที่นี่เลย [cmd/mcpserver](cmd/mcpserver/README.md) แยกออกไปเพราะเดี๋ยวจะยาวเกินไปสำหรับใน README นี้

---------------------------------------------------------------------------

## API Endpoints Spec(ส่วนเสริม)

### การยืนยันตัวตน

- `POST /api/auth/register` - ลงทะเบียนผู้ใช้ใหม่
- `POST /api/auth/login` - เข้าสู่ระบบ

### ผู้ใช้ทั่วไป

- `GET /api/me` - ดึงข้อมูลผู้ใช้ปัจจุบัน
- `PUT /api/me` - อัพเดทข้อมูลผู้ใช้ปัจจุบัน
- `POST /api/me/profile-image` - อัพโหลดรูปโปรไฟล์
- `GET /api/me/login-history` - ดึงประวัติการเข้าสู่ระบบ

### แอดมิน

- `GET /api/admin/users` - ดึงรายการผู้ใช้ทั้งหมด
- `POST /api/admin/users` - สร้างผู้ใช้ใหม่
- `GET /api/admin/users/:id` - ดึงผู้ใช้ด้วย ID
- `PUT /api/admin/users/:id` - อัปเดตผู้ใช้
- `DELETE /api/admin/users/:id` - ลบผู้ใช้
- `GET /api/admin/users/:id/login-history` - ดึงประวัติการเข้าสู่ระบบของผู้ใช้

## บทบาทของผู้ใช้

โปรเจคนี้มีการกำหนดบทบาทด้วย Union Type ซึ่งช่วยให้ผู้ใช้สามารถมีบทบาทหลายบทบาทได้พร้อมกัน:

- **RoleUser** (1): ผู้ใช้ทั่วไป
- **RoleStaff** (2): พนักงาน
- **RoleManager** (4): ผู้จัดการ
- **RoleAdmin** (8): แอดมิน
- **RoleSuperAdmin** (16): ซูเปอร์แอดมิน

ค่าตัวเลขในวงเล็บแสดงค่า bitwise ที่ใช้ในการกำหนดบทบาท ทำให้ผู้ใช้สามารถมีหลายบทบาทได้ เช่น ค่า 9 หมายถึงเป็นทั้งผู้ใช้ทั่วไปและแอดมิน (1 + 8)

## การจัดการรูปโปรไฟล์

ระบบรองรับการอัพโหลดและจัดการรูปโปรไฟล์ของผู้ใช้ด้วย MinIO:

- รองรับไฟล์รูปภาพ JPG, JPEG และ PNG
- จัดเก็บรูปภาพใน MinIO bucket "profiles"
- มีรูปโปรไฟล์ default สำหรับผู้ใช้ที่ยังไม่ได้อัพโหลดรูป
- รูปโปรไฟล์เดิมจะถูกลบเมื่อมีการอัพโหลดรูปใหม่

### การตั้งค่า MinIO

กำหนดค่าการเชื่อมต่อ MinIO ในไฟล์ `.env`:

```env
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_USE_SSL=false
```

เข้าถึง MinIO Console ได้ที่ <http://localhost:9001>

- Username: minioadmin
- Password: minioadmin

## Reference

- [MCP Server collection #1](https://mcp.so/)
- [MCP Server collection #2](https://github.com/modelcontextprotocol/servers)
- [Time MCP Server](https://github.com/modelcontextprotocol/servers/tree/main/src/time)
- [Postgres MCP Server](https://github.com/modelcontextprotocol/servers/tree/main/src/postgres)
