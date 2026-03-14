Project Overview

This project is an open-source alternative to iLovePDF.

The goal is to create a serverless PDF processing platform that organizations can deploy in their own infrastructure to meet security, compliance, and privacy requirements.

The system must be API-first, meaning all functionality must work without a UI.
A web interface may be added later, but it is not required for the MVP.

The platform processes PDF jobs asynchronously using a queue and worker architecture.

Core Principles

The project must remain fully serverless-friendly.

The system must follow Clean Architecture principles.

The API must always be usable without a UI.

All business rules must live in the domain layer.

Infrastructure implementations must be replaceable.

The project must be easy to extend with new PDF operations.

Initial Features (MVP)

The first version must support the following PDF operations:

Merge PDFs

Split PDFs

Protect PDF (add password)

Remove password

New operations must be implemented using the processing strategy pattern.

System Architecture

The system consists of the following components:

Client
↓
API
↓
SQS Queue
↓
Worker
↓
S3 Storage

Responsibilities:

API

authentication

file upload

job creation

job status queries

Worker

consume queue messages

process PDF jobs

store results

update job status

Technology Stack

Backend language:

Go

HTTP framework:

Fiber

Infrastructure services:

Queue
Amazon SQS

Database
Amazon DynamoDB

File storage
Amazon S3

Deployment:

Serverless using

AWS Serverless Application Model
or

Terraform

Project Structure

The project follows Clean Architecture.

cmd/
api/
main.go
worker/
main.go

internal/

core/
domain/
pdf/
entity.go
repository.go

    service/
      pdf_service/
        service.go
        strategy/

router/
pdf_handler.go

infra/
database/
storage/
queue/

Layer responsibilities:

Domain

entities

repository interfaces

business rules

Service

use cases

orchestration

Infrastructure

AWS integrations

database implementations

queue implementations

storage implementations

Router / handlers

HTTP endpoints

API Endpoints

Authentication:

POST /auth/login
POST /auth/register

PDF processing:

POST /pdf/process
Creates a new processing job.

GET /pdf/status/{id}
Returns the job status.

GET /pdf/list
Returns all jobs for the authenticated user.

Job Processing Flow

Client uploads files.

API stores files in S3.

API creates a job record in DynamoDB.

API sends a message to SQS.

Worker consumes the message.

Worker processes the PDF operation.

Worker stores output in S3.

Worker updates job status in DynamoDB.

DynamoDB Model

The system stores job metadata only.

Files are stored in S3.

Job entity example:

JobID
UserID
ProcessType
Status
CreatedAt
DeleteAt
InputFiles
OutputFiles

Status values:

pending
processing
completed
failed
Security Requirements

All endpoints must require authentication except:

/auth/login
/auth/register

Requirements:

JWT authentication

Password hashing using Argon2

Users must only access their own jobs

Testing Requirements

All new features must include automated tests.

Tests should cover:

API endpoints

authentication logic

job creation

job status retrieval

worker processing

Coding Rules

Follow Go best practices.

Do not introduce unnecessary dependencies.

Keep domain code independent from infrastructure.

Avoid tight coupling between layers.

Keep services small and focused.