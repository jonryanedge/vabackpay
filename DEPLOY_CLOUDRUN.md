# Deploy to Google Cloud Run

This guide covers deploying the VA Backpay Calculator to GCP Cloud Run.

## Prerequisites

- Google Cloud SDK installed (`gcloud` command)
- Docker installed (for local testing)
- GCP project with billing enabled

## Quick Deploy

```bash
# 1. Enable required services
gcloud services enable run.googleapis.com containerregistry.googleapis.com

# 2. Configure Docker
gcloud auth configure-docker

# 3. Deploy
gcloud run deploy vabackpay \
  --source . \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars DEBUG=false
```

## Step-by-Step

### 1. Install Google Cloud SDK

Follow instructions at: https://cloud.google.com/sdk/docs/install

### 2. Authenticate

```bash
gcloud auth login
gcloud config set project YOUR_PROJECT_ID
```

### 3. Enable APIs

```bash
gcloud services enable run.googleapis.com
gcloud services enable containerregistry.googleapis.com
```

### 4. Configure Environment Variables

Update your `.env` file for production:

```bash
DEBUG=false
SMTP_HOST=smtp.yourprovider.com
SMTP_PORT=587
SMTP_USERNAME=your-username
SMTP_PASSWORD=your-password
FROM_EMAIL=your-sender@domain.com
FROM_NAME=VA Backpay Calculator
APP_URL=https://your-service.a.run.app
```

### 5. Deploy

```bash
gcloud run deploy vabackpay \
  --source . \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars DEBUG=false
```

- `--source .` builds from current directory
- `--region us-central1` deploys to Iowa (change as needed)
- `--allow-unauthenticated` makes it publicly accessible
- `--set-env-vars` passes environment variables

### 6. Verify

```bash
# Get the URL
gcloud run services describe vabackpay --region us-central1
```

## Custom Domain (Optional)

1. Map custom domain in Cloud Run:
```bash
gcloud run domain-mappings create \
  --service vabackpay \
  --domain yourdomain.com \
  --region us-central1
```

2. Update DNS records as instructed

## Persistent Storage

The current implementation writes to local filesystem. In Cloud Run, this data is ephemeral. For persistent storage:

### Option A: Cloud Storage

Create a bucket and update code to write there:

```bash
gsutil mb -l us-central1 gs://your-bucket-name
```

### Option B: Firestore

Store submissions in Firestore database.

## Docker Build (Alternative)

```bash
# Build image
docker build -t gcr.io/PROJECT_ID/vabackpay .

# Push to Container Registry
docker push gcr.io/PROJECT_ID/vabackpay

# Deploy
gcloud run deploy vabackpay \
  --image gcr.io/PROJECT_ID/vabackpay \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated
```

## Environment Variables in Cloud Run

Set via console or CLI:

```bash
gcloud run deploy vabackpay \
  --source . \
  --region us-central1 \
  --allow-unauthenticated \
  --update-env-vars DEBUG=false,SMTP_HOST=smtp.example.com
```

## Costs

- **Cloud Run**: Pay per use (usually free tier covers low traffic)
- **Cloud Storage** (if used): ~$0.020/GB/month
- **Custom domain SSL**: Free with Cloud Run

## Troubleshooting

### Container won't start

Check logs:
```bash
gcloud logs read --service vabackpay
```

### Email not working

Verify environment variables:
```bash
gcloud run services describe vabackpay --region us-central1 --format=get(spec.template.spec.containers[0].env)
```

### Permission issues

Ensure service account has necessary permissions:
```bash
gcloud run services add-iam-policy-binding vabackpay \
  --member=allUsers \
  --role=roles/run.invoker \
  --region us-central1
```
