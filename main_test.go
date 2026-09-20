package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()

	handler := setupRoutes()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, status)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["status"] != "healthy" {
		t.Errorf("expected status 'healthy', got %v", resp["status"])
	}
}

func TestStatusHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	rr := httptest.NewRecorder()

	handler := setupRoutes()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, status)
	}

	var state PrinterState
	if err := json.Unmarshal(rr.Body.Bytes(), &state); err != nil {
		t.Fatalf("failed to decode status: %v", err)
	}

	if state.ErrorCode != "PC LOAD LETTER" {
		t.Errorf("expected ErrorCode 'PC LOAD LETTER', got %s", state.ErrorCode)
	}
}

func TestJobsPostHandler(t *testing.T) {
	body := []byte(`{"document":"Special_Memo_TPS.pdf","requested_by":"Samir Nagheenanajar","pages":12}`)
	req := httptest.NewRequest(http.MethodPost, "/jobs", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler := setupRoutes()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, status)
	}

	var job PrintJob
	if err := json.Unmarshal(rr.Body.Bytes(), &job); err != nil {
		t.Fatalf("failed to decode job: %v", err)
	}

	if job.Document != "Special_Memo_TPS.pdf" {
		t.Errorf("expected document 'Special_Memo_TPS.pdf', got %s", job.Document)
	}
	if job.Status != "pending_paper_load" {
		t.Errorf("expected status 'pending_paper_load', got %s", job.Status)
	}
}
