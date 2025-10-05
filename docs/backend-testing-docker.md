# Running Backend Middleware Tests in Docker

This document outlines the correct procedure for executing Go middleware tests within the Docker environment using `docker compose`.

## Problem Statement

When attempting to run tests for the backend middleware, it's crucial to ensure they are executed within the designated test environment to avoid failures or unexpected behavior. Running tests against the main application service (`backend`) will not work as intended, as that service is configured to run the application, not the test suite.

## Correct Procedure

To properly run the middleware tests, you must target the `test` service defined in your `docker-compose.yml` file. This service is specifically built using `Dockerfile.test` and is set up with the necessary dependencies and configurations for executing Go tests, including database connectivity.

1.  **Ensure Docker Compose Services are Up (Optional but Recommended):**
    Before running tests, it's often beneficial to ensure your `db` service (and potentially others that your tests might depend on) is running and healthy. You can start all services in detached mode:
    ```bash
    docker compose up -d db
    ```
    Wait for the `db` service to be healthy before proceeding.

2.  **Execute Middleware Tests:**
    Use the `docker compose run` command, specifying the `test` service, to execute your Go middleware tests. The `--rm` flag ensures that the container is removed after the command exits, keeping your system clean.

    ```bash
    docker compose run --rm test sh -c "go test -v ./middleware/... -v"
    ```

    *   `docker compose run`: Executes a one-off command in a new container based on a service.
    *   `--rm`: Automatically remove the container when it exits.
    *   `test`: The name of the service defined in `docker-compose.yml` that is configured for running tests.
    *   `sh -c "go test -v ./middleware/... -v"`: The command to be executed inside the `test` container. This command runs all tests within the `./middleware/...` directory with verbose output.

## Summary of Changes

The key change is to use `docker compose run --rm test ...` instead of `docker compose run --rm backend ...`. This ensures that the tests are run in the environment specifically designed for them, leveraging the `Dockerfile.test` and the associated configurations in `docker-compose.yml`.