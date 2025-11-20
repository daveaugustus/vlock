package vlock

import (
	"fmt"
	"sync"
	"time"

	"github.com/daveaugustus/vlock/pkg/config"
)

// Client represents a Voltage encryption client
// Provides methods for initializing and managing connections to the Voltage service
type Client struct {
	config *config.Config

	// Connection state
	initialized bool
	mu          sync.RWMutex

	// Health monitoring
	lastHealthCheck time.Time
	healthy         bool

	// Session management
	sessionID string
}

// ClientOption is a functional option for configuring the Client
type ClientOption func(*Client) error

// NewClient creates a new Voltage client with the given configuration
// This follows the initialization pattern where configuration is loaded first,
// then the client is created, and finally Initialize() is called to connect
//
// Example:
//
//	cfg, err := config.LoadConfig("./voltageprotector.cfg")
//	if err != nil {
//	    return err
//	}
//
//	client, err := vlock.NewClient(cfg)
//	if err != nil {
//	    return err
//	}
//	defer client.Close()
//
//	if err := client.Initialize(); err != nil {
//	    return err
//	}
func NewClient(cfg *config.Config, opts ...ClientOption) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	client := &Client{
		config:      cfg,
		initialized: false,
		healthy:     false,
	}

	// Apply functional options
	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, fmt.Errorf("failed to apply client option: %w", err)
		}
	}

	return client, nil
}

// Initialize establishes connection to the Voltage service and performs initial setup
// This method must be called before using any encryption/decryption functions
// It initializes the Voltage C library and verifies connectivity
func (c *Client) Initialize() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.initialized {
		return fmt.Errorf("client already initialized")
	}

	// TODO: Initialize Voltage C library via CGO
	// This will call voltage_init() from the C library
	// For now, we'll simulate initialization

	if err := c.initializeVoltageLibrary(); err != nil {
		return fmt.Errorf("failed to initialize Voltage library: %w", err)
	}

	// Perform health check
	if err := c.performHealthCheck(); err != nil {
		return fmt.Errorf("health check failed after initialization: %w", err)
	}

	c.initialized = true
	c.healthy = true
	c.lastHealthCheck = time.Now()

	return nil
}

// initializeVoltageLibrary initializes the Voltage C library
// Implementation is provided by either voltage_cgo.go (with CGO) or voltage_mock.go (without CGO)

// performHealthCheck verifies the Voltage service is accessible
func (c *Client) performHealthCheck() error {
	if c.config == nil {
		return fmt.Errorf("configuration not loaded")
	}

	// Call the C library health check
	return c.performHealthCheckC()
}

// Close gracefully shuts down the Voltage client
// This should be called when the client is no longer needed
// It terminates the Voltage C library connection and cleans up resources
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.initialized {
		return nil // Already closed or never initialized
	}

	// TODO: Terminate Voltage C library via CGO
	// This will call voltage_terminate()

	if err := c.terminateVoltageLibrary(); err != nil {
		return fmt.Errorf("failed to terminate Voltage library: %w", err)
	}

	c.initialized = false
	c.healthy = false

	return nil
}

// terminateVoltageLibrary terminates the Voltage C library
// Implementation is provided by either voltage_cgo.go (with CGO) or voltage_mock.go (without CGO)

// IsInitialized returns whether the client has been initialized
func (c *Client) IsInitialized() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.initialized
}

// IsHealthy returns whether the client is healthy and ready to process requests
func (c *Client) IsHealthy() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.healthy
}

// Config returns the client's configuration (read-only)
func (c *Client) Config() *config.Config {
	return c.config
}

// HealthCheck performs an on-demand health check
// Returns an error if the service is not healthy
func (c *Client) HealthCheck() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.initialized {
		return fmt.Errorf("client not initialized")
	}

	if err := c.performHealthCheck(); err != nil {
		c.healthy = false
		return err
	}

	c.healthy = true
	c.lastHealthCheck = time.Now()

	return nil
}

// LastHealthCheck returns the time of the last health check
func (c *Client) LastHealthCheck() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.lastHealthCheck
}

// Reinitialize attempts to reinitialize the client if initialization fails or connection is lost
// This is useful for recovery scenarios
func (c *Client) Reinitialize() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.initialized {
		// Close existing connection
		if err := c.terminateVoltageLibrary(); err != nil {
			return fmt.Errorf("failed to terminate before reinitialize: %w", err)
		}
		c.initialized = false
	}

	// Reinitialize
	if err := c.initializeVoltageLibrary(); err != nil {
		return fmt.Errorf("failed to reinitialize Voltage library: %w", err)
	}

	c.initialized = true
	c.healthy = true
	c.lastHealthCheck = time.Now()

	return nil
}

// GetSessionID returns the current session ID if available
func (c *Client) GetSessionID() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.sessionID
}

// ClientInfo returns information about the client state
type ClientInfo struct {
	AppName         string
	AppVersion      string
	Environment     string
	Initialized     bool
	Healthy         bool
	LastHealthCheck time.Time
	SessionID       string
}

// Info returns current client information
func (c *Client) Info() ClientInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return ClientInfo{
		AppName:         c.config.AppName,
		AppVersion:      c.config.AppVersion,
		Environment:     c.config.AppEnv,
		Initialized:     c.initialized,
		Healthy:         c.healthy,
		LastHealthCheck: c.lastHealthCheck,
		SessionID:       c.sessionID,
	}
}
