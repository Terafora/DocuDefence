# DocuDefense
## Project Overview

DocuDefense is a secure document management system that enables efficient handling and storage of sensitive files, specifically designed for contract PDFs. Built as a full-stack application with Golang and MongoDB for the backend and React for the frontend, DocuDefense offers user authentication, document uploading, and version control.

The API supports JWT-based authentication, allowing only authorized users to manage their profile and associated documents. The Docker setup enables easy deployment of both frontend and backend services.
***

## Features Completed
### Backend Features:

- **User CRUD Operations:**
  - Create, Retrieve, Update, and Delete user accounts.
- **File Upload:**
  - Upload PDF files and associate them with users with automatic versioning.
- **Authentication:**
  - Basic login functionality using email and password.
  - JWT-based token authentication for protecting routes.
- **JWT Authentication:**
  - Middleware to protect sensitive routes using JWT tokens.
- **User Search, Pagination and Filtering**
  - Fetch all users with pagination and search capabilities.
- **Docker Support:**
  - Backend and frontend Dockerfiles and a Docker Compose file to orchestrate services

### Frontend Features:

- **User Management:**
  - View a list of users.
  - Create new users.
  - Log in to view, update and your profile.
  - Delete your account.
- **Document Management:**
  - Upload PDF files to be associated with your user profile with versioning for each document.
  - Preview, download and delete files.
  - View a list of all users with search functionality and pagination.
- **Dynamic User Interface:**
  - Responsive layout with animations, custom modals, and user feedback messages.


***

## Technology Stack

 - **Backend:**
    - Language: **Golang**
    - Router: [gorilla/mux](https://github.com/gorilla/mux)
    - Testing Framework: Go's built-in testing package
    - Database: **MongoDB** for data storage
    - Authentication: **JWT**
    - File Handling: **Multipart/FormData** for PDF uploads
  
 - **Frontend:**
   - Framework: **React**
   - Styling: **SCSS Modules**
   - PDF Preview: **pdfjs-dist** for in-browser rendering

- **Deployment:**
  - **Docker** for containerization of backend and frontend.
  - **Docker Compose** for service orchestration

***

## Project Structure (As Of Writing)

```bash

.
├── frontend
│   ├── src
│   │   ├── components              # React Components (AuthPanel, UserDashboard, etc.)
│   │   ├── services                # API service functions for frontend
│   │   ├── App.js                  # Main React application file
│   │   └── index.js                # React entry point
├── backend
│   ├── src
│   │   ├── handlers                # Go handlers for CRUD, auth, and file management
│   │   ├── models                  # Go model definitions (User, Document)
│   │   ├── jwtmiddleware.go        # JWT middleware
│   │   ├── middleware.go           # Basic auth middleware
│   │   └── main.go                 # Main backend entry point
│   ├── go.mod                      # Go modules (dependencies)
│   ├── Dockerfile                  # Dockerfile for backend
│   └── .env                        # Environment variables (MongoDB URI, JWT_SECRET)
├── docker-compose.yml              # Compose file for frontend and backend services
├── README.md                       # Project documentation
                    # Project documentation


```

***

## API Endpoints

| Method | Endpoint                       | Description                            | Payload                                      |
|--------|--------------------------------|----------------------------------------|----------------------------------------------|
| GET    | `/users`                       | Retrieves all users                    | N/A                                          |
| POST   | `/users`                       | Creates a new user                     | `{ "first_name": "First", "surname": "Last", "email": "email@example.com", "birthdate": "YYYY-MM-DD", "password": "password" }` |
| PUT    | `/users/{id}`                  | Updates an existing user by ID         | `{ "first_name": "Updated", "surname": "Name", "email": "updated@example.com", "birthdate": "YYYY-MM-DD", "password": "newpassword" }` |
| DELETE | `/users/{id}`                  | Deletes a user by ID                   | N/A                                          |
| POST   | `/login`                       | Logs in a user and retrieves JWT token | `{ "email": "user@example.com", "password": "password" }` |
| POST   | `/users/{id}/upload`           | Uploads a PDF file for a user          | Multipart form with `contract` key           |
| GET    | `/users/{id}/files`            | Retrieves files uploaded by user       | N/A                                          |
| GET    | `/users/{id}/files/{filename}/download` | Downloads a specific file by filename for a user | N/A |
| DELETE | `/users/{id}/files/{filename}/delete` | Deletes a specific file by filename for a user | N/A |

### Example Payloads

#### Creating a New User
```json
{
  "first_name": "John",
  "surname": "Doe",
  "email": "john.doe@example.com",
  "birthdate": "1990-01-01",
  "password": "securepassword123"
}
```

### Updating an Existing User

```json
{
  "first_name": "John",
  "surname": "Smith",
  "email": "john.smith@example.com",
  "birthdate": "1990-01-01",
  "password": "newpassword456"
}

```
## JWT Authentication

1. **Login**: Use the `POST /login` endpoint to authenticate and retrieve a JWT token.
2. **Authorization**: Include the JWT token in the `Authorization` header for protected routes in the format:

```plaintext

Authorization: Bearer <your_token_here>


```

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

### File Upload and Download

- **File Upload**: Use `POST /users/{id}/upload` with a multipart form data containing the key contract for the PDF file.
- **File Download**: Use `GET /users/{id}/files/{filename}/download` to download a specific file.


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
