# DocuDefence
## Project Overview

DocuDefense is a simple CRUD (Create, Read, Update, Delete) API developed in Golang for a fictional user management system. It manages user data, allowing clients to create, retrieve, update, and delete users through API endpoints. The project includes unit testing for each operation and is in the process of adding authentication and file upload functionality.

The API captures basic user details (first name, surname, email, and date of birth), and future enhancements will include file uploads, authentication, and additional security measures.

***

## Features Completed

- CRUD Operations:
  - Create a new user.
  - Retrieve all users.
  - Update an existing user by ID.
  - Delete a user by ID.

- Unit Testing:
  - Tests are provided to ensure each CRUD operation functions as expected.
 
***

## Features In Progress

- File Upload:
  - Functionality to upload a PDF file and associate it with a user.
  - Files will be stored on disk.

- Basic Authentication:
  -  Adding middleware for email and password authentication.
  -  Users will be required to authenticate before accessing endpoints.

- JWT Authentication:
  -  The current plan is to upgrade basic authentication to JWT-based authentication for better security.

- Dockerization:
   - Adding a Dockerfile for easy deployment and running the application in containers.

- Data Storage:
  -  Users are currently stored in memory, but files will be stored on disk once the upload functionality is added.
 
***

## Technology Stack

- Backend Language: Golang
- Router: [gorilla/mux](https://github.com/gorilla/mux)
- Testing Framework: Go's built-in testing package

***

## Project Structure (As Of Writing)

```bash
.
├── src
│   ├── handlers         # Contains the CRUD handler functions
│   │   ├── handlers.go  # Main CRUD logic
│   │   └── handlers_test.go  # Unit tests for the handlers
│   ├── models           # Defines the User struct
│   │   └── user.go      # User data structure
│   └── main.go          # Entry point for the application
├── go.mod               # Go modules (dependencies)
└── README.md            # Project documentation
```

***

## API Endpoints

| Method | Endpoint         | Description                | Payload                              |
|--------|------------------|----------------------------|--------------------------------------|
| GET    | `/users`          | Retrieves all users        | N/A                                  |
| POST   | `/users`          | Creates a new user         | `{ "id": "1", "first_name": "...", ... }` |
| PUT    | `/users/{id}`     | Updates an existing user   | `{ "first_name": "Updated", ... }`   |
| DELETE | `/users/{id}`     | Deletes a user by ID       | N/A       

```JSON
{
  "id": "1",
  "first_name": "Test",
  "surname": "User",
  "email": "testuser@code.com",
  "birthdate": "2024-10-23"
}
```

***

## Unit Testing

Unit tests are provided in the handlers_test.go file to ensure each handler (Create, Read, Update, Delete) is working as expected.
Running the Tests

To run the tests, navigate to the project folder and execute:

```bash
go test ./src/handlers
```

This will run all the unit tests and output the results.

***

## Planned Features

- File Upload: Add functionality for uploading and associating PDF files with users.
- Basic Authentication: Add authentication middleware using email and password.
- JWT Authentication: Secure the API with JWT tokens.
- Dockerization: Add a Dockerfile for easy deployment.
- Additional Features:
  - Pagination, filtering, and search functionality for users.
  - CORS middleware and rate-limiting to enhance security and scalability.

 ***

 ## Installation and Setup
1. Clone the repository

```bash
git clone https://github.com/your-username/DocuDefense.git
```

2. Install dependencies

Ensure that Golang is installed on your machine. Then, navigate to the project folder and run:

```bash
go mod tidy
```

This will install the required dependencies.

3. Run the application

Start the server:

```bash
go run src/main.go
```

The API will be running on http://localhost:8000.

4. Test the API

You can test the API endpoints using tools like Postman or curl. For example, to create a user:

```bash
curl -X POST -d '{"id":"1","first_name":"Test","surname":"User","email":"testuser@code.com","birthdate":"2024-10-23"}' http://localhost:8000/users
```

***

## Contributing

As this is a challenge posed to myself I won't be accepting any outside contributions for the moment.

***

## License

This project is licensed under the MIT License.
