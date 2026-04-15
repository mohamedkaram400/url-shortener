# 🔗 Url Shortener Project

A **scalable backend system** built with **Golang (Gin)** to help users **generate short urls** through an easy, automated, and structured process.  
This project follows a **modular clean architecture** and supports **multiple services** like MySQL, and Redis — all orchestrated with Docker.

---

## 📑 Table of Contents
- [Tech Stack](#-tech-stack)
- [Requirements](#-requirements)
- [Project Structure](#-project-structure)
- [Getting Started](#-getting-started)
- [Usage](#-usage)
- [Example API Request](#-example-api-request)
- [Features](#-features)
- [Endpoint Collection](#-endpoint-collection)
- [License](#-license)


## 🧰 Tech Stack

- 🐹 **Golang (Gin)** — REST API framework  
- 🐬 **MySQL** — Primary relational database  
<!-- - 🧠 **Redis** — Caching and session management   -->
- 🐳 **Docker & Docker Compose** — Containerization and environment orchestration

---

## 📋 Requirements

Before you begin, ensure you have the following installed:

- [Golang](https://go.dev/) `>= 1.21`
- [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)
- Git

---

### 🧠 Project Structure
```
.
├── cmd/                  # Application entry point
├── config/               # Configuration files
├── conn/                 # Database connections (MySQL, Redis)
├── deploy/               # Docker & Nginx configs
├── internal/
│   ├── adapters/         # Repositories, notifiers, schedulers
│   ├── core/             # Business logic & services
│   ├── delivery/         # HTTP Handlers & Routes
├── db/
│   ├── migrations/       # Migrations
│   └── seeders/          # Initial data seeding
├── pkg/                  # Shared utilities & helpers
├── docker-compose.yml
├── Dockerfile
└── README.md
```
---

## 🚀 Getting Started

### 1. Clone the Repository
```bash
git clone https://github.com/mohamedkaram400/url-shortener.git
cd url-shortener
```

## ⚙️ Setup Environment
```bash
cp .env.example .env
```

## 🐳 Docker Setup
```bash
docker -f docker-compose.yml compose up --build
```
---

### Run the Application (Without Docker)

#### If you want to run locally without Docker:

```bash
go mod tidy
go run cmd/main.go
```
---

## 🧪 Usage
Once the application is running:

The server will be available at: http://localhost:9000

You can use tools like Postman or cURL to test the API.

---
## 🧪 Example API Request

### 📝 Create a New short url
```http
POST /api/shorten
Content-Type: application/json
```
#### 📥 Request

```json
{
    "long_url": "https://algomaster.io/learn/lld/what-is-lld",
    "expiration_days": 1,
    "user_id":  1,
    "status":   "Active",
    // "custom_alias": "hello"
}
```

#### 📤 Response

```json
{
    "success": true,
    "message": "URL Generated successfully",
    "data": {
        "short_url": "http://localhost:9000/Zl8xuJf",
        "long_url": "https://algomaster.io/learn/lld/what-is-lld"
    }
}
```

---
## ✨ Features:
- Modular clean architecture
- Multi-database support (MySQL + Redis)
- Easy environment setup with Docker
- Centralized configuration management
- Ready for scaling and service extension

---
## 📚 Endpoint Collection
The full list of API endpoints is included in the project’s main directory for easy reference and testing.

---

## 🛡️ License

This project is licensed under the [MIT License](LICENSE). You are free to use, modify, and share this project with proper attribution.
