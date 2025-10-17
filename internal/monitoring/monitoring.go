package monitoring

import (
	"context"
	"fmt"
	"sync"
)


var (
	// globalService holds the global monitoring service instance
	globalService *MonitoringService
	globalOnce    sync.Once
	globalInit   bool
)

// InitializeFromFile initializes the monitoring service from a configuration file
func InitializeFromFile(configPath string) error {
	var err error
	
	globalOnce.Do(func() {
		configManager := NewConfigManager(configPath)
		if err = configManager.LoadConfig(); err != nil {
			return
		}
		
		globalService, err = NewMonitoringService(configManager.GetConfig())
		if err != nil {
			return
		}
		
		globalInit = true
	})
	
	if err != nil {
		return fmt.Errorf("failed to initialize monitoring service from file: %v", err)
	}
	
	return nil
}

// InitializeWithDefaults initializes the monitoring service with default configuration
func InitializeWithDefaults() error {
	var err error
	
	globalOnce.Do(func() {
		config := DefaultMonitoringConfig()
		globalService, err = NewMonitoringService(config)
		if err != nil {
			return
		}
		
		globalInit = true
	})
	
	if err != nil {
		return fmt.Errorf("failed to initialize monitoring service with defaults: %v", err)
	}
	
	return nil
}

// GetService returns the global monitoring service instance
func GetService() *MonitoringService {
	if !globalInit {
		return nil
	}
	return globalService
}

// IsInitialized returns whether the monitoring service has been initialized
func IsInitialized() bool {
	return globalInit
}

// Start starts the global monitoring service
func Start(ctx context.Context) error {
	if !globalInit {
		return fmt.Errorf("monitoring service is not initialized")
	}
	
	return globalService.Start(ctx)
}

// Shutdown shuts down the global monitoring service
func Shutdown() error {
	if !globalInit {
		return fmt.Errorf("monitoring service is not initialized")
	}
	
	return globalService.Stop()
}