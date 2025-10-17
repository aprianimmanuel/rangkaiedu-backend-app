package mocks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"

	"github.com/aprianimmanuel/rangkaiedu-backend/internal/config"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/monitoring"
)

// MockWhatsAppAPI represents a mock WhatsApp API for testing
type MockWhatsAppAPI struct {
	mu                sync.RWMutex
	sentMessages      []SentMessage
	healthStatus      map[string]bool
	apiCallLatency    time.Duration
	shouldFail        bool
	shouldRateLimit   bool
	rateLimitCount    int
	rateLimitWindow   time.Duration
	lastRateLimitTime time.Time
}

// SentMessage represents a message sent through the mock WhatsApp API
type SentMessage struct {
	To        string    `json:"to"`
	Message   string    `json:"message"`
	Provider  string    `json:"provider"`
	Timestamp time.Time `json:"timestamp"`
	Success   bool      `json:"success"`
	Error     string    `json:"error,omitempty"`
}

// NewMockWhatsAppAPI creates a new mock WhatsApp API
func NewMockWhatsAppAPI() *MockWhatsAppAPI {
	return &MockWhatsAppAPI{
		sentMessages:   make([]SentMessage, 0),
		healthStatus:   make(map[string]bool),
		apiCallLatency: 100 * time.Millisecond, // Default latency
		shouldFail:     false,
		shouldRateLimit: false,
		rateLimitCount: 0,
		rateLimitWindow: 15 * time.Minute,
	}
}

// SetHealthStatus sets the health status for a provider
func (m *MockWhatsAppAPI) SetHealthStatus(provider string, healthy bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.healthStatus[provider] = healthy
}

// SetAPICallLatency sets the simulated API call latency
func (m *MockWhatsAppAPI) SetAPICallLatency(latency time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.apiCallLatency = latency
}

// SetShouldFail sets whether the API should fail
func (m *MockWhatsAppAPI) SetShouldFail(shouldFail bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shouldFail = shouldFail
}

// SetShouldRateLimit sets whether the API should rate limit
func (m *MockWhatsAppAPI) SetShouldRateLimit(shouldRateLimit bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shouldRateLimit = shouldRateLimit
}

// SetRateLimitCount sets the rate limit count
func (m *MockWhatsAppAPI) SetRateLimitCount(count int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.rateLimitCount = count
}

// GetSentMessages returns all sent messages
func (m *MockWhatsAppAPI) GetSentMessages() []SentMessage {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sentMessages
}

// GetSentMessagesByProvider returns messages sent by a specific provider
func (m *MockWhatsAppAPI) GetSentMessagesByProvider(provider string) []SentMessage {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	var messages []SentMessage
	for _, msg := range m.sentMessages {
		if msg.Provider == provider {
			messages = append(messages, msg)
		}
	}
	return messages
}

// ClearSentMessages clears all sent messages
func (m *MockWhatsAppAPI) ClearSentMessages() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sentMessages = make([]SentMessage, 0)
}

// GetHealthStatus returns the health status of all providers
func (m *MockWhatsAppAPI) GetHealthStatus() map[string]bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	status := make(map[string]bool)
	for k, v := range m.healthStatus {
		status[k] = v
	}
	return status
}

