package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetIPStatuses(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/ip", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetIPStatuses)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status code: got %v, expected %v", status, http.StatusOK)
	}
}
