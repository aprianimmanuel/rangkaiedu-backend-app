# Cloud Storage Setup for Rangkai Edu

This document describes how to set up and configure cloud storage for the Rangkai Edu application.

## Overview

The Rangkai Edu application supports both local file storage and cloud storage (Alibaba Cloud OSS). This allows for flexible deployment options depending on your infrastructure requirements.

## Prerequisites

1. Alibaba Cloud account with billing enabled
2. Terraform v1.0 or higher installed
3. Docker and Docker Compose installed
4. Alibaba Cloud CLI installed and configured

## Infrastructure Setup

### 1. Terraform Configuration

The infrastructure is defined in the `terraform/` directory with the following structure:

```
terraform/
├── main.tf              # Provider configuration
├── variables.tf         # Input variables
├── outputs.tf           # Output values
├── terraform.tfvars.example  # Example variables file
├── README.md            # Terraform documentation
├── materials_storage.tf # Storage resources
└── modules/
    └── oss/
        ├── main.tf      # OSS bucket configuration
        ├── variables.tf # Module variables
        └── outputs.tf   # Module outputs
```

### 2. Deploying Infrastructure

1. **Authenticate with Alibaba Cloud:**
   Configure your Alibaba Cloud credentials using one of the following methods:
   
   a. Environment variables:
   ```bash
   export ALICLOUD_ACCESS_KEY="your-access-key-id"
   export ALICLOUD_SECRET_KEY="your-access-key-secret"
   export ALICLOUD_REGION="ap-southeast-1"
   ```
   
   b. Alibaba Cloud CLI configuration:
   ```bash
   aliyun configure
   ```

2. **Configure variables:**
   ```bash
   cp terraform/terraform.tfvars.example terraform/terraform.tfvars
   ```
   
   Edit `terraform/terraform.tfvars` with your project details:
   ```hcl
   project_name = "rangkaiedu"
   region       = "ap-southeast-1"  # Singapore region
   environment  = "dev"
   ```

3. **Initialize Terraform:**
   ```bash
   cd terraform
   terraform init
   ```

4. **Plan the infrastructure:**
   ```bash
   terraform plan
   ```

5. **Apply the infrastructure:**
   ```bash
   terraform apply
   ```

### 3. RAM User Setup

The Terraform configuration creates a RAM user with appropriate permissions to access the storage bucket. The access key ID and secret for this user are output as sensitive values.

To extract the RAM user credentials:

1. After running `terraform apply`, get the access key ID:
   ```bash
   terraform output storage_access_key_id
   ```

2. Get the access key secret:
   ```bash
   terraform output storage_access_key_secret
   ```

## Application Configuration

### Environment Variables

The application uses the following environment variables for storage configuration:

| Variable | Description | Default |
|----------|-------------|---------|
| `STORAGE_PROVIDER` | Storage provider to use (`local` or `oss`) | `local` |
| `OSS_BUCKET_NAME` | Name of the OSS bucket | - |
| `OSS_ACCESS_KEY_ID` | Access key ID for OSS access | - |
| `OSS_ACCESS_KEY_SECRET` | Access key secret for OSS access | - |
| `OSS_REGION` | Alibaba Cloud region | - |
| `OSS_ENDPOINT` | OSS endpoint URL | - |

### Docker Configuration

The `docker-compose.yml` file includes the necessary configuration for cloud storage:

```yaml
environment:
  STORAGE_PROVIDER: ${STORAGE_PROVIDER:-local}
  OSS_BUCKET_NAME: ${OSS_BUCKET_NAME}
  OSS_ACCESS_KEY_ID: ${OSS_ACCESS_KEY_ID}
  OSS_ACCESS_KEY_SECRET: ${OSS_ACCESS_KEY_SECRET}
  OSS_REGION: ${OSS_REGION}
  OSS_ENDPOINT: ${OSS_ENDPOINT}

volumes:
  - ./config:/app/config
```

## Deployment Process

### 1. Local Development with Local Storage

For local development with local file storage:

1. Copy the example environment file:
   ```bash
   cp config/.env.example .env
   ```

2. Start the application:
   ```bash
   docker-compose up -d
   ```

### 2. Deployment with Alibaba Cloud OSS

For deployment with Alibaba Cloud OSS:

1. Set up the infrastructure as described above

2. Copy the example environment file:
   ```bash
   cp config/.env.example .env
   ```

3. Update the environment variables in `.env`:
   ```bash
   STORAGE_PROVIDER=oss
   OSS_BUCKET_NAME=your-bucket-name
   OSS_REGION=ap-southeast-1
   OSS_ENDPOINT=https://oss-ap-southeast-1.aliyuncs.com
   ```

4. Set the access key credentials as environment variables or in the `.env` file:
   ```bash
   OSS_ACCESS_KEY_ID=your-access-key-id
   OSS_ACCESS_KEY_SECRET=your-access-key-secret
   ```

5. Start the application:
   ```bash
   docker-compose up -d
   ```

## Security Considerations

1. **RAM User Permissions**: The RAM user has minimal required permissions (GetObject, PutObject, DeleteObject) for the specific bucket.

2. **Key Management**: Access keys should be stored securely and rotated regularly.

3. **Environment Variables**: Sensitive information should not be hardcoded in configuration files.

4. **File Access**: Files are only accessible to authorized users based on application permissions.

## Troubleshooting

### Common Issues

1. **Authentication Errors**: Ensure the access key ID and secret are correctly set and have the correct permissions.

2. **Bucket Access Errors**: Verify that the RAM user has the necessary permissions for the bucket.

3. **File Upload Errors**: Check that the bucket exists and the RAM user has write permissions.

### Logs and Monitoring

Check the application logs for storage-related errors:
```bash
docker-compose logs backend
```

## Cleanup

To destroy the infrastructure:
```bash
cd terraform
terraform destroy
```

**Warning**: This will permanently delete all stored materials. Ensure you have backups if needed.