// SendWhatsAppMessage simulates sending a WhatsApp message
func (m *MockWhatsAppAPI) SendWhatsAppMessage(cfg *config.Config, to, otp string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Simulate rate limiting
	if m.shouldRateLimit {
		if m.rateLimitCount > 0 {
			if time.Since(m.lastRateLimitTime) < m.rateLimitWindow {
				monitoring.LogAuthFailure(context.Background(), to, "", "", map[string]interface{}{
					"error": "Rate limit exceeded",
					"action": "otp_rate_limit",
				})
				return fmt.Errorf("rate limit exceeded. Please try again later")
			}
		}
		m.lastRateLimitTime = time.Now()
		m.rateLimitCount++
	}
	
	// Simulate API call latency
	time.Sleep(m.apiCallLatency)
	
	// Simulate failure
	if m.shouldFail {
		monitoring.LogAuthFailure(context.Background(), to, "", "", map[string]interface{}{
			"error": "Failed to send WhatsApp OTP",
			"details": "Mock API failure",
		})
		return fmt.Errorf("failed to send WhatsApp OTP: Mock API failure")
	}
	
	// Get the primary WhatsApp provider
	provider, err := cfg.GetPrimaryWhatsAppProvider()
	if err != nil {
		return fmt.Errorf("no WhatsApp provider configured: %w", err)
	}
	
	// Create message
	message := fmt.Sprintf("Your Rangkai Edu OTP code is: %s. This code expires in 10 minutes. Do not share it with anyone. If you didn't request this, please ignore this message.", otp)
	
	// Record the sent message
	sentMsg := SentMessage{
		To:        to,
		Message:   message,
		Provider:  provider.Type,
		Timestamp: time.Now(),
		Success:   true,
	}
	
	m.sentMessages = append(m.sentMessages, sentMsg)
	
	// Log success
	monitoring.LogAuthSuccess(context.Background(), "", to, "", "", map[string]interface{}{
		"action": "otp_sent_whatsapp",
	})
	
	return nil
}

// SendWhatsAppBusinessMessage simulates sending a WhatsApp Business message
func (m *MockWhatsAppAPI) SendWhatsAppBusinessMessage(cfg *config.Config, provider config.WhatsAppProviderConfig, to, otp string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Simulate rate limiting
	if m.shouldRateLimit {
		if m.rateLimitCount > 0 {
			if time.Since(m.lastRateLimitTime) < m.rateLimitWindow {
				monitoring.LogAuthFailure(context.Background(), to, "", "", map[string]interface{}{
					"error": "Rate limit exceeded",
					"action": "otp_rate_limit",
				})
				return fmt.Errorf("rate limit exceeded. Please try again later")
			}
		}
		m.lastRateLimitTime = time.Now()
		m.rateLimitCount++
	}
	
	// Simulate API call latency
	time.Sleep(m.apiCallLatency)
	
	// Simulate failure
	if m.shouldFail {
		monitoring.LogAuthFailure(context.Background(), to, "", "", map[string]interface{}{
			"error": "Failed to send WhatsApp Business OTP",
			"details": "Mock API failure",
		})
		return fmt.Errorf("failed to send WhatsApp Business OTP: Mock API failure")
	}
	
	// Create message
	message := fmt.Sprintf("Your Rangkai Edu OTP code is: %s. This code expires in 10 minutes. Do not share it with anyone. If you didn't request this, please ignore this message.", otp)
	
	// Record the sent message
	sentMsg := SentMessage{
		To:        to,
		Message:   message,
		Provider:  provider.Type,
		Timestamp: time.Now(),
		Success:   true,
	}
	
	m.sentMessages = append(m.sentMessages, sentMsg)
	
	// Log success
	monitoring.LogAuthSuccess(context.Background(), "", to, "", "", map[string]interface{}{
		"action": "otp_sent_whatsapp_business",
	})
	
	return nil
}

