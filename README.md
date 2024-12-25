# E-commerce API

A robust and scalable e-commerce API built with Go, featuring user authentication, product management, order processing, and email notifications. Built using Echo framework and following Domain-Driven Design principles.

## Table of Contents

- [E-commerce API](#e-commerce-api)
  - [Features](#features)
  - [Architecture](#architecture)
  - [Getting Started](#getting-started)
  - [API Documentation](#api-documentation)
  - [Docker Setup](#docker-setup)
  - [Environment Variables](#environment-variables)
  - [Development Commands](#development-commands)
  - [License](#license)

## Features

- **Authentication & Authorization**
  - JWT-based authentication
  - Role-based access control (Admin/User)
  - Secure password hashing

- **Product Management**
  - CRUD operations for products
  - Category management
  - Stock tracking

- **Order Processing**
  - Order creation and management
  - Order status tracking
  - Email notifications

- **Email System**
  - MJML template support
  - Async processing with queue
  - Multiple provider support (SMTP, SendGrid, etc.)
  - Attachment handling

## Architecture

The project follows Domain-Driven Design principles with a clean architecture:

```
├── cmd/
│   └── api/              # Application entry point
├── internal/
│   ├── domain/          # Business logic and entities
│   │   ├── user/
│   │   ├── product/
│   │   └── order/
│   ├── middleware/      # HTTP middleware
│   ├── db/              # Database setup
│   ├── common/          # Shared utilities
│   └── config/          # Configuration
└── templates/           # Email templates
```

## Getting Started

### Prerequisites

- Go 1.22+
- PostgreSQL
- Redis (for queue)
- Docker & Docker Compose (optional)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/nneji123/ecommerce-golang.git
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up environment variables:
   ```bash
   cp .env.sample .env
   # Edit .env with your configuration (refer to .env.sample for example)
   ```

4. Start the server:
   ```bash
   go run cmd/api/main.go
   ```

## API Documentation

API documentation is available at `/swagger/index.html` when running the server. Key endpoints include:

- **Authentication**
  - POST `/auth/register`
  - POST `/auth/login`

- **Products**
  - GET `/api/products`
  - POST `/api/admin/products` (Admin only)
  - PUT `/api/admin/products/{id}` (Admin only)

- **Orders**
  - POST `/api/orders`
  - GET `/api/orders`
  - PUT `/api/orders/{id}/status` (Admin only)

## Docker Setup

You can run the entire stack with Docker Compose for easy local development and testing.

### 1. **Docker Build and Run**

To start the services with Docker Compose, use the following commands:

```bash
# Start services
docker-compose up -d

# Stop services
docker-compose down

# Rebuild services (if you have made changes)
docker-compose up -d --build
```

This will start:
- **API**: The main e-commerce application.
- **PostgreSQL**: The database used for storing user and product data.
- **Mailpit**: A test email server to view emails sent from the application (see below).

### 2. **Viewing Emails with Mailpit**

To view emails sent from the application during development (e.g., registration confirmation, password reset), you can use **Mailpit**. This is a simple, local email testing server that is started with Docker.

Once the services are running, visit [http://localhost:8025](http://localhost:8025) in your browser to view incoming emails. SMTP messages are sent to **Mailpit** using port **1025**, which is automatically configured in the Docker Compose file.

### 3. **Additional Docker Commands**

You can also use Docker to build and run the API container directly:

```bash
# Build the Docker image
docker build -t api .

# Run the Docker container
docker run -p 8080:8080 api
```

## Environment Variables

Required environment variables are specified in the `.env` file. After cloning the repository, copy the sample environment file:

```bash
cp .env.sample .env
```

Refer to `.env.sample` for the necessary environment variables.

### Example `.env` file:

```env

# SERVER
SERVER_PORT=8080
CORS_ALLOWED_ORIGINS=http://localhost:3000, https://example.com, http://localhost:8080

# DATABASE
POSTGRES_DSN=postgres://myuser:mypassword@localhost:5432/database?sslmode=disable
APP_URL="https://example.com"
JWT_SECRET="TEST-SECRET"

# EMAIL CREDENTIALS
EMAIL_FROM_ADDRESS=no-reply@gocommerce.com
EMAIL_FROM_NAME="GoCommerce"
SMTP_SERVER=localhost
SMTP_PORT=1025
SMTP_USERNAME=user1
SMTP_PASSWORD=password1
```

Make sure to configure the SMTP settings to match your email service provider or use Mailpit for local testing.

## Development Commands

You can run the following commands for local development and testing:

```bash
# Run tests
go test ./...

# Generate Swagger docs
swag init -g cmd/api/main.go

# Run with hot reload (using air)
air

# Build for production
go build -o ecommerce-api cmd/api/main.go
```

For local development, we recommend using `air` for hot reloading. Install it with:

```bash
go install github.com/cosmtrek/air@latest
```

## Makefile

### Common Commands

```makefile
# Run tests
make test

# Run the server
make run

# Build the binary
make build

# Build Docker image
make docker-build

# Run Docker container
make docker-run
```

You can refer to the `Makefile` for additional commands that simplify common tasks such as building, testing, and running the application.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

