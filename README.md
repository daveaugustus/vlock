# Voltage Go Wrapper# Voltage Go Wrapper# VLock - Voltage C Library Wrapper for Go



A Go wrapper for the Fiserv Voltage SecureData Format-Preserving Encryption (FPE) library, providing idiomatic Go APIs for secure data encryption and protection.



## OverviewA Go wrapper for the Fiserv Voltage SecureData Format-Preserving Encryption (FPE) library, providing idiomatic Go APIs for secure data encryption and protection.A CGO-based wrapper that provides a clean, developer-friendly Go interface for the Voltage C (Protector) library, enabling Format-Preserving Encryption (FPE) and tokenization in Go services.



The Voltage Go Wrapper provides a clean, type-safe interface to the Voltage SecureData C library, following Go best practices and patterns. This implementation focuses on **configuration and initialization** as specified in ARCH-16835.



### Key Features## Overview## Overview



- ✅ **Type-safe Configuration**: Environment variable and file-based configuration with validation

- ✅ **Idiomatic Go API**: Clean interfaces following Go conventions

- ✅ **Thread-safe Operations**: Concurrent-safe client with proper mutex protectionThe Voltage Go Wrapper provides a clean, type-safe interface to the Voltage SecureData C library, following Go best practices and patterns. This implementation focuses on **configuration and initialization** as specified in ARCH-16835.VLock simplifies the integration of Voltage's C-based SDK into Go applications by providing a clean wrapper that eliminates the need for developers to work directly with C code. This wrapper provides one-time standardized integration for all Go services, reducing onboarding time and preventing repeated C-level integration work across teams.

- ✅ **Health Monitoring**: Built-in health checks and connection management

- ✅ **Error Handling**: Go-friendly error types with retry logic support

- ✅ **CGO Integration**: Seamless integration with Voltage C library

- ✅ **Comprehensive Testing**: Full test coverage with mock support### Key Features## Features



## Installation



```bash- ✅ **Type-safe Configuration**: Environment variable and file-based configuration with validation- **Text Encryption/Decryption** - Protect and access plain text values

go get github.com/daveaugustus/vlock

```- ✅ **Idiomatic Go API**: Clean interfaces following Go conventions- **Binary Encryption/Decryption** - Support for raw data or files



### Prerequisites- ✅ **Thread-safe Operations**: Concurrent-safe client with proper mutex protection- **Masked Access** - Partial redaction of sensitive fields based on patterns



- Go 1.21 or higher- ✅ **Health Monitoring**: Built-in health checks and connection management- **Flexible Configuration** - Support for both file-based and environment variable configuration

- Voltage SecureData SimpleAPI library (for production use)

- GCC (for CGO compilation on macOS/Linux)- ✅ **Error Handling**: Go-friendly error types with retry logic support- **Multi-Environment Support** - Manage separate configurations for DEV, QA, CAT, and PROD



## Project Structure- ✅ **CGO Integration**: Seamless integration with Voltage C library- **Key Rotation** - Automatic key material updates with zero downtime



```- ✅ **Comprehensive Testing**: Full test coverage with mock support- **Error Handling** - Descriptive error codes and Go-friendly error messages

vlock/

├── pkg/- **Thread Safety** - Safe multi-threaded usage after initialization

│   ├── vlock/              # Main client package

│   │   ├── voltage.go       # Client implementation## Installation- **Container-Friendly** - Environment variable injection for Docker/Kubernetes deployments

│   │   ├── voltage_cgo.go   # CGO bindings (production)

│   │   ├── voltage_mock.go  # Mock implementation (testing)- **Memory Management** - Automatic buffer allocation to prevent C-level memory issues

│   │   ├── errors.go        # Error handling

│   │   └── voltage_test.go  # Test suite```bash

│   └── config/             # Configuration package

│       ├── config.go        # Configuration managementgo get github.com/daveaugustus/vlock## Architecture

│       └── config_test.go   # Configuration tests

├── examples/               # Example programs```

│   ├── main.go             # Usage example

│   └── voltage.cfg         # Example configuration```

├── go.mod

├── go.sum### Prerequisites┌─────────────────────────────────┐

└── README.md

```│     Go Service                  │



## Quick Start- Go 1.21 or higher│  (Business Logic / APIs)        │



### Basic Usage- Voltage SecureData SimpleAPI library (for production use)└────────────┬────────────────────┘



```go- GCC/MinGW (for CGO compilation)             │

package main

┌────────────▼────────────────────┐

import (

    "log"## Quick Start│  Go Wrapper (CGO)               │

    "github.com/daveaugustus/vlock/pkg/vlock"

    "github.com/daveaugustus/vlock/pkg/config"│  - Init / Encrypt / Decrypt     │

)

### Basic Usage│  - Error & Memory Handling      │

func main() {

    // Load configuration└────────────┬────────────────────┘

    cfg, err := config.LoadConfig("voltage.cfg")

    if err != nil {```go             │

        log.Fatalf("Failed to load config: %v", err)

    }package main┌────────────▼────────────────────┐



    // Create and initialize client│  Voltage C Library              │

    client, err := vlock.NewClient(cfg)

    if err != nil {import (│  (voltage_protect_text, etc.)   │

        log.Fatalf("Failed to create client: %v", err)

    }    "log"└─────────────────────────────────┘

    defer client.Close()

    "github.com/daveaugustus/vlock"```

    if err := client.Initialize(); err != nil {

        log.Fatalf("Failed to initialize: %v", err)    "github.com/daveaugustus/vlock/config"

    }

)## Installation

    // Client is ready for use

    log.Println("Voltage client initialized successfully")

}

