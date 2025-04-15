# New API Endpoint Generation

Your goal is to generate a new API endpoint following the project's architecture and testing standards.

Requirements:

- Follow the layered architecture: handler -> service -> repository.
- Handler layer must handle all validations and auth checks.
- Use proper dependency injection pattern.
- Include GORM repository implementation.
- Define interfaces in the same file as the implementing struct, prefixed with 'I'.
- Generate corresponding test files using `testify/suite`.
- Use `testify/assert` for assertions in tests.
- Generate mocks using `uber-go/mock` with `//go:generate` comments in interface files.
- Update router configuration in `/internal/router/router.go`.

Required components:

1. Model struct (`/internal/models`)
2. Repository interface and implementation (`/internal/repository`)
3. Service interface and implementation (`/internal/service`)
4. Handler with validation (`/internal/handler/{driver}handler`)
5. Router configuration update
6. Unit tests using `testify/suite`, `testify/assert`, and `uber-go/mock`

Architecture rules:

- Validation completes at handler layer.
- Use dependency injection initialized in `cmd/{appname}/main.go`.
- Follow existing patterns in `/internal/{handler,service,repository}`.
- Add proper error handling and logging (`uber-go/zap`) at each layer.
- Place interfaces in the same file as their corresponding struct.
