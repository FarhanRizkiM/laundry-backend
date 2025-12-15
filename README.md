# Laundry Management System — Backend

Backend service for managing laundry operations including order intake, production workflow, delivery handling, and status tracking.

## Project Scope

This project focuses on building the backend system for a laundry management application.

Core responsibilities:

- Order management
- Laundry workflow tracking
- Role-based access (owner, cashier, staff, courier)
- Customer order status lookup

## Business Workflow Overview

1. Cashier creates a new laundry order.
2. Staff processes the order through laundry stages.
3. Orders are marked ready for pickup or delivery.
4. Courier handles delivery if applicable.
5. Customer can check order status using order code.
6. Owner monitors overall business performance.

## System Roles

- **Owner**: monitors system activity and business summary.
- **Cashier**: creates orders and manages payments.
- **Staff**: processes laundry orders.
- **Courier**: delivers laundry orders.
- **Customer**: checks order status without login.

## Order Status Lifecycle

Internal order statuses:

- pending
- in-progress
- ready-pickup
- ready-delivery
- being-delivered
- finished-delivery
- picked-up
- cancelled

## Database Design

The database schema is designed around order-centric workflow and includes:

- users
- services
- orders
- order_items
- payments
- status_history
- courier_assignments
- receipts

## Project Structure

````text
laundry-backend/
├─ cmd/
│  └─ server/
│     └─ main.go        # application entry point
├─ internal/
│  ├─ config/           # configuration loader
│  ├─ db/               # database connection & helpers
│  ├─ handlers/         # HTTP handlers (request/response)
│  ├─ services/         # business logic
│  ├─ repositories/     # database access layer
│  ├─ models/           # data models / entities
│  └─ middlewares/      # HTTP middlewares
├─ migrations/          # SQL database migrations
├─ docs/                # design & documentation
├─ spec/                # API specifications
└─ README.md
## Project Structure

```text
laundry-backend/
├─ cmd/
│  └─ server/
│     └─ main.go        # application entry point
├─ internal/
│  ├─ config/           # configuration loader
│  ├─ db/               # database connection & helpers
│  ├─ handlers/         # HTTP handlers (request/response)
│  ├─ services/         # business logic
│  ├─ repositories/     # database access layer
│  ├─ models/           # data models / entities
│  └─ middlewares/      # HTTP middlewares
├─ migrations/          # SQL database migrations
├─ docs/                # design & documentation
├─ spec/                # API specifications
└─ README.md
````

## Tech Stack

- Language: Go
- Database: MySQL
- Architecture: Layered (handler, service, repository)

## Development Status

Current phase:

- Database schema and migrations completed
- API contracts and implementation are in progress
