# CRM Backend Demo

A backend server built with Go for managing customer data in a CRM (Customer Relationship Management) system. This project allows users to perform CRUD (Create, Read, Update, Delete) operations on customer data through a RESTful API.

## Features

- **CRUD Operations**: Perform Create, Read, Update, and Delete operations on customer data.
- **RESTful API**: Built using Go and Gorilla Mux to handle RESTful routing.
- **JSON Responses**: All responses from the API are in JSON format for easy integration with front-end applications.
- **In-Memory Database**: Customer data is stored in an in-memory database (a map in Go).

## API Reference

### Get all customers

```http
GET /customers
```

#### Response

- `200 OK`: Returns a list of all customers.

### Get customer by ID

```http
GET /customers/{id}
```

#### Path Parameters

- `id` (integer): The ID of the customer.

#### Response

- `200 OK`: Returns the customer data.
- `404 Not Found`: If the customer with the specified ID does not exist.

### Add a new customer

```http
POST /customers
```

#### Request Body

- JSON object containing customer details:
  - `name` (string): The customer's name.
  - `role` (string): The customer's role (e.g., "Manager").
  - `email` (string): The customer's email address.
  - `phone` (string): The customer's phone number.
  - `address` (string): The customer's address.
  - `contacted` (boolean): Indicates whether the customer has been contacted.

#### Response

- `201 Created`: Returns the created customer data.

### Update customer by ID

```http
PUT /customers/{id}
```

#### Path Parameters

- `id` (integer): The ID of the customer.

#### Request Body

- JSON object containing updated customer details:
  - `name` (string): The customer's name.
  - `role` (string): The customer's role (e.g., "Manager").
  - `email` (string): The customer's email address.
  - `phone` (string): The customer's phone number.
  - `address` (string): The customer's address.
  - `contacted` (boolean): Indicates whether the customer has been contacted.

#### Response

- `200 OK`: Returns the updated customer data.
- `404 Not Found`: If the customer with the specified ID does not exist.

### Delete customer by ID

```http
DELETE /customers/{id}
```

#### Path Parameters

- `id` (integer): The ID of the customer.

#### Response

- `200 OK`: Confirms the deletion of the customer.
- `404 Not Found`: If the customer with the specified ID does not exist.

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/crm-backend-demo.git
   ```

2. Navigate to the project directory:

   ```bash
   cd crm-backend-demo
   ```

3. Install dependencies (if any) and build the project:

   ```bash
   go build
   ```

4. Run the application:

   ```bash
   go run main.go
   ```

## Run Locally

1. Start the server:

   ```bash
   go run main.go
   ```

2. The server will run on `http://localhost:3000`. You can use Postman, cURL, or any other HTTP client to interact with the API.

## Usage/Examples

### Adding a Customer

```bash
curl -X POST http://localhost:3000/customers \
-H "Content-Type: application/json" \
-d '{"name": "John Doe", "role": "Manager", "address": "123 Main St", "phone": "555-555-5555", "email": "john@example.com", "contacted": true}'
```

### Getting All Customers

```bash
curl http://localhost:3000/customers
```

### Getting a Customer by ID

```bash
curl http://localhost:3000/customers/1
```

### Updating a Customer

```bash
curl -X PUT http://localhost:3000/customers/1 \
-H "Content-Type: application/json" \
-d '{"name": "Jane Doe", "role": "Supervisor", "address": "456 Elm St", "phone": "555-555-5555", "email": "jane@example.com", "contacted": false}'
```

### Deleting a Customer

```bash
curl -X DELETE http://localhost:3000/customers/1
```

## Tech Stack

- **Go**: The main programming language used.
- **Gorilla Mux**: Router package for handling API routes.

## Running Tests

To run tests for the application, you can use the following command:

```bash
go test
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

## Acknowledgements

- [Gorilla Mux](https://github.com/gorilla/mux) - For providing a robust routing library for Go.
- **Udacity** - For providing the excellent Go course that inspired and guided the development of this project.
