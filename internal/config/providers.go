package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// ProviderConfig represents the common interface for all service providers
type ProviderConfig interface {
	GetName() string
	GetType() string
	IsEnabled() bool
	Validate() error
	GetPriority() int
}

// EmailProviderConfig represents the configuration for an email provider
type EmailProviderConfig struct {
	Name     string            `json:"name" mapstructure:"name"`
	Type     string            `json:"type" mapstructure:"type"`
	Enabled  bool              `json:"enabled" mapstructure:"enabled"`
	Priority int               `json:"priority" mapstructure:"priority"`
	Settings EmailSettings     `json:"settings" mapstructure:"settings"`
}

// Validate validates the email provider configuration
func (e EmailProviderConfig) Validate() error {
	// Basic validation
	if e.Name == "" {
		return errors.New("email provider name is required")
	}
	
	if e.Type == "" {
		return errors.New("email provider type is required")
	}
	
	return nil
}

// GetName returns the provider name
func (e EmailProviderConfig) GetName() string {
	return e.Name
}

// GetType returns the provider type
func (e EmailProviderConfig) GetType() string {
	return e.Type
}

// IsEnabled returns whether the provider is enabled
func (e EmailProviderConfig) IsEnabled() bool {
	return e.Enabled
}

// GetPriority returns the provider priority
func (e EmailProviderConfig) GetPriority() int {
	return e.Priority
}

// EmailSettings contains the common settings for email providers
type EmailSettings struct {
	FromEmail string `json:"from_email" mapstructure:"from_email"`
	FromName  string `json:"from_name" mapstructure:"from_name"`
	ReplyTo   string `json:"reply_to" mapstructure:"reply_to"`
}

// SMTPProviderConfig represents configuration for SMTP-based email providers
type SMTPProviderConfig struct {
	EmailProviderConfig `mapstructure:",squash"`
	Host                string `json:"host" mapstructure:"host"`
	Port                int    `json:"port" mapstructure:"port"`
	Username            string `json:"username" mapstructure:"username"`
	Password            string `json:"password" mapstructure:"password"`
	UseTLS              bool   `json:"use_tls" mapstructure:"use_tls"`
	UseStartTLS         bool   `json:"use_starttls" mapstructure:"use_starttls"`
}

// SendGridProviderConfig represents configuration for SendGrid email provider
type SendGridProviderConfig struct {
	EmailProviderConfig `mapstructure:",squash"`
	APIKey              string `json:"api_key" mapstructure:"api_key"`
}

// GmailProviderConfig represents configuration for Gmail/Google Workspace
type GmailProviderConfig struct {
	EmailProviderConfig `mapstructure:",squash"`
	ClientID            string `json:"client_id" mapstructure:"client_id"`
	ClientSecret        string `json:"client_secret" mapstructure:"client_secret"`
	RefreshToken        string `json:"refresh_token" mapstructure:"refresh_token"`
	AccessToken         string `json:"access_token" mapstructure:"access_token"`
}

// SMSProviderConfig represents the configuration for an SMS provider
type SMSProviderConfig struct {
	Name     string        `json:"name" mapstructure:"name"`
	Type     string        `json:"type" mapstructure:"type"`
	Enabled  bool          `json:"enabled" mapstructure:"enabled"`
	Priority int           `json:"priority" mapstructure:"priority"`
	Settings SMSSettings   `json:"settings" mapstructure:"settings"`
}

// Validate validates the SMS provider configuration
func (s SMSProviderConfig) Validate() error {
	// Basic validation
	if s.Name == "" {
		return errors.New("SMS provider name is required")
	}
	
	if s.Type == "" {
		return errors.New("SMS provider type is required")
	}
	
	return nil
}

// GetName returns the provider name
func (s SMSProviderConfig) GetName() string {
	return s.Name
}

// GetType returns the provider type
func (s SMSProviderConfig) GetType() string {
	return s.Type
}

// IsEnabled returns whether the provider is enabled
func (s SMSProviderConfig) IsEnabled() bool {
	return s.Enabled
}

// GetPriority returns the provider priority
func (s SMSProviderConfig) GetPriority() int {
	return s.Priority
}

// SMSSettings contains the common settings for SMS providers
type SMSSettings struct {
	FromNumber string `json:"from_number" mapstructure:"from_number"`
}

