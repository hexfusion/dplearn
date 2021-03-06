#!/usr/bin/env bash
set -e

if ! [[ "$0" =~ "./scripts/gcp/ubuntu-python3-gpu.gcp.sh" ]]; then
  echo "must be from repository root"
  exit 255
fi

if [[ "${GCP_KEY_PATH}" ]]; then
  echo GCP_KEY_PATH is defined: \""${GCP_KEY_PATH}"\"
else
  echo GCP_KEY_PATH is not defined!
  exit 255
fi

gcloud config set project dplearn

gcloud beta compute instances create dplearn-gpu \
  --custom-cpu=4 \
  --custom-memory=16 \
  --zone us-west1-b \
  --image-family=ubuntu-1604-lts \
  --image-project=ubuntu-os-cloud \
  --boot-disk-size=60 \
  --boot-disk-type="pd-ssd" \
  --network default \
  --tags=dplearn,http-server,https-server \
  --maintenance-policy=MIGRATE \
  --restart-on-failure \
  --accelerator type=nvidia-tesla-k80,count=1 \
  --metadata-from-file gcp-key-dplearn=${GCP_KEY_PATH},startup-script=./scripts/gcp/ubuntu-python3-gpu.ansible.sh
