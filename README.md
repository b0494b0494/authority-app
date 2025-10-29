# Authority App

This is a full-stack application that demonstrates a simple authentication and authorization system using Go, Next.js, and OpenFGA.

## Features

- User registration and login with JWT authentication
- Protected routes that require authentication
- File upload with authorization checks using OpenFGA
- A simple and clean UI built with Next.js and Material-UI

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Node.js and npm

### Installation

1. Clone the repository
2. Install frontend dependencies
   ```bash
   cd frontend
   npm install
   ```
3. Run the application
   ```bash
   docker-compose up --build
   ```

## Usage

1. Open your browser and navigate to `http://localhost:3000`
2. Register a new user
3. Login with the registered user
4. Access the protected route and upload a file

## Technologies Used

- **Backend**: Go, Gin, GORM, OpenFGA
- **Frontend**: Next.js, React, Material-UI, Axios
- **Database**: SQLite
- **Containerization**: Docker, Docker Compose
