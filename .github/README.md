# GitHub Actions Deployment Setup

This repository uses GitHub Actions to automatically deploy both the frontend and backend to Google Cloud Platform (GCP).

## Frontend Deployment

The frontend is deployed to GCP Storage buckets for static website hosting.

### Required GitHub Secrets for Frontend

To enable the frontend deployment workflow, you need to add the following secrets to your GitHub repository:

1. **GCP_SA_KEY**: A JSON service account key with permissions to create and manage GCP Storage buckets
2. **GCP_PROJECT_ID**: Your Google Cloud project ID
3. **GCP_BUCKET_NAME** (optional): A specific bucket name to use for deployment. If not provided, a name will be generated automatically.

## Backend Deployment

The backend is deployed to GCP Cloud Run as a containerized service.

### Required GitHub Secrets for Backend

To enable the backend deployment workflow, you need to add the following secrets to your GitHub repository:

1. **GCP_SA_KEY**: A JSON service account key with permissions to create and manage GCP Cloud Run services and Container Registry
2. **GCP_PROJECT_ID**: Your Google Cloud project ID
3. **POSTGRES_USER**: PostgreSQL database username
4. **POSTGRES_PASSWORD**: PostgreSQL database password
5. **POSTGRES_HOST**: PostgreSQL database host
6. **POSTGRES_PORT**: PostgreSQL database port
7. **POSTGRES_DB**: PostgreSQL database name
8. **JWT_SECRET**: Secret key for JWT token generation and validation

## How to Set Up GCP Service Account

1. Go to the [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the necessary APIs:
   - Cloud Storage API
   - Cloud Run API
   - Container Registry API
4. Go to IAM & Admin > Service Accounts
5. Create a new service account with the following roles:
   - Storage Admin
   - Storage Object Admin
   - Cloud Run Admin
   - Service Account User
6. Create a new key for this service account (JSON format)
7. Download the key file

## How to Add Secrets to GitHub

1. Go to your GitHub repository
2. Click on "Settings" > "Secrets and variables" > "Actions"
3. Click "New repository secret"
4. Add each of the required secrets:
   - Name: `GCP_SA_KEY`, Value: *paste the entire contents of the JSON key file*
   - Name: `GCP_PROJECT_ID`, Value: *your GCP project ID*
   - Name: `GCP_BUCKET_NAME`, Value: *your desired bucket name* (optional)
   - Name: `POSTGRES_USER`, Value: *your PostgreSQL database username*
   - Name: `POSTGRES_PASSWORD`, Value: *your PostgreSQL database password*
   - Name: `POSTGRES_HOST`, Value: *your PostgreSQL database host*
   - Name: `POSTGRES_PORT`, Value: *your PostgreSQL database port*
   - Name: `POSTGRES_DB`, Value: *your PostgreSQL database name*
   - Name: `JWT_SECRET`, Value: *your secret key for JWT token generation and validation*

## Workflow Behavior

### Frontend Workflow
The frontend deployment workflow will:
- Trigger on pushes to the main branch that affect files in the frontend directory
- Build the Vue.js application
- Create a GCP Storage bucket if it doesn't exist
- Upload the built files to the bucket
- Configure the bucket for web hosting
- Make the bucket contents publicly accessible
- Print the deployment URL

### Backend Workflow
The backend deployment workflow will:
- Trigger on pushes to the main branch that affect files in the backend directory
- Build a Docker image for the Go application
- Push the image to Google Container Registry
- Deploy the image to Cloud Run
- Configure environment variables for the service
- Print the deployment URL

### Test Workflow
The test workflow will:
- Trigger on all pull requests to the main branch
- Run Go tests for the backend
- Run type checking and linting for the frontend
- This helps ensure code quality before merging changes

## Manual Deployment

You can also manually trigger both workflows from the "Actions" tab in your GitHub repository.
