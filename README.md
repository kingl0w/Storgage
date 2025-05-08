# Storage WebApp

A self-hosted storage web application for securely sharing files between your work and personal machines.

This project is built using:
- **Go** for the backend API
- **Svelte** for the frontend interface
- **Azure Blob Storage** for file storage
- An invitation-based authentication system

---

## Features

- Upload and access files through a clean, responsive UI
- Admin-controlled access via invitation tokens
- Persistent file storage with Azure Blob Storage
- Full-stack Dockerized setup

---

##  Setup & Installation

### Prerequisites

- [Docker](https://www.docker.com/)
- An [Azure Blob Storage](https://azure.microsoft.com/en-us/products/storage/blobs/) account

---

### Environment Variables

Create two `.env` files — one in the `backend/` folder and one in the `frontend/` folder.

#### `backend/.env`

env
PORT=8080
JWT_SECRET=your_jwt_secret

# Azure Blob Storage
- AZURE_STORAGE_ACCOUNT=your_azure_account_name
- AZURE_STORAGE_ACCESS_KEY=your_azure_access_key
- AZURE_STORAGE_CONTAINER=your_container_name

# Admin Login
- ADMIN_USERNAME=your_admin_username
- ADMIN_PASSWORD=your_admin_password

# PostgreSQL
- POSTGRES_USER=your_db_username
- POSTGRES_PASSWORD=your_db_password
- POSTGRES_DB=cloudstorage
- DATABASE_URL=postgres://your_db_username:your_db_password@db:5432/cloudstorage?sslmode=disable

#### 'frontend/.env'
VITE_API_URL=http://localhost:8080/api

## Running the app

Once your .env files are configured:

git clone https://github.com/yourusername/storage-webapp.git
cd storage-webapp
docker compose up --build

The frontend will be available at http://localhost:3000

The backend API will be running at http://localhost:8080/api

The PostgreSQL service will be running internally on port 5432


