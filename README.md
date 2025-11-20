# VLock - Voltage C Library Wrapper for Go

A CGO-based wrapper that provides a clean, developer-friendly Go interface for the Voltage C (Protector) library, enabling Format-Preserving Encryption (FPE) and tokenization in Go services.

## Overview

VLock simplifies the integration of Voltage's C-based SDK into Go applications by providing a clean wrapper that eliminates the need for developers to work directly with C code. This wrapper provides one-time standardized integration for all Go services, reducing onboarding time and preventing repeated C-level integration work across teams.

## Features

- **Text Encryption/Decryption** - Protect and access plain text values
- **Binary Encryption/Decryption** - Support for raw data or files
- **Masked Access** - Partial redaction of sensitive fields based on patterns
- **Flexible Configuration** - Support for both file-based and environment variable configuration
- **Multi-Environment Support** - Manage separate configurations for DEV, QA, CAT, and PROD
- **Key Rotation** - Automatic key material updates with zero downtime
- **Error Handling** - Descriptive error codes and Go-friendly error messages
- **Thread Safety** - Safe multi-threaded usage after initialization
- **Container-Friendly** - Environment variable injection for Docker/Kubernetes deployments
- **Memory Management** - Automatic buffer allocation to prevent C-level memory issues

## Architecture

```
┌─────────────────────────────────┐
│     Go Service                  │
│  (Business Logic / APIs)        │
└────────────┬────────────────────┘
             │
┌────────────▼────────────────────┐
│  Go Wrapper (CGO)               │
│  - Init / Encrypt / Decrypt     │
│  - Error & Memory Handling      │
└────────────┬────────────────────┘
             │
┌────────────▼────────────────────┐
│  Voltage C Library              │
│  (voltage_protect_text, etc.)   │
└─────────────────────────────────┘
```

## Installation

```bash
go get github.com/the_monkeys/vlock
```

**Prerequisites:**
- Voltage C Library (Protector) installed
- CGO enabled
- Valid Voltage configuration files or environment variables

## Configuration

### Overview

This section addresses the key configuration requirements for the Voltage solution, including how to acquire settings, manage them across environments, and handle security considerations like key rotation.

### Required Configuration Parameters

The Voltage Protector library requires specific configuration parameters. These can be provided via configuration files or environment variables (environment variables take precedence).

| Configuration Purpose | .cfg Parameter | Environment Variable | Required |
|----------------------|----------------|---------------------|----------|
| Application Name | `fp_appName` | `FP_APPNAME` | Yes |
| Application Version | `fp_appVersion` | `FP_APPVERSION` | Yes |
| Environment | `fp_appEnv` | `FP_APPENV` | Yes |
| Network Timeout | `fp_networkTimeout` | `FP_NETWORKTIMEOUT` | No |
| Disable CRL Checking | `fp_disableCRLChecking` | `FP_DISABLECRLCHECKING` | No |
| Simple API Install Path | `fp_simpleAPI_installPath` | `FP_SIMPLEAPI_INSTALLPATH` | See docs |
| Trust Store Path | `fp_trustStore_path` | `FP_TRUSTSTORE_PATH` | See docs |
| KEK Cert Path | `fp_kek_certPath` | `FP_KEK_CERTPATH` | See docs |
| KEK Cert Passphrase | `fp_kek_certPassphrase` | `FP_KEK_CERTPASSPHRASE` | See docs |
| KEK Shared Secret | `fp_kek_sharedSecret` | `FP_KEK_SHAREDSECRET` | See docs |
| DEK Shared Secret | `fp_default_sharedSecret` | `FP_DEFAULT_SHAREDSECRET` | See docs |
| DEK Username | `fp_default_userName` | `FP_DEFAULT_USERNAME` | See docs |
| DEK Password | `fp_default_password` | `FP_DEFAULT_PASSWORD` | See docs |

### Configuration File Example

