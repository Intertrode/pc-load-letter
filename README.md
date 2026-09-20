# PC LOAD LETTER (`pc-load-letter`)

> *"PC Load Letter? What the f\*\*\* does that mean?!"* — Michael Bolton, *Office Space*

`pc-load-letter` is an enterprise print-spooler microservice engineered for **Intertrode** (formerly Initech). It provides queue management, error telemetry, and continuous integration pipeline automation.

---

## 🖨️ Microservice Overview

This Go microservice simulates the notorious office printer queue:
* **`GET /healthz`**: Liveness probe returning service health and UTC timestamp.
* **`GET /status`**: Returns real-time printer hardware state, tray status, and the infamous `PC LOAD LETTER` error code.
* **`GET /jobs`**: Lists queued documents (including Peter Gibbons' 42-page TPS report cover sheet).
* **`POST /jobs`**: Enqueues a new print job.

---

## 🚀 Running Locally

```bash
# Run unit tests
go test -v ./...

# Start the server (default port 8080)
go run .

# Test status endpoint
curl http://localhost:8080/status

# Enqueue a test document
curl -X POST http://localhost:8080/jobs \
  -H "Content-Type: application/json" \
  -d '{"document":"2026_Audit_Report.pdf","requested_by":"Donnie Page","pages":10}'
```

---

## 🏗️ Architecture & CI/CD Pipeline

The repository includes a GitHub Actions continuous integration pipeline (`.github/workflows/ci.yml`):
* Automated Go build and comprehensive unit testing across modules.
* Cloud deployment authentication hooks for staging delivery.
* Automated release artifact packaging (`dist/pc-load-letter-linux-amd64.tar.gz`).
