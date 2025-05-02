// Package integration_test содержит интеграционные тесты для проверки создания клиента и корректной работы лимитера
package integration_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"load_balancer/testutils"
)

// TestCreateClientIntegration проверяет корректное создание клиента
func TestCreateClientIntegration(t *testing.T) {
	handler := testutils.NewHandlerWithMockedDeps()

	body := bytes.NewBufferString(`{"id": "test-client", "capacity": 10, "refill_rate": 5}`)
	req := httptest.NewRequest(http.MethodPost, "/clients", body)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}
}

// TestRateLimitExceededIntegration проверяет работу лимитера
func TestRateLimitExceededIntegration(t *testing.T) {
	handler := testutils.NewHandlerWithMockedDeps()

	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Client-ID", "limited-client")

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if i >= 5 && rr.Code != http.StatusTooManyRequests {
			t.Errorf("Expected 429 Too Many Requests, got %d at request %d", rr.Code, i)
		}
	}
}
