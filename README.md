# DocuDefense
## Project Overview

DocuDefense is a secure document management system that enables efficient handling and storage of sensitive files, specifically designed for contract PDFs. Built as a full-stack application with Golang and MongoDB for the backend and React for the frontend, DocuDefense offers user authentication, document uploading, and version control.

The API supports JWT-based authentication, allowing only authorized users to manage their profile and associated documents. The Docker setup enables easy deployment of both frontend and backend services.
***

## Features
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

### User Management

| Method | Endpoint              | Description                          | Auth       |
|--------|------------------------|--------------------------------------|------------|
| GET    | `/users`              | Retrieve all users (pagination)      | No         |
| POST   | `/users`              | Create a new user                    | No         |
| GET    | `/users/email`        | Get user ID by email                 | Yes (JWT)  |
| PUT    | `/users/{id}`         | Update user by ID                    | Yes (JWT)  |
| DELETE | `/users/{id}`         | Delete user by ID                    | Yes (JWT)  |

### Authenitcation

| Method | Endpoint | Description               | Auth |
|--------|----------|---------------------------|------|
| POST   | `/login` | Log in and get JWT token  | No   |

### Dociment Management

| Method | Endpoint                           | Description                     | Auth       |
|--------|------------------------------------|---------------------------------|------------|
| POST   | `/users/{id}/upload`               | Upload a PDF file               | Yes (JWT)  |
| GET    | `/users/{id}/files`                | Get all files for a user        | Yes (JWT)  |
| GET    | `/users/{id}/files/{filename}/download` | Download a file                | Yes (JWT)  |
| DELETE | `/users/{id}/files/{filename}/delete`   | Delete a file                   | Yes (JWT)  |



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

#### Updating an Existing User

```json
{
  "first_name": "John",
  "surname": "Smith",
  "email": "john.smith@example.com",
  "birthdate": "1990-01-01",
  "password": "newpassword456"
}

```
#### Login

```JSON
{
  "email": "john.doe@example.com",
  "password": "password123"
}
```

***

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
     - Admin Actions for user management and audit logs.

***

## Running the Project
### Prerequisites:

 - **Golang** installed (1.20 or newer)
 - **MongoDB** instance running locally or remotely
 - **Node.js** (v18+) for the React frontend
 - **Docker** (optional, for containerization)

### Environment Variables

Create a `.env` file in the backend directory with the following variables:

```plaintext

MONGODB_URI=<your_mongodb_uri>
JWT_SECRET=<your_jwt_secret>

```

### Running Locally

**1. Clone the Repository**
  ```bash
  git clone https://github.com/your-username/DocuDefense.git
  cd DocuDefense
  ```
**2. Backend Setup**
   - Navigate to the backend directory and install dependencies:
   ```bash
    cd backend
    go mod tidy
   ```
   - Run the backend server:
  ```bash
  go run src/main.go
  ```
**3. Frontend Setup**
  - Navigate to the frontend directory and install dependecies:
  ```bash
    cd ../frontend
    npm install
  ```
  - Start the frontend server:
  ```bash
  cd ../frontend
  npm install
  ```
4. Access the application at `http://localhost:3000`.

### Running with Docker
**1. Build and Run with Docker Compose**
   ```bash
    docker-compose up --build
   ```
   This command starts both frontend and backend services, exposing the frontend on port `3000` and backend on port `8000`.

***

## Usage

DocuDefense is a secure document management platform that offers a range of features to help users manage, secure, and organize their files. Below are the instructions on how to use each feature within the application.

### 1. Registering and Logging In

- **Register**: New users can create an account by clicking the "New User?" link in the navigation bar. This opens the registration form where users should provide their first name, surname, email, birthdate, and a password. Once registered, an alert confirms account creation, allowing the user to log in.
- **Login**: Registered users can log in by entering their email and password. A JWT token is generated and stored in the browser’s local storage upon successful login, enabling access to restricted features.
- **Logout**: To log out, simply click the "Logout" button in the navigation bar. This clears the JWT token from local storage and restricts access to protected features.

### 2. User Dashboard

- After logging in, users can access their **User Dashboard**, which provides options to upload, manage, and preview files as well as update or delete their account.
- **File Upload**: Users can upload PDF files securely. To upload a file, select a PDF file using the file input, then click "Upload PDF". The application version controls each file, saving new versions with each upload.
- **File Management**:
    - **Preview**: Click the "Preview" button to open a modal with a PDF preview of the file.
    - **Download**: Download any version of a file by clicking "Download". This will download the specific file version as a PDF.
    - **Delete**: To delete a file, click the "Delete" button. A confirmation prompt appears before deletion. Previous versions of the file are also listed, with options to preview or download specific versions.