// SendTwilioWhatsAppMessage simulates sending a Twilio WhatsApp message
func (m *MockWhatsAppAPI) SendTwilioWhatsAppMessage(cfg *config.Config, provider config.WhatsAppProviderConfig, to, otp string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Simulate rate limiting
	if m.shouldRateLimit {
		if m.rateLimitCount > 0 {
			if time.Since(m.lastRateLimitTime) < m.rateLimitWindow {
				monitoring.LogAuthFailure(context.Background(), to, "", "", map[string]interface{}{
					"error": "Rate limit exceeded",
					"action": "otp_rate_limit",
				})
				return fmt.Errorf("rate limit exceeded. Please try again later")
			}
		}
		m.lastRateLimitTime = time.Now()
		m.rateLimitCount++
	}
	
	// Simulate API call latency
	time.Sleep(m.apiCallLatency)
	
	// Simulate failure
	if m.shouldFail {
		monitoring.LogAuthFailure(context.Background(), to, "", "", map[string]interface{}{
			"error": "Failed to send Twilio WhatsApp OTP",
			"details": "Mock API failure",
		})
		return fmt.Errorf("failed to send Twilio WhatsApp OTP: Mock API failure")
	}
	
	// Create message
	message := fmt.Sprintf("Your Rangkai Edu OTP code is: %s. This code expires in 10 minutes. Do not share it with anyone. If you didn't request this, please ignore this message.", otp)
	
	// Record the sent message
	sentMsg := SentMessage{
		To:        to,
		Message:   message,
		Provider:  provider.Type,
		Timestamp: time.Now(),
		Success:   true,
	}
	
	m.sentMessages = append(m.sentMessages, sentMsg)
	
	// Log success
	monitoring.LogAuthSuccess(context.Background(), "", to, "", "", map[string]interface{}{
		"action": "otp_sent_twilio_whatsapp",
	})
	
	return nil
}

// SendTestWhatsAppMessage simulates sending a test WhatsApp message
func (m *MockWhatsAppAPI) SendTestWhatsAppMessage(cfg *config.Config, to string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Simulate API call latency
	time.Sleep(m.apiCallLatency)
	
	// Simulate failure
	if m.shouldFail {
		return fmt.Errorf("failed to send test WhatsApp message: Mock API failure")
	}
	
	// Get the primary WhatsApp provider
	provider, err := cfg.GetPrimaryWhatsAppProvider()
	if err != nil {
		return fmt.Errorf("no WhatsApp provider configured: %w", err)
	}
	
	// Create test message
	message := "This is a test message from Rangkai Edu WhatsApp service."
	
	// Record the sent message
	sentMsg := SentMessage{
		To:        to,
		Message:   message,
		Provider:  provider.Type,
		Timestamp: time.Now(),
		Success:   true,
	}
	
	m.sentMessages = append(m.sentMessages, sentMsg)
	
	return nil
}

// ValidatePhoneFormat validates that a phone number is in the correct format
func (m *MockWhatsAppAPI) ValidatePhoneFormat(phone string) bool {
	// Basic validation - phone number should start with + and contain only digits and spaces
	if len(phone) < 10 {
		return false
	}
	
	if phone[0] != '+' {
		return false
	}
	
	// Check remaining characters are digits or spaces
	for _, char := range phone[1:] {
		if char != ' ' && (char < '0' || char > '9') {
			return false
		}
	}
	
	return true
}

// FormatPhoneNumber formats a phone number to E.164 format
func (m *MockWhatsAppAPI) FormatPhoneNumber(phone string) string {
	// Remove all non-digit characters except the leading +
	var formatted string
	for _, char := range phone {
		if char == '+' {
			formatted += "+"
		} else if char >= '0' && char <= '9' {
			formatted += string(char)
		}
	}
	
	return formatted
}

// GetRateLimitInfo returns rate limit information
func (m *MockWhatsAppAPI) GetRateLimitInfo() (count int, window time.Duration, lastTime time.Time) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.rateLimitCount, m.rateLimitWindow, m.lastRateLimitTime
}

// ResetRateLimit resets the rate limiter
func (m *MockWhatsAppAPI) ResetRateLimit() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.rateLimitCount = 0
	m.lastRateLimitTime = time.Time{}
}

// CreateHTTPServer creates an HTTP server for mocking WhatsApp API endpoints
func (m *MockWhatsAppAPI) CreateHTTPServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate API call latency
		time.Sleep(m.apiCallLatency)
		
		// Simulate rate limiting
		if m.shouldRateLimit {
			if m.rateLimitCount > 0 {
				if time.Since(m.lastRateLimitTime) < m.rateLimitWindow {
					w.WriteHeader(http.StatusTooManyRequests)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"error": "Rate limit exceeded",
					})
					return
				}
			}
			m.lastRateLimitTime = time.Now()
			m.rateLimitCount++
		}
		
		// Simulate failure
		if m.shouldFail {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "Mock API failure",
			})
			return
		}
		
		// Simulate successful response
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Message sent successfully",
		})
	}))
}

