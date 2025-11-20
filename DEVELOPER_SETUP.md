# VLock - Developer Setup Guide

## Overview

This guide explains how to set up and use the VLock Voltage wrapper library on your local development machine.

---

## Prerequisites

### 1. Go Installation
```bash
# Check if Go is installed
go version

# Should show Go 1.21 or higher
```

If not installed, download from [golang.org](https://golang.org/dl/)

### 2. CGO Enabled

VLock requires CGO to interface with the Voltage C library.

```bash
# Linux/Mac
export CGO_ENABLED=1

# Windows PowerShell
$env:CGO_ENABLED=1

# Verify
go env CGO_ENABLED  # Should output: 1
```

### 3. Voltage C Library (Protector)

The Voltage SimpleAPI C library must be installed on your system.

**Installation locations:**
- **Linux/Mac**: `/opt/voltage/simpleapi`
- **Windows**: `C:\voltage\simpleapi` (or custom location)

**Required libraries:**
- `libfiservprotect` - Core protection library
- `libvibesimple` - Voltage SimpleAPI library

---

## Installation

### Option 1: Using `go get` (Recommended)

```bash
go get github.com/the_monkeys/vlock
```

### Option 2: Clone Repository

```bash
git clone https://github.com/the_monkeys/vlock.git
cd vlock
go mod download
```

---

## Dependencies Beyond `go install`

### 1. Voltage SimpleAPI C Library

**What it is:** The underlying C library that performs the actual encryption/decryption.

**Where to get it:**
- Contact your Voltage platform team
- Request access to Voltage SimpleAPI for your platform
- Download from internal Voltage distribution

**Installation:**
```bash
# Linux/Mac
tar -xzf voltage-simple-api-c-6.22.0.2-[Platform].tar.gz
sudo mv voltage-simple-api-c-6.22.0.2 /opt/voltage/simpleapi

# Windows
# Extract zip and place in C:\voltage\simpleapi
```

### 2. Trust Store

**What it is:** Certificate trust store for Voltage server communication.

**Where to get it:**
- Provided by Voltage platform team
- Environment-specific (DEV, QA, PROD)

**Installation:**
```bash
# Linux/Mac
sudo mkdir -p /opt/voltage/trustStore
sudo cp -r trustStore/* /opt/voltage/trustStore/

# Windows
# Create C:\voltage\trustStore
# Copy trustStore files there
```

### 3. C Compiler

CGO requires a C compiler:

**Linux:**
```bash
# Ubuntu/Debian
sudo apt-get install build-essential

# RHEL/CentOS
sudo yum groupinstall "Development Tools"
```

**Mac:**
```bash
# Install Xcode Command Line Tools
xcode-select --install
```

**Windows:**
- Install [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) or
- Install [MinGW-w64](https://www.mingw-w64.org/) or
- Use Visual Studio C++ Build Tools

---

## Configuration Settings

### Where Developers Get Config Settings

#### 1. Request from Voltage Platform Team

**What to request:**
- Configuration file template
- Environment-specific credentials
- Certificate files (`.pfx`)
- Shared secrets
- Trust store files

**How to request:**
- Open ticket with Voltage platform team
- Specify your environment (DEV, QA, PROD)
- Provide your application name and version

#### 2. Configuration File Location

Developers should place configuration in a secure location:

**Recommended structure:**
```
/home/username/voltage/config/    (Linux/Mac)
C:\Users\Username\voltage\config\  (Windows)
```

**Example:**
```
~/voltage/config/
├── voltageprotector.cfg
├── vsconfig.xml
├── cert.pfx
└── trustStore/
```

#### 3. Using Environment Variables (Alternative)

Instead of config files, developers can use environment variables:

**Linux/Mac (~/.bashrc or ~/.zshrc):**
```bash
export FP_APPNAME=VLockDev
export FP_APPVERSION=1.0.0
export FP_APPENV=DEV
export FP_SIMPLEAPI_INSTALLPATH=/opt/voltage/simpleapi
export FP_TRUSTSTORE_PATH=/opt/voltage/trustStore
export FP_KEK_CERTPATH=/home/username/voltage/config/cert.pfx
export FP_KEK_CERTPASSPHRASE="your_passphrase"
export FP_DEFAULT_SHAREDSECRET="your_secret"
```

**Windows (PowerShell profile):**
```powershell
$env:FP_APPNAME="VLockDev"
$env:FP_APPVERSION="1.0.0"
$env:FP_APPENV="DEV"
$env:FP_SIMPLEAPI_INSTALLPATH="C:\voltage\simpleapi"
$env:FP_TRUSTSTORE_PATH="C:\voltage\trustStore"
$env:FP_KEK_CERTPATH="C:\Users\Username\voltage\config\cert.pfx"
$env:FP_KEK_CERTPASSPHRASE="your_passphrase"
$env:FP_DEFAULT_SHAREDSECRET="your_secret"
```

---

## Platform-Specific Setup

### Windows Setup

#### Differences on Windows:
1. **Path separators:** Use backslashes `\` or forward slashes `/`
2. **Library paths:** Different from Linux/Mac
3. **CGO flags:** May need different linker flags

#### Windows-Specific Configuration:

**voltageprotector.cfg:**
```ini
fp_simpleAPI_installPath=C:/voltage/simpleapi
fp_trustStore_path=C:/voltage/trustStore
fp_kek_certPath=C:/Users/Username/voltage/config/cert.pfx
```

#### CGO Configuration for Windows:

**go.mod or build command:**
```bash
# Set library path
set CGO_LDFLAGS=-LC:\voltage\simpleapi\lib -lfiservprotect -lvibesimple

# Build
go build
```

#### Common Windows Issues:
- ✅ **DLL not found:** Add library path to PATH environment variable
- ✅ **Permission denied:** Run as Administrator
- ✅ **Path spaces:** Use quotes around paths with spaces

**Testing Windows Setup:**

Contact **Sarika** or **Paul** - they use Windows machines and can help evaluate Windows-specific issues.

---

### macOS Setup

#### Differences on macOS:
1. **Library extension:** `.dylib` instead of `.so`
2. **Installation paths:** `/opt/voltage/` (same as Linux)
3. **Security:** May need to allow libraries in System Preferences

#### Mac-Specific Configuration:

**CGO Configuration for Mac:**
```bash
# Set library paths
export CGO_LDFLAGS="-L/opt/voltage/simpleapi/lib -lfiservprotect -lvibesimple"
export DYLD_LIBRARY_PATH=/opt/voltage/simpleapi/lib:$DYLD_LIBRARY_PATH

# Build
go build
```

#### Common macOS Issues:
- ✅ **Library not loaded:** Check `DYLD_LIBRARY_PATH`
- ✅ **Untrusted developer:** Go to System Preferences → Security & Privacy
- ✅ **Permission denied:** Use `sudo` for installation

---

### Linux Setup

#### Standard Linux Configuration:

**CGO Configuration for Linux:**
```bash
# Set library paths
export CGO_LDFLAGS="-L/opt/voltage/simpleapi/lib -lfiservprotect -lvibesimple"
export LD_LIBRARY_PATH=/opt/voltage/simpleapi/lib:$LD_LIBRARY_PATH

# Build
go build
```

#### Common Linux Issues:
- ✅ **Library not found:** Check `LD_LIBRARY_PATH`
- ✅ **Permission denied:** Use `sudo` for installation
- ✅ **Missing gcc:** Install build-essential package

---

## Quick Start for Developers

### Step 1: Install Dependencies

```bash
# 1. Install Go (if not already installed)
# Download from golang.org

# 2. Enable CGO
export CGO_ENABLED=1  # Linux/Mac
$env:CGO_ENABLED=1    # Windows

# 3. Install Voltage libraries (contact Voltage team)
```

### Step 2: Get Configuration

```bash
# Option A: Request from Voltage platform team
# - voltageprotector.cfg
# - vsconfig.xml
# - cert.pfx
# - trustStore files

# Option B: Copy from team shared location
# (if your team has a shared config repo)
```

### Step 3: Install VLock

```bash
go get github.com/the_monkeys/vlock
```

### Step 4: Test Your Setup

```go
package main

import (
    "fmt"
    "log"
    "github.com/the_monkeys/vlock"
)

func main() {
    // Initialize with your config
    config, err := vlock.LoadConfig("/path/to/voltageprotector.cfg")
    if err != nil {
        log.Fatal("Config error:", err)
    }
    
    fmt.Printf("Loaded config for: %s (%s)\n", 
        config.AppName, config.AppEnv)
    fmt.Println("✅ Setup successful!")
}
```

Run:
```bash
go run main.go
```

Expected output:
```
Loaded config for: VLockDev (DEV)
✅ Setup successful!
```

---

## Testing Your Setup

### 1. Verify Configuration Loading

```bash
# Run configuration tests
cd /path/to/vlock
go test -v -run TestLoadConfig
```

### 2. Check Environment Variables

```bash
# Linux/Mac
env | grep FP_

# Windows PowerShell
Get-ChildItem Env: | Where-Object {$_.Name -like "FP_*"}
```

### 3. Verify Library Paths

```bash
# Linux
ldd /opt/voltage/simpleapi/lib/libfiservprotect.so

# Mac
otool -L /opt/voltage/simpleapi/lib/libfiservprotect.dylib

# Windows
# Use Dependency Walker or similar tool
```

---

## Troubleshooting

### Issue: "cannot find -lfiservprotect"

**Solution:**
```bash
# Add library path to CGO_LDFLAGS
export CGO_LDFLAGS="-L/opt/voltage/simpleapi/lib -lfiservprotect -lvibesimple"
```

### Issue: "AppName is required"

**Solution:** Ensure config file has required fields or set environment variables:
```bash
export FP_APPNAME=YourAppName
export FP_APPVERSION=1.0.0
export FP_APPENV=DEV
```

### Issue: "Error loading shared library"

**Solution:**
```bash
# Linux
export LD_LIBRARY_PATH=/opt/voltage/simpleapi/lib:$LD_LIBRARY_PATH

# Mac
export DYLD_LIBRARY_PATH=/opt/voltage/simpleapi/lib:$DYLD_LIBRARY_PATH

# Windows
# Add C:\voltage\simpleapi\lib to PATH
```

### Issue: "Permission denied" on config files

**Solution:**
```bash
chmod 600 /path/to/voltageprotector.cfg
chmod 600 /path/to/cert.pfx
```

---

## Security Best Practices

### 1. Protect Configuration Files

```bash
# Set restrictive permissions
chmod 600 voltageprotector.cfg
chmod 600 cert.pfx
```

### 2. Never Commit Credentials

- ✅ Use `.gitignore` to exclude config files
- ✅ Store real credentials in secure vault
- ✅ Use environment variables for CI/CD

### 3. Separate Environments

- ✅ Use different configs for DEV/QA/PROD
- ✅ Never mix credentials across environments
- ✅ Test in DEV before promoting to higher environments

---

## Getting Help

### Internal Resources
- **Voltage Platform Team** - For library installation and credentials
- **Security Team** - For certificate and secret management
- **Sarika & Paul** - Windows-specific issues
- **@John Farley** - Architecture and design questions

### Documentation
- Main README.md - Usage and API reference
- config/README.md - Configuration guide
- IMPLEMENTATION_PLAN.md - Full implementation details

### Common Questions

**Q: Where do I get the Voltage libraries?**
A: Contact your Voltage platform team or check internal distribution.

**Q: Can I use this without the Voltage C library?**
A: No, the C library is required. VLock is a wrapper around it.

**Q: How do I get config for my environment?**
A: Request from Voltage platform team with your app name and environment.

**Q: Is this tested on Windows?**
A: Yes, contact Sarika or Paul who use Windows machines.

---

## Next Steps

1. ✅ Complete setup following this guide
2. ✅ Run configuration tests
3. ✅ Review main README.md for API usage
4. ✅ Check config/README.md for configuration details
5. ✅ Start integrating VLock into your application

---

**Need help?** Reach out to @John Farley or the Voltage platform team.

**Found an issue?** Create a ticket in the Finxact Architecture JIRA.