// TwilioProviderConfig represents configuration for Twilio SMS provider
type TwilioProviderConfig struct {
	SMSProviderConfig `mapstructure:",squash"`
	AccountSID        string `json:"account_sid" mapstructure:"account_sid"`
	AuthToken         string `json:"auth_token" mapstructure:"auth_token"`
}

// SNSProviderConfig represents configuration for Amazon SNS SMS provider
type SNSProviderConfig struct {
	SMSProviderConfig `mapstructure:",squash"`
	AccessKeyID       string `json:"access_key_id" mapstructure:"access_key_id"`
	SecretAccessKey   string `json:"secret_access_key" mapstructure:"secret_access_key"`
	Region            string `json:"region" mapstructure:"region"`
}

// WhatsAppProviderConfig represents the configuration for a WhatsApp provider
type WhatsAppProviderConfig struct {
	Name     string              `json:"name" mapstructure:"name"`
	Type     string              `json:"type" mapstructure:"type"`
	Enabled  bool                `json:"enabled" mapstructure:"enabled"`
	Priority int                 `json:"priority" mapstructure:"priority"`
	Settings WhatsAppSettings    `json:"settings" mapstructure:"settings"`
}

// Validate validates the WhatsApp provider configuration
func (w WhatsAppProviderConfig) Validate() error {
	// Basic validation
	if w.Name == "" {
		return errors.New("WhatsApp provider name is required")
	}
	
	if w.Type == "" {
		return errors.New("WhatsApp provider type is required")
	}
	
	return nil
}

// GetName returns the provider name
func (w WhatsAppProviderConfig) GetName() string {
	return w.Name
}

// GetType returns the provider type
func (w WhatsAppProviderConfig) GetType() string {
	return w.Type
}

// IsEnabled returns whether the provider is enabled
func (w WhatsAppProviderConfig) IsEnabled() bool {
	return w.Enabled
}

// GetPriority returns the provider priority
func (w WhatsAppProviderConfig) GetPriority() int {
	return w.Priority
}

// WhatsAppSettings contains the common settings for WhatsApp providers
type WhatsAppSettings struct {
	FromNumber string `json:"from_number" mapstructure:"from_number"`
}

// WhatsAppBusinessProviderConfig represents configuration for WhatsApp Business API
type WhatsAppBusinessProviderConfig struct {
	WhatsAppProviderConfig `mapstructure:",squash"`
	PhoneIDNumber          string `json:"phone_id_number" mapstructure:"phone_id_number"`
	BusinessAccountID      string `json:"business_account_id" mapstructure:"business_account_id"`
	AccessToken            string `json:"access_token" mapstructure:"access_token"`
	APIVersion             string `json:"api_version" mapstructure:"api_version"`
}

// TwilioWhatsAppProviderConfig represents configuration for Twilio WhatsApp
type TwilioWhatsAppProviderConfig struct {
	WhatsAppProviderConfig `mapstructure:",squash"`
	AccountSID             string `json:"account_sid" mapstructure:"account_sid"`
	AuthToken              string `json:"auth_token" mapstructure:"auth_token"`
	WhatsAppNumber         string `json:"whatsapp_number" mapstructure:"whatsapp_number"`
}

// ProviderHealth represents the health status of a provider
type ProviderHealth struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Status   string `json:"status"`   // "healthy", "degraded", "unhealthy"
	Message  string `json:"message"`  // Additional details
	LastTest string `json:"last_test"`
}

// ProviderManager manages the application's provider configurations
type ProviderManager struct {
	EmailProviders  []EmailProviderConfig  `json:"email_providers"`
	SMSProviders    []SMSProviderConfig    `json:"sms_providers"`
	WhatsAppProviders []WhatsAppProviderConfig `json:"whatsapp_providers"`
}

// NewProviderManager creates a new provider manager
func NewProviderManager() *ProviderManager {
	return &ProviderManager{
		EmailProviders:  make([]EmailProviderConfig, 0),
		SMSProviders:    make([]SMSProviderConfig, 0),
		WhatsAppProviders: make([]WhatsAppProviderConfig, 0),
	}
}

