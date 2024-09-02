package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

// customers is a map that acts as our in-memory database,
// where the key is an integer ID and the value is a Customer.
var customers map[int]Customer

// lastID keeps track of the last assigned ID for customers.
var lastID int

// addCustomer adds a new customer to the customers map.
func addCustomer(c Customer) int {
	lastID++
	c.ID = lastID
	customers[c.ID] = c
	return c.ID
}

// getCustomer retrieves a customer by their ID from the customers map.
func getCustomer(id int) (Customer, bool) {
	customer, exists := customers[id]
	return customer, exists
}

// updateCustomer updates an existing customer's data in the customers map.
func updateCustomer(id int, c Customer) error {
	if _, exists := customers[id]; !exists {
		return fmt.Errorf("customer with ID %d not found", id)
	}
	c.ID = id // Ensure the ID is maintained
	customers[id] = c
	return nil
}

// deleteCustomer removes a customer from the customers map by their ID.
func deleteCustomer(id int) error {
	if _, exists := customers[id]; !exists {
		return fmt.Errorf("customer with ID %d not found", id)
	}
	delete(customers, id)
	return nil
}

// getCustomers returns the entire customers map.
func getCustomers() map[int]Customer {
	return customers
}

// setJSONResponseHeaders sets the Content-Type to application/json
// and writes the provided HTTP status code to the response.
func setJSONResponseHeaders(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
}

// addCustomerHandler handles POST requests to add a new customer.
// It reads the request body, decodes the JSON data, adds the customer,
// and responds with the new customer's data.
func addCustomerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newEntry Customer

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(reqBody, &newEntry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := addCustomer(newEntry)

	w.Header().Set("Location", fmt.Sprintf("/customers/%d", id))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newEntry)
}

// getCustomerHandler handles GET requests to retrieve a specific customer by ID.
func getCustomerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	customerID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusNotFound)
		return
	}

	customer, exists := getCustomer(customerID)
	if !exists {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customer)
}

// updateCustomerHandler handles PUT requests to update an existing customer's data.
func updateCustomerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	customerID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	_, exists := getCustomer(customerID)
	if !exists {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	var updatedCustomer Customer
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(reqBody, &updatedCustomer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = updateCustomer(customerID, updatedCustomer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedCustomer)
}

// deleteCustomerHandler handles DELETE requests to remove a customer by ID.
func deleteCustomerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	customerID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusNotFound)
		return
	}

	err = deleteCustomer(customerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Customer deleted successfully"})
}

// getCustomersHandler handles GET requests to retrieve all customers.
func getCustomersHandler(w http.ResponseWriter, r *http.Request) {
	setJSONResponseHeaders(w, http.StatusOK)

	customers := getCustomers()
	json.NewEncoder(w).Encode(customers)
}

// serveHTMLHandler serves the index.html file at the "/" route.
func serveHTMLHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

// seedMockDB populates the customers map with initial data.
func seedMockDB() {
	customers[1] = Customer{
		ID:        1,
		Name:      "John Doe",
		Role:      "Manager",
		Address:   "123 Main St",
		Phone:     "555-555-5555",
		Email:     "johndoe@example.com",
		Contacted: true,
	}
	customers[2] = Customer{
		ID:        2,
		Name:      "Jane Doe",
		Role:      "Employee",
		Address:   "456 Main St",
		Phone:     "555-555-5555",
		Email:     "janedoe@example.com",
		Contacted: false,
	}
	customers[3] = Customer{
		ID:        3,
		Name:      "Jim Doe",
		Role:      "Supervisor",
		Address:   "789 Main St",
		Phone:     "555-555-5555",
		Email:     "jimdoe@example.com",
		Contacted: true,
	}
	customers[4] = Customer{
		ID:        4,
		Name:      "Jill Doe",
		Role:      "HR",
		Address:   "1011 Main St",
		Phone:     "555-555-5555",
		Email:     "jilldoe@example.com",
		Contacted: false,
	}

	// Update lastID based on seeded data
	lastID = len(customers)
}

func main() {
	customers = make(map[int]Customer)

	seedMockDB()

	router := mux.NewRouter()

	router.HandleFunc("/", serveHTMLHandler).Methods("GET")
	router.HandleFunc("/customers", getCustomersHandler).Methods("GET")
	router.HandleFunc("/customers", addCustomerHandler).Methods("POST")
	router.HandleFunc("/customers/{id}", getCustomerHandler).Methods("GET")
	router.HandleFunc("/customers/{id}", updateCustomerHandler).Methods("PUT")
	router.HandleFunc("/customers/{id}", deleteCustomerHandler).Methods("DELETE")

	fmt.Println("Server is starting on port 3000...")
	if err := http.ListenAndServe(":3000", router); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
