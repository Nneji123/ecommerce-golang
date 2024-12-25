# E-commerce API

A robust and scalable e-commerce API built with Go, featuring user authentication, product management, order processing, and email notifications. Built using Echo framework and following Domain-Driven Design principles.

## TODO

### Core Features
- [x] User Authentication with JWT
- [x] Product Management
- [x] Order Processing
- [x] Email Notifications
- [ ] Payment Processing Integration
- [ ] Inventory Management System
- [ ] Shopping Cart Functionality

### Enhancements
- [x] Async Email Processing
- [x] MJML Email Templates
- [x] File Upload for Product Images
- [ ] Webhook Support for Order Updates
- [ ] Cache Layer Implementation
- [ ] Real-time Inventory Updates

### Documentation
- [x] API Documentation with Swagger
- [x] Email Template Documentation
- [ ] Deployment Guide
- [ ] Contributing Guidelines

## Table of Contents

- [E-commerce API](#e-commerce-api)
  - [Features](#features)
  - [Architecture](#architecture)
  - [Getting Started](#getting-started)
  - [API Documentation](#api-documentation)
  - [Docker Setup](#docker-setup)
  - [Environment Variables](#environment-variables)
  - [Development Commands](#development-commands)

## Features

- **Authentication & Authorization**
  - JWT-based authentication
  - Role-based access control (Admin/User)
  - Secure password hashing

- **Product Management**
  - CRUD operations for products
  - Category management
  - Image upload support
  - Stock tracking

- **Order Processing**
  - Order creation and management
  - Order status tracking
  - Email notifications
  - Stock validation

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

- Go 1.23+
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
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. Start the server:
   ```bash
   go run cmd/api/main.go
   ```

## API Documentation

API documentation is available at `/swagger/index.html` when running the server. Key endpoints include:

- **Authentication**
  - POST `/api/auth/register`
  - POST `/api/auth/login`

- **Products**
  - GET `/api/products`
  - POST `/api/admin/products` (Admin only)
  - PUT `/api/admin/products/{id}` (Admin only)

- **Orders**
  - POST `/api/orders`
  - GET `/api/orders`
  - PUT `/api/orders/{id}/status` (Admin only)

## Docker Setup

Run the entire stack with Docker Compose:

```bash
# Start services
docker-compose up -d

# Stop services
docker-compose down

# Rebuild services
docker-compose up -d --build
```

## Environment Variables

Required environment variables:

```env
# Server
PORT=8080
ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=ecommerce
DB_USER=postgres
DB_PASSWORD=password

# JWT
JWT_SECRET=your-secret-key
JWT_DURATION=24h

# Email
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=your-email
SMTP_PASSWORD=your-password
```

## Development Commands

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

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.