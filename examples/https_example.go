// Package example demonstrates how to use the HTTPS utilities
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aprianimmanuel/rangkaiedu-backend/config"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/config"
	"github.com/gin-gonic/gin"
)

func main() {
	// Example 1: Load HTTPS configuration
	// fmt.Println("=== Loading HTTPS Configuration ===")
	// httpsConfig := config.LoadHTTPSConfig()
	// fmt.Printf("HTTPS Enabled: %v\n", httpsConfig.Enabled)
	// fmt.Printf("HTTPS Port: %d\n", httpsConfig.Port)
	// fmt.Printf("HTTP Port: %d\n", httpsConfig.HTTPPort)
	// fmt.Printf("Certificate Type: %s\n", httpsConfig.Certificate.Type)

	// Example 2: Create HTTPS manager
	// fmt.Println("\n=== Creating HTTPS Manager ===")
	// manager, err := https.NewManager(httpsConfig)
	// if err != nil {
	// 	log.Printf("Warning: Failed to create HTTPS manager: %v", err)
	// } else {
	// 	fmt.Println("HTTPS manager created successfully")
	// }

	// Example 3: Create a simple Gin server with HSTS middleware
	// fmt.Println("\n=== Creating Gin Server with HSTS ===")
	// r := gin.Default()

	// Add HSTS middleware if enabled
	// if httpsConfig.HSTS.Enabled {
	// 	r.Use(https.HSTSMiddleware(&httpsConfig.HSTS))
	// 	fmt.Println("HSTS middleware added")
	// }

	// Add a simple route
	// r.GET("/", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "Hello HTTPS World!",
	// 		"https":   https.IsHTTPSRequest(c.Request),
	// 		"scheme":  https.GetScheme(c.Request),
	// 	})
	// })

	// Add health check route
	// r.GET("/health", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"status": "healthy",
	// 		"time":   time.Now().Format(time.RFC3339),
	// 	})
	// })

	// Example 4: Demonstrate HTTPS detection
	// r.GET("/detect", func(c *gin.Context) {
	// 	isHTTPS := https.IsHTTPSRequest(c.Request)
	// 	scheme := https.GetScheme(c.Request)
		
	// 	c.JSON(200, gin.H{
	// 		"is_https": isHTTPS,
	// 		"scheme":   scheme,
	// 		"headers":  c.Request.Header,
	// 	})
	// })

	// fmt.Println("Server routes configured")

	// Note: In a real application, you would start the server here
	// For this example, we'll just show what would be done
	// fmt.Println("\n=== Server Startup (Example) ===")
	// fmt.Println("To start the server, you would typically call:")
	// fmt.Println("  r.RunTLS(\":443\", \"cert.pem\", \"key.pem\")")
	// fmt.Println("or use the HTTPS manager to start both HTTPS and HTTP redirect servers")

	// Example 5: Show environment variable configuration
	// fmt.Println("\n=== Environment Configuration Example ===")
	// fmt.Println("To enable HTTPS, set these environment variables:")
	// fmt.Println("  HTTPS_ENABLED=true")
	// fmt.Println("  HTTPS_PORT=8443")
	// fmt.Println("  HTTPS_CERTIFICATE_FILE_CERT_FILE=./certs/server.crt")
	// fmt.Println("  HTTPS_CERTIFICATE_FILE_KEY_FILE=./certs/server.key")

	// Example 6: Show certificate generation command
	// fmt.Println("\n=== Certificate Generation ===")
	// fmt.Println("For development, generate self-signed certificates with:")
	// fmt.Println("  mkdir -p certs")
	// fmt.Println("  openssl req -x509 -newkey rsa:4096 -keyout certs/server.key -out certs/server.crt -days 365 -nodes -subj \"/CN=localhost\"")

	// fmt.Println("\n=== HTTPS Example Complete ===")

	// Keep the example running for a short time to show output
	// time.Sleep(100 * time.Millisecond)

	// In a real application, you would start the server here:
	// r.RunTLS(":8443", "./certs/server.crt", "./certs/server.key")
	
	// For demonstration purposes, we'll exit here
	// os.Exit(0)
	
	fmt.Println("HTTPS example")
}