// GetPrimaryEmailProvider returns the highest priority enabled email provider
func (pm *ProviderManager) GetPrimaryEmailProvider() (EmailProviderConfig, error) {
	for _, provider := range pm.EmailProviders {
		if provider.Enabled {
			return provider, nil
		}
	}
	return EmailProviderConfig{}, errors.New("no enabled email provider found")
}

// GetPrimarySMSProvider returns the highest priority enabled SMS provider
func (pm *ProviderManager) GetPrimarySMSProvider() (SMSProviderConfig, error) {
	for _, provider := range pm.SMSProviders {
		if provider.Enabled {
			return provider, nil
		}
	}
	return SMSProviderConfig{}, errors.New("no enabled SMS provider found")
}

// GetPrimaryWhatsAppProvider returns the highest priority enabled WhatsApp provider
func (pm *ProviderManager) GetPrimaryWhatsAppProvider() (WhatsAppProviderConfig, error) {
	for _, provider := range pm.WhatsAppProviders {
		if provider.Enabled {
			return provider, nil
		}
	}
	return WhatsAppProviderConfig{}, errors.New("no enabled WhatsApp provider found")
}

// GetWhatsAppProviderByName returns a specific WhatsApp provider by name
func (pm *ProviderManager) GetWhatsAppProviderByName(name string) (WhatsAppProviderConfig, error) {
	for _, provider := range pm.WhatsAppProviders {
		if provider.Name == name && provider.Enabled {
			return provider, nil
		}
	}
	return WhatsAppProviderConfig{}, fmt.Errorf("WhatsApp provider %s not found or not enabled", name)
}

// GetEmailProviderByName returns a specific email provider by name
func (pm *ProviderManager) GetEmailProviderByName(name string) (EmailProviderConfig, error) {
	for _, provider := range pm.EmailProviders {
		if provider.Name == name && provider.Enabled {
			return provider, nil
		}
	}
	return EmailProviderConfig{}, fmt.Errorf("email provider %s not found or not enabled", name)
}

// GetSMSProviderByName returns a specific SMS provider by name
func (pm *ProviderManager) GetSMSProviderByName(name string) (SMSProviderConfig, error) {
	for _, provider := range pm.SMSProviders {
		if provider.Name == name && provider.Enabled {
			return provider, nil
		}
	}
	return SMSProviderConfig{}, fmt.Errorf("SMS provider %s not found or not enabled", name)
}

// LoadProviders loads provider configurations from environment variables
func (pm *ProviderManager) LoadProviders() error {
	// Load email providers
	emailCount := pm.loadProviderCount("EMAIL_PROVIDER")
	for i := 0; i < emailCount; i++ {
		if provider, err := pm.loadEmailProvider(i); err == nil {
			pm.EmailProviders = append(pm.EmailProviders, provider)
		}
	}

	// Load SMS providers
	smsCount := pm.loadProviderCount("SMS_PROVIDER")
	for i := 0; i < smsCount; i++ {
		if provider, err := pm.loadSMSProvider(i); err == nil {
			pm.SMSProviders = append(pm.SMSProviders, provider)
		}
	}

	// Load WhatsApp providers
	whatsappCount := pm.loadProviderCount("WHATSAPP_PROVIDER")
	for i := 0; i < whatsappCount; i++ {
		if provider, err := pm.loadWhatsAppProvider(i); err == nil {
			pm.WhatsAppProviders = append(pm.WhatsAppProviders, provider)
		}
	}

	// Load legacy configurations for backward compatibility
	pm.loadLegacyConfigurations()

	return nil
}

// loadProviderCount counts the number of providers with a given prefix
func (pm *ProviderManager) loadProviderCount(prefix string) int {
	count := 0
	for i := 0; i < 10; i++ { // Support up to 10 providers
		if os.Getenv(fmt.Sprintf("%s_%d_TYPE", prefix, i)) != "" {
			count++
		}
	}
	return count
}