- **Account Management**:
    - **Update Profile**: Users can update their personal information, including first name, surname, email, and password. This can be done by clicking the "Edit Profile" button, filling in the fields, and saving changes.
    - **Delete Account**: If a user wants to delete their account and all associated files, they can click the "Delete Account" button. A confirmation prompt ensures that this action is intentional.

### 3. User Directory and Search

- **User List**: The "All Users" page provides a directory of all users. Users can browse through the list using pagination controls or search for specific users.
- **Search Users**: Enter a first name or surname in the search bar to find specific users. Search results are paginated, and users can navigate through them using "Next" and "Previous" buttons.

### 4. File Version Control

- DocuDefense offers version control for uploaded files. Each time a user uploads a file with the same name, the system saves it as a new version rather than overwriting the original.
- **Viewing Previous Versions**: Users can expand a file’s entry in the "My Files" section to view and manage previous versions. Each version can be previewed or downloaded.
- **Deleting Versions**: When deleting a file, all associated versions are also deleted from the system.

### 5. JWT Authentication and Authorization

- **Token Storage**: The application uses JWT tokens for session management. Tokens are stored securely in the browser's local storage and are automatically attached to authenticated requests.
- **Token Expiry**: Each token expires after 1 hour. Users will need to re-login once their session expires.
- **Protected Routes**: Routes such as `/dashboard`, `/allusers`, and user-specific file management routes require a valid JWT token, ensuring only authorized users can access these features.

### 6. Document Management Features

- **File Upload & Storage**: Users can securely upload files, with each file stored as a document entry in MongoDB. The files themselves are stored on disk.
- **Advanced Search**: Users can search by first name or surname to locate specific users.
- **Tagging and Categorizing**: Users can categorize files during upload to make organization simpler. (Note: If not currently implemented, consider adding this feature in the future.)
- **File Permissions**: By default, users can only view, edit, or delete their own files.

### 7. Front-End Navigation

- The front-end is organized for seamless navigation between the homepage, about section, user dashboard, and user directory.
- **Navbar**: The navbar provides quick access to key sections of the application, with login/logout buttons dynamically updating based on the user's session status.
- **Animation and Responsiveness**: The front-end is designed to be responsive, with animations that enhance user experience, such as scroll-based transformations.

### 8. Admin Actions (Future Enhancements)

- **Admin User Management**: Consider implementing an admin feature for managing all users, roles, and permissions across the platform.
- **Audit Logs**: Adding logging features to track user activity (such as file uploads, downloads, and deletions) could be a valuable enhancement.

### Quick Summary of Commands and Endpoints

Below is a summary of key commands and endpoints available in the application for easy reference.

- **File Management**:
    - Upload a file: `/users/{id}/upload`
    - View user files: `/users/{id}/files`
    - Download file: `/users/{id}/files/{filename}/download`
    - Delete file: `/users/{id}/files/{filename}/delete`

- **User Management**:
    - Register a new user: `/users` (POST)
    - Login user: `/login` (POST)
    - Update user details: `/users/{id}` (PUT)
    - Delete user: `/users/{id}` (DELETE)
    - Get all users or search: `/api/users` (GET with search query)

***

## Error Handling

The API provides meaningful HTTP status codes and error messages:

  - `401 Unauthorized`: Invalid or missing JWT token.
  - `404 Not Found`: Resource does not exist.
  - `500 Internal Server Error`: Server-side error.

### Example Response for Unauthorized Access

```json

{
  "error": "Invalid or expired token"
}

```

***
### Bugs

- **PDF Preview** shows images rotated 180 degrees seemingly randomly and hasn't been addressed just yet.
- **Token Expirey** Won't log users out but will stop them from being able to view account related information and documents. Currently the token expires after and hour.
- **Styling Inconsistencies** appear between modals and such and need to be touched up in the future.

***

### Future Improvements

- **Unit Testing**: Comprehensive test coverage for all API endpoints.
- **Enhanced Security**: Rate limiting and stricter CORS policy.
- **Additional Features**: Pagination for file versions, advanced search filters.

***

## License

This project is licensed under the MIT License.