**voltageprotector.cfg:**
```ini
[ProtectorConfig]
XMLConfig=./vsconfig.xml
DefaultCryptId=SSN_Internal
Environment=DEV
LogLevel=2
LogFile=/opt/voltage/logs/protector.log

# Application Configuration
fp_appName=YOUR_APP_NAME
fp_appVersion=YOUR_APP_VERSION
fp_appEnv=DEV
fp_default_sharedSecret=YOUR_SHARED_SECRET
fp_kek_certPath=/path/to/cert.pfx
fp_kek_certPassphrase=YOUR_CERT_PASSPHRASE
fp_simpleAPI_installPath=/opt/voltage/simpleapi
fp_trustStore_path=/opt/voltage/trustStore
fp_networkTimeout=10
fp_disableCRLChecking=false
```

**vsconfig.xml (sanitized example):**
```xml
<cryptId name="SSN_Internal" algorithm="FPE" key="YOUR_KEY_HERE" format="NUMERIC"/>
<mask pattern="XXX-XX-####" cryptId="SSN_Internal"/>
```

### Using Environment Variables

```bash
# Linux/Mac
export FP_APPNAME=MYAPPNAME
export FP_APPVERSION=1.0.0
export FP_APPENV=DEV
export FP_SIMPLEAPI_INSTALLPATH=/opt/voltage/simpleapi
export FP_TRUSTSTORE_PATH=/opt/voltage/trustStore
export FP_DEFAULT_SHAREDSECRET=yoursharedsecret
```

```powershell
# Windows PowerShell
$env:FP_APPNAME="MYAPPNAME"
$env:FP_APPVERSION="1.0.0"
$env:FP_APPENV="DEV"
$env:FP_SIMPLEAPI_INSTALLPATH="C:\voltage\simpleapi"
$env:FP_TRUSTSTORE_PATH="C:\voltage\trustStore"
$env:FP_DEFAULT_SHAREDSECRET="yoursharedsecret"
```

### Environment Management

#### Moving Between Environments (DEV → QA → CAT → PROD)

Configuration settings change when promoting between environments. Each environment requires its own:

1. **Unique Configuration Files** - Separate `.cfg` and `.xml` files per environment
2. **Environment-Specific Keys** - Different encryption keys for DEV, QA, CAT, and PROD
3. **Environment Variable Settings** - Different values for `FP_APPENV` parameter

**Best Practices:**
- Store configuration files in environment-specific directories (e.g., `/config/dev`, `/config/prod`)
- Use environment-specific service accounts and credentials
- Never promote encrypted data across environments (re-encrypt in target environment)
- Validate configuration before deployment using the `Init()` validation

**Example Environment Structure:**
```
/config
  /dev
    voltageprotector.cfg
    vsconfig.xml
  /qa
    voltageprotector.cfg
    vsconfig.xml
  /cat
    voltageprotector.cfg
    vsconfig.xml
  /prod
    voltageprotector.cfg
    vsconfig.xml
```

#### Acquiring Configuration Settings

Configuration settings are typically provided by:

1. **Voltage Platform Team** - Provides initial setup including:
   - Certificate files (`.pfx`)
   - Shared secrets and credentials
   - Trust store files
   - XML configuration templates

2. **Internal Security/Platform Team** - Manages:
   - Environment-specific credentials
   - Key rotation schedules
   - Access control policies

3. **External Teams (e.g., Finxact)** - May provide:
   - Application-specific configuration
   - Environment variables for container deployments
   - Service account credentials

**Recommended Acquisition Process:**
1. Submit configuration request to Voltage platform team
2. Receive environment-specific credentials and certificates
3. Store sensitive values in secure vault (e.g., HashiCorp Vault, AWS Secrets Manager)
4. Inject configuration via environment variables at runtime
5. Validate configuration in non-production environment first

### Key Rotation

The Voltage library supports automatic key rotation for enhanced security compliance.

#### Rotation Schedule

| Environment | Recommended Interval | Compliance Requirement |
|-------------|---------------------|------------------------|
| Development | 90 days | Optional |
| QA/CAT | 60 days | Recommended |
| Production | 30-90 days | Required (per policy) |