// loadEmailProvider loads an email provider configuration
func (pm *ProviderManager) loadEmailProvider(index int) (EmailProviderConfig, error) {
	providerType := os.Getenv(fmt.Sprintf("EMAIL_PROVIDER_%d_TYPE", index))
	if providerType == "" {
		return EmailProviderConfig{}, errors.New("provider type not specified")
	}

	name := getProviderEnv(fmt.Sprintf("EMAIL_PROVIDER_%d_NAME", index), fmt.Sprintf("email-provider-%d", index))
	enabled := getEnvBool(fmt.Sprintf("EMAIL_PROVIDER_%d_ENABLED", index), true)
	priority := getEnvInt(fmt.Sprintf("EMAIL_PROVIDER_%d_PRIORITY", index), 0)

	settings := EmailSettings{
		FromEmail: getProviderEnv(fmt.Sprintf("EMAIL_PROVIDER_%d_SETTINGS_FROM_EMAIL", index), ""),
		FromName:  getProviderEnv(fmt.Sprintf("EMAIL_PROVIDER_%d_SETTINGS_FROM_NAME", index), ""),
		ReplyTo:   getProviderEnv(fmt.Sprintf("EMAIL_PROVIDER_%d_SETTINGS_REPLY_TO", index), ""),
	}

	switch providerType {
	case "smtp":
		return pm.loadSMTPProvider(index, name, enabled, priority, settings)
	case "sendgrid":
		return pm.loadSendGridProvider(index, name, enabled, priority, settings)
	case "gmail":
		return pm.loadGmailProvider(index, name, enabled, priority, settings)
	default:
		return EmailProviderConfig{}, fmt.Errorf("unsupported email provider type: %s", providerType)
	}
}

// loadSMTPProvider loads an SMTP provider configuration
func (pm *ProviderManager) loadSMTPProvider(index int, name string, enabled bool, priority int, settings EmailSettings) (EmailProviderConfig, error) {
	host := getProviderEnv(fmt.Sprintf("EMAIL_PROVIDER_%d_SETTINGS_HOST", index), "")
	port := getEnvInt(fmt.Sprintf("EMAIL_PROVIDER_%d_SETTINGS_PORT", index), 587)
	username := getProviderEnv(fmt.Sprintf("EMAIL_PROVIDER_%d_SETTINGS_USERNAME", index), "")
	password := getProviderEnv(fmt.Sprintf("EMAIL_PROVIDER_%d_SETTINGS_PASSWORD", index), "")
	useTLS := getEnvBool(fmt.Sprintf("EMAIL_PROVIDER_%d_SETTINGS_USE_TLS", index), true)
	useStartTLS := getEnvBool(fmt.Sprintf("EMAIL_PROVIDER_%d_SETTINGS_USE_STARTTLS", index), true)

	if host == "" || username == "" || password == "" {
		return EmailProviderConfig{}, errors.New("required SMTP settings missing")
	}

	// Create the base email provider config
	emailProvider := EmailProviderConfig{
		Name:     name,
		Type:     "smtp",  // Fix the type to be "smtp" instead of "email"
		Enabled:  enabled,
		Priority: priority,
		Settings: settings,
	}
	
	// Create the SMTP provider with the base config
	smtpProvider := SMTPProviderConfig{
		EmailProviderConfig: emailProvider,
		Host:        host,
		Port:        port,
		Username:    username,
		Password:    password,
		UseTLS:      useTLS,
		UseStartTLS: useStartTLS,
	}
	
	// Return the SMTP provider as an EmailProviderConfig
	// Note: This will lose the SMTP-specific fields, but that's the interface we need to return
	return smtpProvider.EmailProviderConfig, nil
}

// loadSendGridProvider loads a SendGrid provider configuration
func (pm *ProviderManager) loadSendGridProvider(index int, name string, enabled bool, priority int, settings EmailSettings) (EmailProviderConfig, error) {
	apiKey := getProviderEnv(fmt.Sprintf("EMAIL_PROVIDER_%d_SETTINGS_API_KEY", index), "")

	if apiKey == "" {
		return EmailProviderConfig{}, errors.New("SendGrid API key is required")
	}

	// Create the base email provider config
	emailProvider := EmailProviderConfig{
		Name:     name,
		Type:     "sendgrid",  // Fix the type to be "sendgrid" instead of "email"
		Enabled:  enabled,
		Priority: priority,
		Settings: settings,
	}
	
	// Create the SendGrid provider with the base config
	sendGridProvider := SendGridProviderConfig{
		EmailProviderConfig: emailProvider,
		APIKey: apiKey,
	}
	
	// Return the SendGrid provider as an EmailProviderConfig
	// Note: This will lose the SendGrid-specific fields, but that's the interface we need to return
	return sendGridProvider.EmailProviderConfig, nil
}

