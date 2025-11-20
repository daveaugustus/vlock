package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/daveaugustus/vlock/pkg/config"
	"github.com/daveaugustus/vlock/pkg/vlock"
)

func main() {
	fmt.Println("=== Voltage Go Wrapper - Initialization Example ===\n")

	// Step 1: Load configuration
	fmt.Println("Step 1: Loading configuration...")
	cfg, err := loadConfiguration()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	fmt.Printf("✓ Configuration loaded successfully\n")
	fmt.Printf("  - App Name: %s\n", cfg.AppName)
	fmt.Printf("  - App Version: %s\n", cfg.AppVersion)
	fmt.Printf("  - Environment: %s\n\n", cfg.AppEnv)

	// Step 2: Create Voltage client
	fmt.Println("Step 2: Creating Voltage client...")
	client, err := vlock.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	fmt.Println("✓ Client created successfully\n")

	// Step 3: Initialize the client
	fmt.Println("Step 3: Initializing Voltage library...")
	if err := client.Initialize(); err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}
	fmt.Println("✓ Voltage library initialized successfully\n")

	// Step 4: Verify client is healthy
	fmt.Println("Step 4: Performing health check...")
	if err := client.HealthCheck(); err != nil {
		log.Printf("WARNING: Health check failed: %v\n", err)
	} else {
		fmt.Println("✓ Health check passed\n")
	}

	// Step 5: Get client information
	fmt.Println("Step 5: Client information...")
	info := client.Info()
	fmt.Printf("  - Initialized: %v\n", info.Initialized)
	fmt.Printf("  - Healthy: %v\n", info.Healthy)
	if !info.LastHealthCheck.IsZero() {
		fmt.Printf("  - Last Health Check: %v\n", info.LastHealthCheck.Format(time.RFC3339))
	}
	fmt.Println()

	// Step 6: Simulate some work
	fmt.Println("Step 6: Simulating application work...")
	fmt.Println("  (In a real application, you would perform encryption/decryption here)")
	time.Sleep(2 * time.Second)
	fmt.Println("✓ Work completed\n")

	// Step 7: Perform another health check
	fmt.Println("Step 7: Final health check...")
	if err := client.HealthCheck(); err != nil {
		log.Printf("WARNING: Health check failed: %v\n", err)
	} else {
		fmt.Println("✓ Health check passed\n")
	}

	// Step 8: Clean shutdown
	fmt.Println("Step 8: Closing Voltage client...")
	if err := client.Close(); err != nil {
		log.Printf("WARNING: Error during shutdown: %v\n", err)
	} else {
		fmt.Println("✓ Client closed successfully\n")
	}

	fmt.Println("=== Example completed successfully ===")
}

// loadConfiguration demonstrates different ways to load configuration
func loadConfiguration() (*config.Config, error) {
	// Try to load from configuration file first
	configFile := os.Getenv("VOLTAGE_CONFIG_FILE")
	if configFile == "" {
		configFile = "voltage.cfg"
	}

	if _, err := os.Stat(configFile); err == nil {
		fmt.Printf("  Attempting to load from configuration file: %s...\n", configFile)
		cfg, err := config.LoadConfig(configFile)
		if err != nil {
			// Config file exists but has validation errors
			// Fall back to programmatic configuration for this example
			fmt.Printf("  Configuration file has errors, using programmatic config instead\n")
			fmt.Printf("  (Error was: %v)\n", err)
		} else {
			return cfg, nil
		}
	}

	// Create configuration programmatically
	// This is useful for testing or when config comes from another source
	fmt.Println("  Creating configuration programmatically...")
	cfg := &config.Config{
		AppName:         "VoltageExampleApp",
		AppVersion:      "1.0.0",
		AppEnv:          "DEV",
		DEKSharedSecret: "example_shared_secret",
		ConfigFilePath:  "voltage.cfg",
		NetworkTimeout:  30,
		LogLevel:        2,
	}

	return cfg, nil
}

// Example of error handling with Voltage errors
func demonstrateErrorHandling(client *vlock.Client) {
	fmt.Println("\n=== Error Handling Examples ===\n")

	// Example 1: Handle initialization errors
	err := client.Initialize()
	if err != nil {
		if voltageErr, ok := err.(*vlock.VoltageError); ok {
			fmt.Printf("Voltage Error Code: %d\n", voltageErr.Code)
			fmt.Printf("Error Message: %s\n", voltageErr.Message)
			fmt.Printf("Is Retryable: %v\n", voltageErr.IsRetryable())
			fmt.Printf("Error Category: %s\n", voltageErr.Category())

			// Retry logic for transient errors
			if voltageErr.IsRetryable() {
				fmt.Println("This error is retryable. Implementing retry logic...")
				// Implement exponential backoff retry
			}
		}
	}

	// Example 2: Handle health check errors
	if err := client.HealthCheck(); err != nil {
		fmt.Printf("Health check failed: %v\n", err)
		// Decide whether to continue or abort based on error type
	}
}

// Example of using functional options with NewClient
func demonstrateClientOptions() {
	fmt.Println("\n=== Client Options Examples ===\n")

	cfg := &config.Config{
		AppName:         "VoltageExampleApp",
		AppVersion:      "1.0.0",
		AppEnv:          "DEV",
		DEKSharedSecret: "example_shared_secret",
		ConfigFilePath:  "examples/voltage.cfg",
	}

	// Create client with functional options
	client, err := vlock.NewClient(cfg)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return
	}
	defer client.Close()

	// Initialize and use the client
	if err := client.Initialize(); err != nil {
		log.Printf("Failed to initialize: %v", err)
		return
	}

	fmt.Println("✓ Client created with custom options")
}
