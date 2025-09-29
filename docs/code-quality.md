# Code Quality Tools Configuration

## Overview
This document outlines the configuration of code quality tools for both the backend (Go) and frontend (JavaScript/React) components of the Rangkai Edu application.

## Backend Code Quality Tools

### Go Linting Tools

#### 1. golint
golint is a linter for Go source code that checks for style and correctness issues.

**Installation:**
```bash
go install golang.org/x/lint/golint@latest
```

**Usage:**
```bash
# Lint specific package
golint ./package

# Lint entire project
golint ./...
```

**Configuration:**
golint doesn't require configuration files. It follows Go's official style guide.

#### 2. go vet
go vet is a tool that finds suspicious constructs in Go code.

**Usage:**
```bash
# Vet specific package
go vet ./package

# Vet entire project
go vet ./...
```

**Configuration:**
go vet is part of the Go toolchain and doesn't require separate installation.

#### 3. golangci-lint
golangci-lint is a fast Go linters runner that includes golint, go vet, and many other linters.

**Installation:**
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

**Configuration:**
Create a `.golangci.yml` file in the project root:

```yaml
run:
  timeout: 5m
  tests: true

linters:
  enable:
    - gofmt
    - golint
    - govet
    - errcheck
    - staticcheck
    - unused
    - gosimple
    - structcheck
    - varcheck
    - ineffassign
    - deadcode

linters-settings:
  golint:
    min-confidence: 0.8

issues:
  exclude-use-default: false
```

**Usage:**
```bash
# Run all linters
golangci-lint run

# Run linters and fix issues where possible
golangci-lint run --fix
```

### Code Formatting

#### gofmt
gofmt formats Go code according to the standard Go format.

**Usage:**
```bash
# Format specific file
gofmt -w file.go

# Format entire project
gofmt -w .
```

#### goimports
goimports updates Go import lines, adding missing ones and removing unreferenced ones.

**Installation:**
```bash
go install golang.org/x/tools/cmd/goimports@latest
```

**Usage:**
```bash
# Format imports and code
goimports -w .
```

## Frontend Code Quality Tools

### ESLint
ESLint is a static code analysis tool for identifying problematic patterns in JavaScript and JSX code.

**Configuration:**
The project already has an `eslint.config.js` file with basic configuration. We'll enhance it with additional rules and plugins.

**Enhanced ESLint Configuration:**
```javascript
import js from '@eslint/js'
import globals from 'globals'
import reactHooks from 'eslint-plugin-react-hooks'
import reactRefresh from 'eslint-plugin-react-refresh'
import react from 'eslint-plugin-react'
import { defineConfig, globalIgnores } from 'eslint/config'

export default defineConfig([
  globalIgnores(['dist']),
  {
    files: ['**/*.{js,jsx}'],
    extends: [
      js.configs.recommended,
      react.configs['recommended'],
      reactHooks.configs['recommended-latest'],
      reactRefresh.configs.vite,
    ],
    languageOptions: {
      ecmaVersion: 2020,
      globals: globals.browser,
      parserOptions: {
        ecmaVersion: 'latest',
        ecmaFeatures: { jsx: true },
        sourceType: 'module',
      },
    },
    rules: {
      'no-unused-vars': ['error', { varsIgnorePattern: '^[A-Z_]' }],
      'react/prop-types': 'off',
      'react/react-in-jsx-scope': 'off',
    },
    settings: {
      react: {
        version: 'detect',
      },
    },
  },
])
```

### Prettier
Prettier is an opinionated code formatter that enforces consistent code style.

**Installation:**
```bash
npm install -D prettier eslint-config-prettier eslint-plugin-prettier
```

**Configuration:**
Create a `.prettierrc` file in the project root:

```json
{
  "semi": false,
  "trailingComma": "es5",
  "singleQuote": true,
  "printWidth": 80,
  "tabWidth": 2,
  "useTabs": false,
  "bracketSpacing": true,
  "arrowParens": "avoid"
}
```

**Integration with ESLint:**
Update `eslint.config.js` to include Prettier:

