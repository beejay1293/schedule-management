### Schedule Management — gRPC Based Appointment Booking

This project implements a simple appointment scheduling system using **Go (backend)**, **React + TypeScript (frontend)**, and **gRPC** for communication.
It demonstrates handling concurrent booking, proper database locking, and frontend integration via **gRPC-Web**.

---

### 🏗️ Architecture Overview

```
Frontend (React + TypeScript)
   │
   ▼
gRPC-Web Proxy (grpcwebproxy)
   │
   ▼
Backend (Go + gRPC + PostgreSQL)
```

* **Frontend:** React + TypeScript using `protobuf-ts` for generated gRPC client stubs
* **Backend:** Go gRPC service exposing endpoints for creating, listing, and deleting appointments
* **Database:** PostgreSQL
* **Proxy:** `grpcwebproxy` for bridging gRPC-Web requests to native gRPC server
* **Migrations:** Automatically run at app startup if `RUN_DB_MIGRATIONS=true`

---

### ⚙️ Setup Instructions

#### 1. Clone and Build the Proxy

```bash
git clone https://github.com/improbable-eng/grpc-web.git
cd grpc-web/go/grpcwebproxy
go build -o grpcwebproxy
sudo mv grpcwebproxy /usr/local/bin/
```

This installs the `grpcwebproxy` binary globally.

#### 2. Environment Variables

Create a `.env` file in the backend root:

```bash
POSTGRES_DB_URL=postgres://postgres:postgres@localhost:5432/schedule_db?sslmode=disable
RUN_DB_MIGRATIONS=true
```

Create a `.env` file in the frontend root:

```bash
VITE_GRPC_PROXY_URL=http://localhost:8080
```

#### 3. Start PostgreSQL

- Make sure PostgreSQL is running and accessible on port `5432`.
- Also, create a database named `schedule_db` if it doesn’t already exist.

#### 4. Run Backend

```bash
make start-app
```

> 🧩 The application automatically runs database migrations if `RUN_DB_MIGRATIONS` is set to `true`.

#### 5. Run gRPC-Web Proxy

```bash
grpcwebproxy \
  --backend_addr=localhost:50051 \
  --run_tls_server=false \
  --allow_all_origins \
  --server_http_debug_port=8080
```

This exposes the gRPC service over HTTP on port `8080` for your web client.

#### 6. Run Frontend

```bash
cd frontend
npm install
npm run dev
```

---

### 🧬 Code Generation

The `proto` files are already included in both backend and frontend directories, and the generated code is checked into version control.

You **don’t need to regenerate** the code during setup, but for reference:

* **Backend (Go):**

  ```bash
  cd backend
  make gen-proto
  ```

* **Frontend (TypeScript using protobuf-ts):**

  ```bash
  cd frontend
  npm run gen:protobuf-ts
  ```

  (Command is defined in `package.json`)

---

### 🧪 Running Tests

The backend includes tests for:

* Appointment creation
* Conflict prevention
* Concurrency handling
* Input validation

Run all tests with:

```bash
cd backend
make test
```

---

### 💡 Implementation Summary

- Built a **gRPC service** in Go to handle appointments.  
- Used **PostgreSQL** with a **unique constraint on appointment datetime** and **transactional inserts** to safely prevent double bookings.  
- Implemented **database-level concurrency control** instead of in-memory locks (like `sync.Mutex`), ensuring consistency across multiple service instances in distributed environments.  
- Frontend communicates through **gRPC-Web**, bridged using `grpcwebproxy`.  
- Used **protobuf-ts** for strongly typed client stubs in TypeScript.  
- Automatically runs migrations on startup for a smoother local setup.  
- Organized a comprehensive test suite to verify **CRUD operations**, **conflict prevention**, and **concurrency safety**.



## 🧠 Assumptions & Limitations

* Appointments are validated to prevent overlapping time ranges
* All time values stored in UTC
* No authentication (out of scope for this assessment)
* Real-time updates are implemented via gRPC streaming; however, in local testing the stream can occasionally break, especially when creating new appointments. Delete events are generally more stable, as they trigger fewer chunked HTTP responses.
* Mutex-based synchronization not used in the service layer because it wouldn't provide safety in a multi-instance deployment; proper concurrency control should rely on the database
