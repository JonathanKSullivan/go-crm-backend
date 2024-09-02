package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Customer struct represents the data model for a customer.
// JSON tags are used to map struct fields to JSON keys.
type Customer struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Contacted bool   `json:"contacted"`
}

// CustomerStore defines the interface for customer storage.
type CustomerStore interface {
	AddCustomer(c Customer) int
	GetCustomer(id int) (Customer, bool)
	UpdateCustomer(id int, c Customer) error
	DeleteCustomer(id int) error
	GetCustomers() map[int]Customer
}

// InMemoryCustomerStore implements CustomerStore with an in-memory map.
type InMemoryCustomerStore struct {
	customers map[int]Customer
	lastID    int
}

// NewInMemoryCustomerStore initializes a new InMemoryCustomerStore with seed data.
func NewInMemoryCustomerStore() *InMemoryCustomerStore {
	store := &InMemoryCustomerStore{
		customers: make(map[int]Customer),
	}
	store.seedMockDB()
	return store
}

func (s *InMemoryCustomerStore) AddCustomer(c Customer) int {
	s.lastID++
	c.ID = s.lastID
	s.customers[c.ID] = c
	return c.ID
}

func (s *InMemoryCustomerStore) GetCustomer(id int) (Customer, bool) {
	customer, exists := s.customers[id]
	return customer, exists
}

func (s *InMemoryCustomerStore) UpdateCustomer(id int, c Customer) error {
	if _, exists := s.customers[id]; !exists {
		return fmt.Errorf("customer with ID %d not found", id)
	}
	c.ID = id // Ensure the ID is maintained
	s.customers[id] = c
	return nil
}

func (s *InMemoryCustomerStore) DeleteCustomer(id int) error {
	if _, exists := s.customers[id]; !exists {
		return fmt.Errorf("customer with ID %d not found", id)
	}
	delete(s.customers, id)
	return nil
}

func (s *InMemoryCustomerStore) GetCustomers() map[int]Customer {
	return s.customers
}

func (s *InMemoryCustomerStore) seedMockDB() {
	s.customers[1] = Customer{ID: 1, Name: "John Doe", Role: "Manager", Address: "123 Main St", Phone: "555-555-5555", Email: "johndoe@example.com", Contacted: true}
	s.customers[2] = Customer{ID: 2, Name: "Jane Doe", Role: "Employee", Address: "456 Main St", Phone: "555-555-5555", Email: "janedoe@example.com", Contacted: false}
	s.customers[3] = Customer{ID: 3, Name: "Jim Doe", Role: "Supervisor", Address: "789 Main St", Phone: "555-555-5555", Email: "jimdoe@example.com", Contacted: true}
	s.customers[4] = Customer{ID: 4, Name: "Jill Doe", Role: "HR", Address: "1011 Main St", Phone: "555-555-5555", Email: "jilldoe@example.com", Contacted: false}
	s.lastID = len(s.customers)
}

// JSONResponse writes a JSON response with a given status code.
func JSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func addCustomerHandler(store CustomerStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newCustomer Customer
		if err := json.NewDecoder(r.Body).Decode(&newCustomer); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id := store.AddCustomer(newCustomer)
		w.Header().Set("Location", fmt.Sprintf("/customers/%d", id))
		JSONResponse(w, http.StatusCreated, newCustomer)
	}
}

func getCustomerHandler(store CustomerStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["id"]
		customerID, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid customer ID", http.StatusBadRequest)
			return
		}

		customer, exists := store.GetCustomer(customerID)
		if !exists {
			http.Error(w, "Customer not found", http.StatusNotFound)
			return
		}

		JSONResponse(w, http.StatusOK, customer)
	}
}

func updateCustomerHandler(store CustomerStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["id"]
		customerID, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid customer ID", http.StatusBadRequest)
			return
		}

		var updatedCustomer Customer
		if err := json.NewDecoder(r.Body).Decode(&updatedCustomer); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := store.UpdateCustomer(customerID, updatedCustomer); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		JSONResponse(w, http.StatusOK, updatedCustomer)
	}
}

func deleteCustomerHandler(store CustomerStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["id"]
		customerID, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid customer ID", http.StatusBadRequest)
			return
		}

		if err := store.DeleteCustomer(customerID); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		JSONResponse(w, http.StatusOK, map[string]string{"message": "Customer deleted successfully"})
	}
}

func getCustomersHandler(store CustomerStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		JSONResponse(w, http.StatusOK, store.GetCustomers())
	}
}

func serveHTMLHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func main() {
	store := NewInMemoryCustomerStore()
	router := mux.NewRouter()

	router.HandleFunc("/", serveHTMLHandler).Methods("GET")
	router.HandleFunc("/customers", getCustomersHandler(store)).Methods("GET")
	router.HandleFunc("/customers", addCustomerHandler(store)).Methods("POST")
	router.HandleFunc("/customers/{id}", getCustomerHandler(store)).Methods("GET")
	router.HandleFunc("/customers/{id}", updateCustomerHandler(store)).Methods("PUT")
	router.HandleFunc("/customers/{id}", deleteCustomerHandler(store)).Methods("DELETE")

	fmt.Println("Server is starting on port 3000...")
	if err := http.ListenAndServe(":3000", router); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
