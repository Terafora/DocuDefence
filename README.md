# DocuDefense
## Project Overview

DocuDefense is a CRUD (Create, Read, Update, Delete) API developed in Golang and functions as a secure way to upload and sace PDF documents, primarily contracts. It manages user data using MongoDB, allowing clients to create, retrieve, update, and delete users through API endpoints. The project includes unit testing for each operation and implements authentication functionality.

The API captures basic user details (first name, surname, email, and date of birth), with data securely stored in MongoDB. Future enhancements will include additional security measures.
***

## Features Completed

- CRUD Operations:
   - Create a new user.
   - Retrieve all users.
   - Update an existing user by ID.
   - Delete a user by ID.

- Unit Testing:
   - Tests are provided to ensure each CRUD operation functions as expected.
 
- File Upload:
       - Functionality to upload a PDF file and associate it with a user.
       - Note: Files will be stored on MongoDB, enhancing data management and retrieval.

    - Basic Authentication:
        - Adding middleware for email and password authentication.
        - Users will be required to authenticate before accessing endpoints.

     - JWT Authentication:
        - The current plan is to upgrade basic authentication to JWT-based authentication for better security.

- Features In Progress
   
   - Dockerization:
       - Adding a Dockerfile for easy deployment and running the application in containers.


***

## Technology Stack

 - Backend Language: Golang
 - Router: [gorilla/mux](https://github.com/gorilla/mux)
 - Testing Framework: Go's built-in testing package
 - Database: MongoDB for data storage

***

## Project Structure (As Of Writing)

```bash

.
├── src
│   ├── handlers         # Contains the CRUD handler functions
│   │   ├── handlers.go  # Main CRUD logic
│   │   ├── jwtmiddleware.go  # JWT Authentication middleware
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

### Example Payload for Creating a User

```JSON
{
  "first_name": "UpdatedName",
  "surname": "UpdatedSurname",
  "email": "newuser@example.com",
  "birthdate": "2024-11-01"
}
```

### Example Payload for Updating a User

```JSON
{
  "first_name": "UpdatedName",
  "surname": "UpdatedSurname",
  "email": "newuser@example.com",
  "birthdate": "2024-11-01"
}
```

***

## Error Handling

The API responds with appropriate status codes and messages when errors occur. Here are the common errors:

  - 401 Unauthorized: When the Authorization header is missing, invalid, or the token is expired.
  - 404 Not Found: When the user or resource you're trying to access is not found.
  - 500 Internal Server Error: For server-side issues.

Example response for an invalid JWT token:

```json

{
  "error": "Invalid or expired token"
}

```

***

## JWT Authentication Testing

To test JWT authentication:

  1. **Login to get the JWT Token**
     
Use the  `POST /login` endpoint to log in with your user credentials. You’ll receive a JWT token in the response.

  2. **Use the token in requests**

For protected routes like Update, Delete, and File Upload, add the JWT token to the Authorization header in this format:

```plaintext

Bearer <your_token_here>

```
  Example in Postman:
      - Add a new Header called `Authorization`
      - Set the value to: `Bearer <your_token_here>`

  3. **Testing a Protected Route**

Now, send a request to a protected endpoint like `DELETE /users/{id}` using the token from step 1.

If your token is valid, you’ll get a response like this:

```json

{
  "message": "User deleted successfully"
}

```

***

## File Upload Example

To upload a PDF file for a user, use the `/users/{id}/upload` endpoint. Here’s how you can do it in Postman:

  1. **Select POST method**
  
  Use `POST /users/{id}/upload` to target a specific user.

  2. **Add Authorization**
     
  Remember to include your JWT token in the Authorization header.

  3. **Upload the file**
      - In Postman, go to the **Body** tab.
      - Select **form-data**.
      - Add a new key named `contract` and set the type to **File**.
      - Upload a PDF file.

Example Response:

```json

{
  "message": "File uploaded successfully",
  "filename": "contract.pdf"
}

```

***

## Environment Variables

To keep sensitive data like the JWT secret key secure, use environment variables. Here’s how you can set it up:

  1. **Create a `.env` file** in the project root.
  2. Add the secret key:

```plaintext

JWT_SECRET=your_secret_key

```

3. Access the variable in your code:

In Go, use os.Getenv("JWT_SECRET") to retrieve the secret key.

Example:

```go

jwtKey := []byte(os.Getenv("JWT_SECRET"))

```

***

## JWT Token Expiry

Tokens in this API expire after a set time, currently 1 hour. If a token is expired, users will receive the following error:

```json

{
  "error": "Invalid or expired token"
}

```

To get a new token, the user must log in again using their email and password. Once logged in, a new token will be issued.

***

## Unit Testing

Unit tests are provided in the handlers_test.go file to ensure each handler (Create, Read, Update, Delete) is working as expected.

### Running the Tests

To run the tests, navigate to the project folder and execute:

```bash

go test ./src/handlers

```

This will run all the unit tests and output the results.

***

## Dependencies

Here’s a list of Go libraries used in the project:

  - [gorilla/mux](https://github.com/gorilla/mux): For routing and handling HTTP requests.
  - [JWT-go](https://github.com/golang-jwt/jwt): For generating and parsing JWT tokens for authentication.
  - [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver): For MongoDB database interactions.

***

## Planned Features

  - Dockerization: Add a Dockerfile for easy deployment.
  - Additional Features:
     - Pagination, filtering, and search functionality for users.
     - CORS middleware and rate-limiting to enhance security and scalability.

***

## Installation and Setup

1.Clone the repository:

```bash

git clone https://github.com/your-username/DocuDefense.git

```

2. Install dependencies:

Ensure that Golang is installed on your machine. Then, navigate to the project folder and run:

```bash

go mod tidy

```

This will install the required dependencies.

3. Run the application:

Start the server:

```bash

go run src/main.go

The API will be running on http://localhost:8000.

```

4. Test the API:

You can test the API endpoints using tools like Postman or curl. For example, to create a user:

```bash

curl -X POST -d '{"id":"1","first_name":"Test","surname":"User","email":"testuser@code.com","birthdate":"2024-10-23"}' http://localhost:8000/users

```

## Contributing

As this is a challenge posed to myself, I won't be accepting any outside contributions for the moment.

***

## License

This project is licensed under the MIT License.
