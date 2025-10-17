
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/aprianimmanuel/rangkaiedu-backend/config"
	"github.com/aprianimmanuel/rangkaiedu-backend/utils/providers"
)

func main() {
	// Define command line flags
	providerType := flag.String("type", "", "Provider type (smtp, sendgrid, gmail, twilio, sns)")
	providerName := flag.String("name", "", "Provider name")
	list := flag.Bool("list", false, "List all configured providers")
	health := flag.Bool("health", false, "Check provider health status")
	validate := flag.Bool("validate", false, "Validate all provider configurations")
	encrypt := flag.Bool("encrypt", false, "Encrypt all credentials")
	decrypt := flag.Bool("decrypt", false, "Decrypt all credentials")
	generate := flag.Bool("generate", false, "Generate provider configuration template")
	
	flag.Parse()

	// Load configuration
	cfg := config.Load()
	
	// Create provider manager
	pm := providers.NewProviderManager(cfg)

	// Execute commands based on flags
	if *list {
		printProviderConfiguration(cfg)
		return
	}

	if *health {
		printProviderHealth(cfg)
		return
	}

	if *validate {
		if err := pm.ValidateAllProviders(); err != nil {
			log.Fatalf("Provider validation failed: %v", err)
		}
		fmt.Println("All provider configurations are valid")
		return
	}

	if *encrypt {
		if err := pm.EncryptCredentials(); err != nil {
			log.Fatalf("Failed to encrypt credentials: %v", err)
		}
		fmt.Println("All credentials have been encrypted")
		return
	}

	if *decrypt {
		if err := pm.DecryptCredentials(); err != nil {
			log.Fatalf("Failed to decrypt credentials: %v", err)
		}
		fmt.Println("All credentials have been decrypted")
		return
	}

	if *generate {
		if *providerType == "" || *providerName == "" {
			log.Fatal("Both --type and --name are required for generating configuration")
		}
		
		config, err := providers.GenerateProviderConfig(*providerType, *providerName)
		if err != nil {
			log.Fatalf("Failed to generate configuration: %v", err)
		}
		
		fmt.Println(config)
		return
	}

	// If no specific command, show help
	if flag.NArg() == 0 {
		printHelp()
		return
	}

	// Handle specific provider operations
	if *providerType != "" && *providerName != "" {
		handleProviderOperation(cfg, *providerType, *providerName)
		return
	}

	printHelp()
}

func printHelp() {
	fmt.Println("Provider Management CLI")
	fmt.Println("========================")
	fmt.Println("Usage:")
	fmt.Println("  providers [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -config <path>           Path to configuration file")
	fmt.Println("  -type <type>             Provider type (smtp, sendgrid, gmail, twilio, sns)")
	fmt.Println("  -name <name>             Provider name")
	fmt.Println("  -list                    List all configured providers")
	fmt.Println("  -health                  Check provider health status")
	fmt.Println("  -validate                Validate all provider configurations")
	fmt.Println("  -encrypt                 Encrypt all credentials")
	fmt.Println("  -decrypt                 Decrypt all credentials")
	fmt.Println("  -generate                Generate provider configuration template")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  providers -list")
	fmt.Println("  providers -health")
	fmt.Println("  providers -validate")
	fmt.Println("  providers -encrypt")
	fmt.Println("  providers -decrypt")
	fmt.Println("  providers -generate -type smtp -name my-smtp")
	fmt.Println("  providers -type smtp -name my-smtp -test")
}

func printProviderConfiguration(cfg *config.Config) {
	fmt.Println("=== Provider Configuration ===")
	fmt.Println()
	
	fmt.Println("Email Providers:")
	fmt.Println("----------------")
	for i, provider := range cfg.ProviderManager.EmailProviders {
		fmt.Printf("%d. %s\n", i+1, provider.Name)
		fmt.Printf("   Type: %s\n", provider.Type)
		fmt.Printf("   Enabled: %t\n", provider.Enabled)
		fmt.Printf("   Priority: %d\n", provider.Priority)
		fmt.Printf("   From Email: %s\n", provider.From)
		fmt.Println()
	}
	
	fmt.Println("SMS Providers:")
	fmt.Println("---------------")
	for i, provider := range cfg.ProviderManager.SMSProviders {
		fmt.Printf("%d. %s\n", i+1, provider.Name)
		fmt.Printf("   Type: %s\n", provider.Type)
		fmt.Printf("   Enabled: %t\n", provider.Enabled)
		fmt.Printf("   Priority: %d\n", provider.Priority)
		fmt.Printf("   API Key: %s\n", provider.APIKey)
		fmt.Printf("   Endpoint: %s\n", provider.Endpoint)
		fmt.Println()
	}
}