```func main() {```bash



### Programmatic Configuration    // Load configurationgo get github.com/the_monkeys/vlock



```go    cfg, err := config.LoadConfig("voltage.cfg")```

cfg := &config.Config{

    AppName:         "MyApp",    if err != nil {

    AppVersion:      "1.0.0",

    AppEnv:          "PROD",        log.Fatalf("Failed to load config: %v", err)**Prerequisites:**

    DEKSharedSecret: "your_shared_secret",

    ConfigFilePath:  "/path/to/voltage.cfg",    }- Voltage C Library (Protector) installed

    NetworkTimeout:  30,

    LogLevel:        2,- CGO enabled

}

    // Create and initialize client- Valid Voltage configuration files or environment variables

client, err := vlock.NewClient(cfg)

// ... use client    client, err := vlock.NewClient(cfg)

```

    if err != nil {## Configuration

## Configuration

        log.Fatalf("Failed to create client: %v", err)

### Configuration File Format

    }### Overview

The Voltage wrapper supports INI-style configuration files:

    defer client.Close()

```ini

# Application SettingsThis section addresses the key configuration requirements for the Voltage solution, including how to acquire settings, manage them across environments, and handle security considerations like key rotation.

FP_APPNAME=MyApplication

FP_APPVERSION=1.0.0    if err := client.Initialize(); err != nil {

FP_APPENV=PROD

        log.Fatalf("Failed to initialize: %v", err)### Required Configuration Parameters

# Voltage Installation Paths

FP_SIMPLEAPI_INSTALLPATH=/opt/voltage/simpleapi    }

FP_TRUSTSTORE_PATH=/opt/voltage/truststore/cacerts.jks

FP_XMLCONFIG=/opt/voltage/config/voltage.xmlThe Voltage Protector library requires specific configuration parameters. These can be provided via configuration files or environment variables (environment variables take precedence).



# Key Encryption Key (KEK) Settings    // Client is ready for use

FP_KEK_CERTPATH=/opt/voltage/certs/kek.p12

FP_KEK_CERTPASSPHRASE=secret_passphrase    log.Println("Voltage client initialized successfully")| Configuration Purpose | .cfg Parameter | Environment Variable | Required |

FP_KEK_SHAREDSECRET=kek_shared_secret

}|----------------------|----------------|---------------------|----------|

# Data Encryption Key (DEK) Settings

FP_DEFAULT_SHAREDSECRET=your_shared_secret```| Application Name | `fp_appName` | `FP_APPNAME` | Yes |

FP_DEFAULT_USERNAME=voltage_user

FP_DEFAULT_PASSWORD=voltage_password| Application Version | `fp_appVersion` | `FP_APPVERSION` | Yes |



# Optional Settings### Programmatic Configuration| Environment | `fp_appEnv` | `FP_APPENV` | Yes |

FP_NETWORKTIMEOUT=30

FP_DISABLECRLCHECKING=false| Network Timeout | `fp_networkTimeout` | `FP_NETWORKTIMEOUT` | No |

FP_DEFAULT_CRYPTID=FPE_SSN

FP_LOGLEVEL=2```go| Disable CRL Checking | `fp_disableCRLChecking` | `FP_DISABLECRLCHECKING` | No |

FP_LOGFILE=/var/log/voltage/wrapper.log

```cfg := &config.Config{| Simple API Install Path | `fp_simpleAPI_installPath` | `FP_SIMPLEAPI_INSTALLPATH` | See docs |



### Environment Variables    AppName:         "MyApp",| Trust Store Path | `fp_trustStore_path` | `FP_TRUSTSTORE_PATH` | See docs |



All configuration values can be overridden using environment variables. Environment variables take precedence over configuration file values.    AppVersion:      "1.0.0",| KEK Cert Path | `fp_kek_certPath` | `FP_KEK_CERTPATH` | See docs |



```bash    AppEnv:          "PROD",| KEK Cert Passphrase | `fp_kek_certPassphrase` | `FP_KEK_CERTPASSPHRASE` | See docs |

export FP_APPNAME="MyApp"

export FP_APPVERSION="1.0.0"    DEKSharedSecret: "your_shared_secret",| KEK Shared Secret | `fp_kek_sharedSecret` | `FP_KEK_SHAREDSECRET` | See docs |

export FP_APPENV="PROD"

export FP_DEFAULT_SHAREDSECRET="your_shared_secret"    ConfigFilePath:  "/path/to/voltage.cfg",| DEK Shared Secret | `fp_default_sharedSecret` | `FP_DEFAULT_SHAREDSECRET` | See docs |

```

    NetworkTimeout:  30,| DEK Username | `fp_default_userName` | `FP_DEFAULT_USERNAME` | See docs |

### Configuration Precedence

    LogLevel:        2,| DEK Password | `fp_default_password` | `FP_DEFAULT_PASSWORD` | See docs |

1. **Defaults** - Built-in default values

2. **Configuration File** - Values from .cfg file}

3. **Environment Variables** - Highest priority

### Configuration File Example

## API Reference

client, err := vlock.NewClient(cfg)

### Client Initialization

// ... use client**voltageprotector.cfg:**

#### NewClient

``````ini

Creates a new Voltage client with the provided configuration.

[ProtectorConfig]

```go

