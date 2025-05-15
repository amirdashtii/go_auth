# Go Auth Service

A robust authentication and authorization service built with Go, following clean architecture principles.

## Features

- User authentication with phone number and password
- JWT token-based authentication
- Role-based access control (RBAC)
- Admin panel for user management
- Redis for token storage and OTP
- PostgreSQL for data persistence
- Comprehensive error handling with bilingual messages (English and Persian)
- Input validation
- Unit tests
- Logging to standard output (stdout)

## Project Structure

```
.
├── cmd/
│   └── main.go            # Application entry point
├── config/                # Configuration files (config.go, *.yaml, .env.example)
├── controller/            # HTTP handlers
│   ├── dto/               # Data Transfer Objects
│   ├── middleware/        # Gin middleware (auth, logger)
│   └── validators/        # Request validation logic
├── docs/                  # API documentation (Swagger/OpenAPI files: docs.go, swagger.json, swagger.yaml)
├── infrastructure/
│   ├── logger/            # Logging implementations (file, zerolog)
│   └── repository/        # Data persistence implementations (Postgres, Redis, InMemory)
├── internal/
│   └── core/
│       ├── entities/      # Core domain entities (user, token)
│       ├── errors/        # Custom error types and messages
│       ├── ports/         # Interfaces for services and repositories
│       └── service/       # Business logic (application services) and their mocks
│           └── mocks/     # Mock implementations for testing
├── migrations/            # Database migration scripts (.sql files)
├── docker-compose.yml     # Docker Compose configuration
├── go.mod                 # Go module definition
├── go.sum                 # Go module checksums
├── LICENSE                # Project license file
├── README.md              # This file
```

## Getting Started

### Prerequisites

- Go 1.21 or higher (check `go.mod` for the exact version, currently it is 1.21)
- PostgreSQL
- Redis

### Installation

1.  Clone the repository:

    ```bash
    git clone https://github.com/amirdashtii/go_auth.git # Replace with your actual repository URL if different
    cd go_auth
    ```

2.  Install dependencies:

    ```bash
    go mod download
    ```

3.  Set up configuration:
    The application's configuration is loaded with the following priority:

    1.  **Environment Variables (`.env` file):** Highest priority. Values here override YAML and default settings.
    2.  **YAML Configuration Files (`config/*.yaml`):** Middle priority. Loaded if corresponding environment variables are not set.
    3.  **Default Values (in code):** Lowest priority. Used if no configuration is provided via `.env` or YAML files.

    - **Environment Variables:**
      Copy the example environment file and customize it:

      ```bash
      cp config/.env.example .env
      ```

      Then, edit `.env` with your database credentials, Redis info, JWT secrets, server port, etc.

    - **YAML Configuration:**
      An example YAML configuration is provided: `config/development.yaml.example`. You can copy it to `config/development.yaml` (or other environment-specific names like `config/production.yaml`) and customize it.
      ```bash
      cp config/development.yaml.example config/development.yaml
      ```
      _Note: `config/_.yaml`files (except`_.example.yaml`files) are configured to be ignored by Git via`.gitignore`._

4.  Run the application:
    ```bash
    go run cmd/main.go
    ```

## API Endpoints

### Authentication (`/auth`)

- `POST /auth/register`: Register a new user.
  - Request Body: `dto.RegisterRequest`
  - Response: Success message or error.
- `POST /auth/login`: Login with phone number and password.
  - Request Body: `dto.LoginRequest`
  - Response: Access and refresh tokens or error.
- `POST /auth/logout`: Logout user (requires authentication).
  - Invalidates user's tokens.
  - Response: Success message or error.
- `POST /auth/refresh-token`: Refresh access token using a valid refresh token.
  - Request Body: `dto.RefreshTokenRequest`
  - Response: New access and refresh tokens or error.

### User Management (`/users`) - Authenticated User

All endpoints in this section require user authentication.

- `GET /users/me`: Get current authenticated user's profile.
  - Response: `dto.UserProfileResponse` or error.
- `PUT /users/me`: Update current authenticated user's profile.
  - Request Body: `dto.UserUpdateRequest`
  - Response: Success message or error.
- `PUT /users/me/change-password`: Change current authenticated user's password.
  - Request Body: `dto.ChangePasswordRequest`
  - Response: Success message or error.
- `DELETE /users/me`: Delete current authenticated user's profile.
  - Response: Success message or error.

### Admin User Management (`/users`) - Admin Only

All endpoints in this section require admin-level authentication and authorization. These admin endpoints are grouped under the `/users` path but are distinct due to admin middleware.

- `GET /users`: List all users (admin/super-admin only).
  - Query Parameters:
    - `status`: Filter by status (e.g., `active`, `inactive`). Default: `active`.
    - `role`: Filter by role (e.g., `user`, `admin`). Default: `user`.
    - `sort`: Sort field (e.g., `created_at`, `first_name`). Default: `created_at`.
    - `order`: Sort order (`asc` or `desc`). Default: `desc`.
  - Response: List of `dto.AdminUserResponse` or error.
- `GET /users/:id`: Get user details by ID (admin/super-admin only).
  - Path Parameter: `id` (User UUID)
  - Response: `dto.AdminUserResponse` or error.
- `PUT /users/:id`: Update user details by ID (admin/super-admin only).
  - Path Parameter: `id` (User UUID)
  - Request Body: `dto.AdminUserUpdateRequest`
  - Response: Success message or error.
- `PUT /users/:id/role`: Change a user's role by ID (admin/super-admin only).
  - Path Parameter: `id` (User UUID)
  - Request Body: `dto.ChangeRoleRequest`
  - Response: Success message or error.
- `PUT /users/:id/status`: Change a user's status by ID (admin/super-admin only).
  - Path Parameter: `id` (User UUID)
  - Request Body: `dto.ChangeStatusRequest`
  - Response: Success message or error.
- `DELETE /users/:id`: Delete a user by ID (admin/super-admin only).
  - Path Parameter: `id` (User UUID)
  - Response: Success message or error.

## Error Handling

The service implements a comprehensive error handling system with:

- Custom error types defined in `internal/core/errors`.
- Bilingual error messages (English and Persian).
- Proper HTTP status codes for API responses.
- Detailed error responses in JSON format.

## Logging

Application logs are written to standard output (stdout) in JSON format (powered by Zerolog). This facilitates easy log collection and processing by containerization platforms (like Docker, Kubernetes) or external log management systems.

## Testing

Run tests with verbose output:

```bash
go test ./... -v
```

## Contributing

1.  Create a new branch for your feature (e.g., `git checkout -b feature/your-feature-name`).
2.  Make your changes and add/update tests accordingly.
3.  Ensure all tests pass: `go test ./... -v`.
4.  Format your code: `go fmt ./...`.
5.  Consider linting your code if a linter (e.g., `golangci-lint`) is part of the project setup.
6.  Submit a pull request with a clear description of your changes.

## License

This project is licensed under the MIT License - see the `LICENSE` file for details. (Ensure a `LICENSE` file with MIT content exists in your project root).
