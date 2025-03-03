#!/bin/bash
set -e

# Configuration
PROJECT_ID="single-bulwark-450909-q4"
SERVICE_ACCOUNT_NAME="crowdfund-deployer"
SERVICE_ACCOUNT_EMAIL="$SERVICE_ACCOUNT_NAME@$PROJECT_ID.iam.gserviceaccount.com"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}This script must be run with your GCP owner account, not the service account${NC}"
echo -e "${YELLOW}It will enable the necessary APIs and grant permissions to your service account${NC}"

# Check if gcloud is installed
if ! command -v gcloud &> /dev/null; then
    echo -e "${RED}gcloud CLI is not installed. Please install it first:${NC}"
    echo "https://cloud.google.com/sdk/docs/install"
    exit 1
fi

# Ensure user is logged in with owner account
echo -e "${YELLOW}Please make sure you're logged in with your owner account, not the service account${NC}"
echo -e "${YELLOW}Current account:${NC}"
gcloud auth list --filter=status:ACTIVE --format="value(account)"

read -p "Continue with this account? (y/n) " continue_setup
if [[ $continue_setup != "y" && $continue_setup != "Y" ]]; then
    echo -e "${YELLOW}Please login with your owner account:${NC}"
    echo "gcloud auth login"
    exit 1
fi

# Set the project
echo -e "${YELLOW}Setting project to $PROJECT_ID...${NC}"
gcloud config set project $PROJECT_ID

# Enable required APIs
echo -e "${YELLOW}Enabling required APIs...${NC}"
gcloud services enable artifactregistry.googleapis.com containerregistry.googleapis.com compute.googleapis.com

# Grant permissions to the service account
echo -e "${YELLOW}Granting permissions to service account $SERVICE_ACCOUNT_EMAIL...${NC}"

# Add Artifact Registry Admin role
echo -e "${YELLOW}Adding Artifact Registry Admin role...${NC}"
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:$SERVICE_ACCOUNT_EMAIL" \
  --role="roles/artifactregistry.admin"

# Add Storage Admin role
echo -e "${YELLOW}Adding Storage Admin role...${NC}"
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:$SERVICE_ACCOUNT_EMAIL" \
  --role="roles/storage.admin"

# Add Container Registry Service Agent role
echo -e "${YELLOW}Adding Container Registry Service Agent role...${NC}"
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:$SERVICE_ACCOUNT_EMAIL" \
  --role="roles/containerregistry.ServiceAgent"

# Add Compute Admin role for VM instance access
echo -e "${YELLOW}Adding Compute Admin role...${NC}"
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:$SERVICE_ACCOUNT_EMAIL" \
  --role="roles/compute.admin"

# Create Artifact Registry repository if it doesn't exist
echo -e "${YELLOW}Creating Artifact Registry repository if it doesn't exist...${NC}"
if ! gcloud artifacts repositories describe crowdfund --location=us-central1 &> /dev/null; then
    gcloud artifacts repositories create crowdfund \
        --repository-format=docker \
        --location=us-central1 \
        --description="Docker repository for Crowdfund application"
    echo -e "${GREEN}Repository created successfully${NC}"
else
    echo -e "${GREEN}Repository already exists${NC}"
fi

echo -e "${GREEN}Setup complete!${NC}"
echo -e "${YELLOW}Now you can run the fix-auth-issue.sh script with your service account${NC}"
echo -e "And then build and push your Docker image"
