# 🛒 Mini E-Commerce API

A robust RESTful API for e-commerce services built with **Go (Golang)** and **Clean Architecture**.

## 🛠 Tech Stack
- **Language:** Go (Golang)
- **Framework:** Gin
- **Database:** PostgreSQL
- **Infrastructure:** Docker, Docker Compose
- **Testing:** Gomock, Testify
- **Security:** JWT, Bcrypt

## ✨ Features Overview

### 🏆 Key Highlights (Technical Challenges)
- **Architecture:** Designed using **Clean Architecture** to ensure separation of concerns and testability.
- **Concurrency Control:** Implemented **Pessimistic Locking (`SELECT ... FOR UPDATE`)** to prevent inventory race conditions and overselling.
- **Atomic Transactions:** Ensures data consistency across Order, OrderItems, and Inventory using **Database Transactions (ACID)**.
- **Testing:** Achieved high test coverage for Business Logic Layer using **Unit Testing**, **Table-Driven Tests**, and **Mocking**.
- **Security:** Secured endpoints using **JWT Authentication** and **Role-Based Access Control (RBAC)** middleware.

### 📦 Core Modules
- **🔐 Authentication & Security**
  - Secure JWT Authentication (Access & Refresh Tokens).
  - RBAC Middleware for Admin, Seller, and User roles.
  
- **🛒 Shopping Cart**
  - **Smart Cart System:** Supports both Logged-in Users and **Guest Users** (Session-based).
  - Real-time stock and price validation before checkout.

- **📦 Product Catalog**
  - Product management with ownership authorization (Seller can only edit their own products).
  - Efficient filtering and pagination support.

- **📝 Order Management**
  - Full order lifecycle: Create, View History, Cancel.
  - Automatic stock restoration upon order cancellation.
  - Admin controls for order status updates.

## 🚀 How to Run

```
## Clone Repository
git clone [https://github.com/codepnw/mini-ecommerce.git](https://github.com/codepnw/mini-ecommerce.git)
cd mini-ecommerce

## Setup Environment
cp -n .env.example .env

## Run with Docker
docker compose up --build -d
```
The API Server will start at: **http://localhost:8080**

**NOTE**: The database will be automatically initialized with schema and seed data located in **./scripts**

## 📨 API Documentation
You can import the **Postman** collection file to test the endpoints:

File: **./docs/postman_collection.json**

## 🔐 Default Credentials

The database comes pre-filled with the following accounts for testing:

| Role      | Email                | Password | Permissions |
| :---      | :---                 | :---     | :---        |
| **Admin** | `admin@example.com`  | `123456` | Full Access (Manage Products/Orders) |
| **Seller**| `seller@example.com` | `123456` | Manage Own Products & Orders |
| **User** | `user@example.com`   | `123456` | Buy Products, View History |
