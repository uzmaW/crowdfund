#!/bin/bash
set -e  # Exit on error

# Configuration
INSTANCE_NAME="github-actions-vm"
ZONE="us-central1-a"
MACHINE_TYPE="e2-micro"
IMAGE_FAMILY="debian-11"
IMAGE_PROJECT="debian-cloud"
DISK_SIZE="10GB"
TAGS="http-server"

# Deploy VM
gcloud compute instances create $INSTANCE_NAME \
  --zone=$ZONE \
  --machine-type=$MACHINE_TYPE \
  --image-family=$IMAGE_FAMILY \
  --image-project=$IMAGE_PROJECT \
  --boot-disk-size=$DISK_SIZE \
  --tags=$TAGS \
  --quiet  # Disable interactive prompts

# Open firewall (if needed)
if [[ $TAGS == *"http-server"* ]]; then
  gcloud compute firewall-rules create allow-http \
    --allow=tcp:80 \
    --source-ranges=0.0.0.0/0 \
    --target-tags=http-server \
    --quiet
fi