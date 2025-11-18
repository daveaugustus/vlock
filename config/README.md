# VLock Configuration Guide

This directory contains configuration files for the VLock Voltage wrapper library across different environments.

## Directory Structure

```
config/
├── dev/
│   ├── fiservprotector20.cfg   # Development configuration
│   └── vsconfig.xml             # Development encryption rules
├── qa/
│   ├── fiservprotector20.cfg   # QA configuration
│   └── vsconfig.xml             # QA encryption rules
├── prod/
│   ├── fiservprotector20.cfg   # Production configuration
│   └── vsconfig.xml             # Production encryption rules
└── README.md                    # This file
```

## Configuration File Types

### 1. fiservprotector20.cfg
- **Format**: INI-style key-value pairs
- **Purpose**: Primary configuration for application settings, paths, and credentials
- **Contains**: App info, paths, authentication, network settings, logging

### 2. vsconfig.xml
- **Format**: XML
- **Purpose**: Encryption rules, cryptIds, and masking patterns
- **Contains**: Encryption algorithms, keys, masking patterns, security settings

## Using Configuration Files

### Option 1: Direct File Path
```go
err := voltage.Init("/path/to/config/dev/fiservprotector20.cfg")
```

### Option 2: Environment Variables (Recommended)
```bash
export FP_APPNAME=VLockDev
export FP_APPVERSION=1.0.0
export FP_APPENV=DEV
export FP_DEFAULT_SHAREDSECRET=your_secret_here
```

### Option 3: Hybrid Approach
Use config file for paths and structure, override sensitive values with environment variables:
```bash
# Base configuration from file
export CONFIG_PATH=/path/to/config/dev/fiservprotector20.cfg

# Override secrets via environment
export FP_KEK_CERTPASSPHRASE=actual_passphrase
export FP_DEFAULT_SHAREDSECRET=actual_secret
```

## Environment-Specific Settings

| Setting | DEV | QA | PROD |
|---------|-----|-----|------|
| AppEnv | DEV | QA | PROD |
| CRL Checking | Disabled | Enabled | Enabled |
| Log Level | 3 (Detailed) | 2 (Normal) | 2 (Normal) |
| Network Timeout | 10s | 20s | 30s |
| Key Rotation | Optional | 60 days | 30-60 days |

## Important Security Notes

### ⚠️ DO NOT
- Commit real credentials to version control
- Use production keys in non-production environments
- Share configuration files containing secrets
- Reuse keys across environments

### ✅ DO
- Store sensitive values in secrets management systems (Vault, AWS Secrets Manager, etc.)
- Use environment variables for credential injection
- Rotate keys regularly according to security policy
- Restrict file permissions (chmod 600 for .cfg files)
- Use separate keys for each environment
- Validate configuration before deployment

## Acquiring Configuration

Configuration values are typically obtained from:

1. **Fiserv/Voltage Team**
   - Certificate files (.pfx)
   - Initial shared secrets
   - Trust store files
   - XML configuration templates

2. **Internal Security Team**
   - Environment-specific credentials
   - Key rotation schedules
   - Compliance requirements

3. **Platform Team**
   - Installation paths
   - Service account credentials
   - Network configurations

## Configuration Parameters Reference

### Required Parameters
- `fp_appName` - Your application name
- `fp_appVersion` - Application version
- `fp_appEnv` - Environment (DEV, QA, CAT, PROD)

### Authentication (at least one required)
- KEK: `fp_kek_certPath` + `fp_kek_certPassphrase` OR `fp_kek_sharedSecret`
- DEK: `fp_default_sharedSecret` OR (`fp_default_userName` + `fp_default_password`)

### Paths
- `fp_simpleAPI_installPath` - Voltage SimpleAPI installation path
- `fp_trustStore_path` - Trust store directory path
- `XMLConfig` - Path to vsconfig.xml

### Optional Parameters
- `fp_networkTimeout` - Network timeout in seconds (default: 10)
- `fp_disableCRLChecking` - Disable CRL checking (default: false)
- `DefaultCryptId` - Default encryption algorithm ID
- `LogLevel` - Logging verbosity (0-3)
- `LogFile` - Log file path

## Example Usage

### Development
```go
package main

import (
    "log"
    "github.com/the_monkeys/vlock"
)

func main() {
    // Load dev configuration
    err := voltage.Init("./config/dev/fiservprotector20.cfg")
    if err != nil {
        log.Fatal("Init failed:", err)
    }
    defer voltage.Close()
    
    // Use encryption
    encrypted, _ := voltage.Encrypt("123-45-6789", "SSN_Internal")
    log.Println("Encrypted:", encrypted)
}
```

### Production (with Environment Variables)
```bash
# Set environment variables
export FP_APPNAME=VLockProd
export FP_APPVERSION=1.0.0
export FP_APPENV=PROD
export FP_KEK_CERTPASSPHRASE=$(vault read -field=passphrase secret/voltage/prod)
export FP_DEFAULT_SHAREDSECRET=$(vault read -field=secret secret/voltage/prod)

# Run application (config path optional if all required vars are set)
./your-app
```

## Troubleshooting

### Configuration Not Loading
- Check file path is correct and accessible
- Verify file permissions (readable by application user)
- Check for syntax errors in .cfg file

### Validation Errors
- Ensure all required parameters are set
- Verify AppEnv is one of: DEV, QA, CAT, PROD
- Check that at least one authentication method is configured

### Environment Variables Not Working
- Ensure variables are exported before running application
- Check variable names match exactly (case-sensitive)
- Environment variables override file values

## Support

For configuration issues or questions:
- Review [Fiserv Protector Configuration Parameters](https://enterprise-confluence.onefiserv.net/pages/viewpage.action?pageId=494549691)
- Contact Voltage platform team
- Review main README.md for detailed documentation