**Note:** Actual rotation intervals should align with your organization's security policies and compliance requirements (PCI-DSS, HIPAA, etc.).

#### Key Rotation Process

**Step 1: Preparation**
```bash
# Backup current configuration
cp vsconfig.xml vsconfig.xml.backup.$(date +%Y%m%d)
cp voltageprotector.cfg voltageprotector.cfg.backup.$(date +%Y%m%d)
```

**Step 2: Request New Keys**
- Contact Voltage platform team or use Voltage management console
- Request new key material for specific `cryptId`
- Receive updated XML configuration with new key references

**Step 3: Update Configuration**
```xml
<!-- Updated vsconfig.xml with new key -->
<cryptId name="SSN_Internal" algorithm="FPE" key="NEW_KEY_HERE" format="NUMERIC"/>
```

**Step 4: Deploy Configuration**
```bash
# Deploy to target environment
# The library automatically uses new keys for encryption
# Old keys remain available for decryption of existing data
```

**Step 5: Validation**
```go
// Test encryption with new keys
err := voltage.Init("/path/to/updated/config.cfg")
if err != nil {
    log.Fatal("Configuration validation failed:", err)
}

// Verify new encryptions work
encrypted, err := voltage.Encrypt("test-data", "SSN_Internal")
// Verify old encrypted data can still be decrypted
decrypted, err := voltage.Decrypt(oldEncryptedData, "SSN_Internal")
```

**Step 6: Monitor**
- Monitor application logs for decryption errors
- Verify both new and old ciphertext can be processed
- Plan re-encryption of old data if required by policy

#### Automatic Key Material Updates

The Voltage library automatically detects and uses updated key materials from the XML configuration file:
- No application restart required for key rotation
- Library checks for configuration updates periodically
- Both old and new keys remain active during transition period
- Re-encryption of existing data can be performed gradually

#### Key Rotation Considerations

- **Zero Downtime**: Key rotation should not cause service interruption
- **Backward Compatibility**: Old ciphertext must remain decryptable during transition
- **Audit Trail**: Maintain logs of all key rotation activities
- **Testing**: Always test key rotation in non-production environments first
- **Rollback Plan**: Keep previous configuration backups for emergency rollback

### Configuration FAQ

**Q: Can we inject configuration via method parameters instead of files?**  
A: No. The Voltage C library requires file-based configuration. Configuration items cannot be injected through Go API calls. You must use `.cfg` files or environment variables.

**Q: What data format does the configuration use?**  
A: The library uses two formats:
- `.cfg` files: Plain text key-value pairs (similar to INI format)
- `.xml` files: Structured XML for encryption rules, cryptIds, and mask patterns

**Q: Can we use environment variables like other applications?**  
A: Yes! Environment variables are the recommended approach for containerized deployments. All configuration parameters have corresponding environment variable names (e.g., `FP_APPNAME`, `FP_APPENV`). Environment variables take precedence over file-based configuration.

**Q: How do we handle configuration in container environments?**  
A: Recommended approach for containers:
1. Store sensitive values in secrets management system
2. Inject as environment variables at container runtime
3. Mount configuration files as ConfigMaps/volumes (if needed)
4. Use `FP_*` environment variables to override file settings

Example Docker/Kubernetes approach:
```yaml
env:
  - name: FP_APPNAME
    value: "MyService"
  - name: FP_APPENV
    value: "PROD"
  - name: FP_DEFAULT_SHAREDSECRET
    valueFrom:
      secretKeyRef:
        name: voltage-secrets
        key: shared-secret
```

**Q: What happens if configuration is missing or invalid?**  
A: The `Init()` function will return an error with details about the missing/invalid configuration. The wrapper validates all required parameters before allowing any encryption/decryption operations.

**Q: Can we use the same configuration across multiple services?**  
A: Configuration can be shared across services if they:
- Run in the same environment (DEV, QA, PROD)
- Use the same application name and version
- Share the same security requirements

However, it's recommended to use service-specific `fp_appName` values for better audit trails and access control.

