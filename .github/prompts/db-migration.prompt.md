# Database Migration Generation

Your goal is to generate database migration files following the project's standards.

Requirements:

- Use GORM migrations
- Include both up and down migrations
- Follow Postgres best practices
- Consider database constraints and indexes
- Handle large tables appropriately

Rules:

- Always include transaction support
- Add proper documentation for migration steps
- Include rollback procedures
- Consider data preservation
- Follow naming convention: YYYYMMDDHHMMSS_description.go

Reference existing migrations in cmd/migrate/main.go
