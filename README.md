# DocuDefense
## Project Overview

DocuDefense is a full-stack application that functions as a secure document management system, specifically designed for handling contracts in PDF format. The backend, built in Golang, provides a secure API for CRUD operations on users and file uploads, leveraging MongoDB for data storage. The frontend, built in React, offers a simple interface for users to manage their profile and documents.

The API supports user authentication using JWT, ensuring that only authenticated users can perform actions like updating or deleting profiles, and uploading files.
***

## Features Completed
### Backend Features:

- **User CRUD Operations:**
  - Create, Retrieve, Update, and Delete users.
- **File Upload:**
  - Upload PDF files and associate them with users.
- **Authentication:**
  - Basic login functionality using email and password.
  - JWT-based token authentication for protecting routes.
- **JWT Authentication:**
  - Middleware to protect sensitive routes using JWT tokens.
- **Docker**

### Frontend Features:

- **User Management:**
  - View a list of users.
  - Create new users.
  - Log in to view and update your profile.
  - Delete your account.
- **File Upload:**
  - Upload PDF files to be associated with your user profile.


***

## Technology Stack

 - **Backend:**
    - Language: **Golang**
    - Router: [gorilla/mux](https://github.com/gorilla/mux)
    - Testing Framework: Go's built-in testing package
    - Database: **MongoDB** for data storage
    - Authentication: **JWT**
    - File Handling: **Multipart/FormData** file uploads
  
 - **Frontend:**
   - Framework: **React**
   - Styling: **SCSS Modules**

***

## Project Structure (As Of Writing)

```bash

.
├── frontend
│   ├── src
│   │   ├── components              # React Components (Login, UserList, etc.)
│   │   ├── services                # API service functions for frontend
│   │   ├── App.js                  # Main application entry for React
│   │   └── index.js                # Entry point for React rendering
├── backend
│   ├── src
│   │   ├── handlers                # Golang CRUD handler functions
│   │   ├── models                  # Golang model definitions (User)
│   │   ├── jwtmiddleware.go        # JWT middleware for protected routes
│   │   ├── middleware.go           # Basic auth middleware (to be deprecated)
│   │   └── main.go                 # Main application entry point
│   ├── go.mod                      # Go modules (dependencies)
│   └── .env                        # Environment variables (e.g., MongoDB URI, JWT_SECRET)
├── Dockerfile                      # Dockerfile for Docker support (coming soon)
├── docker-compose.yml              # Compose file for running MongoDB and the backend
├── README.md                       # Project documentation


```

***

## API Endpoints
| Method | Endpoint         | Description                | Payload                              |
|--------|------------------|----------------------------|--------------------------------------|
| GET    | `/users`          | Retrieves all users        | N/A                                  |
| POST   | `/users`          | Creates a new user         | `{ "id": "1", "first_name": "...", ... }` |
| PUT    | `/users/{id}`     | Updates an existing user   | `{ "first_name": "Updated", ... }`   |
| DELETE | `/users/{id}`     | Deletes a user by ID       | N/A 
| POST    | `/login`    | Logs in and retrieves a JWT token   | `{ "email": "...", "password": "..."}`   |
| POST | `/users/{id}/upload`     | Uploads a PDF file for a user      | Multipart form with `contract` key 

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

  - Additional Features:
     - Pagination, filtering, and search functionality for users.
     - CORS middleware and rate-limiting to enhance security and scalability.

***

## Running the Project
Prerequisites:

 - **Golang** installed (1.16+)
 - **MongoDB** instance running locally or remotely
 - **Node.js** (v14+) for the React frontend
 - **Docker** (optional, for containerization)

Steps to Run:

 1. **Clone the repository:**

 ```bash

git clone https://github.com/your-username/DocuDefense.git
cd DocuDefense
```

2. **Set up Backend:**

- Ensure you have a .env file with your MongoDB URI and JWT secret.
- Install Go dependencies:

```bash

go mod tidy

```

- Run the backend:

```bash

    go run src/main.go
```

3. **Set up Frontend:**

- Navigate to the frontend directory:

```bash

        cd frontend
        npm install
        npm start
```

- The frontend should now be running at `http://localhost:3000`.

***

## Docker Setup (Optional)

If you prefer to use Docker for containerizing the application, you can follow these steps:

**1. Dockerfile**

The Dockerfile is already provided in the project for both the backend and frontend. The Dockerfile is located at the root of the `backend`  directory.

**2. Build Docker Image**

To build the Docker image for the backend, navigate to the backend directory and run:

```bash

docker build -t docudefense-backend .

```
**Work In Progress**

***

## Contributing

As this is a challenge posed to myself, I won't be accepting any outside contributions for the moment.

***

## License

This project is licensed under the MIT License.
