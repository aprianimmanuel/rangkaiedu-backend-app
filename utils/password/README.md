# Password Utility Package

This package provides secure password hashing and verification functionality using the bcrypt algorithm for the Rangkai Edu authentication system.

## Security Features

- Uses bcrypt algorithm, which is specifically designed for password hashing
- Implements a cost factor of 12, providing a good balance between security and performance
- Automatically handles salt generation for each password hash
- Resistant to rainbow table and brute force attacks

## Functions

### HashPassword

```go
func HashPassword(password string) (string, error)
```

Generates a bcrypt hash for the given password.

**Parameters:**
- `password` (string): The plain text password to hash

**Returns:**
- `string`: The bcrypt hash of the password
- `error`: Any error that occurred during hashing

**Example:**
```go
hashedPassword, err := password.HashPassword("mySecretPassword")
if err != nil {
    // Handle error
    log.Fatal("Error hashing password:", err)
}
```

### CheckPasswordHash

```go
func CheckPasswordHash(password, hash string) bool
```

Compares a plain text password with a bcrypt hash.

**Parameters:**
- `password` (string): The plain text password to verify
- `hash` (string): The bcrypt hash to compare against

**Returns:**
- `bool`: True if the password matches the hash, false otherwise

**Example:**
```go
isValid := password.CheckPasswordHash("mySecretPassword", hashedPassword)
if isValid {
    // Password is correct
} else {
    // Password is incorrect
}
```

## Usage in Authentication System

This package will be used by:
1. The user registration handler to hash passwords before storing in the database
2. The user login handler to verify passwords against stored hashes

## Security Considerations

1. **Cost Factor**: The bcrypt cost factor is set to 12, which provides a good balance between security and performance. Higher values increase security but also increase computation time exponentially.

2. **Salt Generation**: The package automatically generates a unique salt for each password, preventing rainbow table attacks.

3. **Error Handling**: Proper error handling should be implemented when using these functions to prevent information leakage.

4. **Password Policies**: While this package securely hashes passwords, it's important to enforce strong password policies in the application layer.