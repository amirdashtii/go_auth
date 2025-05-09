# Go Auth Service

A robust authentication and authorization service built with Go, following clean architecture principles.

## Features

- User authentication with phone number and password
- JWT token-based authentication
- Role-based access control (RBAC)
- Admin panel for user management
- Redis for token storage
- PostgreSQL for data persistence
- Comprehensive error handling
- Input validation
- Unit tests

## Project Structure

```
.
├── cmd/
│   └── api/                 # Application entry point
├── internal/
│   ├── core/               # Core business logic
│   │   ├── entities/       # Domain entities
│   │   ├── errors/         # Custom error definitions
│   │   ├── ports/          # Interface definitions
│   │   └── service/        # Business logic implementation
│   ├── infrastructure/     # External implementations
│   │   ├── repository/     # Database implementations
│   │   └── token/          # Token service implementation
│   └── controller/         # HTTP handlers and validators
├── pkg/                    # Shared utilities
└── configs/                # Configuration files
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL
- Redis
- Make (optional)

### Installation

1. Clone the repository:

```bash
git clone https://github.com/amirdashtii/go_auth.git
cd go_auth
```

2. Install dependencies:

```bash
go mod download
```

3. Set up environment variables:

```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Run the application:

```bash
make run
# or
go run cmd/api/main.go
```

## API Endpoints

### Authentication

- `POST /auth/register` - Register a new user
- `POST /auth/login` - Login with phone number and password
- `POST /auth/refresh-token` - Refresh access token
- `POST /auth/logout` - Logout user (requires authentication)

### User Management

- `GET /users/me` - Get current user profile (requires authentication)
- `PUT /users/me` - Update user profile (requires authentication)
- `PUT /users/me/change-password` - Change password (requires authentication)
- `DELETE /users/me` - Delete user profile (requires authentication)

### Admin Endpoints

- `GET /admin/users` - List all users (requires admin authentication)
  - Query parameters:
    - `status`: Filter by status (active/inactive)
    - `role`: Filter by role (user/admin/super_admin)
    - `sort`: Sort field (created_at, etc.)
    - `order`: Sort order (asc/desc)
- `GET /admin/users/:id` - Get user details (requires admin authentication)
- `PUT /admin/users/:id` - Update user (requires admin authentication)
- `PUT /admin/users/:id/role` - Change user role (requires admin authentication)
- `PUT /admin/users/:id/status` - Change user status (requires admin authentication)
- `DELETE /admin/users/:id` - Delete user (requires admin authentication)

## Error Handling

The service implements a comprehensive error handling system with:

- Custom error types
- Bilingual error messages (English and Persian)
- Proper HTTP status codes
- Detailed error responses

## Testing

Run tests:

```bash
make test
# or
go test ./...
```

## Contributing

1. Create a new branch for your feature
2. Make your changes
3. Write or update tests
4. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
