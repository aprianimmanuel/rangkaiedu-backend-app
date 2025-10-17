package utils_test

import (
	"testing"
	"time"
	"net/http"

	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/utils/mfa"
)

func TestGenerateTOTPSecret(t *testing.T) {
	secret, err := mfa.GenerateTOTPSecret()
	
	assert.NoError(t, err)
	assert.NotEmpty(t, secret)
	// Base32 encoded 20-byte secret should be 32 characters
	assert.Len(t, secret, 32)
}

func TestGenerateBackupCodes(t *testing.T) {
	codes, err := mfa.GenerateBackupCodes()
	
	assert.NoError(t, err)
	assert.Len(t, codes, mfa.BackupCodeCount)
	
	// Check that all codes are unique
	codeMap := make(map[string]bool)
	for _, code := range codes {
		assert.Len(t, code, mfa.BackupCodeLength)
		assert.NotContains(t, codeMap, code)
		codeMap[code] = true
	}
}

func TestValidateTOTP(t *testing.T) {
	// Generate a secret for testing
	secret, err := mfa.GenerateTOTPSecret()
	assert.NoError(t, err)
	
	// Generate a valid TOTP token
	token, err := totp.GenerateCode(secret, time.Now())
	assert.NoError(t, err)
	
	// Test valid token
	assert.True(t, mfa.ValidateTOTP(secret, token))
	
	// Test invalid token
	assert.False(t, mfa.ValidateTOTP(secret, "invalidtoken"))
	
	// Test empty secret
	assert.False(t, mfa.ValidateTOTP("", token))
	
	// Test empty token
	assert.False(t, mfa.ValidateTOTP(secret, ""))
}

func TestValidateBackupCode(t *testing.T) {
	codes := []string{"ABC123", "DEF456", "GHI789", "JKL012", "MNO345", "PQR678"}
	
	// Test valid backup code
	assert.True(t, mfa.ValidateBackupCode(codes, "ABC123"))
	
	// Test invalid backup code
	assert.False(t, mfa.ValidateBackupCode(codes, "INVALID"))
	
	// Test empty codes
	assert.False(t, mfa.ValidateBackupCode([]string{}, "ABC123"))
	
	// Test empty code
	assert.False(t, mfa.ValidateBackupCode(codes, ""))
}

func TestCreateQRCode(t *testing.T) {
	secret, err := mfa.GenerateTOTPSecret()
	assert.NoError(t, err)
	
	qrCode, err := mfa.CreateQRCode(secret, "test@example.com")
	
	assert.NoError(t, err)
	assert.NotEmpty(t, qrCode)
	
	// The QR code should be a valid base64 string
	// We can't easily test the actual QR code content, but we can verify it's not empty
}

func TestGenerateRandomCode(t *testing.T) {
	// Test generating codes of different lengths
	for length := 1; length <= 10; length++ {
		code, err := mfa.GenerateRandomCode(length)
		assert.NoError(t, err)
		assert.Len(t, code, length)
		
		// Verify all characters are alphanumeric
		for _, char := range code {
			assert.Contains(t, "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ", string(char))
		}
	}
	
	// Test edge cases
	_, err := mfa.GenerateRandomCode(0)
	assert.Error(t, err)
	
	_, err = mfa.GenerateRandomCode(-1)
	assert.Error(t, err)
}

func TestGetClientIP(t *testing.T) {
	// Test cases for different IP header scenarios
	testCases := []struct {
		name           string
		headers        map[string]string
		remoteAddr     string
		expectedIP     string
	}{
		{
			name:       "X-Forwarded-For header present",
			headers:    map[string]string{"X-Forwarded-For": "192.168.1.100, 10.0.0.1"},
			remoteAddr: "172.16.0.1:12345",
			expectedIP: "192.168.1.100",
		},
		{
			name:       "X-Real-IP header present",
			headers:    map[string]string{"X-Real-IP": "10.0.0.1"},
			remoteAddr: "172.16.0.1:12345",
			expectedIP: "10.0.0.1",
		},
		{
			name:       "No headers, use RemoteAddr",
			headers:    map[string]string{},
			remoteAddr: "192.168.1.100:12345",
			expectedIP: "192.168.1.100",
		},
		{
			name:       "X-Forwarded-For with multiple IPs",
			headers:    map[string]string{"X-Forwarded-For": "203.0.113.1, 198.51.100.1, 192.0.2.1"},
			remoteAddr: "172.16.0.1:12345",
			expectedIP: "203.0.113.1",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a proper http.Request
			req := &http.Request{
				Header: http.Header{},
				RemoteAddr: tc.remoteAddr,
			}
			
			// Set headers properly
			for key, value := range tc.headers {
				req.Header.Set(key, value)
			}
			
			ip := mfa.GetClientIP(req)
			assert.Equal(t, tc.expectedIP, ip)
		})
	}
}

func TestGetUserAgent(t *testing.T) {
	userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
	
	req := &http.Request{
		Header: http.Header{
			"User-Agent": []string{userAgent},
		},
	}
	
	assert.Equal(t, userAgent, mfa.GetUserAgent(req))
}

// mockRequest is a mock implementation of http.Request for testing
type mockRequest struct {
	headers    map[string]string
	remoteAddr string
}

func (m *mockRequest) Header() map[string][]string {
	headers := make(map[string][]string)
	for k, v := range m.headers {
		headers[k] = []string{v}
	}
	return headers
}

func (m *mockRequest) RemoteAddr() string {
	return m.remoteAddr
}