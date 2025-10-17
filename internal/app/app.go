package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/config"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/handlers"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/monitoring"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/utils/db"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/utils/https"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/routes"
)

// App represents the application structure
type App struct {
	Config     *config.Config
	DB         *db.DB
	Router     *gin.Engine
	Monitoring *monitoring.MonitoringService
	Container  *Container
}

// New creates a new application instance
func New() (*App, error) {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	database, err := db.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize monitoring
	monitor, err := monitoring.NewMonitoringService(monitoring.DefaultMonitoringConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize monitoring: %w", err)
	}

	app := &App{
		Config:     cfg,
		DB:         database,
		Monitoring: monitor,
		Router:     gin.Default(),
	}

	return app, nil
}

// Initialize sets up the application components
func (a *App) Initialize() error {
	// Initialize database connection pool
	if err := a.DB.Init(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Create dependency injection container
	a.Container = NewContainer(a.Config, a.DB.GetConnection())

	// Initialize monitoring system
	monitoringConfigPath := "./config/monitoring.json"
	if err := monitoring.InitializeFromFile(monitoringConfigPath); err != nil {
		log.Printf("Warning: Failed to initialize monitoring from file: %v", err)
		log.Printf("Initializing monitoring with default configuration")
		if err := monitoring.InitializeWithDefaults(); err != nil {
			log.Printf("Warning: Failed to initialize monitoring with defaults: %v", err)
		}
	}

	// Initialize security logging system
	securityConfigPath := "./config/security-logging-staging.json"
	// Check if staging config exists, otherwise use example config
	if _, err := os.Stat(securityConfigPath); os.IsNotExist(err) {
		securityConfigPath = "./config/security-logging-example.json"
	}

	log.Printf("[DEBUG] Security config path: %s", securityConfigPath)

	// Load security configuration
	log.Printf("[DEBUG] Calling monitoring.NewSecurityConfigManager...")
	securityConfigManager := monitoring.NewSecurityConfigManager(securityConfigPath)
	log.Printf("[DEBUG] SecurityConfigManager created: %+v", securityConfigManager)
	
	if err := securityConfigManager.LoadConfig(); err != nil {
		log.Printf("Warning: Failed to load security configuration: %v", err)
	} else {
		log.Printf("Security configuration loaded successfully")
		// Log some configuration details for debugging
		securityConfig := securityConfigManager.GetConfig()
		if securityConfig != nil && securityConfig.Logging != nil {
			log.Printf("Security logging enabled: %v, level: %v",
				securityConfig.Logging.Enabled, securityConfig.Logging.Level)

			// Initialize security logger with loaded configuration
			if securityConfig.Logging.Enabled {
				log.Printf("[DEBUG] Calling monitoring.NewSecurityEventLogger...")
				securityLogger, err := monitoring.NewSecurityEventLogger(securityConfig.Logging)
				if err != nil {
					log.Printf("Warning: Failed to create security logger: %v", err)
				} else {
					log.Printf("[DEBUG] Security logger created: %+v", securityLogger)
					// Set the security logger on the global monitoring service
					log.Printf("[DEBUG] Checking if monitoring is initialized...")
					if monitoring.IsInitialized() {
						log.Printf("[DEBUG] Monitoring is initialized, getting service...")
						service := monitoring.GetService()
						if service != nil {
							log.Printf("[DEBUG] Service obtained: %+v", service)
							log.Printf("[DEBUG] Calling service.SetSecurityLogger...")
							service.SetSecurityLogger(securityLogger)
							log.Printf("Security logger initialized successfully")
						} else {
							log.Printf("Warning: Monitoring service is nil")
						}
					} else {
						log.Printf("Warning: Monitoring is not initialized")
					}
				}
			}
		}
	}

	return nil
}

// SetupRoutes sets up all application routes
func (a *App) SetupRoutes() {
	// Add HSTS middleware if enabled
	log.Printf("[DEBUG] HSTS enabled: %v", a.Config.HTTPS.HSTS.Enabled)
	if a.Config.HTTPS.HSTS.Enabled {
		log.Printf("[DEBUG] Calling https.HSTSMiddleware...")
		a.Router.Use(https.HSTSMiddleware(&a.Config.HTTPS.HSTS))
		log.Printf("[DEBUG] HSTS middleware applied")
	}

	// Welcome route
	a.Router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Rangkai Edu Backend API",
		})
	})

	// Create handlers
	authHandler := handlers.NewAuthHandler()
	
	// Setup all routes with handler
	routes.SetupAllRoutes(a.Router, authHandler)
}

