# Frontend Testing Framework Setup

## Overview
This document outlines the setup and configuration of testing frameworks for the Rangkai Edu frontend application. We will be implementing both unit testing with Vitest and end-to-end testing with Cypress.

## Unit Testing with Vitest

### Introduction
Vitest is a fast unit test framework powered by Vite. It's designed to be a drop-in replacement for Jest with better performance and integration with Vite projects.

### Installation
Vitest can be installed as a development dependency:
```bash
npm install -D vitest jsdom @testing-library/react @testing-library/jest-dom
```

### Configuration
Vitest can be configured in the `vite.config.js` file or through a separate `vitest.config.js` file.

Example configuration in `vite.config.js`:
```javascript
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  test: {
    environment: 'jsdom',
    setupFiles: './tests/setup.js',
    globals: true,
  }
})
```

### Directory Structure
```
frontend-app/
├── src/
│   ├── components/
│   ├── pages/
│   └── utils/
├── tests/
│   ├── components/
│   ├── pages/
│   ├── utils/
│   └── setup.js
└── vite.config.js
```

### Writing Tests
Tests are written in files with the `.test.jsx` or `.spec.jsx` extension.

Example test:
```javascript
import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import Button from '../src/components/ui/Button'

describe('Button', () => {
  it('renders with correct text', () => {
    render(<Button>Click me</Button>)
    expect(screen.getByText('Click me')).toBeInTheDocument()
  })
})
```

### Running Tests
Tests can be run using the following commands:
```bash
# Run all tests
npm run test

# Run tests in watch mode
npm run test:watch

# Run tests with coverage
npm run test:coverage
```

## End-to-End Testing with Cypress

### Introduction
Cypress is a next-generation front-end testing tool built for the modern web. It provides a complete end-to-end testing framework with a rich API for writing tests.

### Installation
Cypress can be installed as a development dependency:
```bash
npm install -D cypress @testing-library/cypress
```

### Configuration
Cypress is configured through the `cypress.config.js` file.

Example configuration:
```javascript
import { defineConfig } from 'cypress'

export default defineConfig({
  e2e: {
    baseUrl: 'http://localhost:3000',
    setupNodeEvents(on, config) {
      // implement node event listeners here
    },
  },
})
```

### Directory Structure
```
frontend-app/
├── src/
│   ├── components/
│   └── pages/
├── cypress/
│   ├── e2e/
│   ├── fixtures/
│   ├── support/
│   └── downloads/
└── cypress.config.js
```

### Writing Tests
Cypress tests are written in files with the `.cy.js` extension in the `cypress/e2e` directory.

Example test:
```javascript
describe('Login Flow', () => {
  it('successfully logs in', () => {
    cy.visit('/login')
    cy.get('[data-cy="phone-input"]').type('+6281234567890')
    cy.get('[data-cy="send-otp-button"]').click()
    cy.get('[data-cy="otp-input"]').type('123456')
    cy.get('[data-cy="verify-otp-button"]').click()
    cy.url().should('include', '/dashboard')
  })
})
```

### Running Tests
Tests can be run using the following commands:
```bash
# Open Cypress UI
npm run cypress:open

# Run Cypress tests in headless mode
npm run cypress:run
```

## Test Categories

### Unit Tests
Unit tests focus on testing individual components and functions in isolation:
- Component rendering tests
- Utility function tests
- Hook tests
- State management tests

### Integration Tests
Integration tests verify that components work together correctly:
- Component composition tests
- Context provider tests
- API integration tests

### End-to-End Tests
End-to-end tests verify complete user workflows:
- Authentication flows
- Navigation tests
- Form submissions
- Data persistence tests

## Testing Best Practices

### 1. Test Structure
Follow the AAA pattern:
- **Arrange**: Set up the test data and conditions
- **Act**: Execute the functionality being tested
- **Assert**: Verify the expected outcomes

### 2. Test Descriptions
Use clear, descriptive test names that explain what is being tested and what the expected outcome is.

### 3. Test Data
Use realistic test data that represents actual usage scenarios.

### 4. Test Isolation
Ensure tests are independent and can run in any order.

### 5. Mocking
Use mocking for external dependencies:
- API calls
- Browser APIs
- Third-party libraries

## Code Coverage

Both Vitest and Cypress provide code coverage reporting:
- Vitest uses Istanbul for coverage reporting
- Cypress has built-in coverage reporting

## Continuous Integration

Tests will be integrated into the CI/CD pipeline:
- Unit tests run on every commit
- End-to-end tests run on deployment
- Coverage reports generated for each test run

## Next Steps

1. Install Vitest and related dependencies
2. Configure Vitest for the project
3. Install Cypress and related dependencies
4. Configure Cypress for the project
5. Write sample tests for existing components
6. Configure CI/CD pipeline for automated testing