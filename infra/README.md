# Reel Farm AWS Infrastructure

This Terraform configuration sets up AWS infrastructure for the Reel Farm application with the following components:

## Architecture

- **VPC**: Custom VPC with public subnet
- **Internet Gateway**: For public internet access
- **Application Load Balancer (ALB)**: Routes traffic to ECS tasks
- **ECS Fargate**: Serverless container platform for running your application
- **Security Groups**: Separate security groups for ALB and Fargate tasks
- **S3 Bucket**: For storing application data
- **ECR Repository**: For storing Docker container images

## Components

### Networking

- VPC with CIDR `10.0.0.0/16`
- Public subnet with CIDR `10.0.1.0/24`
- Internet Gateway for public access
- Route table with default route to Internet Gateway

### Security

- ALB Security Group: Allows HTTP (80) and HTTPS (443) from anywhere
- Fargate Security Group: Allows traffic from ALB only
- S3 bucket with encryption and public access blocked
- IAM roles for ECS execution and task permissions

### Application Load Balancer

- Listens on port 80 (HTTP)
- Health checks on configured path
- Routes traffic to ECS Fargate tasks

### ECS Fargate

- Serverless container platform
- Task definition with configurable CPU/memory
- Service with desired count of 1
- CloudWatch logging enabled

### Storage

- S3 bucket with versioning and encryption
- ECR repository for Docker images
- IAM policies for S3 access from ECS tasks

## Prerequisites

1. AWS CLI configured with appropriate credentials
2. Terraform installed (>= 1.0)
3. Docker installed (for building/pushing images)

## Usage

1. **Initialize Terraform**:

   ```bash
   terraform init
   ```

2. **Configure variables**:

   ```bash
   cp terraform.tfvars.example terraform.tfvars
   # Edit terraform.tfvars with your specific values
   ```

3. **Plan the deployment**:

   ```bash
   terraform plan
   ```

4. **Apply the configuration**:

   ```bash
   terraform apply
   ```

5. **Build and push your Docker image**:

   ```bash
   # Get the ECR repository URL from terraform output
   ECR_URL=$(terraform output -raw ecr_repository_url)

   # Login to ECR
   aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin $ECR_URL

   # Build your image
   docker build -t reel-farm-app .

   # Tag and push
   docker tag reel-farm-app:latest $ECR_URL:latest
   docker push $ECR_URL:latest
   ```

6. **Update ECS service** (if needed):
   ```bash
   aws ecs update-service --cluster reel-farm-cluster --service reel-farm-service --force-new-deployment
   ```

## Outputs

After deployment, Terraform will output:

- `alb_dns_name`: The DNS name of your load balancer
- `s3_bucket_name`: Name of your S3 bucket
- `ecr_repository_url`: URL of your ECR repository
- `ecs_cluster_name`: Name of your ECS cluster

## Accessing Your Application

Once deployed, you can access your application via the ALB DNS name:

```bash
ALB_DNS=$(terraform output -raw alb_dns_name)
curl http://$ALB_DNS
```

## Security Notes

- The Fargate tasks are assigned public IPs and can access the internet via the Internet Gateway
- S3 bucket has public access blocked for security
- Security groups follow least-privilege principle
- ECR repository has image scanning enabled

## Cost Optimization

- Uses Fargate Spot for cost savings (can be configured)
- Configurable CPU/memory allocation
- CloudWatch logs retention set to 30 days
- Single AZ deployment (can be expanded for HA)

## Cleanup

To destroy all resources:

```bash
terraform destroy
```

**Note**: This will delete all resources including S3 bucket contents and ECR images.
