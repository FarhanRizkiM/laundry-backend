# API Contract - Laundry Management System

## General Rules

- Base URL: /api/v1
- All requests and responses use JSON
- Authentication via Authorization header
- Timestamps use ISO 8601 format

## Authentication

Authenticated endpoints require:
Authorization: Bearer <token>

## Roles:

- owner
- kasir
- staff
- courier

## Roles & Access Rules

- Owner: read-only access to summaries
- Kasir: create orders, manage payments
- Staff: update production status
- Courier: handle delivery status
- Customer: public access (no auth) for order status lookup

## Endpoints

### Auth

- POST /auth/login
- POST /auth/logout
- GET /auth/me

### Users

- POST /users
- GET /users
- GET /users/{id}
- PUT /users/{id}
- DELETE /users/{id}

### Service Categories

- POST /service-categories
- GET /service-categories
- GET /service-categories/{id}
- PUT /service-categories/{id}
- DELETE /service-categories/{id}

### Services

- POST /services
- GET /services
- GET /services/{id}
- PUT /services/{id}
- DELETE /services/{id}

### Orders

- POST /orders
- GET /orders
- GET /orders/{id}
- PUT /orders/{id}
- PATCH /orders/{id}/status

### Payments

- POST /orders/{id}/payments

### Courier

- POST /orders/{id}/assign-courier
- PATCH /orders/{id}/delivery-status

### Customer

- GET /public/orders/{kode_nota}
