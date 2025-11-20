package vlock

import (
	"fmt"
)

// ErrorCode represents Voltage C library error codes
type ErrorCode int

// Voltage C Library Error Codes
// Based on Voltage C Protector error code documentation
const (
	ErrSuccess              ErrorCode = 0
	ErrInvalidParameter     ErrorCode = 1
	ErrMemoryAllocation     ErrorCode = 2
	ErrConfigNotFound       ErrorCode = 3
	ErrConfigInvalid        ErrorCode = 4
	ErrInitializationFailed ErrorCode = 5
	ErrNotInitialized       ErrorCode = 6
	ErrAlreadyInitialized   ErrorCode = 7
	ErrConnectionFailed     ErrorCode = 8
	ErrAuthenticationFailed ErrorCode = 9
	ErrCryptIDNotFound      ErrorCode = 10
	ErrEncryptionFailed     ErrorCode = 11
	ErrDecryptionFailed     ErrorCode = 12
	ErrInvalidData          ErrorCode = 13
	ErrBufferTooSmall       ErrorCode = 14
	ErrNetworkTimeout       ErrorCode = 15
	ErrCertificateError     ErrorCode = 16
	ErrKeyNotFound          ErrorCode = 17
	ErrPermissionDenied     ErrorCode = 18
	ErrServiceUnavailable   ErrorCode = 19
	ErrUnknown              ErrorCode = 999
)

// VoltageError represents an error from the Voltage library
type VoltageError struct {
	Code    ErrorCode
	Message string
	Detail  string
	CError  int // Original C error code
}

// Error implements the error interface
func (e *VoltageError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("voltage error [%d]: %s (%s)", e.Code, e.Message, e.Detail)
	}
	return fmt.Sprintf("voltage error [%d]: %s", e.Code, e.Message)
}

// Is checks if the error is of a specific type
func (e *VoltageError) Is(target error) bool {
	t, ok := target.(*VoltageError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// IsRetryable returns true if the error is transient and the operation can be retried
func (e *VoltageError) IsRetryable() bool {
	switch e.Code {
	case ErrNetworkTimeout, ErrConnectionFailed, ErrServiceUnavailable:
		return true
	default:
		return false
	}
}

// NewVoltageError creates a new VoltageError from a C error code
func NewVoltageError(cErrorCode int, detail string) *VoltageError {
	code, message := mapCErrorCode(cErrorCode)
	return &VoltageError{
		Code:    code,
		Message: message,
		Detail:  detail,
		CError:  cErrorCode,
	}
}

// mapCErrorCode maps C library error codes to Go error codes with messages
func mapCErrorCode(cErrorCode int) (ErrorCode, string) {
	switch cErrorCode {
	case 0:
		return ErrSuccess, "success"
	case 1:
		return ErrInvalidParameter, "invalid parameter provided"
	case 2:
		return ErrMemoryAllocation, "memory allocation failed"
	case 3:
		return ErrConfigNotFound, "configuration file not found"
	case 4:
		return ErrConfigInvalid, "configuration is invalid"
	case 5:
		return ErrInitializationFailed, "initialization failed"
	case 6:
		return ErrNotInitialized, "library not initialized"
	case 7:
		return ErrAlreadyInitialized, "library already initialized"
	case 8:
		return ErrConnectionFailed, "failed to connect to Voltage service"
	case 9:
		return ErrAuthenticationFailed, "authentication failed"
	case 10:
		return ErrCryptIDNotFound, "crypt ID not found"
	case 11:
		return ErrEncryptionFailed, "encryption operation failed"
	case 12:
		return ErrDecryptionFailed, "decryption operation failed"
	case 13:
		return ErrInvalidData, "invalid data format"
	case 14:
		return ErrBufferTooSmall, "output buffer too small"
	case 15:
		return ErrNetworkTimeout, "network operation timed out"
	case 16:
		return ErrCertificateError, "certificate validation error"
	case 17:
		return ErrKeyNotFound, "encryption key not found"
	case 18:
		return ErrPermissionDenied, "permission denied"
	case 19:
		return ErrServiceUnavailable, "Voltage service unavailable"
	default:
		return ErrUnknown, fmt.Sprintf("unknown error (code: %d)", cErrorCode)
	}
}

// Predefined errors for common scenarios
var (
	ErrClientNotInitialized = &VoltageError{
		Code:    ErrNotInitialized,
		Message: "client not initialized",
		Detail:  "call Initialize() before using client methods",
	}

	ErrClientAlreadyInitialized = &VoltageError{
		Code:    ErrAlreadyInitialized,
		Message: "client already initialized",
		Detail:  "call Close() before reinitializing",
	}

	ErrInvalidConfig = &VoltageError{
		Code:    ErrConfigInvalid,
		Message: "invalid configuration",
		Detail:  "check configuration file and required fields",
	}

	ErrNilConfig = &VoltageError{
		Code:    ErrInvalidParameter,
		Message: "configuration cannot be nil",
		Detail:  "provide a valid configuration object",
	}
)

// WrapError wraps a Go error as a VoltageError if it isn't already one
func WrapError(err error, code ErrorCode, detail string) error {
	if err == nil {
		return nil
	}

	// If already a VoltageError, return as is
	if _, ok := err.(*VoltageError); ok {
		return err
	}

	return &VoltageError{
		Code:    code,
		Message: err.Error(),
		Detail:  detail,
	}
}

// ErrorCategory represents a high-level category of errors
type ErrorCategory int

const (
	CategoryConfiguration ErrorCategory = iota
	CategoryInitialization
	CategoryConnection
	CategoryAuthentication
	CategoryEncryption
	CategoryDecryption
	CategoryNetwork
	CategoryCertificate
	CategoryUnknown
)

// Category returns the high-level category of the error
func (e *VoltageError) Category() ErrorCategory {
	switch e.Code {
	case ErrConfigNotFound, ErrConfigInvalid, ErrInvalidParameter:
		return CategoryConfiguration
	case ErrInitializationFailed, ErrNotInitialized, ErrAlreadyInitialized:
		return CategoryInitialization
	case ErrConnectionFailed, ErrServiceUnavailable:
		return CategoryConnection
	case ErrAuthenticationFailed, ErrPermissionDenied:
		return CategoryAuthentication
	case ErrEncryptionFailed, ErrCryptIDNotFound, ErrKeyNotFound:
		return CategoryEncryption
	case ErrDecryptionFailed:
		return CategoryDecryption
	case ErrNetworkTimeout:
		return CategoryNetwork
	case ErrCertificateError:
		return CategoryCertificate
	default:
		return CategoryUnknown
	}
}

// String returns the string representation of the error category
func (c ErrorCategory) String() string {
	switch c {
	case CategoryConfiguration:
		return "Configuration"
	case CategoryInitialization:
		return "Initialization"
	case CategoryConnection:
		return "Connection"
	case CategoryAuthentication:
		return "Authentication"
	case CategoryEncryption:
		return "Encryption"
	case CategoryDecryption:
		return "Decryption"
	case CategoryNetwork:
		return "Network"
	case CategoryCertificate:
		return "Certificate"
	default:
		return "Unknown"
	}
}