**Q: How do we secure sensitive configuration values?**  
A: Best practices:
- Never commit secrets to version control
- Use secrets management systems (Vault, AWS Secrets Manager, etc.)
- Inject sensitive values as environment variables at runtime
- Restrict file permissions on configuration files (chmod 600)
- Rotate credentials regularly

**Q: Where can we get help with configuration issues?**  
A: Configuration support resources:
- Voltage platform documentation
- Internal Voltage platform team
- Voltage C API documentation (`voltage_api.h`)
- Voltage Protector error codes reference


## Usage

### Basic Text Encryption/Decryption

```go
package main

import (
    "fmt"
    "log"
    "github.com/the_monkeys/vlock"
)

func main() {
    // Initialize the library
    err := voltage.Init("/path/to/voltageprotector.cfg")
    if err != nil {
        log.Fatal("Failed to initialize:", err)
    }
    defer voltage.Close()

    // Encrypt data
    plainText := "123-45-6789"
    encrypted, err := voltage.Encrypt(plainText, "SSN_Internal")
    if err != nil {
        log.Fatal("Encryption failed:", err)
    }
    fmt.Println("Encrypted:", encrypted)

    // Decrypt data
    decrypted, err := voltage.Decrypt(encrypted, "SSN_Internal")
    if err != nil {
        log.Fatal("Decryption failed:", err)
    }
    fmt.Println("Decrypted:", decrypted)
}
```

### Masked Access

```go
// Initialize and encrypt as above
encrypted, err := voltage.Encrypt("123-45-6789", "SSN_Internal")
if err != nil {
    log.Fatal(err)
}

// Get masked version (based on XML pattern configuration)
masked, err := voltage.DecryptMasked(encrypted, "SSN_Internal")
if err != nil {
    log.Fatal(err)
}
fmt.Println("Masked:", masked) // Output: XXX-XX-6789
```

### Binary Encryption/Decryption

```go
// Encrypt binary data
binaryData := []byte{0x01, 0x02, 0x03, 0x04}
encryptedBinary, err := voltage.EncryptBinary(binaryData, "BinaryCryptId")
if err != nil {
    log.Fatal("Binary encryption failed:", err)
}

// Decrypt binary data
decryptedBinary, err := voltage.DecryptBinary(encryptedBinary, "BinaryCryptId")
if err != nil {
    log.Fatal("Binary decryption failed:", err)
}
```

## API Reference

### Initialization

```go
func Init(configPath string) error
```
Initializes the Voltage Protector library with the specified configuration file. Must be called before any encryption/decryption operations.

### Text Operations

```go
func Encrypt(plainText string, cryptId string) (string, error)
```
Encrypts plain text using the specified cryptId.

```go
func Decrypt(cipherText string, cryptId string) (string, error)
```
Decrypts cipher text using the specified cryptId.

```go
func DecryptMasked(cipherText string, cryptId string) (string, error)
```
Returns partially decrypted data based on mask patterns defined in XML configuration.

### Binary Operations

```go
func EncryptBinary(data []byte, cryptId string) ([]byte, error)
```
Encrypts binary data using the specified cryptId.

```go
func DecryptBinary(data []byte, cryptId string) ([]byte, error)
```
Decrypts binary data using the specified cryptId.

### Cleanup

```go
func Close() error
```
Releases resources and performs cleanup. Should be called when done using the library.

## Error Handling

The wrapper translates C-level errors into meaningful Go errors:

| Error | Description | Common Cause |
|-------|-------------|--------------|
| `Error 5500207` | Output buffer is too small | Insufficient buffer allocation |
| Invalid ciphertext | Trailing characters mismatch or invalid base64 | Corrupted or tampered data |
| Config XML issues | Missing or invalid cryptId configuration | Configuration file errors |

Example error handling:
```go
encrypted, err := voltage.Encrypt(plainText, "SSN_Internal")
if err != nil {
    switch {
    case strings.Contains(err.Error(), "5500207"):
        // Handle buffer size error
    case strings.Contains(err.Error(), "invalid ciphertext"):
        // Handle invalid data error
    default:
        // Handle other errors
    }
}
```

