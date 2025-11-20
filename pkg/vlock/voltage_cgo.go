package vlock

/*
#cgo CFLAGS: -I${SRCDIR}/include
#cgo LDFLAGS: -L${SRCDIR}/lib -lvoltage

#include <stdlib.h>
#include "voltage.h"

// Wrapper functions for C API
int voltage_go_init(const char* config_file, char** error_msg) {
    return voltage_init(config_file, error_msg);
}

int voltage_go_terminate(char** error_msg) {
    return voltage_terminate(error_msg);
}

int voltage_go_health_check(char** error_msg) {
    return voltage_health_check(error_msg);
}

void voltage_go_free_string(char* str) {
    if (str != NULL) {
        free(str);
    }
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// initializeVoltageLibrary initializes the Voltage C library with the configuration
// This replaces the placeholder implementation in voltage.go
func (c *Client) initializeVoltageLibrary() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.voltageInitialized {
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

	// Convert Go string to C string
	cConfigPath := C.CString(configPath)
	defer C.free(unsafe.Pointer(cConfigPath))

	// Error message pointer
	var cErrorMsg *C.char
	defer func() {
		if cErrorMsg != nil {
			C.voltage_go_free_string(cErrorMsg)
		}
	}()

	// Call C API
	result := C.voltage_go_init(cConfigPath, &cErrorMsg)

	// Convert C error code to Go error
	if result != 0 {
		errorMsg := "unknown initialization error"
		if cErrorMsg != nil {
			errorMsg = C.GoString(cErrorMsg)
		}
		return mapCErrorCode(int(result), errorMsg, nil)
	}

	c.voltageInitialized = true
	return nil
}

// terminateVoltageLibrary terminates the Voltage C library
// This replaces the placeholder implementation in voltage.go
func (c *Client) terminateVoltageLibrary() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.voltageInitialized {
		// Not an error - already terminated or never initialized
		return nil
	}

	// Error message pointer
	var cErrorMsg *C.char
	defer func() {
		if cErrorMsg != nil {
			C.voltage_go_free_string(cErrorMsg)
		}
	}()

	// Call C API
	result := C.voltage_go_terminate(&cErrorMsg)

	// Convert C error code to Go error
	if result != 0 {
		errorMsg := "unknown termination error"
		if cErrorMsg != nil {
			errorMsg = C.GoString(cErrorMsg)
		}
		return mapCErrorCode(int(result), errorMsg, nil)
	}

	c.voltageInitialized = false
	return nil
}

// performHealthCheckC performs a health check against the Voltage library
func (c *Client) performHealthCheckC() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.voltageInitialized {
		return ErrClientNotInitialized
	}

	// Error message pointer
	var cErrorMsg *C.char
	defer func() {
		if cErrorMsg != nil {
			C.voltage_go_free_string(cErrorMsg)
		}
	}()

	// Call C API
	result := C.voltage_go_health_check(&cErrorMsg)

	// Convert C error code to Go error
	if result != 0 {
		errorMsg := "health check failed"
		if cErrorMsg != nil {
			errorMsg = C.GoString(cErrorMsg)
		}
		return mapCErrorCode(int(result), errorMsg, nil)
	}

	return nil
}

// GetVoltageVersion returns the version of the Voltage C library
// This is a utility function for debugging and logging
func GetVoltageVersion() string {
	// This would call a C function to get the version
	// For now, return a placeholder
	// In production, you would call: C.voltage_get_version()
	return "1.0.0-placeholder"
}
