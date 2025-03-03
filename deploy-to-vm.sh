#!/bin/bash

# Use GitHub Secrets for configuration
PROJECT_ID=$PROJECT_ID               # Passed via GitHub Actions
INSTANCE_NAME="github-actions-vm"    # Name your VM
ZONE="us-central1-a"
MACHINE_TYPE="e2-micro"
IMAGE_FAMILY="debian-11"
IMAGE_PROJECT="debian-cloud"
DISK_SIZE="10GB"
TAGS="http-server"

# Disable interactive prompts
gcloud config set disable_prompts true

# Deploy the VM
gcloud compute instances create $INSTANCE_NAME \
  --project=$PROJECT_ID \
  --zone=$ZONE \
  --machine-type=$MACHINE_TYPE \
  --image-family=$IMAGE_FAMILY \
  --image-project=$IMAGE_PROJECT \
  --boot-disk-size=$DISK_SIZE \
  --tags=$TAGS

# Open firewall (if needed)
if [[ $TAGS == *"http-server"* ]]; then
  gcloud compute firewall-rules create allow-http \
    --allow=tcp:80 \
    --source-ranges=0.0.0.0/0 \
    --target-tags=http-server \
    --project=$PROJECT_ID \
    --quiet
fi