func printProviderHealth(cfg *config.Config) {
	fmt.Println("=== Provider Health Status ===")
	fmt.Println()
	
	health := cfg.GetProviderHealth()
	if len(health) == 0 {
		fmt.Println("No providers configured")
		return
	}
	
	for _, h := range health {
		fmt.Printf("Provider: %s (%s)\n", h.Name, h.Type)
		fmt.Printf("Status: %s\n", h.Status)
		fmt.Printf("Message: %s\n", h.Message)
		fmt.Println()
	}
}

func handleProviderOperation(cfg *config.Config, providerType, providerName string) {
	switch providerType {
	case "smtp":
		handleSMTPProvider(cfg, providerName)
	case "sendgrid":
		handleSendGridProvider(cfg, providerName)
	case "gmail":
		handleGmailProvider(cfg, providerName)
	case "twilio":
		handleTwilioProvider(cfg, providerName)
	case "sns":
		handleSNSProvider(cfg, providerName)
	default:
		log.Fatalf("Unsupported provider type: %s", providerType)
	}
}

func handleSMTPProvider(cfg *config.Config, providerName string) {
	// Find the provider in the EmailProviders list
	var provider config.EmailProvider
	for _, p := range cfg.ProviderManager.EmailProviders {
		if p.Name == providerName {
			provider = p
			break
		}
	}
	
	if provider.Name == "" {
		log.Fatalf("Failed to get SMTP provider: %s not found", providerName)
	}
	
	fmt.Printf("=== SMTP Provider: %s ===\n", providerName)
	fmt.Printf("Type: %s\n", provider.Type)
	fmt.Printf("Host: %s\n", provider.Host)
	fmt.Printf("Port: %d\n", provider.Port)
	fmt.Printf("Username: %s\n", provider.Username)
	fmt.Printf("Password: %s\n", provider.Password)
	fmt.Printf("From: %s\n", provider.From)
	fmt.Printf("Enabled: %t\n", provider.Enabled)
	fmt.Printf("Priority: %d\n", provider.Priority)
}

func handleSendGridProvider(cfg *config.Config, providerName string) {
	// Find the provider in the EmailProviders list
	var provider config.EmailProvider
	for _, p := range cfg.ProviderManager.EmailProviders {
		if p.Name == providerName {
			provider = p
			break
		}
	}
	
	if provider.Name == "" {
		log.Fatalf("Failed to get SendGrid provider: %s not found", providerName)
	}
	
	fmt.Printf("=== SendGrid Provider: %s ===\n", providerName)
	fmt.Printf("Type: %s\n", provider.Type)
	fmt.Printf("From: %s\n", provider.From)
	fmt.Printf("Enabled: %t\n", provider.Enabled)
	fmt.Printf("Priority: %d\n", provider.Priority)
}

func handleGmailProvider(cfg *config.Config, providerName string) {
	// Find the provider in the EmailProviders list
	var provider config.EmailProvider
	for _, p := range cfg.ProviderManager.EmailProviders {
		if p.Name == providerName {
			provider = p
			break
		}
	}
	
	if provider.Name == "" {
		log.Fatalf("Failed to get Gmail provider: %s not found", providerName)
	}
	
	fmt.Printf("=== Gmail Provider: %s ===\n", providerName)
	fmt.Printf("Type: %s\n", provider.Type)
	fmt.Printf("From: %s\n", provider.From)
	fmt.Printf("Enabled: %t\n", provider.Enabled)
	fmt.Printf("Priority: %d\n", provider.Priority)
}

func handleTwilioProvider(cfg *config.Config, providerName string) {
	// Find the provider in the SMSProviders list
	var provider config.SMSProvider
	for _, p := range cfg.ProviderManager.SMSProviders {
		if p.Name == providerName {
			provider = p
			break
		}
	}
	
	if provider.Name == "" {
		log.Fatalf("Failed to get Twilio provider: %s not found", providerName)
	}
	
	fmt.Printf("=== Twilio Provider: %s ===\n", providerName)
	fmt.Printf("Type: %s\n", provider.Type)
	fmt.Printf("API Key: %s\n", provider.APIKey)
	fmt.Printf("Endpoint: %s\n", provider.Endpoint)
	fmt.Printf("Enabled: %t\n", provider.Enabled)
	fmt.Printf("Priority: %d\n", provider.Priority)
}

func handleSNSProvider(cfg *config.Config, providerName string) {
	// Find the provider in the SMSProviders list
	var provider config.SMSProvider
	for _, p := range cfg.ProviderManager.SMSProviders {
		if p.Name == providerName {
			provider = p
			break
		}
	}
	
	if provider.Name == "" {
		log.Fatalf("Failed to get SNS provider: %s not found", providerName)
	}
	
	fmt.Printf("=== SNS Provider: %s ===\n", providerName)
	fmt.Printf("Type: %s\n", provider.Type)
	fmt.Printf("API Key: %s\n", provider.APIKey)
	fmt.Printf("Endpoint: %s\n", provider.Endpoint)
	fmt.Printf("Enabled: %t\n", provider.Enabled)
	fmt.Printf("Priority: %d\n", provider.Priority)
}
