package httpapi_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/sample-custody/internal/app"
)

func TestHTTPWorkflow(t *testing.T) {
	handler := app.New().Handler
	response := postJSON(t, handler, "/samples", `{"id":"http-sample-001","code":"HTTP-001","material":"serum","origin":"test-lab","quantity":3,"unit":"ml"}`)
	if response.StatusCode != http.StatusCreated || !strings.Contains(response.Body, "registered") {
		t.Fatalf("register response: %+v", response)
	}
	response = postJSON(t, handler, "/samples/http-sample-001/transfers", `{"from":"intake","to":"review","location":"cold-room","operator":"http-test"}`)
	if response.StatusCode != http.StatusCreated || !strings.Contains(response.Body, "review") {
		t.Fatalf("transfer response: %+v", response)
	}
	response = postJSON(t, handler, "/batches", `{"id":"http-batch-001","name":"http batch","purpose":"smoke test","sample_ids":["http-sample-001"]}`)
	if response.StatusCode != http.StatusCreated || !strings.Contains(response.Body, "open") {
		t.Fatalf("batch response: %+v", response)
	}
	response = postJSON(t, handler, "/batches/http-batch-001/complete", `{}`)
	if response.StatusCode != http.StatusOK || !strings.Contains(response.Body, "completed") {
		t.Fatalf("completion response: %+v", response)
	}
}

func TestHTTPHealth(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	recorder := httptest.NewRecorder()
	app.New().Handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("health status: %d", recorder.Code)
	}
}

// TestHTTPDuplicateSampleIsConflict 验证重复登记同一个样本编码时被识别为明确的冲突（409），
// 而非系统失败（500），以便客户端能据此判断无需重试。
func TestHTTPDuplicateSampleIsConflict(t *testing.T) {
	handler := app.New().Handler
	body := `{"id":"dup-sample-001","code":"DUP-001","material":"serum","origin":"test-lab","quantity":1,"unit":"ml"}`
	first := postJSON(t, handler, "/samples", body)
	if first.StatusCode != http.StatusCreated {
		t.Fatalf("first register status: %d", first.StatusCode)
	}
	duplicate := postJSON(t, handler, "/samples", `{"id":"dup-sample-002","code":"DUP-001","material":"serum","origin":"test-lab","quantity":1,"unit":"ml"}`)
	if duplicate.StatusCode != http.StatusConflict {
		t.Fatalf("duplicate register status = %d, want %d (body: %s)", duplicate.StatusCode, http.StatusConflict, duplicate.Body)
	}
	if !strings.Contains(duplicate.Body, `"error":"conflict"`) {
		t.Fatalf("duplicate register body missing conflict kind: %s", duplicate.Body)
	}
}

type jsonResponse struct {
	StatusCode int
	Body       string
}

func postJSON(t *testing.T, handler http.Handler, path, body string) jsonResponse {
	t.Helper()
	request := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	return jsonResponse{StatusCode: recorder.Code, Body: recorder.Body.String()}
}