```javascript
import js from '@eslint/js'
import globals from 'globals'
import reactHooks from 'eslint-plugin-react-hooks'
import reactRefresh from 'eslint-plugin-react-refresh'
import react from 'eslint-plugin-react'
import prettier from 'eslint-plugin-prettier'
import prettierConfig from 'eslint-config-prettier'
import { defineConfig, globalIgnores } from 'eslint/config'

export default defineConfig([
  globalIgnores(['dist']),
  {
    files: ['**/*.{js,jsx}'],
    extends: [
      js.configs.recommended,
      react.configs['recommended'],
      reactHooks.configs['recommended-latest'],
      reactRefresh.configs.vite,
      prettierConfig,
    ],
    languageOptions: {
      ecmaVersion: 2020,
      globals: globals.browser,
      parserOptions: {
        ecmaVersion: 'latest',
        ecmaFeatures: { jsx: true },
        sourceType: 'module',
      },
    },
    plugins: {
      prettier,
    },
    rules: {
      'no-unused-vars': ['error', { varsIgnorePattern: '^[A-Z_]' }],
      'react/prop-types': 'off',
      'react/react-in-jsx-scope': 'off',
      'prettier/prettier': 'error',
    },
    settings: {
      react: {
        version: 'detect',
      },
    },
  },
])
```

**Prettier Ignore File:**
Create a `.prettierignore` file:

```
dist
node_modules
.env
.env.local
.env.development.local
.env.test.local
.env.production.local
```

## Running Code Quality Tools

### Backend
```bash
# Run linters
golangci-lint run

# Format code
go fmt ./...
goimports -w .

# Vet code
go vet ./...
```

### Frontend
```bash
# Run ESLint
npm run lint

# Fix ESLint issues
npm run lint -- --fix

# Format with Prettier
npx prettier --write .
```

## Git Hooks Integration

To ensure code quality is maintained, we'll integrate these tools with Git hooks using husky and lint-staged.

**Installation:**
```bash
npm install -D husky lint-staged
```

**Configuration:**
Add to `package.json`:

```json
{
  "scripts": {
    "prepare": "husky install"
  },
  "lint-staged": {
    "*.{js,jsx}": [
      "eslint --fix",
      "prettier --write"
    ],
    "*.go": [
      "gofmt -w",
      "goimports -w"
    ]
  }
}
```

Create `.husky/pre-commit`:

```bash
#!/bin/sh
. "$(dirname "$0")/_/husky.sh"

npx lint-staged
```

## Continuous Integration

Code quality checks will be integrated into the CI/CD pipeline:

1. **Backend Pipeline:**
   - Run `golangci-lint` on all Go files
   - Run `go vet` on all Go files
   - Check code formatting with `gofmt`

2. **Frontend Pipeline:**
   - Run ESLint on all JavaScript/JSX files
   - Run Prettier to check formatting
   - Fail build if any issues are found

## Editor Integration

### VS Code
Install the following extensions:
- Go (for Go language support)
- ESLint (for JavaScript/JSX linting)
- Prettier - Code formatter (for code formatting)

### Configuration
Create `.vscode/settings.json`:

```json
{
  "editor.formatOnSave": true,
  "editor.defaultFormatter": "esbenp.prettier-vscode",
  "editor.codeActionsOnSave": {
    "source.fixAll.eslint": true
  },
  "go.formatTool": "goimports",
  "go.lintTool": "golangci-lint",
  "[go]": {
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
      "source.organizeImports": true
    }
  }
}
```

## Best Practices

1. **Consistent Code Style:**
   - Use the same formatting rules across the project
   - Enforce style rules through automated tools

2. **Early Detection:**
   - Run linters and formatters in development
   - Integrate with Git hooks to prevent bad code from being committed

3. **Performance:**
   - Configure tools to only check changed files
   - Use caching where possible

4. **Documentation:**
   - Document the code quality standards
   - Provide clear instructions for running tools

## Next Steps

1. Install and configure golangci-lint for backend
2. Install and configure Prettier for frontend
3. Set up Git hooks for automated code quality checks
4. Configure editor integrations
5. Document the code quality standards