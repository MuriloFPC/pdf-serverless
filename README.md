# PDF Serverless

A self-hosted, serverless PDF processing platform, inspired by iLovePDF.

## Project Overview

PDF Serverless is designed to provide organizations with a secure, private, and compliant way to process PDF files within their own infrastructure. It is an API-first platform that uses an asynchronous queue and worker architecture, making it ideal for high-scale, serverless deployments on AWS.

## Features (MVP)

The following operations are currently supported:

- **Merge PDFs**: Combine multiple PDF documents into a single file.
- **Split PDFs**: Break a PDF document into multiple single-page files.
- **Protect PDF**: Secure your PDF with a password.
- **Remove Password**: Decrypt and remove protection from a PDF (requires password).

## Architecture

The system follows **Clean Architecture** principles and is written in **Go 1.26**.

- **API Layer**: Handles authentication, file uploads, job creation, and status monitoring.
- **Service Layer**: Orchestrates the business logic and uses the **Strategy Pattern** for PDF operations.
- **Domain Layer**: Contains the core entities and repository interfaces.
- **Infrastructure Layer**: Implements storage (S3), database (DynamoDB), and queue (SQS) providers.

Currently, the project includes memory-based implementations for all infrastructure components, allowing for easy local testing and demonstration.

## Getting Started

### Prerequisites

- Go 1.26 or later

### Installation

1. Clone the repository.
2. Install dependencies:
   ```bash
   go mod download
   ```

### Running the API (Local Demo)

The API includes an integrated background worker for local testing with memory-based storage and queue.

```bash
go run cmd/api/main.go
```

The server will start on `http://localhost:3000`.

## API Endpoints

### Authentication

- `POST /auth/register`: Create a new user account.
- `POST /auth/login`: Authenticate and receive a JWT token.

### PDF Processing (Requires JWT)

- `POST /pdf/process`: Upload files and start a processing job.
- `GET /pdf/status/:id`: Check the status of a specific job.
- `GET /pdf/list`: List all jobs for the authenticated user.

## Tech Stack

- **Backend**: Go 1.26
- **Web Framework**: Fiber
- **PDF Engine**: pdfcpu
- **Authentication**: JWT & Argon2
- **Infrastructure (Planned)**: AWS S3, DynamoDB, SQS
