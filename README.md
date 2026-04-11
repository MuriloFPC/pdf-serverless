# PDF Serverless

An open-source, serverless PDF processing platform, inspired by iLovePDF.

## 🚀 Project Overview

PDF Serverless is designed to provide organizations with a secure, private, and compliant way to process PDF files within their own AWS infrastructure. It is an API-first platform that uses an asynchronous queue and worker architecture, making it ideal for high-scale, serverless deployments.

### 🌐 Project URLs

- **Frontend**: [http://pdf-serverless-opensource-frontend.s3-website-us-east-1.amazonaws.com/](http://pdf-serverless-opensource-frontend.s3-website-us-east-1.amazonaws.com/)
- **API**: [https://sjyu4p7dsf.execute-api.us-east-1.amazonaws.com](https://sjyu4p7dsf.execute-api.us-east-1.amazonaws.com)
- **Swagger Documentation**: [https://sjyu4p7dsf.execute-api.us-east-1.amazonaws.com/swagger/index.html](https://sjyu4p7dsf.execute-api.us-east-1.amazonaws.com/swagger/index.html)

---

## 🎯 Core Focus

The project is built with three main pillars:
1. **Security**: Complete control over your data. Files are processed in your own infrastructure and automatically deleted after a configurable TTL.
2. **Governance**: Ensure compliance with data privacy regulations (GDPR, LGPD) by keeping the processing pipeline within your organization's boundaries.
3. **Easy Implementation**: Simple deployment using AWS SAM and a clean, API-first architecture that is easy to extend and integrate.

---

## ✨ Features (MVP)

The following operations are currently supported:

- **Merge PDFs**: Combine multiple PDF documents into a single file.
- **Split PDFs**: Break a PDF document into multiple single-page files.
- **Protect PDF**: Secure your PDF with a password.
- **Remove Password**: Decrypt and remove protection from a PDF (requires current password).

---

## 🏗️ Architecture

The system follows **Clean Architecture** principles and is written in **Go 1.26**.

- **API Layer**: Handles authentication, file uploads, job creation, and status monitoring.
- **Service Layer**: Orchestrates business logic and uses the **Strategy Pattern** for PDF operations.
- **Domain Layer**: Contains core entities and repository interfaces.
- **Infrastructure Layer**: Implements storage (S3), database (DynamoDB), and queue (SQS) providers.

Currently, the project supports memory-based implementations for local testing and AWS (S3, DynamoDB, SQS) for production.

---

## 🛠️ Getting Started

### Prerequisites

To run or deploy this project, you will need:

- **Go 1.26** or higher.
- **Make** (to run build commands and shortcuts).
- **AWS CLI** configured with your credentials.
- **AWS SAM CLI** installed for infrastructure deployment.

### Installation

1. Clone the repository.
2. Install dependencies:
   ```bash
   go mod download
   ```

### Running Locally

To run the API locally using memory-based storage and queue (no AWS required):

```bash
make local
```

The server will start at `http://localhost:3000`. You can access Swagger at `http://localhost:3000/swagger/index.html` to test the endpoints.

---

## 📦 Deployment to AWS

The project uses **AWS SAM (Serverless Application Model)** for infrastructure as code.

### Manual Deployment

If you need to perform a manual deployment, follow these steps:

1. **Configure your AWS credentials** via `aws configure`.
2. **Run the deploy command**:
   ```bash
   make deploy
   ```
   This command will run tests, build binaries for Linux/ARM64, and start `sam deploy --guided`.
3. Follow the terminal instructions to configure the stack name, region, and parameters (like `JwtSecret`).

### Frontend Deployment

The frontend is located in the `/frontend` folder. To deploy manually:
1. Go to the frontend folder: `cd frontend`
2. Install dependencies and build:
   ```bash
   npm install
   npm run build
   ```
3. Sync the build with the S3 bucket (replace with your bucket name):
   ```bash
   aws s3 sync dist/ s3://pdf-serverless-opensource-frontend --acl public-read
   ```

---

## 🔒 Security

- All files uploaded and processed in S3 have Lifecycle Rules for automatic deletion based on the selected TTL (e.g., 24h, 72h).
- Authentication is handled via JWT (JSON Web Tokens).
- User passwords are protected using the Argon2 hashing algorithm.

---

## 💻 Tech Stack

- **Backend**: Go 1.26
- **Web Framework**: Fiber
- **PDF Engine**: pdfcpu
- **Authentication**: JWT & Argon2
- **Infrastructure**: AWS S3, DynamoDB, SQS, Lambda, API Gateway
- **Frontend**:  Vite (in the `/frontend` folder)