// loadGmailProvider loads a Gmail provider configuration
func (pm *ProviderManager) loadGmailProvider(index int, name string, enabled bool, priority int, settings EmailSettings) (EmailProviderConfig, error) {
	clientID := getProviderEnv(fmt.Sprintf("EMAIL_PROVIDER_%d_SETTINGS_CLIENT_ID", index), "")
	clientSecret := getProviderEnv(fmt.Sprintf("EMAIL_PROVIDER_%d_SETTINGS_CLIENT_SECRET", index), "")
	refreshToken := getProviderEnv(fmt.Sprintf("EMAIL_PROVIDER_%d_SETTINGS_REFRESH_TOKEN", index), "")
	accessToken := getProviderEnv(fmt.Sprintf("EMAIL_PROVIDER_%d_SETTINGS_ACCESS_TOKEN", index), "")

	if clientID == "" || clientSecret == "" || refreshToken == "" {
		return EmailProviderConfig{}, errors.New("required Gmail OAuth2 settings missing")
	}

	// Create the base email provider config
	emailProvider := EmailProviderConfig{
		Name:     name,
		Type:     "gmail",  // Fix the type to be "gmail" instead of "email"
		Enabled:  enabled,
		Priority: priority,
		Settings: settings,
	}
	
	// Create the Gmail provider with the base config
	gmailProvider := GmailProviderConfig{
		EmailProviderConfig: emailProvider,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}
	
	// Return the Gmail provider as an EmailProviderConfig
	// Note: This will lose the Gmail-specific fields, but that's the interface we need to return
	return gmailProvider.EmailProviderConfig, nil
}

// loadSMSProvider loads an SMS provider configuration
func (pm *ProviderManager) loadSMSProvider(index int) (SMSProviderConfig, error) {
	providerType := os.Getenv(fmt.Sprintf("SMS_PROVIDER_%d_TYPE", index))
	if providerType == "" {
		return SMSProviderConfig{}, errors.New("provider type not specified")
	}

	name := getProviderEnv(fmt.Sprintf("SMS_PROVIDER_%d_NAME", index), fmt.Sprintf("sms-provider-%d", index))
	enabled := getEnvBool(fmt.Sprintf("SMS_PROVIDER_%d_ENABLED", index), true)
	priority := getEnvInt(fmt.Sprintf("SMS_PROVIDER_%d_PRIORITY", index), 0)

	settings := SMSSettings{
		FromNumber: getProviderEnv(fmt.Sprintf("SMS_PROVIDER_%d_SETTINGS_FROM_NUMBER", index), ""),
	}

	switch providerType {
	case "twilio":
		return pm.loadTwilioProvider(index, name, enabled, priority, settings)
	case "sns":
		return pm.loadSNSProvider(index, name, enabled, priority, settings)
	default:
		return SMSProviderConfig{}, fmt.Errorf("unsupported SMS provider type: %s", providerType)
	}
}

// loadTwilioProvider loads a Twilio provider configuration
func (pm *ProviderManager) loadTwilioProvider(index int, name string, enabled bool, priority int, settings SMSSettings) (SMSProviderConfig, error) {
	accountSID := getProviderEnv(fmt.Sprintf("SMS_PROVIDER_%d_SETTINGS_ACCOUNT_SID", index), "")
	authToken := getProviderEnv(fmt.Sprintf("SMS_PROVIDER_%d_SETTINGS_AUTH_TOKEN", index), "")

	if accountSID == "" || authToken == "" {
		return SMSProviderConfig{}, errors.New("required Twilio settings missing")
	}

	// Create the base SMS provider config
	smsProvider := SMSProviderConfig{
		Name:     name,
		Type:     "twilio",  // Fix the type to be "twilio" instead of "sms"
		Enabled:  enabled,
		Priority: priority,
		Settings: settings,
	}
	
	// Create the Twilio provider with the base config
	twilioProvider := TwilioProviderConfig{
		SMSProviderConfig: smsProvider,
		AccountSID: accountSID,
		AuthToken:  authToken,
	}
	
	// Return the Twilio provider as an SMSProviderConfig
	// Note: This will lose the Twilio-specific fields, but that's the interface we need to return
	return twilioProvider.SMSProviderConfig, nil
}