// CreateWhatsAppBusinessHTTPServer creates an HTTP server for mocking WhatsApp Business API
func (m *MockWhatsAppAPI) CreateWhatsAppBusinessHTTPServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate API call latency
		time.Sleep(m.apiCallLatency)
		
		// Simulate rate limiting
		if m.shouldRateLimit {
			if m.rateLimitCount > 0 {
				if time.Since(m.lastRateLimitTime) < m.rateLimitWindow {
					w.WriteHeader(http.StatusTooManyRequests)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"error": "Rate limit exceeded",
					})
					return
				}
			}
			m.lastRateLimitTime = time.Now()
			m.rateLimitCount++
		}
		
		// Simulate failure
		if m.shouldFail {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "Mock WhatsApp Business API failure",
			})
			return
		}
		
		// Parse request body
		var requestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "Invalid request body",
			})
			return
		}
		
		// Record the sent message
		if to, ok := requestBody["to"].(string); ok {
			if text, ok := requestBody["text"].(map[string]interface{}); ok {
				if body, ok := text["body"].(string); ok {
					m.mu.Lock()
					provider := "whatsapp_business" // Default provider type
					if providerType, ok := requestBody["messaging_product"].(string); ok {
						provider = providerType
					}
					
					sentMsg := SentMessage{
						To:        to,
						Message:   body,
						Provider:  provider,
						Timestamp: time.Now(),
						Success:   true,
					}
					m.sentMessages = append(m.sentMessages, sentMsg)
					m.mu.Unlock()
				}
			}
		}
		
		// Simulate successful response
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Message sent successfully",
		})
	}))
}

// CreateTwilioHTTPServer creates an HTTP server for mocking Twilio WhatsApp API
func (m *MockWhatsAppAPI) CreateTwilioHTTPServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate API call latency
		time.Sleep(m.apiCallLatency)
		
		// Simulate rate limiting
		if m.shouldRateLimit {
			if m.rateLimitCount > 0 {
				if time.Since(m.lastRateLimitTime) < m.rateLimitWindow {
					w.WriteHeader(http.StatusTooManyRequests)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"error": "Rate limit exceeded",
					})
					return
				}
			}
			m.lastRateLimitTime = time.Now()
			m.rateLimitCount++
		}
		
		// Simulate failure
		if m.shouldFail {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "Mock Twilio API failure",
			})
			return
		}
		
		// Parse request body
		var requestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "Invalid request body",
			})
			return
		}
		
		// Record the sent message
		if to, ok := requestBody["To"].(string); ok {
			if body, ok := requestBody["Body"].(string); ok {
				m.mu.Lock()
				sentMsg := SentMessage{
					To:        to,
					Message:   body,
					Provider:  "twilio_whatsapp",
					Timestamp: time.Now(),
					Success:   true,
				}
				m.sentMessages = append(m.sentMessages, sentMsg)
				m.mu.Unlock()
			}
		}
		
		// Simulate successful response
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Message sent successfully",
		})
	}))
}

// GetStatistics returns statistics about the mock API usage
func (m *MockWhatsAppAPI) GetStatistics() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	return map[string]interface{}{
		"total_messages_sent": len(m.sentMessages),
		"messages_by_provider": func() map[string]int {
			providerCounts := make(map[string]int)
			for _, msg := range m.sentMessages {
				providerCounts[msg.Provider]++
			}
			return providerCounts
		}(),
		"health_status": m.healthStatus,
		"api_call_latency_ms": m.apiCallLatency.Milliseconds(),
		"should_fail": m.shouldFail,
		"should_rate_limit": m.shouldRateLimit,
		"rate_limit_count": m.rateLimitCount,
		"rate_limit_window_ms": m.rateLimitWindow.Milliseconds(),
		"last_rate_limit_time": m.lastRateLimitTime,
	}
}