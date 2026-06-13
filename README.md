# KanBan Auth Service - MimoGen

Authentication microservice for KanBan SaaS platform.

## Features

- JWT authentication (access + refresh tokens)
- User registration and login
- OAuth2 support (Google, GitHub)
- PostgreSQL database

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | /api/v1/auth/register | Register new user |
| POST | /api/v1/auth/login | Login |
| POST | /api/v1/auth/refresh | Refresh tokens |
| POST | /api/v1/auth/logout | Logout |
| GET | /api/v1/auth/me | Get profile |
| PUT | /api/v1/auth/me | Update profile |

## Tech Stack

- Go 1.22
- Chi Router
- PostgreSQL
- JWT (golang-jwt)

## Quick Start

```bash
# Set environment variables
export DB_HOST=localhost
export DB_NAME=kanban_auth
export JWT_SECRET=your-secret

# Run
go run ./cmd/main.go
```

---

<details>
<summary><strong>IMPORTANT NOTICE</strong></summary>

<br>

**This repository was entirely generated using [MiMoCode](https://github.com/xiaomi/mimocode) - an AI-powered coding assistant by Xiaomi.**

All code, tests, documentation, and infrastructure configuration in this repository were created through AI-assisted development with MiMoCode. The codebase demonstrates the capabilities of AI in generating production-ready microservices architecture.

---

**Этот репозиторий был полностью сгенерирован с помощью [MiMoCode](https://github.com/xiaomi/mimocode) - AI-ассистента для программирования от Xiaomi.**

Весь код, тесты, документация и инфраструктурная конфигурация в этом репозитории были созданы с помощью AI-ассистированной разработки MiMoCode. Кодовая база демонстрирует возможности ИИ в генерации готовых к производству микросервисных архитектур.

</details>