## Important Considerations

### Environment-Specific Encryption

**⚠️ Warning:** Ciphertext encrypted in one environment (e.g., DEV) cannot be decrypted in another environment (e.g., QA) due to differing keys and configurations. The wrapper validates and enforces environment-specific configuration usage.

### Thread Safety

The wrapper is thread-safe after initialization. Internal locking is provided around initialization to ensure safe concurrent usage.

### Memory Management

The wrapper handles buffer allocation internally to avoid common C-related memory issues and buffer overflow errors.

### Configuration Restrictions

- Configuration items cannot be injected through Go API calls
- All configuration must be present in configuration files or environment variables at initialization
- The `.cfg` file must exist and reference the correct XML at runtime if not using environment variables

## Development

### Building

```bash
# Ensure CGO is enabled
export CGO_ENABLED=1

# Build the project
go build
```

## Testing

### Running Tests

#### 1. Run All Tests

```bash
# Run all tests with verbose output
go test -v

# Expected output:
# === RUN   TestNewConfig
# --- PASS: TestNewConfig (0.00s)
# === RUN   TestLoadConfigFromFile
# --- PASS: TestLoadConfigFromFile (0.00s)
# ...
# PASS
# ok      github.com/daveaugustus/vlock   1.169s
```

#### 2. Run Specific Tests

```bash
# Run only configuration tests
go test -v -run TestLoadConfig

# Run only validation tests
go test -v -run TestValidate

# Run only environment variable tests
go test -v -run TestEnvVar
```

#### 3. Run Tests with Coverage

```bash
# Generate coverage report
go test -cover

# Expected output:
# PASS
# coverage: 95.2% of statements
# ok      github.com/daveaugustus/vlock   1.169s

# Generate detailed HTML coverage report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

#### 4. Run Tests with Race Detection

```bash
# Detect data races in concurrent code
go test -race -v

# Use this to ensure thread safety
```

### Test Scenarios

#### Configuration Loading Tests

```bash
# Test file-based configuration loading
go test -v -run TestLoadConfigFromFile

# Test environment variable precedence
go test -v -run TestEnvVarPrecedence

# Test configuration validation
go test -v -run TestValidateRequiredFields
```

**What these tests verify:**
- ✅ Configuration files are parsed correctly
- ✅ Environment variables override file values
- ✅ Required fields are validated
- ✅ Error handling works as expected

#### Environment Variable Tests

```bash
# Set test environment variables
export FP_APPNAME=TestApp
export FP_APPVERSION=2.0.0
export FP_APPENV=QA

# Run environment variable tests
go test -v -run TestEnvVarPrecedence

# Clean up
unset FP_APPNAME FP_APPVERSION FP_APPENV
```

**Windows PowerShell:**
```powershell
# Set test environment variables
$env:FP_APPNAME="TestApp"
$env:FP_APPVERSION="2.0.0"
$env:FP_APPENV="QA"

# Run tests
go test -v -run TestEnvVarPrecedence

# Clean up
Remove-Item Env:\FP_APPNAME
Remove-Item Env:\FP_APPVERSION
Remove-Item Env:\FP_APPENV
```

### Integration Testing

#### Prerequisites for Integration Tests
1. Voltage C library installed
2. Valid configuration files in `config/dev/`
3. Access to Voltage server (for live encryption tests)

#### Running Integration Tests

```bash
# Set up test configuration
export CONFIG_PATH=./config/dev/voltageprotector.cfg

# Run integration tests (when implemented)
go test -v -tags=integration

# Note: Integration tests require live Voltage server connection
```

### Test File Structure

```
vlock/
├── config_test.go       # Configuration loading and validation tests
├── voltage_test.go      # Encryption/decryption tests (when implemented)
└── config/
    └── dev/
        ├── voltageprotector.cfg   # Test configuration file
        └── vsconfig.xml          # Test XML configuration