func NewClient(cfg *config.Config, opts ...ClientOption) (*Client, error)## ConfigurationXMLConfig=./vsconfig.xml

```

DefaultCryptId=SSN_Internal

**Parameters:**

- `cfg`: Configuration object (required, cannot be nil)### Configuration File FormatEnvironment=DEV

- `opts`: Optional client options (variadic)

LogLevel=2

**Returns:**

- `*Client`: Initialized client instanceThe Voltage wrapper supports INI-style configuration files:LogFile=/opt/voltage/logs/protector.log

- `error`: Error if configuration is invalid



**Example:**

```go```ini# Application Configuration

cfg := &config.Config{

    AppName:         "MyApp",# Application Settingsfp_appName=YOUR_APP_NAME

    AppVersion:      "1.0.0",

    AppEnv:          "PROD",FP_APPNAME=MyApplicationfp_appVersion=YOUR_APP_VERSION

    DEKSharedSecret: "secret",

    ConfigFilePath:  "voltage.cfg",FP_APPVERSION=1.0.0fp_appEnv=DEV

}

FP_APPENV=PRODfp_default_sharedSecret=YOUR_SHARED_SECRET

client, err := vlock.NewClient(cfg)

if err != nil {fp_kek_certPath=/path/to/cert.pfx

    log.Fatalf("Failed to create client: %v", err)

}# Voltage Installation Pathsfp_kek_certPassphrase=YOUR_CERT_PASSPHRASE

```

FP_SIMPLEAPI_INSTALLPATH=/opt/voltage/simpleapifp_simpleAPI_installPath=/opt/voltage/simpleapi

#### Initialize

FP_TRUSTSTORE_PATH=/opt/voltage/truststore/cacerts.jksfp_trustStore_path=/opt/voltage/trustStore

Initializes the Voltage library and establishes connection.

FP_XMLCONFIG=/opt/voltage/config/voltage.xmlfp_networkTimeout=10

```go

func (c *Client) Initialize() errorfp_disableCRLChecking=false

```

# Key Encryption Key (KEK) Settings```

**Returns:**

- `error`: Error if initialization failsFP_KEK_CERTPATH=/opt/voltage/certs/kek.p12



**Example:**FP_KEK_CERTPASSPHRASE=secret_passphrase**vsconfig.xml (sanitized example):**

```go

if err := client.Initialize(); err != nil {FP_KEK_SHAREDSECRET=kek_shared_secret```xml

    log.Fatalf("Initialization failed: %v", err)

}<cryptId name="SSN_Internal" algorithm="FPE" key="YOUR_KEY_HERE" format="NUMERIC"/>

```

# Data Encryption Key (DEK) Settings<mask pattern="XXX-XX-####" cryptId="SSN_Internal"/>

**Note:** Must be called before performing any encryption/decryption operations.

FP_DEFAULT_SHAREDSECRET=your_shared_secret```

#### Close

FP_DEFAULT_USERNAME=voltage_user

Gracefully shuts down the client and terminates the Voltage library connection.

FP_DEFAULT_PASSWORD=voltage_password### Using Environment Variables

```go

func (c *Client) Close() error

```

# Optional Settings```bash

**Returns:**

- `error`: Error if shutdown failsFP_NETWORKTIMEOUT=30# Linux/Mac



**Example:**FP_DISABLECRLCHECKING=falseexport FP_APPNAME=MYAPPNAME

```go

defer client.Close()FP_DEFAULT_CRYPTID=FPE_SSNexport FP_APPVERSION=1.0.0

```

FP_LOGLEVEL=2export FP_APPENV=DEV

**Best Practice:** Always defer `Close()` after successful client creation.

FP_LOGFILE=/var/log/voltage/wrapper.logexport FP_SIMPLEAPI_INSTALLPATH=/opt/voltage/simpleapi

### Health Monitoring

```export FP_TRUSTSTORE_PATH=/opt/voltage/trustStore

#### HealthCheck

export FP_DEFAULT_SHAREDSECRET=yoursharedsecret

Performs a health check on the Voltage library connection.

### Environment Variables```

```go

func (c *Client) HealthCheck() error

```

All configuration values can be overridden using environment variables. Environment variables take precedence over configuration file values.```powershell

**Returns:**

- `error`: Error if health check fails, nil if healthy# Windows PowerShell



#### IsHealthy```bash$env:FP_APPNAME="MYAPPNAME"



Returns the current health status without performing a new check.export FP_APPNAME="MyApp"$env:FP_APPVERSION="1.0.0"



```goexport FP_APPVERSION="1.0.0"$env:FP_APPENV="DEV"

func (c *Client) IsHealthy() bool

```export FP_APPENV="PROD"$env:FP_SIMPLEAPI_INSTALLPATH="C:\voltage\simpleapi"



**Returns:**export FP_DEFAULT_SHAREDSECRET="your_shared_secret"$env:FP_TRUSTSTORE_PATH="C:\voltage\trustStore"

