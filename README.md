# PC LOAD LETTER (`pc-load-letter`)

> *"PC Load Letter? What the f\*\*\* does that mean?!"* — Michael Bolton, *Office Space*

`pc-load-letter` is an enterprise print-spooler microservice engineered for **Intertrode** (formerly Initech). It provides simulated queue management, error telemetry, and automated CI/CD pipeline protection powered by **Tracebit Threat Deception**.

---

## 🖨️ Microservice Overview

This Go microservice simulates the notorious office printer queue:
* **`GET /healthz`**: Liveness probe returning service health and UTC timestamp.
* **`GET /status`**: Returns real-time printer hardware state, tray status, and the infamous `PC LOAD LETTER` error.
* **`GET /jobs`**: Lists queued documents (including Peter Gibbons' 42-page TPS report cover sheet).
* **`POST /jobs`**: Enqueues a new print job.

### Running Locally

```bash
# Run unit tests
go test -v ./...

# Start the server (default port 8080)
go run .

# Test status endpoint
curl http://localhost:8080/status
```

---

## 🛡️ Security Architecture: Tracebit CI/CD Deception

This repository serves as a live testbed for **Tracebit Threat Deception** in modern cloud CI/CD pipelines.

### The Threat Model
Modern CI/CD pipelines are susceptible to **Software Supply Chain Compromises**:
1. A developer imports a compromised third-party GitHub Action or malicious dependency (`npm`, `pip`, Go module).
2. The malicious code executes during the runner build phase, attempting to harvest environment variables, cloud keys, and secrets.

### The Tracebit Deception Defense
In `.github/workflows/ci.yml`, the pipeline invokes `tracebit-com/tracebit-community-action@v1`:
* **Dynamic Canary Generation:** On every pipeline execution, Tracebit mints fresh, ephemeral AWS honeytokens.
* **Contextual Tagging:** These credentials carry invisible metadata identifying this repository (`Intertrode/pc-load-letter`), workflow run ID, commit SHA, and job name.
* **Tripwire Alerting:** If an attacker or compromised action exfiltrates these credentials and attempts to use them (e.g., calling `aws sts get-caller-identity`), Tracebit detects the call via CloudTrail and alerts within seconds with the exact build and attacker IP.

---

## 🧪 Testing the Canary Tripwire

1. Trigger a workflow run on GitHub Actions (**Actions → CI & Tracebit Supply-Chain Defense → Run workflow**).
2. Inspect the **Verify Tracebit Canary Injection** job log to view the dynamically issued AWS canary key ID.
3. Simulate an adversary from your local terminal:
   ```bash
   AWS_ACCESS_KEY_ID="<CANARY_KEY_ID>" \
   AWS_SECRET_ACCESS_KEY="<CANARY_SECRET_KEY>" \
   aws sts get-caller-identity
   ```
4. Check the **Tracebit Portal** under **Alerts** to see the high-fidelity intrusion detection.
