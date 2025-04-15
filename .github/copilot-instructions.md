# Project Technology Stack and Coding Standards

## Go Backend Development

- Always use the latest stable Go version. Check <https://go.dev/doc/devel/release> for updates
- Use idiomatic Go code and follow official Go style guidelines
- Follow the Standard Go Project Layout: <https://github.com/golang-standards/project-layout>
- All packages must have proper documentation and error handling
- Use `github.com/go-resty/resty` for HTTP client interactions
- Use `github.com/uber-go/zap` for logging
- All code comments must be written in English only

## Technology Stack

- Postgres as the primary database (check <https://hub.docker.com/_/postgres> for latest version)
- GORM for database operations and ORM
- Echo framework for HTTP routing and middleware (Backend API)
- JWT for authentication (Backend API)
- Next.js (App Router, TypeScript) and Tailwind CSS for frontend development.
  - Frontend will be configured for **Static HTML Export** (`output: 'export'` in `next.config.js`).
  - Dynamic data interactions (login, notes, profile) will use client-side JavaScript to fetch data from the Go backend API.

## Project Structure Requirements

- Backend applications must be placed in `/cmd/{appname}` (e.g., /cmd/api, /cmd/migrate).
- Frontend application (Next.js) should be placed in `/frontend` directory.
- Follow clear layer separation (3-Tier) for the backend: handler -> service -> repository with dependency injection.
- Handler layer (Backend): Responsible for receiving input from drivers (e.g., HTTP), validating input, and handling Authentication/Authorization. Structure handlers by type, e.g., `/internal/handler/{httphandler}`.
- Service layer (Backend): Contains the core business logic. Should not access data storage or external services directly, but only through repository interfaces.
- Repository layer (Backend): Implements Facade pattern to handle all data storage logic and external services access. A repository may combine multiple storage technologies (e.g., database + cache + file storage) for a domain-specific entity. For example, a UserRepository might access both Postgres for user data and MinIO for profile images, keeping service layer isolated from implementation details.
- Validation, authentication, and request header to context management must be handled in the `/internal/handler` layer (Backend).
- Use singular form for folder names (model, config, etc.).
- Interface Naming and Location (Backend): Define interfaces in the same file as the struct they describe. Prefix interface names with 'I' based on the implementing struct (e.g., `IUserRepository` for `UserRepository`, `IUserService` for `UserService`).
- Common libraries for the backend go in `/pkg` directory. These should be purely technical utilities and helpers without any business logic. Business-specific features (such as visitor counting) should be implemented in the appropriate business layer (repository/service).
- Middleware Exception: While most business logic should be in the appropriate business layer, HTTP middleware that enforces business rules (such as business hours restrictions, rate limiting by business tier, or similar cross-cutting concerns) can be placed in `/pkg/middleware`. This is an exception to the general rule, as these middleware components can be reusable across different parts of the application.

## Docker & Deployment

- Use multi-stage Dockerfile for optimized builds (for both backend and potentially the static frontend server like Nginx).
- The output of the frontend build (`next build` with `output: 'export'`) will be the static files located in the `/frontend/out` directory. These files are ready to be deployed to a static web server (e.g., Nginx, Apache, CDN, static hosting providers).
- Docker image versions must be specified with major.minor (e.g., 1.2), without patch version. The 'latest' tag should only be used if no specific version tags are available for an image.
- Environment configs in .env file with support for sit/uat/prod environments.

## Code Quality & Standards

- All new code must have corresponding tests.
- Go Testing:
  - Use `github.com/stretchr/testify/assert` for assertions.
  - Use `github.com/stretchr/testify/suite` for organizing test suites.
  - Use `github.com/uber-go/mock` for generating mocks. Include `//go:generate mockgen -source=./your_interface_file.go -destination=./mocks/mock_your_interface.go -package=mocks` comment in the interface file.
- Use GORM for database operations with proper transaction handling.
- Echo framework for HTTP routing with proper middleware chain.
- JWT for authentication with secure practices.
- Makefile targets must have individual .PHONY declarations.
- All comments in code must be in English only, following Go standard documentation format.

## Development Workflow

- Always test code changes before deployment.
- Document APIs and additional information in /doc directory.
- Frontend (Next.js with Static Export): Follow Next.js best practices. Client-side JavaScript will handle API calls to the Go backend for dynamic features.

## Architecture Patterns

- Backend: Strict layered architecture (3-Tier: handler -> service -> repository) with clear dependency injection.
- Backend: Request/response validation completes at handler layer.
- Backend: Middleware and utilities belong in /pkg.
- Backend: Repository pattern for data access with Facade design pattern implementation.
- Backend: Service layer for business logic and orchestration, never directly accessing storage.
- Frontend: Interacts with the backend API via client-side requests.