- `bool`: true if client is healthy, false otherwise

```$env:FP_DEFAULT_SHAREDSECRET="yoursharedsecret"

#### LastHealthCheck

```

Returns the timestamp of the last health check.

### Configuration Precedence

```go

func (c *Client) LastHealthCheck() time.Time### Environment Management

```

1. **Defaults** - Built-in default values

**Returns:**

- `time.Time`: Timestamp of last health check2. **Configuration File** - Values from .cfg file#### Moving Between Environments (DEV → QA → CAT → PROD)



### Client State3. **Environment Variables** - Highest priority



#### IsInitializedConfiguration settings change when promoting between environments. Each environment requires its own:



Checks if the client has been initialized.## API Reference



```go1. **Unique Configuration Files** - Separate `.cfg` and `.xml` files per environment

func (c *Client) IsInitialized() bool

```### Client Initialization2. **Environment-Specific Keys** - Different encryption keys for DEV, QA, CAT, and PROD



**Returns:**3. **Environment Variable Settings** - Different values for `FP_APPENV` parameter

- `bool`: true if initialized, false otherwise

#### NewClient

#### Config

**Best Practices:**

Returns the current configuration.

Creates a new Voltage client with the provided configuration.- Store configuration files in environment-specific directories (e.g., `/config/dev`, `/config/prod`)

```go

func (c *Client) Config() *config.Config- Use environment-specific service accounts and credentials

```

```go- Never promote encrypted data across environments (re-encrypt in target environment)

**Returns:**

- `*config.Config`: Current configurationfunc NewClient(cfg *config.Config, opts ...ClientOption) (*Client, error)- Validate configuration before deployment using the `Init()` validation



#### Info```



Returns detailed client information.**Example Environment Structure:**



```go**Parameters:**```

func (c *Client) Info() ClientInfo

```- `cfg`: Configuration object (required, cannot be nil)/config



**Returns:**- `opts`: Optional client options (variadic)  /dev

- `ClientInfo`: Struct containing client state information

    voltageprotector.cfg

### Advanced Operations

**Returns:**    vsconfig.xml

#### Reinitialize

- `*Client`: Initialized client instance  /qa

Reinitializes the client (performs Close followed by Initialize).

- `error`: Error if configuration is invalid    voltageprotector.cfg

```go

func (c *Client) Reinitialize() error    vsconfig.xml

```

**Example:**  /cat

**Returns:**

- `error`: Error if reinitialization fails```go    voltageprotector.cfg



## Error Handlingcfg := &config.Config{    vsconfig.xml



The wrapper provides rich error types for better error handling:    AppName:         "MyApp",  /prod



```go    AppVersion:      "1.0.0",    voltageprotector.cfg

type VoltageError struct {

    Code    ErrorCode    AppEnv:          "PROD",    vsconfig.xml

    Message string

    Detail  string    DEKSharedSecret: "secret",```

    CError  int // Original C error code

}    ConfigFilePath:  "voltage.cfg",

```

}#### Acquiring Configuration Settings

### Error Categories



- **Configuration Errors**: Invalid or missing configuration

- **Initialization Errors**: Library initialization failuresclient, err := vlock.NewClient(cfg)Configuration settings are typically provided by:

- **Connection Errors**: Network or connectivity issues

- **Operation Errors**: Encryption/decryption failuresif err != nil {



### Retryable Errors    log.Fatalf("Failed to create client: %v", err)1. **Voltage Platform Team** - Provides initial setup including:



Check if an error is retryable:}   - Certificate files (`.pfx`)



```go```   - Shared secrets and credentials

if err != nil {

    if voltageErr, ok := err.(*vlock.VoltageError); ok {   - Trust store files

        if voltageErr.IsRetryable() {

            // Implement retry logic#### Initialize   - XML configuration templates

            time.Sleep(time.Second * 2)

            err = client.Reinitialize()

        }

    }Initializes the Voltage library and establishes connection.2. **Internal Security/Platform Team** - Manages:

}

```   - Environment-specific credentials



## Testing```go   - Key rotation schedules



### Running Testsfunc (c *Client) Initialize() error   - Access control policies



```bash```

# Run all tests

go test -v ./...3. **External Teams (e.g., Finxact)** - May provide:



# Run with coverage**Returns:**   - Application-specific configuration

go test -v -cover ./...

- `error`: Error if initialization fails   - Environment variables for container deployments

# Run specific package tests

go test -v ./pkg/vlock   - Service account credentials

go test -v ./pkg/config

```**Example:**



### Mock Mode```go**Recommended Acquisition Process:**



