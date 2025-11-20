# VLock Implementation Plan - Based on Extracted Documentation

## Discovered Resources

From the extracted zip file, we now have:

### 1. C API Headers
- `fiserv_api.h` - Main API functions
- `fiserv_config.h` - Configuration and initialization
- `fiserv_kek.h` - Key Encryption Key functions
- `fiserv_rotation.h` - Key rotation functions
- `fsvapi.h` - API definitions

### 2. Example Code
- `test_golang.go` - Working CGO example showing how to call C functions
- `example.txt` - Quick reference for API functions and cryptIds

### 3. Documentation
- Error code documentation (multiple versions)
- API quick start guide
- Example XML configuration files

### 4. Configuration Example
- `fiservprotector 20.cfg` - Real configuration file example

---

## Core C API Functions (from fiserv_api.h)

### Text Operations
```c
int fiserv_protect_text(char *output, int outputBufferSize, 
                        const char *input, const char *cryptId);

int fiserv_access_text(char *output, int outputBufferSize, 
                       const char *input, const char *cryptId);

int fiserv_access_masked(char *output, int outputBufferSize, 
                         const char *input, const char *cryptId);
```

### Binary Operations
```c
int fiserv_protect_binary(char *output, int *outputLength, 
                          const char *input, int inputLength, const char *cryptId);

int fiserv_access_binary(char *output, int *outputLength,
                         const char  *input, int inputLength, 
                         const char *cryptId);
```

### Error Handling
```c
const char *fiserv_get_last_error(void);
```

---

## Initialization Functions (from fiserv_config.h)

### Primary Initialization Methods
```c
// Initialize from config file path
int fiserv_init_protector(const char *fiserv_config_path);

// Initialize from key-value config + resource path
int fiserv_init_protector_key_value(const char *fiserv_config_path, 
                                     const char *resource_path);

// Initialize from environment variables
int fiserv_init_protector_env(const char *resource_path);

// Initialize from memory buffer
int fiserv_init_protector_mem(char *buf, const char *resource_path);
```

### Configuration Management
```c
int fiserv_set_config_path(const char *path);
int fiserv_set_config_xml(const char *auth_buf, const char *config_buf, 
                          const char *fiserv_buf);
int read_fiserv_config(const char *fname);
```

### Logging and Debugging
```c
int fiserv_set_log_file(const char *fname);
void fiserv_set_log_handler(fiserv_log_handler_t a_handler);
void fiserv_set_debug_fp(FILE *fp);
int fiserv_set_debug_file(const char *fname);
void fiserv_turn_off_debug(void);
```

### Version Information
```c
void fiserv_get_protector_version(char *buf, int size);
void fiserv_get_protector_version_short(char *buf, int size);
```

---

## Available CryptIds (from example.txt)

Standard cryptIds available:
- `Card_Internal` - Credit card encryption
- `SSN_Internal` - Social Security Number encryption
- `CustomText_Internal` - Generic text encryption
- `GenericText_Trailing4_Internal` - Text with last 4 visible (from test_golang.go)

---

## Implementation Roadmap

### Phase 1: Core Wrapper (Priority: HIGH)
✅ **DONE** - Configuration management
- [x] Config struct
- [x] File parsing
- [x] Environment variable support
- [x] Validation

🔲 **TODO** - CGO Bindings
- [ ] Create `voltage.go` with CGO directives
- [ ] Wrap `fiserv_init_protector()`
- [ ] Wrap `fiserv_protect_text()`
- [ ] Wrap `fiserv_access_text()`
- [ ] Wrap `fiserv_access_masked()`
- [ ] Wrap `fiserv_get_last_error()`

🔲 **TODO** - Error Handling
- [ ] Create Go error types
- [ ] Map C error codes to Go errors
- [ ] Parse error messages from `fiserv_get_last_error()`

### Phase 2: Binary Operations (Priority: MEDIUM)
🔲 **TODO** - Binary Support
- [ ] Wrap `fiserv_protect_binary()`
- [ ] Wrap `fiserv_access_binary()`
- [ ] Handle binary buffer management
- [ ] Add tests for binary operations

### Phase 3: Advanced Features (Priority: LOW)
🔲 **TODO** - Advanced Initialization
- [ ] Support `fiserv_init_protector_env()` for pure env var init
- [ ] Support `fiserv_init_protector_mem()` for in-memory config

🔲 **TODO** - Logging Integration
- [ ] Wrap logging functions
- [ ] Integrate with Go logging frameworks
- [ ] Add debug mode support

🔲 **TODO** - Version Information
- [ ] Add `GetVersion()` function
- [ ] Add version compatibility checks

🔲 **TODO** - Key Rotation (if fiserv_rotation.h has functions)
- [ ] Review key rotation API
- [ ] Implement rotation wrappers if needed

---

## Go Wrapper Design (Based on test_golang.go)

