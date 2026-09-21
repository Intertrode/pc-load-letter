# Intertrode Print Services (`pc-load-letter`)

`pc-load-letter` is an enterprise print-spooler microservice engineered for **Intertrode Core Services**. It provides resilient queue management, hardware tray telemetry, and automated asynchronous document routing across regional office infrastructure.

---

## 🖨️ Service Overview

The microservice exposes a lightweight REST API for document dispatch and queue state inspection:

* **`GET /healthz`**: Standard liveness probe returning service health and UTC timestamp for Kubernetes/load-balancer health checks.
* **`GET /status`**: Returns real-time hardware status, paper tray occupancy levels, and internal diagnostic error codes.
* **`GET /jobs`**: Lists currently queued document jobs awaiting spooler dispatch.
* **`POST /jobs`**: Submits a new document to the asynchronous spooler queue.

---

## ⚙️ Configuration & Environment Variables

| Variable | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `PORT` | integer | `8080` | Port on which the HTTP server listens |
| `ENVIRONMENT` | string | `development` | Runtime environment (`development`, `staging`, `production`) |
| `LOG_LEVEL` | string | `info` | Logging verbosity (`debug`, `info`, `warn`, `error`) |
| `AWS_REGION` | string | `us-east-1` | Target AWS region for artifact delivery and cloud telemetry |

---

## 🚀 Local Development

```bash
# Run unit test suite with coverage
go test -v -race ./...

# Compile and run service locally
go run .

# Verify liveness probe
curl http://localhost:8080/healthz

# Inspect hardware tray status
curl http://localhost:8080/status

# Enqueue a batch document
curl -X POST http://localhost:8080/jobs \
  -H "Content-Type: application/json" \
  -d '{"document":"Executive_Payroll_Summary_2026.pdf","requested_by":"Finance Ops","pages":12}'
```

---

## 🏗️ Architecture & Continuous Delivery

The repository integrates an automated GitHub Actions CI/CD pipeline (`.github/workflows/ci.yml`):
* **Build & Validation:** Multi-stage Go 1.22 compilation, unit test execution, and static analysis.
* **Cloud Authentication:** Automated short-lived cloud provider IAM session configuration for staging artifact delivery.
* **Distribution Artifacts:** Automated release packaging generating standardized deployment archives (`dist/pc-load-letter-linux-amd64.tar.gz`).