// StartHTTPRedirectServer starts the HTTP redirect server if enabled
func (a *App) StartHTTPRedirectServer() {
	if a.Config.HTTPS.Redirect.Enabled && a.Config.HTTPS.HTTPPort > 0 && a.Config.HTTPS.Enabled {
		go func() {
			redirectHandler := http.NewServeMux()
			redirectHandler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				// Get the target HTTPS URL
				target := "https://" + r.Host
				if a.Config.HTTPS.Port != 443 {
					target = "https://" + r.Host + ":" + fmt.Sprintf("%d", a.Config.HTTPS.Port)
				}
				target += r.URL.Path
				if r.URL.RawQuery != "" {
					target += "?" + r.URL.RawQuery
				}

				// Set HSTS header if enabled
				if a.Config.HTTPS.HSTS.Enabled {
					hstsValue := fmt.Sprintf("max-age=%d", a.Config.HTTPS.HSTS.MaxAge)
					if a.Config.HTTPS.HSTS.IncludeSubDomains {
						hstsValue += "; includeSubDomains"
					}
					if a.Config.HTTPS.HSTS.Preload {
						hstsValue += "; preload"
					}
					w.Header().Set("Strict-Transport-Security", hstsValue)
				}

				// Perform redirect
				statusCode := a.Config.HTTPS.Redirect.StatusCode
				if statusCode == 0 {
					statusCode = http.StatusMovedPermanently
				}
				http.Redirect(w, r, target, statusCode)
			})

			log.Printf("Starting HTTP redirect server on port %d", a.Config.HTTPS.HTTPPort)
			if err := http.ListenAndServe(fmt.Sprintf(":%d", a.Config.HTTPS.HTTPPort), redirectHandler); err != nil {
				log.Printf("HTTP redirect server error: %v", err)
			}
		}()
	}
}

// Run starts the application
func (a *App) Run() error {
	// Start monitoring service
	if monitoring.IsInitialized() {
		ctx := context.Background()
		if err := monitoring.Start(ctx); err != nil {
			log.Printf("Warning: Failed to start monitoring service: %v", err)
		} else {
			log.Printf("Monitoring service started successfully")
		}
	}

	defer func() {
		// Close database connection
		a.DB.Close()

		// Shutdown monitoring service
		if monitoring.IsInitialized() {
			if err := monitoring.Shutdown(); err != nil {
				log.Printf("Warning: Failed to shutdown monitoring service: %v", err)
			} else {
				log.Printf("Monitoring service shutdown successfully")
			}
		}
	}()

	// Start HTTP redirect server if enabled
	a.StartHTTPRedirectServer()

	// Start HTTPS server if enabled
	if a.Config.HTTPS.Enabled {
		log.Printf("Starting HTTPS server on port %d", a.Config.HTTPS.Port)
		if err := a.Router.RunTLS(fmt.Sprintf(":%d", a.Config.HTTPS.Port), a.Config.HTTPS.Certificate.File.CertFile, a.Config.HTTPS.Certificate.File.KeyFile); err != nil {
			return fmt.Errorf("failed to start HTTPS server: %w", err)
		}
	} else {
		// Fallback to HTTP for development
		log.Printf("Starting HTTP server on port %d", 8080)
		if err := a.Router.Run(":8080"); err != nil {
			return fmt.Errorf("failed to start HTTP server: %w", err)
		}
	}

	return nil
}