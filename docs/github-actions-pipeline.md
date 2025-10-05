# GitHub Actions CI/CD Pipeline Implementation

## Overview
This document details the implementation of the GitHub Actions CI/CD pipeline for the Rangkai Edu project. The pipeline automates testing, building, and deployment of both frontend and backend applications.

## Pipeline Structure

### Workflow Files
1. **`.github/workflows/ci.yml`** - Continuous Integration pipeline
2. **`.github/workflows/cd.yml`** - Continuous Deployment pipeline
3. **`.github/workflows/security.yml`** - Scheduled security scanning

## Continuous Integration (CI) Pipeline

### Trigger Events
- Push to `develop`, `staging`, or `main` branches
- Pull requests to `develop`, `staging`, or `main` branches

### Backend CI Jobs
1. **Code Checkout**: Retrieves the latest source code
2. **Setup Go Environment**: Configures Go 1.24 with dependency caching
3. **Install Dependencies**: Downloads Go modules
4. **Code Quality Checks**: Runs `golangci-lint` for linting
5. **Security Scanning**: Runs `gosec` for Go security vulnerabilities
6. **Unit Testing**: Executes Go tests with coverage reporting
7. **Build Backend**: Compiles the Go application
8. **Build Docker Image**: Creates backend Docker image
9. **Scan Docker Image**: Uses Trivy to scan for vulnerabilities

### Frontend CI Jobs
1. **Code Checkout**: Retrieves the latest source code
2. **Setup Node.js Environment**: Configures Node.js 20 with npm caching
3. **Install Dependencies**: Runs `npm ci` for clean dependency installation
4. **Code Quality Checks**: Runs ESLint for code linting
5. **Security Scanning**: Runs `npm audit` for dependency vulnerabilities
6. **Build Frontend**: Creates production build with Vite
7. **Build Docker Image**: Creates frontend Docker image
8. **Scan Docker Image**: Uses Trivy to scan for vulnerabilities

### Docker Image Management
- Images are built and tagged with the commit SHA
- Images are pushed to GitHub Container Registry (GHCR)
- Tags are created based on branch names:
  - `develop` branch: `develop` tag
  - `staging` branch: `staging` tag
  - `main` branch: `latest` tag

## Continuous Deployment (CD) Pipeline

### Trigger Events
- Successful completion of CI workflow

### Staging Deployment
- Triggered when CI workflow completes successfully on `develop` branch
- Deploys to staging environment at https://staging.rangkaiedu.com
- Environment variables are configured for staging

### Production Deployment
- Triggered when CI workflow completes successfully on `main` branch
- Deploys to production environment at https://rangkaiedu.com
- Environment variables are configured for production

## Security Scanning Pipeline

### Trigger Events
- Scheduled daily at 2 AM UTC
- Manual trigger via workflow dispatch

### Security Checks
1. **Go Security Scan**: Runs `gosec` on backend code
2. **NPM Audit**: Checks frontend dependencies for vulnerabilities
3. **Docker Image Scanning**: Scans container images for vulnerabilities

## Environment Configuration

### GitHub Secrets Required
- `GITHUB_TOKEN` - Automatically provided by GitHub Actions
- `CODECOV_TOKEN` - For code coverage reporting (optional)

### Environment Variables
- `BACKEND_IMAGE_NAME` - Backend Docker image name
- `FRONTEND_IMAGE_NAME` - Frontend Docker image name

### Environments
1. **Staging** - Pre-production testing environment
2. **Production** - Live production environment

## Implementation Details

### Docker Image Building
- Uses Docker Buildx for improved caching
- Multi-stage builds for smaller image sizes
- Non-root user for security
- Health checks included in Docker images

### Caching Strategy
- Go module caching
- npm dependency caching
- Docker layer caching
- GitHub Actions cache for faster builds

### Security Features
- Trivy vulnerability scanning for Docker images
- gosec for Go code security analysis
- npm audit for frontend dependencies
- Non-root user in container images

### Quality Assurance
- Code coverage reporting with Codecov
- Linting for both Go and JavaScript/JSX
- Parallel execution of CI jobs
- Comprehensive test suite execution

## Testing Workflow

### Backend Testing
- Unit tests with Go testing framework
- Race condition detection
- Code coverage reporting
- Integration with Codecov

### Frontend Testing
- ESLint for code quality
- npm audit for security
- Production build verification

## Deployment Process

### Staging Deployment Steps
1. Checkout code
2. Login to GitHub Container Registry
3. Deploy Docker images to staging infrastructure
4. Perform health checks
5. Notify team of deployment status

### Production Deployment Steps
1. Checkout code
2. Login to GitHub Container Registry
3. Deploy Docker images to production infrastructure
4. Perform health checks
5. Notify team and stakeholders of deployment status

## Monitoring and Notifications

### Pipeline Monitoring
- Real-time status through GitHub Actions UI
- Build time metrics
- Success/failure tracking

### Deployment Notifications
- Slack notifications (planned)
- Email notifications (planned)
- Deployment status reporting

## Rollback Procedures

### Automated Rollback
- Previous Docker images can be redeployed
- Database migrations should be reversible
- Health checks monitor rollback success

### Manual Rollback
- Deploy specific image tags to revert changes
- Coordinate with team for database state
- Update documentation with rollback steps

## Best Practices Implemented

### Security
- No hardcoded credentials
- Regular security scanning
- Non-root containers
- Dependency vulnerability scanning

### Performance
- Parallel job execution
- Caching strategies
- Efficient Docker layering
- BuildKit optimizations

### Reliability
- Comprehensive error handling
- Health checks for services
- Environment-specific configurations
- Proper logging and monitoring hooks

## Future Enhancements

### Planned Improvements
1. Integration with monitoring systems (Prometheus, Grafana)
2. Automated performance testing
3. Enhanced notification system
4. Blue-green deployment strategy
5. Database migration automation
6. Integration testing in CI pipeline

### Scalability Considerations
1. Matrix testing for multiple environments
2. Parallel deployment strategies
3. Infrastructure as Code integration
4. Multi-region deployment support

## Troubleshooting

### Common Issues
1. **Docker build failures**: Check Dockerfile syntax and dependencies
2. **Test failures**: Review test output and fix code issues
3. **Security scan failures**: Address reported vulnerabilities
4. **Deployment failures**: Check environment configurations

### Debugging Steps
1. Review GitHub Actions logs
2. Check Docker image scanning results
3. Verify environment variables
4. Confirm GitHub Secrets are properly configured

## Maintenance

### Regular Tasks
1. Update workflow dependencies
2. Review security scan results
3. Monitor build performance
4. Update documentation

### Updates and Upgrades
1. Go version updates
2. Node.js version updates
3. Docker base image updates
4. GitHub Actions version updates