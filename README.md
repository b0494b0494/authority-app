# Authority App

This project is a sample application demonstrating an authorization system for a hierarchical directory structure. It utilizes OpenFGA for authorization, gRPC for communication, and PostgreSQL with `ltree` for managing directory hierarchies. The frontend is built with Next.js and Material-UI.

## Features

-   **Hierarchical Directory Management**: Organize directories in a tree-like structure.
-   **Role-Based Access Control (RBAC)**: Define roles (owner, editor, viewer) for users within directories.
-   **Permission Inheritance**: Permissions are inherited from parent directories to child directories.
-   **gRPC API**: High-performance communication between frontend and backend.
-   **OpenFGA Integration**: Externalized authorization checks using OpenFGA.

## Technologies Used

### Backend
-   **Go**: For the gRPC server implementation.
-   **OpenFGA**: An open-source authorization engine.
-   **PostgreSQL**: Relational database with `ltree` extension for hierarchical data.
-   **gRPC**: For inter-service communication.

### Frontend
-   **Next.js**: React framework for building user interfaces.
-   **TypeScript**: For type-safe JavaScript development.
-   **Material-UI (MUI)**: React UI framework.

### Development & Deployment
-   **Docker & Docker Compose**: For containerization and orchestration of services.

## Getting Started

### Prerequisites

-   Docker and Docker Compose installed.
-   Go (for backend development, if not using Docker for development).
-   Node.js and npm/yarn (for frontend development, if not using Docker for development).

### Setup

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/your-username/authority-app.git
    cd authority-app
    ```

2.  **Start the services**:
    The project uses Docker Compose to orchestrate the `postgres`, `openfga`, and `backend` services.
    ```bash
    docker-compose up -d --build
    ```
    This command will:
    -   Build the `backend` service Docker image.
    -   Pull `postgres` and `openfga` Docker images.
    -   Start all services in detached mode.
    -   Ensure all services pass their health checks before dependent services start.

3.  **Verify services status**:
    ```bash
    docker-compose ps
    ```
    All services (`postgres`, `openfga`, `backend`) should show `(healthy)` status.

### Accessing the Application

-   **Backend gRPC API**: Accessible on port `50051`.
-   **OpenFGA HTTP API**: Accessible on port `8080`.
-   **OpenFGA Playground**: Accessible on port `3000` (e.g., `http://localhost:3000`).
-   **PostgreSQL**: Accessible on port `5432`.

*(Frontend access details will be added once the frontend development is complete.)*

## Project Structure

```
.
├── backend/                # Go gRPC backend service
├── frontend/               # Next.js frontend application
├── proto/                  # Protocol buffer definitions
├── docker-compose.yml      # Docker Compose configuration
├── DESIGN.md               # Project design document
├── PROGRESS.md             # Development progress and troubleshooting log
└── README.md               # This file
```

## Development Notes

### Docker Health Check Troubleshooting

During initial setup, several issues were encountered with Docker Compose health checks, particularly for the `openfga` and `backend` services. These issues were resolved by:

-   Ensuring `grpc-health-probe` was correctly installed and available within the `backend` service container.
-   Configuring `grpc-health-probe` to target the correct addresses as reported by the services' logs.
-   Adjusting `start_period` for services to allow sufficient initialization time.
-   Implementing a gRPC health check service in `backend/main.go`.
-   Removing conflicting volume mounts in `docker-compose.yml` that prevented the backend executable from being found.

For detailed troubleshooting steps and resolutions, refer to `PROGRESS.md`.

## Contributing

(To be added)

## License

(To be added)
