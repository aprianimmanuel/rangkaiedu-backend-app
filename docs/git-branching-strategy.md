# Git Branching Strategy for Rangkai Edu

## Overview

This document outlines the Git branching strategy for the Rangkai Edu project. The strategy is designed to support continuous integration and deployment while maintaining code quality and stability.

## Branch Structure

### Main Branches

1. **main** - Production-ready code
   - Contains code that is currently running in production
   - Only accepts merges from `staging` branch via pull requests
   - Protected with strict rules (see below)

2. **staging** - Pre-production testing
   - Contains code that is ready for production but undergoing final testing
   - Only accepts merges from `develop` branch via pull requests
   - Protected with moderate rules (see below)

3. **develop** - Ongoing development
   - Integration branch for features
   - Accepts merges from feature branches
   - Less restrictive than main and staging

### Supporting Branches

1. **Feature Branches** - `feature/*`
   - Branch from: `develop`
   - Merge back to: `develop`
   - Naming convention: `feature/JIRA-123-short-description`

2. **Release Branches** - `release/*`
   - Branch from: `develop`
   - Merge back to: `staging` and `develop`
   - Naming convention: `release/v1.2.0`

3. **Hotfix Branches** - `hotfix/*`
   - Branch from: `main`
   - Merge back to: `main`, `staging`, and `develop`
   - Naming convention: `hotfix/JIRA-456-short-description`

## Branch Protection Rules

### Main Branch Protection

The `main` branch has the following protection rules:

1. **Require pull request reviews before merging**
   - At least 2 approved reviews required
   - Dismiss stale pull request approvals when new commits are pushed
   - Require review from Code Owners

2. **Require status checks to pass before merging**
   - Continuous Integration (CI) checks must pass
   - Required status checks:
     - Backend unit tests
     - Frontend unit tests
     - Integration tests
     - Security scans
     - Code quality checks

3. **Require branches to be up to date before merging**
   - Ensure branch is up-to-date with base branch before merging

4. **Include administrators**
   - Apply these rules to repository administrators as well

5. **Restrict who can push to matching branches**
   - Only specific teams can push directly (restricted)

### Staging Branch Protection

The `staging` branch has the following protection rules:

1. **Require pull request reviews before merging**
   - At least 1 approved review required
   - Dismiss stale pull request approvals when new commits are pushed

2. **Require status checks to pass before merging**
   - Continuous Integration (CI) checks must pass
   - Required status checks:
     - Backend unit tests
     - Frontend unit tests
     - Integration tests

3. **Require branches to be up to date before merging**
   - Ensure branch is up-to-date with base branch before merging

4. **Include administrators**
   - Apply these rules to repository administrators as well

## Workflows

### Feature Development Workflow

1. Create a feature branch from `develop`
   ```bash
   git checkout develop
   git pull origin develop
   git checkout -b feature/JIRA-123-short-description
   ```

2. Develop the feature and commit changes
   ```bash
   git add .
   git commit -m "feat: implement user authentication"
   ```

3. Push the feature branch to remote
   ```bash
   git push origin feature/JIRA-123-short-description
   ```

4. Create a Pull Request from feature branch to `develop`
5. After review and approval, merge the Pull Request
6. Delete the feature branch

### Release Workflow

1. Create a release branch from `develop`
   ```bash
   git checkout develop
   git pull origin develop
   git checkout -b release/v1.2.0
   ```

2. Perform release preparations (version bump, final testing, etc.)
3. Push the release branch to remote
   ```bash
   git push origin release/v1.2.0
   ```

4. Create a Pull Request from release branch to `staging`
5. After testing on staging environment, create a Pull Request from release branch to `main`
6. After review and approval, merge both Pull Requests
7. Create a release tag on `main`
   ```bash
   git checkout main
   git pull origin main
   git tag -a v1.2.0 -m "Release version 1.2.0"
   git push origin v1.2.0
   ```

8. Merge release branch back to `develop` to incorporate any changes
9. Delete the release branch

### Hotfix Workflow