### Basic Structure
```go
package vlock

/*
#cgo LDFLAGS: -L/path/to/voltage/lib -lfiservprotect -lvibesimple
#include <stdlib.h>
extern int fiserv_init_protector(const char *fiserv_config_path);
extern int fiserv_protect_text(char *output, int outputBufferSize, 
                                const char *input, const char *cryptId);
extern int fiserv_access_text(char *output, int outputBufferSize,
                               const char *input, const char *cryptId);
extern int fiserv_access_masked(char *output, int outputBufferSize,
                                 const char *input, const char *cryptId);
extern const char *fiserv_get_last_error(void);
*/
import "C"

import (
    "fmt"
    "unsafe"
)

// Init initializes the Voltage protector with config file
func Init(configPath string) error {
    cPath := C.CString(configPath)
    defer C.free(unsafe.Pointer(cPath))
    
    result := C.fiserv_init_protector(cPath)
    if result != 0 {
        return fmt.Errorf("initialization failed: %s", GetLastError())
    }
    return nil
}

// Encrypt protects plaintext using specified cryptId
func Encrypt(plaintext, cryptId string) (string, error) {
    cInput := C.CString(plaintext)
    defer C.free(unsafe.Pointer(cInput))
    
    cCryptId := C.CString(cryptId)
    defer C.free(unsafe.Pointer(cCryptId))
    
    outputSize := 1024 // Adjust based on needs
    cOutput := (*C.char)(C.malloc(C.size_t(outputSize)))
    defer C.free(unsafe.Pointer(cOutput))
    
    result := C.fiserv_protect_text(cOutput, C.int(outputSize), cInput, cCryptId)
    if result != 0 {
        return "", fmt.Errorf("encryption failed: %s", GetLastError())
    }
    
    return C.GoString(cOutput), nil
}

// Decrypt accesses ciphertext using specified cryptId
func Decrypt(ciphertext, cryptId string) (string, error) {
    cInput := C.CString(ciphertext)
    defer C.free(unsafe.Pointer(cInput))
    
    cCryptId := C.CString(cryptId)
    defer C.free(unsafe.Pointer(cCryptId))
    
    outputSize := 1024
    cOutput := (*C.char)(C.malloc(C.size_t(outputSize)))
    defer C.free(unsafe.Pointer(cOutput))
    
    result := C.fiserv_access_text(cOutput, C.int(outputSize), cInput, cCryptId)
    if result != 0 {
        return "", fmt.Errorf("decryption failed: %s", GetLastError())
    }
    
    return C.GoString(cOutput), nil
}

// DecryptMasked returns masked version of ciphertext
func DecryptMasked(ciphertext, cryptId string) (string, error) {
    cInput := C.CString(ciphertext)
    defer C.free(unsafe.Pointer(cInput))
    
    cCryptId := C.CString(cryptId)
    defer C.free(unsafe.Pointer(cCryptId))
    
    outputSize := 1024
    cOutput := (*C.char)(C.malloc(C.size_t(outputSize)))
    defer C.free(unsafe.Pointer(cOutput))
    
    result := C.fiserv_access_masked(cOutput, C.int(outputSize), cInput, cCryptId)
    if result != 0 {
        return "", fmt.Errorf("masked access failed: %s", GetLastError())
    }
    
    return C.GoString(cOutput), nil
}

// GetLastError retrieves the last error message from C library
func GetLastError() string {
    cErr := C.fiserv_get_last_error()
    if cErr == nil {
        return "unknown error"
    }
    return C.GoString(cErr)
}

// Close performs cleanup (if needed)
func Close() error {
    // Add cleanup code if library provides a shutdown function
    return nil
}
```

---

## Key Reliability Requirements (from example.txt)

1. **Always initialize before use**
   - Call `Init()` before any encryption/decryption
   - Consider warmup calls for production

2. **Validate all inputs**
   - Don't depend on Voltage library for validation
   - Check for empty strings, null values, etc.

3. **Environment consistency**
   - ONLY decrypt ciphertext from the same environment
   - DEV → DEV, QA → QA, PROD → PROD
   - Cross-environment decryption will fail!

---

## Buffer Size Management

From test_golang.go, we see 1024 bytes is used as default. Consider:
- Dynamic sizing based on input length
- Error handling for "buffer too small" (Error 5500207)
- Retry logic with larger buffer if needed

---

## Next Steps

1. **Create `voltage.go`** - Main CGO wrapper file
2. **Create `voltage_test.go`** - Integration tests
3. **Update `go.mod`** - Add CGO build tags and requirements
4. **Create `Makefile`** - Build and test automation
5. **Document library paths** - Where to find libfiservprotect and libvibesimple
6. **Add examples/** - Usage examples for common scenarios

---

## Environment Setup Required

For building and testing, developers will need:
- Voltage SimpleAPI C library installed
- `libfiservprotect` library
- `libvibesimple` library
- CGO enabled: `export CGO_ENABLED=1`
- Proper library paths set in `#cgo LDFLAGS`

---

## Testing Strategy

1. **Unit Tests** - Test individual functions with mocked C calls (if possible)
2. **Integration Tests** - Real encryption/decryption with test config
3. **Environment Tests** - Verify env-specific behavior
4. **Error Tests** - Test error conditions and error messages
5. **Buffer Tests** - Test various buffer sizes and edge cases
6. **Concurrent Tests** - Test thread safety

---

## Documentation Needed

1. **Installation Guide** - How to install Voltage libraries
2. **Quick Start** - 5-minute getting started guide
3. **API Reference** - Complete function documentation
4. **Configuration Guide** - (Already done ✅)
5. **Error Handling** - Common errors and solutions
6. **Examples** - Real-world usage patterns
7. **Migration Guide** - For teams moving from direct C usage

---

This implementation plan provides a complete roadmap for building the VLock wrapper based on the actual C API and working example code!
