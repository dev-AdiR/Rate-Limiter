# Rate Limiter

A lightweight, JWT-authenticated **reverse proxy with rate limiting** built in Go. It sits in front of your backend services, validates incoming JWT tokens, and forwards requests to a configurable target URL — all containerized and ready to deploy with Docker.

---

## Features

- **Reverse Proxy** — Transparently forwards HTTP requests to a configurable upstream service
- **JWT Authentication** — Validates `Authorization: Bearer <token>` headers using a configurable secret
- **Rate Limiting** — Controls the number of requests allowed per client (in-progress module)
- **Docker Support** — First-class Docker and Docker Compose setup for easy deployment
- **Task Runner** — Uses [Task](https://taskfile.dev) for consistent build/run commands across environments

---

## Project Structure

```
Rate-Limiter/
├── Env/                  # Environment configuration files
├── proxy/                # Reverse proxy logic and HTTP handler registration
├── rate_limit/           # Rate limiting middleware (in progress)
├── main.go               # Application entrypoint
├── Dockerfile            # Docker image definition
├── docker-compose.yml    # Multi-container orchestration
├── Taskfile.yml          # Task runner commands (build, dev, run)
├── go.mod                # Go module definition
└── go.sum                # Go dependency checksums
```

---

## Prerequisites

- [Go 1.25+](https://golang.org/dl/)
- [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)
- [Task](https://taskfile.dev/installation/) (optional, for local dev)

---

## Getting Started

### Run with Docker Compose (recommended)

1. Clone the repository:
   ```bash
   git clone https://github.com/dev-AdiR/Rate-Limiter.git
   cd Rate-Limiter
   ```

2. Configure environment variables in `docker-compose.yml` (see [Configuration](#configuration)).

3. Start the service:
   ```bash
   docker-compose up --build
   ```

The proxy will be available at `http://localhost:8081`.

---

### Run Locally

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Build and run using Task:
   ```bash
   task run
   ```

   Or manually:
   ```bash
   go build -o ./dist/app
   ./dist/app
   ```

3. For live-reload during development (requires [Air](https://github.com/air-verse/air)):
   ```bash
   task dev
   ```

---

## Configuration

The service is configured via environment variables:

| Variable     | Description                                      | Example                                      |
|--------------|--------------------------------------------------|----------------------------------------------|
| `Target_Url` | The upstream service URL to proxy requests to    | `http://host.docker.internal:8080`           |
| `JWT_SECRET` | Secret key used to validate incoming JWT tokens  | `KbzVX1ZlViuzPHUbzm9VyiQGCLvLeDCdnQDbN4kAUam` |

> **Note:** In production, replace the example `JWT_SECRET` with a strong, randomly generated secret and store it securely (e.g. via a secrets manager or `.env` file not committed to version control).

When running locally, you can place these in an `.env` file inside the `Env/` directory or export them in your shell.

---

## How It Works

```
Client Request
     │
     ▼
Rate Limiter Proxy (port 8081)
     │  1. Validate JWT token
     │  2. Apply rate limit rules
     │  3. Forward request
     ▼
Target Service (Target_Url)
     │
     ▼
Response returned to client
```

1. An incoming request hits the proxy on port `8081`.
2. The proxy validates the `Authorization: Bearer <token>` JWT header against `JWT_SECRET`.
3. If valid (and within rate limits), the request is forwarded to `Target_Url`.
4. The upstream response is returned to the original caller.

---

## Task Commands

| Command      | Description                          |
|--------------|--------------------------------------|
| `task build` | Compile the binary to `./dist/app`   |
| `task run`   | Build and run the application        |
| `task dev`   | Start with live-reload via Air       |
| `task start` | Run the already-compiled binary      |

---

## Tech Stack

| Technology | Purpose |
|---|---|

| [quic-go](https://github.com/quic-go/quic-go) | QUIC/HTTP3 support |
| [godotenv](https://github.com/joho/godotenv) | `.env` file loading |
| [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver) | Database connectivity |
| Docker | Containerization |
| [Task](https://taskfile.dev) | Build automation |

---

## Contributing

Pull requests are welcome! For significant changes, please open an issue first to discuss what you'd like to change.

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Commit your changes: `git commit -m 'Add my feature'`
4. Push to the branch: `git push origin feature/my-feature`
5. Open a Pull Request

---

## License

This project does not currently specify a license. All rights reserved by the author unless otherwise stated.