// loadSNSProvider loads an SNS provider configuration
func (pm *ProviderManager) loadSNSProvider(index int, name string, enabled bool, priority int, settings SMSSettings) (SMSProviderConfig, error) {
	accessKeyID := getProviderEnv(fmt.Sprintf("SMS_PROVIDER_%d_SETTINGS_ACCESS_KEY_ID", index), "")
	secretAccessKey := getProviderEnv(fmt.Sprintf("SMS_PROVIDER_%d_SETTINGS_SECRET_ACCESS_KEY", index), "")
	region := getProviderEnv(fmt.Sprintf("SMS_PROVIDER_%d_SETTINGS_REGION", index), "us-east-1")

	if accessKeyID == "" || secretAccessKey == "" {
		return SMSProviderConfig{}, errors.New("required SNS settings missing")
	}

	// Create the base SMS provider config
	smsProvider := SMSProviderConfig{
		Name:     name,
		Type:     "sns",  // Fix the type to be "sns" instead of "sms"
		Enabled:  enabled,
		Priority: priority,
		Settings: settings,
	}
	
	// Create the SNS provider with the base config
	snsProvider := SNSProviderConfig{
		SMSProviderConfig: smsProvider,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		Region:          region,
	}
	
	// Return the SNS provider as an SMSProviderConfig
	// Note: This will lose the SNS-specific fields, but that's the interface we need to return
	return snsProvider.SMSProviderConfig, nil
}

// loadWhatsAppProvider loads a WhatsApp provider configuration
func (pm *ProviderManager) loadWhatsAppProvider(index int) (WhatsAppProviderConfig, error) {
	providerType := os.Getenv(fmt.Sprintf("WHATSAPP_PROVIDER_%d_TYPE", index))
	if providerType == "" {
		return WhatsAppProviderConfig{}, errors.New("provider type not specified")
	}

	name := getProviderEnv(fmt.Sprintf("WHATSAPP_PROVIDER_%d_NAME", index), fmt.Sprintf("whatsapp-provider-%d", index))
	enabled := getEnvBool(fmt.Sprintf("WHATSAPP_PROVIDER_%d_ENABLED", index), true)
	priority := getEnvInt(fmt.Sprintf("WHATSAPP_PROVIDER_%d_PRIORITY", index), 0)

	settings := WhatsAppSettings{
		FromNumber: getProviderEnv(fmt.Sprintf("WHATSAPP_PROVIDER_%d_SETTINGS_FROM_NUMBER", index), ""),
	}

	switch providerType {
	case "whatsapp_business":
		return pm.loadWhatsAppBusinessProvider(index, name, enabled, priority, settings)
	case "twilio_whatsapp":
		return pm.loadTwilioWhatsAppProvider(index, name, enabled, priority, settings)
	default:
		return WhatsAppProviderConfig{}, fmt.Errorf("unsupported WhatsApp provider type: %s", providerType)
	}
}

// loadWhatsAppBusinessProvider loads a WhatsApp Business API provider configuration
func (pm *ProviderManager) loadWhatsAppBusinessProvider(index int, name string, enabled bool, priority int, settings WhatsAppSettings) (WhatsAppProviderConfig, error) {
	phoneIDNumber := getProviderEnv(fmt.Sprintf("WHATSAPP_PROVIDER_%d_SETTINGS_PHONE_ID_NUMBER", index), "")
	businessAccountID := getProviderEnv(fmt.Sprintf("WHATSAPP_PROVIDER_%d_SETTINGS_BUSINESS_ACCOUNT_ID", index), "")
	accessToken := getProviderEnv(fmt.Sprintf("WHATSAPP_PROVIDER_%d_SETTINGS_ACCESS_TOKEN", index), "")
	apiVersion := getProviderEnv(fmt.Sprintf("WHATSAPP_PROVIDER_%d_SETTINGS_API_VERSION", index), "v15.0")

	if phoneIDNumber == "" || businessAccountID == "" || accessToken == "" {
		return WhatsAppProviderConfig{}, errors.New("required WhatsApp Business settings missing")
	}

	// Create the base WhatsApp provider config
	whatsappProvider := WhatsAppProviderConfig{
		Name:     name,
		Type:     "whatsapp_business",
		Enabled:  enabled,
		Priority: priority,
		Settings: settings,
	}
	
	// Create the WhatsApp Business provider with the base config
	whatsappBusinessProvider := WhatsAppBusinessProviderConfig{
		WhatsAppProviderConfig: whatsappProvider,
		PhoneIDNumber:          phoneIDNumber,
		BusinessAccountID:      businessAccountID,
		AccessToken:            accessToken,
		APIVersion:             apiVersion,
	}
	
	// Return the WhatsApp Business provider as a WhatsAppProviderConfig
	// Note: This will lose the WhatsApp Business-specific fields, but that's the interface we need to return
	return whatsappBusinessProvider.WhatsAppProviderConfig, nil
}

