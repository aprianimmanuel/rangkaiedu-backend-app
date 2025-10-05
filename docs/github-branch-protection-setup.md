# GitHub Branch Protection Rules Setup

This document provides step-by-step instructions for setting up branch protection rules for the Rangkai Edu project in the GitHub UI.

## Prerequisites

- Repository administrator access
- GitHub account with appropriate permissions

## Setting up Branch Protection for Main Branch

1. Navigate to the repository on GitHub
2. Go to **Settings** tab
3. In the left sidebar, click **Branches**
4. Under **Branch protection rules**, click **Add rule**
5. In the **Branch name pattern** field, enter `main`
6. Configure the following settings:

### Protection Settings for Main Branch

- [x] **Require pull request reviews before merging**
  - Set **Required approving reviews** to `2`
  - [x] **Dismiss stale pull request approvals when new commits are pushed**
  - [x] **Require review from Code Owners**

- [x] **Require status checks to pass before merging**
  - [x] **Require branches to be up to date before merging**
  - In the **Status checks found in the last week for this repository** section, select:
    - `build (ubuntu-latest)` (from CI workflow)

- [x] **Include administrators**
  - This ensures that branch protection rules apply to repository administrators as well

- [x] **Restrict who can push to matching branches**
  - Add teams or users who should have direct push access (typically none for main branch)

7. Click **Create** to save the branch protection rule

## Setting up Branch Protection for Staging Branch

1. While still in the **Branches** settings page, click **Add rule**
2. In the **Branch name pattern** field, enter `staging`
3. Configure the following settings:

### Protection Settings for Staging Branch

- [x] **Require pull request reviews before merging**
  - Set **Required approving reviews** to `1`
  - [x] **Dismiss stale pull request approvals when new commits are pushed**

- [x] **Require status checks to pass before merging**
  - [x] **Require branches to be up to date before merging**
  - In the **Status checks found in the last week for this repository** section, select:
    - `build (ubuntu-latest)` (from CI workflow)

- [x] **Include administrators**
  - This ensures that branch protection rules apply to repository administrators as well

4. Click **Create** to save the branch protection rule

## Verification

After setting up the branch protection rules, verify they are working correctly:

1. Try to push directly to `main` or `staging` (should be rejected)
2. Try to merge a pull request without required approvals (should be blocked)
3. Try to merge a pull request with failing status checks (should be blocked)

## Additional Recommendations

1. **Enable required linear history** to prevent merge commits if your team prefers a linear history
2. **Enable allow force pushes** only for specific administrators if needed for emergency situations
3. **Enable allow deletions** should generally be left unchecked to prevent accidental branch deletion
4. **Lock branch** can be used temporarily to prevent changes during critical periods

## Troubleshooting

If you encounter issues with branch protection rules:

1. Check that the CI workflow is running correctly and reporting status
2. Verify that Code Owners are correctly configured in the `.github/CODEOWNERS` file
3. Ensure that users have the appropriate permissions to approve pull requests
4. Check that the branch name patterns match exactly (case-sensitive)

## References

- [GitHub Documentation: About protected branches](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/defining-the-mergeability-of-pull-requests/about-protected-branches)
- [GitHub Documentation: Configuring protected branches](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/defining-the-mergeability-of-pull-requests/configuring-protected-branches)