The wrapper includes a mock implementation for testing without the actual Voltage library:if err := client.Initialize(); err != nil {1. Submit configuration request to Voltage platform team



```go    log.Fatalf("Initialization failed: %v", err)2. Receive environment-specific credentials and certificates

// Mock mode is automatically used when CGO is not available

// or when building with -tags="!cgo"}3. Store sensitive values in secure vault (e.g., HashiCorp Vault, AWS Secrets Manager)



go test -tags="!cgo" ./...```4. Inject configuration via environment variables at runtime

```

5. Validate configuration in non-production environment first

## Examples

**Note:** Must be called before performing any encryption/decryption operations.

See the [examples](./examples) directory for complete working examples:

### Key Rotation

- **[main.go](./examples/main.go)**: Complete initialization example

- **[voltage.cfg](./examples/voltage.cfg)**: Example configuration file#### Close



### Running ExamplesThe Voltage library supports automatic key rotation for enhanced security compliance.



```bashGracefully shuts down the client and terminates the Voltage library connection.

# From project root

go run examples/main.go#### Rotation Schedule



# With custom config file```go

VOLTAGE_CONFIG_FILE=/path/to/custom.cfg go run examples/main.go

```func (c *Client) Close() error| Environment | Recommended Interval | Compliance Requirement |



## Best Practices```|-------------|---------------------|------------------------|



### 1. Always Close the Client| Development | 90 days | Optional |



```go**Returns:**| QA/CAT | 60 days | Recommended |

client, err := vlock.NewClient(cfg)

if err != nil {- `error`: Error if shutdown fails| Production | 30-90 days | Required (per policy) |

    return err

}

defer client.Close() // Ensures cleanup

```**Example:****Note:** Actual rotation intervals should align with your organization's security policies and compliance requirements (PCI-DSS, HIPAA, etc.).



### 2. Check Initialization State```go



```godefer client.Close()#### Key Rotation Process

if !client.IsInitialized() {

    if err := client.Initialize(); err != nil {```

        return err

    }**Step 1: Preparation**

}

```**Best Practice:** Always defer `Close()` after successful client creation.```bash



### 3. Implement Health Monitoring# Backup current configuration



```go### Health Monitoringcp vsconfig.xml vsconfig.xml.backup.$(date +%Y%m%d)

ticker := time.NewTicker(5 * time.Minute)

defer ticker.Stop()cp voltageprotector.cfg voltageprotector.cfg.backup.$(date +%Y%m%d)



go func() {#### HealthCheck```

    for range ticker.C {

        if err := client.HealthCheck(); err != nil {

            log.Printf("Health check failed: %v", err)

            client.Reinitialize()Performs a health check on the Voltage library connection.**Step 2: Request New Keys**

        }

    }- Contact Voltage platform team or use Voltage management console

}()

``````go- Request new key material for specific `cryptId`



### 4. Thread Safetyfunc (c *Client) HealthCheck() error- Receive updated XML configuration with new key references



The client is thread-safe and can be shared across goroutines:```



```go**Step 3: Update Configuration**

var client *vlock.Client // Shared client

**Returns:**```xml

// Goroutine 1

go func() {- `error`: Error if health check fails, nil if healthy<!-- Updated vsconfig.xml with new key -->

    if client.IsHealthy() {

        // Perform operations<cryptId name="SSN_Internal" algorithm="FPE" key="NEW_KEY_HERE" format="NUMERIC"/>

    }

}()**Example:**```



// Goroutine 2```go

go func() {

    info := client.Info()if err := client.HealthCheck(); err != nil {**Step 4: Deploy Configuration**

    log.Printf("Client info: %+v", info)

}()    log.Printf("Health check failed: %v", err)```bash

```

    // Implement retry or failover logic# Deploy to target environment

## Architecture

}# The library automatically uses new keys for encryption

```

┌─────────────────────────────────────────┐```# Old keys remain available for decryption of existing data

│          Application Layer              │

├─────────────────────────────────────────┤```

│    pkg/vlock (Client & Error Handling)  │

│    pkg/config (Configuration Mgmt)      │#### IsHealthy

├─────────────────────────────────────────┤

│           CGO Bindings Layer            │**Step 5: Validation**

│  voltage_cgo.go / voltage_mock.go       │

├─────────────────────────────────────────┤Returns the current health status without performing a new check.```go

│      Voltage C Library (libvoltage)     │

│   Format-Preserving Encryption (FPE)    │// Test encryption with new keys

└─────────────────────────────────────────┘

``````goerr := voltage.Init("/path/to/updated/config.cfg")



## Ticket Referencefunc (c *Client) IsHealthy() boolif err != nil {



This implementation addresses **ARCH-16835: Voltage Wrapper - Create Go wrapper config/init**.```    log.Fatal("Configuration validation failed:", err)



### Scope}



- ✅ Configuration management (file + environment variables)**Returns:**

- ✅ Client initialization and lifecycle

- ✅ Health monitoring and connection management- `bool`: true if client is healthy, false otherwise// Verify new encryptions work

- ✅ Error handling and retry logic

- ✅ Thread-safe operationsencrypted, err := voltage.Encrypt("test-data", "SSN_Internal")

- ✅ Comprehensive testing

**Example:**// Verify old encrypted data can still be decrypted

### Future Work

```godecrypted, err := voltage.Decrypt(oldEncryptedData, "SSN_Internal")

- ⏳ Encryption/decryption operations (separate ticket)

