package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// PrintJob represents a queued print job in the Intertrode office printer system.
type PrintJob struct {
	ID        string    `json:"id"`
	Document  string    `json:"document"`
	Requested string    `json:"requested_by"`
	Pages     int       `json:"pages"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// PrinterState represents the live status of the infamous office printer.
type PrinterState struct {
	Model       string     `json:"model"`
	Status      string     `json:"status"`
	ErrorCode   string     `json:"error_code"`
	Explanation string     `json:"explanation"`
	PaperTray   string     `json:"paper_tray"`
	TonerLevel  int        `json:"toner_level_percent"`
	ActiveJobs  []PrintJob `json:"active_jobs"`
}

var (
	jobsMutex sync.RWMutex
	jobs      = []PrintJob{
		{
			ID:        "job-001",
			Document:  "2026_Q3_TPS_Report_Cover_Sheet.pdf",
			Requested: "Peter Gibbons",
			Pages:     42,
			Status:    "stalled",
			CreatedAt: time.Now().Add(-15 * time.Minute),
		},
		{
			ID:        "job-002",
			Document:  "Intertrode_Acquisition_Analysis.docx",
			Requested: "Bill Lumbergh",
			Pages:     18,
			Status:    "queued",
			CreatedAt: time.Now().Add(-5 * time.Minute),
		},
	}
)

func statusHandler(w http.ResponseWriter, r *http.Request) {
	jobsMutex.RLock()
	defer jobsMutex.RUnlock()

	state := PrinterState{
		Model:       "LaserPrint 4000 (Initech Standard Issue)",
		Status:      "Error",
		ErrorCode:   "PC LOAD LETTER",
		Explanation: "Paper Cassette (PC) empty. Please load 'Letter' size paper.",
		PaperTray:   "Empty",
		TonerLevel:  14,
		ActiveJobs:  jobs,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(state); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "healthy",
		"service": "pc-load-letter",
		"time":    time.Now().UTC().Format(time.RFC3339),
	})
}

func jobsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		jobsMutex.RLock()
		defer jobsMutex.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(jobs)

	case http.MethodPost:
		var req struct {
			Document  string `json:"document"`
			Requested string `json:"requested_by"`
			Pages     int    `json:"pages"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if req.Document == "" {
			http.Error(w, "document name required", http.StatusBadRequest)
			return
		}

		jobsMutex.Lock()
		newJob := PrintJob{
			ID:        fmt.Sprintf("job-%03d", len(jobs)+1),
			Document:  req.Document,
			Requested: req.Requested,
			Pages:     req.Pages,
			Status:    "pending_paper_load",
			CreatedAt: time.Now().UTC(),
		}
		jobs = append(jobs, newJob)
		jobsMutex.Unlock()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(newJob)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func setupRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthHandler)
	mux.HandleFunc("/status", statusHandler)
	mux.HandleFunc("/jobs", jobsHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		_, _ = fmt.Fprintln(w, "PC LOAD LETTER - Intertrode Print Spooler Microservice")
	})
	return mux
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	handler := setupRoutes()
	log.Printf("Starting pc-load-letter print spooler on port %s...", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("server terminated unexpectedly: %v", err)
	}
}
