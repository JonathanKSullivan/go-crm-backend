package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Tests happy path of submitting a well-formed GET /customers request
func TestGetCustomersHandler(t *testing.T) {
	// Setup
	customers = make(map[int]Customer)
	seedMockDB()

	req, err := http.NewRequest("GET", "/customers", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getCustomersHandler)
	handler.ServeHTTP(rr, req)

	// Checks for 200 status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("getCustomersHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Checks for JSON response
	if ctype := rr.Header().Get("Content-Type"); ctype != "application/json" {
		t.Errorf("Content-Type does not match: got %v want %v", ctype, "application/json")
	}
}

// Tests happy path of submitting a well-formed POST /customers request
func TestAddCustomerHandler(t *testing.T) {
	// Setup
	customers = make(map[int]Customer)

	requestBody := strings.NewReader(`
		{
			"name": "Example Name",
			"address": "123 Example St",
			"phone": "555-0199",
			"email": "example@example.com"
		}
	`)

	req, err := http.NewRequest("POST", "/customers", requestBody)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(addCustomerHandler)
	handler.ServeHTTP(rr, req)

	// Checks for 201 status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("addCustomerHandler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Checks for JSON response
	if ctype := rr.Header().Get("Content-Type"); ctype != "application/json" {
		t.Errorf("Content-Type does not match: got %v want %v", ctype, "application/json")
	}
}

// Tests unhappy path of deleting a customer that doesn't exist
func TestDeleteCustomerHandler(t *testing.T) {
	// Setup
	customers = make(map[int]Customer)
	seedMockDB()

	req, err := http.NewRequest("DELETE", "/customers/999", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deleteCustomerHandler)
	handler.ServeHTTP(rr, req)

	// Checks for 404 status code
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("deleteCustomerHandler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

// Tests unhappy path of getting a customer that doesn't exist
func TestGetCustomerHandler(t *testing.T) {
	// Setup
	customers = make(map[int]Customer)
	seedMockDB()

	req, err := http.NewRequest("GET", "/customers/999", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getCustomerHandler)
	handler.ServeHTTP(rr, req)

	// Checks for 404 status code
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("getCustomerHandler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}