- ⏳ Connection pooling optimization (if needed)if client.IsHealthy() {```

- ⏳ Additional format support

    // Proceed with operations

## Support

}**Step 6: Monitor**

For issues or questions:

```- Monitor application logs for decryption errors

1. Check the [examples](./examples) directory

2. Review API documentation above- Verify both new and old ciphertext can be processed

3. Contact: @John Farley (for code review)

#### LastHealthCheck- Plan re-encryption of old data if required by policy

## License



Internal Fiserv project - See company licensing guidelines.

Returns the timestamp of the last health check.#### Automatic Key Material Updates

---



**Version**: 1.0.0  

**Last Updated**: November 2025  ```goThe Voltage library automatically detects and uses updated key materials from the XML configuration file:

**Maintainer**: Dave Augustus

func (c *Client) LastHealthCheck() time.Time- No application restart required for key rotation

```- Library checks for configuration updates periodically

- Both old and new keys remain active during transition period

**Returns:**- Re-encryption of existing data can be performed gradually

- `time.Time`: Timestamp of last health check

#### Key Rotation Considerations

**Example:**

```go- **Zero Downtime**: Key rotation should not cause service interruption

lastCheck := client.LastHealthCheck()- **Backward Compatibility**: Old ciphertext must remain decryptable during transition

if time.Since(lastCheck) > 5*time.Minute {- **Audit Trail**: Maintain logs of all key rotation activities

    client.HealthCheck() // Perform new check- **Testing**: Always test key rotation in non-production environments first

}- **Rollback Plan**: Keep previous configuration backups for emergency rollback

```

### Configuration FAQ

### Client State

**Q: Can we inject configuration via method parameters instead of files?**  

#### IsInitializedA: No. The Voltage C library requires file-based configuration. Configuration items cannot be injected through Go API calls. You must use `.cfg` files or environment variables.



Checks if the client has been initialized.**Q: What data format does the configuration use?**  

A: The library uses two formats:

```go- `.cfg` files: Plain text key-value pairs (similar to INI format)

func (c *Client) IsInitialized() bool- `.xml` files: Structured XML for encryption rules, cryptIds, and mask patterns

```

**Q: Can we use environment variables like other applications?**  

**Returns:**A: Yes! Environment variables are the recommended approach for containerized deployments. All configuration parameters have corresponding environment variable names (e.g., `FP_APPNAME`, `FP_APPENV`). Environment variables take precedence over file-based configuration.

- `bool`: true if initialized, false otherwise

**Q: How do we handle configuration in container environments?**  

#### ConfigA: Recommended approach for containers:

1. Store sensitive values in secrets management system

Returns the current configuration.2. Inject as environment variables at container runtime

3. Mount configuration files as ConfigMaps/volumes (if needed)

```go4. Use `FP_*` environment variables to override file settings

func (c *Client) Config() *config.Config

```Example Docker/Kubernetes approach:

```yaml

**Returns:**env:

- `*config.Config`: Current configuration  - name: FP_APPNAME

    value: "MyService"

#### Info  - name: FP_APPENV

    value: "PROD"

Returns detailed client information.  - name: FP_DEFAULT_SHAREDSECRET

    valueFrom:

```go      secretKeyRef:

func (c *Client) Info() ClientInfo        name: voltage-secrets

```        key: shared-secret

```

**Returns:**

- `ClientInfo`: Struct containing client state information**Q: What happens if configuration is missing or invalid?**  

A: The `Init()` function will return an error with details about the missing/invalid configuration. The wrapper validates all required parameters before allowing any encryption/decryption operations.

**Example:**

```go**Q: Can we use the same configuration across multiple services?**  

info := client.Info()A: Configuration can be shared across services if they:

fmt.Printf("App: %s v%s\n", info.AppName, info.AppVersion)- Run in the same environment (DEV, QA, PROD)

fmt.Printf("Environment: %s\n", info.Environment)- Use the same application name and version

fmt.Printf("Initialized: %v\n", info.Initialized)- Share the same security requirements

fmt.Printf("Healthy: %v\n", info.Healthy)

fmt.Printf("Session ID: %s\n", info.SessionID)However, it's recommended to use service-specific `fp_appName` values for better audit trails and access control.

```

**Q: How do we secure sensitive configuration values?**  

### Advanced OperationsA: Best practices:

- Never commit secrets to version control

#### Reinitialize- Use secrets management systems (Vault, AWS Secrets Manager, etc.)

- Inject sensitive values as environment variables at runtime

Reinitializes the client (performs Close followed by Initialize).- Restrict file permissions on configuration files (chmod 600)

- Rotate credentials regularly

```go

func (c *Client) Reinitialize() error**Q: Where can we get help with configuration issues?**  

```A: Configuration support resources:

- Voltage platform documentation

**Returns:**- Internal Voltage platform team

- `error`: Error if reinitialization fails- Voltage C API documentation (`voltage_api.h`)

- Voltage Protector error codes reference

**Example:**

```go

// Reinitialize after configuration change## Usage

if err := client.Reinitialize(); err != nil {

    log.Printf("Reinitialization failed: %v", err)### Basic Text Encryption/Decryption

}

``````go

package main

**Use Cases:**

- Recovering from connection failuresimport (

- Applying new configuration    "fmt"

- Resetting after errors    "log"

    "github.com/the_monkeys/vlock"

## Error Handling)



The wrapper provides rich error types for better error handling:func main() {

    // Initialize the library

```go    err := voltage.Init("/path/to/voltageprotector.cfg")

type VoltageError struct {    if err != nil {

    Code    ErrorCode        log.Fatal("Failed to initialize:", err)

    Message string    }

    Detail  string    defer voltage.Close()

    CError  int // Original C error code

}    // Encrypt data

```    plainText := "123-45-6789"

    encrypted, err := voltage.Encrypt(plainText, "SSN_Internal")

### Error Categories    if err != nil {

        log.Fatal("Encryption failed:", err)

- **Configuration Errors**: Invalid or missing configuration    }

- **Initialization Errors**: Library initialization failures    fmt.Println("Encrypted:", encrypted)

- **Connection Errors**: Network or connectivity issues

- **Operation Errors**: Encryption/decryption failures    // Decrypt data

    decrypted, err := voltage.Decrypt(encrypted, "SSN_Internal")

### Retryable Errors    if err != nil {

        log.Fatal("Decryption failed:", err)

Check if an error is retryable:    }

    fmt.Println("Decrypted:", decrypted)

```go}

