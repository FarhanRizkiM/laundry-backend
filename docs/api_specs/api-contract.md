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

### Auths

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

- POST /categories

- GET /categories

- GET /categories/{id}

- PUT /categories/{id}

- DELETE /categories/{id}

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

- PATCH /orders/{id}

### Payments

- GET /payments

- GET /payments/{id}

- PATCH /payments/{id}

### Deliveries

- GET /deliveries

- GET /deliveries/{id}

- GET /deliveries/my-tasks

- PATCH /deliveries/{id}

### Customer (Endpoint Public Tracking)

- GET /orders/track/{invoice_number}

### Reports

- GET /reports/dashboard

- GET /reports/revenue

- GET /reports/payments

- GET /reports/employees

- GET /reports/analytics