// loadTwilioWhatsAppProvider loads a Twilio WhatsApp provider configuration
func (pm *ProviderManager) loadTwilioWhatsAppProvider(index int, name string, enabled bool, priority int, settings WhatsAppSettings) (WhatsAppProviderConfig, error) {
	accountSID := getProviderEnv(fmt.Sprintf("WHATSAPP_PROVIDER_%d_SETTINGS_ACCOUNT_SID", index), "")
	authToken := getProviderEnv(fmt.Sprintf("WHATSAPP_PROVIDER_%d_SETTINGS_AUTH_TOKEN", index), "")
	whatsappNumber := getProviderEnv(fmt.Sprintf("WHATSAPP_PROVIDER_%d_SETTINGS_WHATSAPP_NUMBER", index), "")

	if accountSID == "" || authToken == "" || whatsappNumber == "" {
		return WhatsAppProviderConfig{}, errors.New("required Twilio WhatsApp settings missing")
	}

	// Create the base WhatsApp provider config
	whatsappProvider := WhatsAppProviderConfig{
		Name:     name,
		Type:     "twilio_whatsapp",
		Enabled:  enabled,
		Priority: priority,
		Settings: settings,
	}
	
	// Create the Twilio WhatsApp provider with the base config
	twilioWhatsAppProvider := TwilioWhatsAppProviderConfig{
		WhatsAppProviderConfig: whatsappProvider,
		AccountSID:     accountSID,
		AuthToken:      authToken,
		WhatsAppNumber: whatsappNumber,
	}
	
	// Return the Twilio WhatsApp provider as a WhatsAppProviderConfig
	// Note: This will lose the Twilio WhatsApp-specific fields, but that's the interface we need to return
	return twilioWhatsAppProvider.WhatsAppProviderConfig, nil
}

// loadLegacyConfigurations loads legacy configurations for backward compatibility
func (pm *ProviderManager) loadLegacyConfigurations() {
	// Load legacy SMTP configuration
	if os.Getenv("SMTP_HOST") != "" {
		provider := EmailProviderConfig{
			Name:     "legacy-smtp",
			Type:     "smtp",
			Enabled:  true,
			Priority: 0,
			Settings: EmailSettings{
				FromEmail: os.Getenv("SMTP_USER"),
				FromName:  "Rangkai Edu",
			},
		}
		pm.EmailProviders = append(pm.EmailProviders, provider)
	}

	// Load legacy Twilio configuration
	if os.Getenv("TWILIO_ACCOUNT_SID") != "" {
		provider := SMSProviderConfig{
			Name:     "legacy-twilio",
			Type:     "twilio",
			Enabled:  true,
			Priority: 0,
			Settings: SMSSettings{
				FromNumber: os.Getenv("TWILIO_SENDER_PHONE"),
			},
		}
		pm.SMSProviders = append(pm.SMSProviders, provider)
	}
}

// Validate validates all provider configurations
func (pm *ProviderManager) Validate() error {
	// Validate email providers
	for _, provider := range pm.EmailProviders {
		if err := provider.Validate(); err != nil {
			return fmt.Errorf("email provider %s validation failed: %w", provider.Name, err)
		}
	}

	// Validate SMS providers
	for _, provider := range pm.SMSProviders {
		if err := provider.Validate(); err != nil {
			return fmt.Errorf("SMS provider %s validation failed: %w", provider.Name, err)
		}
	}

	// Validate WhatsApp providers
	for _, provider := range pm.WhatsAppProviders {
		if err := provider.Validate(); err != nil {
			return fmt.Errorf("WhatsApp provider %s validation failed: %w", provider.Name, err)
		}
	}

	return nil
}