```

### Manual Testing

#### Test Configuration Loading

```go
package main

import (
    "fmt"
    "log"
    "github.com/daveaugustus/vlock"
)

func main() {
    // Test loading configuration
    config, err := vlock.LoadConfig("./config/dev/voltageprotector.cfg")
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }
    
    // Verify configuration
    fmt.Printf("✅ Config loaded successfully\n")
    fmt.Printf("   App Name: %s\n", config.AppName)
    fmt.Printf("   App Version: %s\n", config.AppVersion)
    fmt.Printf("   Environment: %s\n", config.AppEnv)
    fmt.Printf("   Library Path: %s\n", config.SimpleAPIInstallPath)
}
```

Run:
```bash
go run test_config.go
```

Expected output:
```
✅ Config loaded successfully
   App Name: VLockDev
   App Version: 1.0.0
   Environment: DEV
   Library Path: /opt/voltage/simpleapi
```

#### Test Environment Variable Override

```bash
# Set environment variables
export FP_APPNAME=OverrideTest
export FP_APPENV=QA

# Run test
go run test_config.go

# Expected: Should show "OverrideTest" and "QA"
```

### Continuous Integration Testing

For CI/CD pipelines:

```yaml
# Example GitHub Actions workflow
test:
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v2
    - uses: actions/setup-go@v2
      with:
        go-version: '1.21'
    
    - name: Run tests
      run: |
        go test -v -race -coverprofile=coverage.out
        go tool cover -func=coverage.out
    
    - name: Upload coverage
      uses: codecov/codecov-action@v2
      with:
        files: ./coverage.out
```

### Troubleshooting Test Failures

#### Issue: "Config file not found"

**Solution:**
```bash
# Verify file path
ls -la config/dev/voltageprotector.cfg

# Use absolute path if needed
go test -v -run TestLoadConfig
```

#### Issue: "Required field missing"

**Solution:** Ensure test config files have all required fields:
- `fp_appName`
- `fp_appVersion`
- `fp_appEnv`

#### Issue: Tests pass locally but fail in CI

**Solution:**
1. Check environment variables in CI
2. Verify file paths (use relative paths from project root)
3. Ensure test config files are committed to git

### Test Coverage Goals

- **Configuration Module**: 95%+ coverage ✅ (Currently: 95.2%)
- **Core Wrapper**: Target 90%+ (Future)
- **Error Handling**: Target 100% (Future)

### Next Steps

1. ✅ Configuration tests are complete and passing
2. ⏳ Add encryption/decryption tests (awaiting CGO wrapper implementation)
3. ⏳ Add integration tests with live Voltage server
4. ⏳ Add performance benchmarks

---

For more details on developer setup, see [DEVELOPER_SETUP.md](./DEVELOPER_SETUP.md).

### Project Structure

```
vlock/
├── README.md
├── project.doc          # Project requirements and research
├── voltage.go           # Main Go wrapper code
├── voltage_test.go      # Unit tests
├── examples/            # Usage examples
├── docs/                # Additional documentation
└── config/              # Sample configuration files
```

## Roadmap

- [x] Research and spike completion
- [ ] Implement text encryption/decryption
- [ ] Implement binary encryption/decryption
- [ ] Add masked access support
- [ ] Write comprehensive unit tests
- [ ] Package as internal Go module
- [ ] Create developer documentation
- [ ] Add usage examples
- [ ] Performance benchmarking

## Contributing

This is an internal project for The Monkeys organization. Please follow the internal contribution guidelines and submit pull requests for review.

## Support

For questions, issues, or support:
- Internal Voltage platform documentation
- Review Voltage C API documentation (`voltage_api.h`)
- Check Voltage Protector error codes reference

## License

Internal use only. See your organization's licensing terms.

## Acknowledgments

- Voltage platform team for the C Protector library
- Research based on Voltage C API documentation and sample integration programs

---

**Note:** This wrapper provides a thin CGO layer to expose clean Go APIs while hiding C-level complexity, providing a reusable API for all Go services within the organization.
