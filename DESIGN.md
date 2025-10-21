# Project Design: Authority App

This document outlines the design and technical specifications for the Authority App.

## 1. Backend Design (gRPC, OpenFGA, PostgreSQL)

### 1.1. Overview

This project is a sample application that manages authorization in a hierarchical directory structure using OpenFGA, gRPC, and PostgreSQL (ltree).

Users can have one of the following roles for a directory: `owner`, `editor`, or `viewer`. Permissions of a parent directory are inherited by its child directories.

### 1.2. Architecture

```mermaid
graph TD
    subgraph "Client"
        Client[gRPC Client]
    end

    subgraph "Backend Service (Go)"
        Server[gRPC Server]
    end

    subgraph "Datastore"
        DB[(PostgreSQL)]
        FGA[OpenFGA Server]
    end

    Client -- gRPC Request --> Server
    Server -- Directory Structure Persistence --> DB
    Server -- Permission Check & Relationship Management --> FGA
```

- **gRPC Client**: A client that sends requests from the user.
- **gRPC Server (Go)**: Implements the business logic. It receives requests from the client, performs permission checks with OpenFGA, and persists data to PostgreSQL.
- **PostgreSQL**: Manages the hierarchical structure of directories using the `ltree` type.
- **OpenFGA Server**: Performs authorization checks and stores relationship tuples based on the authorization model.

### 1.3. gRPC API (proto)

`proto/directory/v1/directory.proto`

```protobuf
syntax = "proto3";

package directory.v1;

// Directory Service
service DirectoryService {
  // Create a directory
  rpc CreateDirectory(CreateDirectoryRequest) returns (CreateDirectoryResponse);
  // Share a directory (grant permissions)
  rpc ShareDirectory(ShareDirectoryRequest) returns (ShareDirectoryResponse);
  // Check permissions
  rpc CheckPermission(CheckPermissionRequest) returns (CheckPermissionResponse);
  // List child directories
  rpc ListChildren(ListChildrenRequest) returns (ListChildrenResponse);
}

// --- Messages ---

message User {
  string id = 1;
}

message Directory {
  string id = 1;
  string name = 2;
  string parent_id = 3; // Parent directory ID
}

// Create
message CreateDirectoryRequest {
  string name = 1;
  string parent_id = 2; // Empty for root
  User creator = 3;
}
message CreateDirectoryResponse {
  Directory directory = 1;
}

// Share
message ShareDirectoryRequest {
  User user = 1; // User to be granted permission
  string relation = 2; // "owner", "editor", "viewer"
  string directory_id = 3;
  User granter = 4; // User granting the permission
}
message ShareDirectoryResponse {
  bool success = 1;
}

// Check
message CheckPermissionRequest {
  User user = 1;
  string relation = 2; // "can_view", "can_edit", etc.
  string directory_id = 3;
}
message CheckPermissionResponse {
  bool allowed = 1;
}

// List
message ListChildrenRequest {
  string directory_id = 1;
  User user = 2;
}
message ListChildrenResponse {
  repeated Directory directories = 1;
}
```

### 1.4. DB Schema (PostgreSQL)

`directories` table

| Column Name | Type    | Description                               |
| ----------- | ------- | ----------------------------------------- |
| `id`        | UUID    | Primary Key                               |
| `name`      | TEXT    | Directory Name                            |
| `path`      | LTREE   | Hierarchical path from the root (e.g., `root.dir1.dir2`) |

### 1.5. OpenFGA Authorization Model

`authz.fga.yaml`

```yaml
model:
  schema_version: "1.1"
type_definitions:
  - type: user
  - type: directory
    relations:
      # Reference to the parent directory
      parent: directory
      # Roles directly assigned to each directory
      owner: user
      editor: user
      viewer: user
      # Permission inheritance and definitions
      # Can view: is a viewer, has a higher role, or can view the parent
      can_view: viewer or editor or owner or can_view from parent
      # Can edit: is an editor, has a higher role, or can edit the parent
      can_edit: editor or owner or can_edit from parent
      # Can share (grant permissions): is an owner or can share the parent
      can_share: owner or can_share from parent
      # Can delete/move: is an owner
      can_delete: owner
```