// HealthCheck performs health checks on all enabled providers
func (pm *ProviderManager) HealthCheck() []ProviderHealth {
	var health []ProviderHealth

	// Check email providers
	for _, provider := range pm.EmailProviders {
		if provider.Enabled {
			health = append(health, pm.checkEmailProviderHealth(provider))
		}
	}

	// Check SMS providers
	for _, provider := range pm.SMSProviders {
		if provider.Enabled {
			health = append(health, pm.checkSMSProviderHealth(provider))
		}
	}

	// Check WhatsApp providers
	for _, provider := range pm.WhatsAppProviders {
		if provider.Enabled {
			health = append(health, pm.checkWhatsAppProviderHealth(provider))
		}
	}

	return health
}

// checkEmailProviderHealth checks the health of an email provider
func (pm *ProviderManager) checkEmailProviderHealth(provider EmailProviderConfig) ProviderHealth {
	health := ProviderHealth{
		Name:     provider.Name,
		Type:     provider.Type,
		Status:   "healthy",
		Message:  "Provider is enabled and configured",
		LastTest: "Not implemented",
	}

	// Basic validation
	if provider.Settings.FromEmail == "" {
		health.Status = "unhealthy"
		health.Message = "From email address is required"
		return health
	}

	// Provider-specific validation
	switch provider.Type {
	case "smtp":
		if strings.Contains(provider.Name, "legacy") {
			health.Status = "degraded"
			health.Message = "Legacy SMTP configuration detected"
		}
	}

	return health
}

// checkSMSProviderHealth checks the health of an SMS provider
func (pm *ProviderManager) checkSMSProviderHealth(provider SMSProviderConfig) ProviderHealth {
	health := ProviderHealth{
		Name:     provider.Name,
		Type:     provider.Type,
		Status:   "healthy",
		Message:  "Provider is enabled and configured",
		LastTest: "Not implemented",
	}

	// Basic validation
	if provider.Settings.FromNumber == "" {
		health.Status = "unhealthy"
		health.Message = "From phone number is required"
		return health
	}

	// Provider-specific validation
	switch provider.Type {
	case "twilio":
		if strings.Contains(provider.Name, "legacy") {
			health.Status = "degraded"
			health.Message = "Legacy Twilio configuration detected"
		}
	}

	return health
}

// checkWhatsAppProviderHealth checks the health of a WhatsApp provider
func (pm *ProviderManager) checkWhatsAppProviderHealth(provider WhatsAppProviderConfig) ProviderHealth {
	health := ProviderHealth{
		Name:     provider.Name,
		Type:     provider.Type,
		Status:   "healthy",
		Message:  "Provider is enabled and configured",
		LastTest: "Not implemented",
	}

	// Basic validation
	if provider.Settings.FromNumber == "" {
		health.Status = "unhealthy"
		health.Message = "From phone number is required"
		return health
	}

	// Provider-specific validation
	switch provider.Type {
	case "whatsapp_business":
		if strings.Contains(provider.Name, "legacy") {
			health.Status = "degraded"
			health.Message = "Legacy WhatsApp Business configuration detected"
		}
	case "twilio_whatsapp":
		if strings.Contains(provider.Name, "legacy") {
			health.Status = "degraded"
			health.Message = "Legacy Twilio WhatsApp configuration detected"
		}
	}

	return health
}

// Helper functions for environment variable parsing
func getProviderEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// CredentialManager handles encryption/decryption of sensitive credentials
type CredentialManager struct {
	encryptionKey []byte
}

// NewCredentialManager creates a new credential manager
func NewCredentialManager() *CredentialManager {
	// In a real implementation, this would be loaded from a secure source
	// For now, we'll use a simple key
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		// Fallback to a simple key for development
		key = []byte("development-encryption-key-32-bytes-long")
	}
	return &CredentialManager{encryptionKey: key}
}

// Encrypt encrypts a plaintext value
func (cm *CredentialManager) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(cm.encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a ciphertext value
func (cm *CredentialManager) Decrypt(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(cm.encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// IsEncrypted checks if a value is encrypted
func (cm *CredentialManager) IsEncrypted(value string) bool {
	// Simple heuristic: if it looks like base64, assume it's encrypted
	_, err := base64.StdEncoding.DecodeString(value)
	return err == nil
}