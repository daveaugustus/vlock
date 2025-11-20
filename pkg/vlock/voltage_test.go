package vlock

import (
	"testing"
	"time"

	"github.com/daveaugustus/vlock/pkg/config"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.Config
		shouldError bool
		errorMsg    string
	}{
		{
			name: "Valid configuration",
			config: &config.Config{
				AppName:         "TestApp",
				AppVersion:      "1.0.0",
				AppEnv:          "DEV",
				DEKSharedSecret: "test_secret",
				ConfigFilePath:  "test.cfg",
			},
			shouldError: false,
		},
		{
			name:        "Nil configuration",
			config:      nil,
			shouldError: true,
			errorMsg:    "config cannot be nil",
		},
		{
			name: "Invalid configuration - missing AppName",
			config: &config.Config{
				AppVersion:      "1.0.0",
				AppEnv:          "DEV",
				DEKSharedSecret: "test_secret",
				ConfigFilePath:  "test.cfg",
			},
			shouldError: true,
			errorMsg:    "invalid configuration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)

			if tt.shouldError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
				if client == nil {
					t.Error("Expected client, got nil")
				}
				if client.IsInitialized() {
					t.Error("Client should not be initialized yet")
				}
			}
		})
	}
}

func TestClientInitialize(t *testing.T) {
	cfg := &config.Config{
		AppName:         "TestApp",
		AppVersion:      "1.0.0",
		AppEnv:          "DEV",
		DEKSharedSecret: "test_secret",
		ConfigFilePath:  "test.cfg",
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test initialization
	err = client.Initialize()
	if err != nil {
		t.Errorf("Failed to initialize: %v", err)
	}

	if !client.IsInitialized() {
		t.Error("Client should be initialized")
	}

	if !client.IsHealthy() {
		t.Error("Client should be healthy after initialization")
	}

	// Test double initialization
	err = client.Initialize()
	if err == nil {
		t.Error("Expected error on double initialization")
	}

	// Clean up
	client.Close()
}

func TestClientClose(t *testing.T) {
	cfg := &config.Config{
		AppName:         "TestApp",
		AppVersion:      "1.0.0",
		AppEnv:          "DEV",
		DEKSharedSecret: "test_secret",
		ConfigFilePath:  "test.cfg",
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Initialize client
	if err := client.Initialize(); err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}

	// Close client
	if err := client.Close(); err != nil {
		t.Errorf("Failed to close: %v", err)
	}

	if client.IsInitialized() {
		t.Error("Client should not be initialized after close")
	}

	// Test double close (should not error)
	if err := client.Close(); err != nil {
		t.Errorf("Double close should not error: %v", err)
	}
}

func TestClientHealthCheck(t *testing.T) {
	cfg := &config.Config{
		AppName:         "TestApp",
		AppVersion:      "1.0.0",
		AppEnv:          "DEV",
		DEKSharedSecret: "test_secret",
		ConfigFilePath:  "test.cfg",
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Health check before initialization should fail
	err = client.HealthCheck()
	if err == nil {
		t.Error("Health check should fail before initialization")
	}

	// Initialize
	if err := client.Initialize(); err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}

	// Health check after initialization should succeed
	err = client.HealthCheck()
	if err != nil {
		t.Errorf("Health check failed: %v", err)
	}

	// Verify last health check time
	lastCheck := client.LastHealthCheck()
	if lastCheck.IsZero() {
		t.Error("LastHealthCheck should not be zero after health check")
	}
}

func TestClientReinitialize(t *testing.T) {
	cfg := &config.Config{
		AppName:         "TestApp",
		AppVersion:      "1.0.0",
		AppEnv:          "DEV",
		DEKSharedSecret: "test_secret",
		ConfigFilePath:  "test.cfg",
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Initialize
	if err := client.Initialize(); err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}

	// Reinitialize
	if err := client.Reinitialize(); err != nil {
		t.Errorf("Failed to reinitialize: %v", err)
	}

	if !client.IsInitialized() {
		t.Error("Client should be initialized after reinitialize")
	}
}

func TestClientInfo(t *testing.T) {
	cfg := &config.Config{
		AppName:         "TestApp",
		AppVersion:      "1.0.0",
		AppEnv:          "DEV",
		DEKSharedSecret: "test_secret",
		ConfigFilePath:  "test.cfg",
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	info := client.Info()

	if info.AppName != "TestApp" {
		t.Errorf("Expected AppName 'TestApp', got '%s'", info.AppName)
	}

	if info.AppVersion != "1.0.0" {
		t.Errorf("Expected AppVersion '1.0.0', got '%s'", info.AppVersion)
	}

	if info.Environment != "DEV" {
		t.Errorf("Expected Environment 'DEV', got '%s'", info.Environment)
	}

	if info.Initialized {
		t.Error("Client should not be initialized yet")
	}

	// Initialize and check again
	client.Initialize()
	info = client.Info()

	if !info.Initialized {
		t.Error("Client should be initialized")
	}

	if !info.Healthy {
		t.Error("Client should be healthy")
	}
}

func TestClientConfig(t *testing.T) {
	cfg := &config.Config{
		AppName:         "TestApp",
		AppVersion:      "1.0.0",
		AppEnv:          "DEV",
		DEKSharedSecret: "test_secret",
		ConfigFilePath:  "test.cfg",
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	retrievedCfg := client.Config()
	if retrievedCfg == nil {
		t.Error("Config() should return configuration")
	}

	if retrievedCfg.AppName != "TestApp" {
		t.Errorf("Expected AppName 'TestApp', got '%s'", retrievedCfg.AppName)
	}
}

func TestClientConcurrency(t *testing.T) {
	cfg := &config.Config{
		AppName:         "TestApp",
		AppVersion:      "1.0.0",
		AppEnv:          "DEV",
		DEKSharedSecret: "test_secret",
		ConfigFilePath:  "test.cfg",
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Initialize
	if err := client.Initialize(); err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}

	// Test concurrent reads
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			// These should all be safe concurrent operations
			_ = client.IsInitialized()
			_ = client.IsHealthy()
			_ = client.Config()
			_ = client.Info()
			_ = client.LastHealthCheck()
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestClientLifecycle(t *testing.T) {
	cfg := &config.Config{
		AppName:         "TestApp",
		AppVersion:      "1.0.0",
		AppEnv:          "DEV",
		DEKSharedSecret: "test_secret",
		ConfigFilePath:  "test.cfg",
	}

	// Create client
	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Verify initial state
	if client.IsInitialized() {
		t.Error("New client should not be initialized")
	}
	if client.IsHealthy() {
		t.Error("New client should not be healthy")
	}

	// Initialize
	if err := client.Initialize(); err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}

	// Verify initialized state
	if !client.IsInitialized() {
		t.Error("Client should be initialized")
	}
	if !client.IsHealthy() {
		t.Error("Client should be healthy")
	}

	// Perform health check
	if err := client.HealthCheck(); err != nil {
		t.Errorf("Health check failed: %v", err)
	}

	// Verify last health check time is recent
	lastCheck := client.LastHealthCheck()
	if time.Since(lastCheck) > time.Second {
		t.Error("LastHealthCheck should be recent")
	}

	// Close client
	if err := client.Close(); err != nil {
		t.Errorf("Failed to close: %v", err)
	}

	// Verify closed state
	if client.IsInitialized() {
		t.Error("Client should not be initialized after close")
	}
}