if err != nil {```

    if voltageErr, ok := err.(*vlock.VoltageError); ok {

        if voltageErr.IsRetryable() {### Masked Access

            // Implement retry logic

            time.Sleep(time.Second * 2)```go

            err = client.Reinitialize()// Initialize and encrypt as above

        }encrypted, err := voltage.Encrypt("123-45-6789", "SSN_Internal")

    }if err != nil {

}    log.Fatal(err)

```}



### Error Example// Get masked version (based on XML pattern configuration)

masked, err := voltage.DecryptMasked(encrypted, "SSN_Internal")

```goif err != nil {

err := client.Initialize()    log.Fatal(err)

if err != nil {}

    if voltageErr, ok := err.(*vlock.VoltageError); ok {fmt.Println("Masked:", masked) // Output: XXX-XX-6789

        fmt.Printf("Error Code: %d\n", voltageErr.Code)```

        fmt.Printf("Message: %s\n", voltageErr.Message)

        fmt.Printf("Category: %s\n", voltageErr.Category())### Binary Encryption/Decryption

        fmt.Printf("Retryable: %v\n", voltageErr.IsRetryable())

    }```go

}// Encrypt binary data

```binaryData := []byte{0x01, 0x02, 0x03, 0x04}

encryptedBinary, err := voltage.EncryptBinary(binaryData, "BinaryCryptId")

## Testingif err != nil {

    log.Fatal("Binary encryption failed:", err)

### Running Tests}



```bash// Decrypt binary data

# Run all testsdecryptedBinary, err := voltage.DecryptBinary(encryptedBinary, "BinaryCryptId")

go test -v ./...if err != nil {

    log.Fatal("Binary decryption failed:", err)

# Run with coverage}

go test -v -cover ./...```



# Run specific test## API Reference

go test -v -run TestClientInitialize

```### Initialization



### Mock Mode```go

func Init(configPath string) error

The wrapper includes a mock implementation for testing without the actual Voltage library:```

Initializes the Voltage Protector library with the specified configuration file. Must be called before any encryption/decryption operations.

```go

// Mock mode is automatically used when CGO is not available### Text Operations

// or when building with -tags="!cgo"

```go

go test -tags="!cgo" ./...func Encrypt(plainText string, cryptId string) (string, error)

``````

Encrypts plain text using the specified cryptId.

## Examples

```go

See the [examples](./examples) directory for complete working examples:func Decrypt(cipherText string, cryptId string) (string, error)

```

- **[main.go](./examples/main.go)**: Complete initialization exampleDecrypts cipher text using the specified cryptId.

- **[voltage.cfg](./examples/voltage.cfg)**: Example configuration file

```go

### Running Examplesfunc DecryptMasked(cipherText string, cryptId string) (string, error)

```

```bashReturns partially decrypted data based on mask patterns defined in XML configuration.

# From project root

go run examples/main.go### Binary Operations



# With custom config file```go

VOLTAGE_CONFIG_FILE=/path/to/custom.cfg go run examples/main.gofunc EncryptBinary(data []byte, cryptId string) ([]byte, error)

``````

Encrypts binary data using the specified cryptId.

## Best Practices

```go

### 1. Always Close the Clientfunc DecryptBinary(data []byte, cryptId string) ([]byte, error)

```

```goDecrypts binary data using the specified cryptId.

client, err := vlock.NewClient(cfg)

if err != nil {### Cleanup

    return err

}```go

defer client.Close() // Ensures cleanupfunc Close() error

``````

Releases resources and performs cleanup. Should be called when done using the library.

### 2. Check Initialization State

## Error Handling

```go

if !client.IsInitialized() {The wrapper translates C-level errors into meaningful Go errors:

    if err := client.Initialize(); err != nil {

        return err| Error | Description | Common Cause |

    }|-------|-------------|--------------|

}| `Error 5500207` | Output buffer is too small | Insufficient buffer allocation |

```| Invalid ciphertext | Trailing characters mismatch or invalid base64 | Corrupted or tampered data |

| Config XML issues | Missing or invalid cryptId configuration | Configuration file errors |

### 3. Implement Health Monitoring

Example error handling:

```go```go

ticker := time.NewTicker(5 * time.Minute)encrypted, err := voltage.Encrypt(plainText, "SSN_Internal")

defer ticker.Stop()if err != nil {

    switch {

go func() {    case strings.Contains(err.Error(), "5500207"):

    for range ticker.C {        // Handle buffer size error

        if err := client.HealthCheck(); err != nil {    case strings.Contains(err.Error(), "invalid ciphertext"):

            log.Printf("Health check failed: %v", err)        // Handle invalid data error

            client.Reinitialize()    default:

        }        // Handle other errors

    }    }

}()}

``````



### 4. Handle Retryable Errors## Important Considerations



```go### Environment-Specific Encryption

func initializeWithRetry(client *vlock.Client, maxRetries int) error {

    for i := 0; i < maxRetries; i++ {**⚠️ Warning:** Ciphertext encrypted in one environment (e.g., DEV) cannot be decrypted in another environment (e.g., QA) due to differing keys and configurations. The wrapper validates and enforces environment-specific configuration usage.

        if err := client.Initialize(); err != nil {

            if voltageErr, ok := err.(*vlock.VoltageError); ok && voltageErr.IsRetryable() {### Thread Safety

                time.Sleep(time.Duration(i+1) * time.Second)

                continueThe wrapper is thread-safe after initialization. Internal locking is provided around initialization to ensure safe concurrent usage.

            }

            return err### Memory Management

        }

        return nilThe wrapper handles buffer allocation internally to avoid common C-related memory issues and buffer overflow errors.

    }

    return fmt.Errorf("max retries exceeded")### Configuration Restrictions

}

```- Configuration items cannot be injected through Go API calls

- All configuration must be present in configuration files or environment variables at initialization

### 5. Thread Safety- The `.cfg` file must exist and reference the correct XML at runtime if not using environment variables



The client is thread-safe and can be shared across goroutines:## Development



```go### Building

var client *vlock.Client // Shared client

```bash

// Goroutine 1# Ensure CGO is enabled

go func() {export CGO_ENABLED=1

    if client.IsHealthy() {

        // Perform operations# Build the project

    }go build

}()```



// Goroutine 2## Testing

go func() {

    info := client.Info()### Running Tests

    log.Printf("Client info: %+v", info)

}()#### 1. Run All Tests

```

```bash

## Architecture# Run all tests with verbose output

go test -v

### Component Overview

# Expected output:

```# === RUN   TestNewConfig

┌─────────────────────────────────────────┐# --- PASS: TestNewConfig (0.00s)

│          Application Layer              │# === RUN   TestLoadConfigFromFile

├─────────────────────────────────────────┤# --- PASS: TestLoadConfigFromFile (0.00s)

│        VLock Go Wrapper (vlock)         │# ...

│  ┌────────────┬──────────┬───────────┐  │# PASS

│  │   Client   │  Config  │  Errors   │  │# ok      github.com/daveaugustus/vlock   1.169s

│  └────────────┴──────────┴───────────┘  │```

├─────────────────────────────────────────┤

│           CGO Bindings Layer            │#### 2. Run Specific Tests

│  ┌────────────────────────────────────┐ │

│  │  voltage_cgo.go / voltage_mock.go  │ │```bash

│  └────────────────────────────────────┘ │# Run only configuration tests

├─────────────────────────────────────────┤go test -v -run TestLoadConfig

│      Voltage C Library (libvoltage)     │

│  ┌────────────────────────────────────┐ │# Run only validation tests

│  │   Format-Preserving Encryption     │ │go test -v -run TestValidate

│  └────────────────────────────────────┘ │

└─────────────────────────────────────────┘# Run only environment variable tests

```go test -v -run TestEnvVar

```

### Key Components

#### 3. Run Tests with Coverage

1. **Client** (`voltage.go`): Main client interface with lifecycle management

2. **Config** (`config/`): Configuration loading and validation```bash

3. **Errors** (`errors.go`): Go-friendly error types# Generate coverage report

4. **CGO Bindings** (`voltage_cgo.go`): C library integrationgo test -cover

5. **Mock** (`voltage_mock.go`): Testing without C library

# Expected output:

## Ticket Reference# PASS

# coverage: 95.2% of statements

This implementation addresses **ARCH-16835: Voltage Wrapper - Create Go wrapper config/init**.# ok      github.com/daveaugustus/vlock   1.169s



### Scope# Generate detailed HTML coverage report

go test -coverprofile=coverage.out

- ✅ Configuration management (file + environment variables)go tool cover -html=coverage.out

- ✅ Client initialization and lifecycle```

- ✅ Health monitoring and connection management

- ✅ Error handling and retry logic#### 4. Run Tests with Race Detection

- ✅ Thread-safe operations

- ✅ Comprehensive testing```bash

# Detect data races in concurrent code

### Future Workgo test -race -v



- ⏳ Encryption/decryption operations (separate ticket)# Use this to ensure thread safety

- ⏳ Connection pooling optimization (if needed)```

- ⏳ Additional format support

### Test Scenarios

## Support

#### Configuration Loading Tests

For issues or questions:

```bash

1. Check the [examples](./examples) directory# Test file-based configuration loading

2. Review API documentation abovego test -v -run TestLoadConfigFromFile

3. Contact: @John Farley (for code review)

# Test environment variable precedence

## Licensego test -v -run TestEnvVarPrecedence



Internal Fiserv project - See company licensing guidelines.# Test configuration validation

go test -v -run TestValidateRequiredFields

---```



**Version**: 1.0.0  **What these tests verify:**

**Last Updated**: November 2025  - ✅ Configuration files are parsed correctly

**Maintainer**: Dave Augustus- ✅ Environment variables override file values

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