1. Create a hotfix branch from `main`
   ```bash
   git checkout main
   git pull origin main
   git checkout -b hotfix/JIRA-456-critical-bug-fix
   ```

2. Implement the hotfix
3. Push the hotfix branch to remote
   ```bash
   git push origin hotfix/JIRA-456-critical-bug-fix
   ```

4. Create a Pull Request from hotfix branch to `main`
5. After review and approval, merge the Pull Request
6. Create a release tag on `main`
   ```bash
   git checkout main
   git pull origin main
   git tag -a v1.1.1 -m "Hotfix version 1.1.1"
   git push origin v1.1.1
   ```

7. Create Pull Requests to merge the hotfix into `staging` and `develop`
8. After review and approval, merge both Pull Requests
9. Delete the hotfix branch

## Merge Strategies

### Merge Commit

- Used for: Release and hotfix branches
- Preserves complete history and chronological order
- Creates a merge commit

### Squash and Merge

- Used for: Feature branches
- Creates a single commit with a summary of all changes
- Keeps history clean and linear
- Recommended for most feature branches

### Rebase and Merge

- Used for: When maintaining a linear history is important
- Reapplies commits on top of the target branch
- Creates a linear history without merge commits

## Code Review Process

### Review Requirements

1. All pull requests must be reviewed before merging
2. Feature branches to `develop`: At least 1 reviewer
3. Release branches to `staging`: At least 1 reviewer
4. Release branches to `main`: At least 2 reviewers
5. Hotfix branches to `main`: At least 2 reviewers

### Review Guidelines

1. **Code Quality**
   - Follows project coding standards
   - Proper error handling
   - Adequate logging
   - Efficient algorithms and data structures

2. **Security**
   - No hardcoded credentials
   - Proper input validation
   - Secure handling of sensitive data

3. **Testing**
   - Adequate test coverage
   - Tests are meaningful and pass
   - Edge cases are considered

4. **Documentation**
   - Code is properly documented
   - README files are updated if needed
   - API documentation is updated if needed

### Review Process

1. Reviewer is automatically assigned based on code ownership
2. Reviewer has 24 hours to review during business days
3. Author should address all comments before merging
4. Reviewer should approve the pull request after all comments are resolved

## Contribution Guidelines

### Commit Message Convention

We follow the Conventional Commits specification:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

#### Commit Types

- **feat**: A new feature
- **fix**: A bug fix
- **docs**: Documentation only changes
- **style**: Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
- **refactor**: A code change that neither fixes a bug nor adds a feature
- **perf**: A code change that improves performance
- **test**: Adding missing tests or correcting existing tests
- **build**: Changes that affect the build system or external dependencies
- **ci**: Changes to our CI configuration files and scripts
- **chore**: Other changes that don't modify src or test files
- **revert**: Reverts a previous commit

#### Examples

```
feat: add user authentication module
fix: resolve database connection issue
docs: update API documentation
refactor: improve error handling in user service
```

### Branch Naming Convention

- Feature branches: `feature/JIRA-123-short-description`
- Release branches: `release/v1.2.0`
- Hotfix branches: `hotfix/JIRA-456-short-description`

## Diagram

```
main       o --------------------------------------o------------------------o
           |                                       |                        |
           |                                       |                        |
staging    | o------------------o------------------|------------------------|---o
           | |                  |                  |                        |   |
           | |                  |                  |                        |   |
develop    | | o---o---o---o----|------------------|------------------------|---|---o
           | | |   |   |   |    |                  |                        |   |   |
           | | |   |   |   |    |                  |                        |   |   |
feature    | | o---o   |   o----o                  |                        |   |   o
           | |         |                          |                        |   |
           | |         o--------------------------o                        |   |
           | |                                                            |   |
release    | o------------------------------------------------------------o   |
           |                                                                  |
           |                                                                  |
hotfix     o------------------------------------------------------------------o
```

This strategy ensures a clean, organized workflow that supports continuous integration and deployment while maintaining code quality and stability.