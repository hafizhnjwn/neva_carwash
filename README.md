# Final Project PBKK

## Overview

This web application, built with Go, allows users to create, read, update, and delete code vehicles. It features user authentication, secure password storage, and a responsive interface designed with Tailwind CSS.

## Features

- User Registration and Authentication 
- Create, Read, Update, and Delete (CRUD) Operations for Code vehicles
- Secure Password Hashing with bcrypt
- Session Management with JWT Tokens
- Responsive Design Using Tailwind CSS


## Prerequisites

- Go (version 1.23 or later)
- Git
- SQLite3

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/azarelga/final-project-pbkk.git
cd final-project-pbkk
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Configure Environment Variables

Create a `.env` file in the project root with the following plat:

```
PORT=8080
SESSION_KEY=your_very_long_and_random_secret_key_here
DATABASE_PATH=./vehicles.db
```

Note: Replace `your_very_long_and_random_secret_key_here` with a secure, random string of at least 32 characters. You can generate one using:

```bash
openssl rand -base64 32
```

## Running the Application

### Development Mode

```bash
go run main.go
```

### Build and Run

```bash
go build -o code-vehicles
./code-vehicles
```

## Project Structure

```
.
├── main.go                # Main application entry point
├── .env                   # Environment configuration
├── go.mod                 # Go module dependencies
├── go.sum                 # Go module checksums
├── handlers/              # HTTP request handlers
│   ├── auth.go
│   └── vehicles.go
├── repositories/          # Database access layers
│   ├── user.go
│   └── vehicles.go
├── services/              # Business logic
│   ├── user.go
│   └── vehicles.go
├── middleware/            # Middleware functions
│   └── checkAuth.go
├── database/              # Database initialization and migrations
│   ├── db.go
│   ├── migrate.go
│   └── loadenvs.go
├── templates/             # HTML templates
│   ├── header.html
│   ├── footer.html
│   ├── home.html
│   ├── login.html
│   ├── register.html
│   ├── list.html
│   ├── mylist.html
│   ├── create.html
│   ├── edit.html
│   └── viewvehicle.html
├── .gitignore             # Git ignore file
├── README.md              # Project documentation
└── vehicles.db            # SQLite database (generated at runtime)
```
## CRUD Implementation
The application implements CRUD (Create, Read, Update, Delete) operations for code vehicles:

- **Create**: Users can create new code vehicles using the `Createvehicle` handler in `vehicles.go`. The HTML form for creating vehicles is in `create.html`.

- **Read**: Users can read existing vehicles through several handlers in `vehicles.go`:
  - `GetvehiclesByStatus` retrieves vehicles filtered by status.
  - `GetvehiclesByUserID` lists vehicles created by the current user.
  - `GetvehicleByID` displays detailed information about a specific vehicle.
- **Update**: Users can update their vehicles using the `Updatevehicle` handler in `vehicles.go`. The edit form is provided in `edit.html`.

- **Delete**: Users can delete their vehicles using the `Deletevehicle` handler in `vehicles.go`.

These handlers interact with the `vehicleService` in `vehicles.go`, which uses the `vehicleRepository` in `vehicles.go` for database operations. Authentication and authorization are managed by middleware in `checkAuth.go` to ensure that only authorized users can perform these actions.

## Security Considerations

- Password Hashing: User passwords are securely hashed using bcrypt.
- Session Management: Authentication is managed with JWT tokens stored in cookies.
- Protected Routes: Middleware ensures that certain routes are only accessible to authenticated users.

## Customization

### Changing the Database
By default, the application uses SQLite. To switch to a different database (e.g., MySQL or PostgreSQL), update the database configuration in db.go and adjust the GORM dialect accordingly.

### Styling
The application uses Tailwind CSS for styling. You can customize the styles by modifying the HTML templates in the templates directory.

## Demo Project
https://youtu.be/3ma7tgqgtKg