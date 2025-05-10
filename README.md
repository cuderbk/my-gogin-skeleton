# my-gogin-skeleton Backend Skeleton

A clean, scalable Go backend skeleton built with [Gin](https://github.com/gin-gonic/gin), [Google Wire](https://github.com/google/wire), structured configuration files, and dependency injection.

---

## 🗂️ Project Structure

```
.
├── cmd
│   └── api
│       ├── main.go           # Entry point
│       ├── router.go         # Main router setup
├── pkg/                     # Common packages (e.g., my-gogin-skeleton)
│   ├── logger/                # Custom logger
│   ├── config/                # Config struct definitions
│   ├── models/                # DB models
├── config/                   # YAML configs for each subsystem
├── internal/
│   ├── common/               # Logger, middleware, utilities
│   ├── infra/                # DB, Redis clients, etc.
│   ├── domain/, service/, ... (TBD: your business logic)
├── go.mod / go.sum
└── README.md
```

---

## ⚙️ Prerequisites

- Go 1.23+

---

## 🧪 Quick Start

### 1. Clone and set up

```bash
git clone https://github.com/The-Innovators-DATN/my-gogin-skeleton.git
cd my-gogin-skeleton-backend
```

### 2. Initialize Wire DI

```bash
cd cmd/api
```


### 3. Run the application

```bash
go run cmd/api/*
```

Or directly:

```bash
cd cmd/api
go run .
```

> 🔧 Ensure your `config/*.yaml` files exist and are correctly structured.

---

## 🧰 Configuration

All configuration is stored in the `config/` folder:

- `app.yaml` — app name, host, port
- `database.yaml` — PostgreSQL DSN
- `redis.yaml` — Redis address
- `log.yaml` — log level/output
- ...

All loaded via Viper based on the folder passed in `InitializeServer("config_path")`.

---

## ✅ Health Check

Once running, hit:

```bash
curl http://localhost:<your-port>/healthz
```

Expected response:
```json
{"status": "ok"}
```

---

## 📦 Re-generating Wire (when deps change)

```bash
cd cmd/api
wire
```

## 📌 Troubleshooting

- **`undefined: InitializeServer`**  
  → You're running `go run main.go` instead of `go run .`

- **Config file not found**  
  → Set correct relative path: `../../config` from `cmd/api`

- **Wire errors: "no provider found..."**  
  → You missed a dependency or didn't pass it down to `buildRouter(...)`

---

