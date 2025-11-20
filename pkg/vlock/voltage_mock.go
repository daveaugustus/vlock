//go:build !cgo
// +build !cgo

package vlock

import (
	"fmt"
	"sync"
)

// Mock voltage library state
var (
	mockInitialized bool
	mockMutex       sync.Mutex
	mockConfigPath  string
)

// initializeVoltageLibrary is a mock implementation for systems without CGO
// This replaces the placeholder implementation in voltage.go
func (c *Client) initializeVoltageLibrary() error {
	mockMutex.Lock()
	defer mockMutex.Unlock()

	if mockInitialized {
		return ErrClientAlreadyInitialized
	}

	// Get the configuration file path
	configPath := c.config.ConfigFilePath
	if configPath == "" {
		// Try to construct config file path from XML config path
		if c.config.XMLConfigPath != "" {
			configPath = c.config.XMLConfigPath
		} else {
			return fmt.Errorf("no configuration file path specified")
		}
	}

	// Mock initialization - just store the config path
	mockConfigPath = configPath
	mockInitialized = true

	return nil
}

// terminateVoltageLibrary is a mock implementation for systems without CGO
// This replaces the placeholder implementation in voltage.go
func (c *Client) terminateVoltageLibrary() error {
	mockMutex.Lock()
	defer mockMutex.Unlock()

	if !mockInitialized {
		// Not an error - already terminated or never initialized
		return nil
	}

	// Mock termination
	mockInitialized = false
	mockConfigPath = ""

	return nil
}

// performHealthCheckC is a mock implementation for systems without CGO
func (c *Client) performHealthCheckC() error {
	mockMutex.Lock()
	defer mockMutex.Unlock()

	if !mockInitialized {
		return ErrClientNotInitialized
	}

	// Mock health check - always healthy in mock mode
	return nil
}

// GetVoltageVersion returns the version of the mock Voltage library
func GetVoltageVersion() string {
	return "1.0.0-mock-nocgo"
}

// IsMockMode returns true if running in mock mode (no CGO)
func IsMockMode() bool {
	return true
}