---

### 1.6. Docker Compose Service Configuration

The application services are orchestrated using Docker Compose, defining `postgres`, `openfga`, and `backend` services.

- **`postgres` Service**:
    - Uses `postgres:15-alpine` image.
    - Configured with a health check using `pg_isready` to ensure database readiness.
- **`openfga` Service**:
    - Uses `openfga/openfga:latest` image.
    - Configured to explicitly listen for gRPC on `0.0.0.0:8081` to ensure IPv4 connectivity within the Docker network.
    - Includes a health check using the bundled `grpc-health-probe` targeting the gRPC endpoint. The `start_period` is set to `30s` to allow sufficient time for initialization.
- **`backend` Service**:
    - Custom-built using `backend/Dockerfile`.
    - Configured with a health check using `grpc-health-probe` targeting its gRPC endpoint on `localhost:50051`.
    - The `Dockerfile` ensures `grpc-health-probe` is installed and the main executable is correctly placed for runtime.
    - Depends on `postgres` and `openfga` services being healthy before starting.

This Docker Compose setup ensures that all necessary backend components are properly initialized and communicating, providing a robust development environment.

---

## 2. Frontend Design (Next.js, Material-UI)

### 2.1. Technology Stack

The frontend will be built using a modern, robust, and scalable technology stack.

*   **Framework:** **Next.js (with App Router)**
    *   **Reasoning:** Provides a great developer experience with features like Server-Side Rendering (SSR), file-based routing, and API routes. It's a current industry trend for building high-performance web applications.
*   **Language:** **TypeScript**
    *   **Reasoning:** For type safety, improved code quality, and better developer tooling.
*   **UI Framework:** **Material-UI (MUI)**
    *   **Reasoning:** A comprehensive library of pre-built, accessible, and customizable React components. This will accelerate UI development and ensure a consistent, high-quality user experience.
*   **Styling:** **Emotion**
    *   **Reasoning:** Comes as the default styling engine with MUI. It allows for writing CSS styles with JavaScript.
*   **State Management:**
    *   **Global State:** **Zustand**. It's a simple and unopinionated state management library.
*   **Data Fetching:** **SWR** or **React Query (TanStack Query)**
    *   **Reasoning:** To handle remote data fetching, caching, and revalidation efficiently.
*   **Code Quality:**
    *   **Linting:** **ESLint**
    *   **Formatting:** **Prettier**

### 2.2. Directory Structure

A new `frontend` directory will be created to house the Next.js application.

```
/
├── frontend/
│   ├── app/
│   │   ├── api/          # API Routes for frontend-specific logic
│   │   ├── (pages)/      # Page routes (e.g., dashboard, directories)
│   │   │   ├── layout.tsx
│   │   │   └── page.tsx
│   │   ├── layout.tsx      # Root application layout
│   │   └── globals.css
│   ├── components/
│   │   ├── ui/             # Reusable UI components (Button, Input)
│   │   └── layout/         # Layout components (Header, Sidebar)
│   ├── lib/                # Helper functions, utilities, gRPC client
│   ├── public/             # Static assets
│   └── ... (config files)
├── docker-compose.yml
├── DESIGN.md
└── ...
```

### 2.3. Coding Conventions

*   **Component Naming:** `PascalCase` (e.g., `UserProfile.tsx`).
*   **File Naming:** `kebab-case` (e.g., `user-profile.tsx`), except for Next.js special files.
*   **ESLint & Prettier:** Code will be automatically linted and formatted.
*   **Imports:** Use absolute paths (`@/components/...`) configured via `tsconfig.json`.

### 2.4. Next Steps

- Set up the Next.js project in the `frontend` directory.
- Install and configure Material-UI.
- Set up the gRPC client to communicate with the backend.