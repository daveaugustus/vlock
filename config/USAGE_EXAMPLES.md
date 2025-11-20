# Config Package Usage Examples

## Quick Reference

### Basic Usage

```go
package main

import (
    "log"
    "github.com/daveaugustus/vlock/config"
)

func main() {
    // Load configuration
    cfg, err := config.LoadConfig("./voltageprotector.cfg")
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Config loaded: %s", cfg.String())
}
```

### Loading Strategies

#### 1. Environment Variables Only (Recommended for Containers)
```go
// Set environment variables first
os.Setenv("FP_APPNAME", "MyApp")
os.Setenv("FP_APPVERSION", "1.0.0")
os.Setenv("FP_APPENV", "PROD")
os.Setenv("FP_DEFAULT_SHAREDSECRET", "secret")

cfg, err := config.LoadConfig("") // Empty path = env vars only
```

#### 2. Config File Only
```go
cfg, err := config.LoadConfig("/path/to/voltageprotector.cfg")
```

#### 3. File + Environment Override (Best Practice)
```go
// File provides base config, environment variables override specific values
cfg, err := config.LoadConfig("./voltageprotector.cfg")
// FP_APPNAME env var will override fp_appName from file if set
```

### Environment Checks

```go
if cfg.IsProduction() {
    // Enable production features
    setupMonitoring()
    disableDebugMode()
}

switch cfg.GetEnvironment() {
case "DEV":
    // Development settings
case "QA":
    // QA settings
case "CAT":
    // CAT settings  
case "PROD":
    // Production settings
}
```

### Validation

```go
// Manual validation
cfg := config.NewConfig()
cfg.AppName = "TestApp"
cfg.AppVersion = "1.0.0"
cfg.AppEnv = "DEV"

if err := cfg.Validate(); err != nil {
    // Handle validation error
    log.Fatal(err)
}
```

## Integration Example

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/daveaugustus/vlock/config"
)

type VoltageService struct {
    config *config.Config
}

func NewVoltageService(cfgPath string) (*VoltageService, error) {
    cfg, err := config.LoadConfig(cfgPath)
    if err != nil {
        return nil, fmt.Errorf("failed to load config: %w", err)
    }
    
    return &VoltageService{
        config: cfg,
    }, nil
}

func (s *VoltageService) Initialize() error {
    log.Printf("Initializing Voltage service for: %s", s.config.AppName)
    log.Printf("Environment: %s", s.config.GetEnvironment())
    log.Printf("Library path: %s", s.config.SimpleAPIInstallPath)
    
    // Initialize Voltage C library here...
    
    return nil
}

func (s *VoltageService) Encrypt(data string) (string, error) {
    if s.config.IsProduction() {
        log.Println("Using production encryption settings")
    }
    
    // Encryption logic here...
    
    return "", nil
}

func main() {
    svc, err := NewVoltageService("./voltageprotector.cfg")
    if err != nil {
        log.Fatal(err)
    }
    
    if err := svc.Initialize(); err != nil {
        log.Fatal(err)
    }
    
    // Use the service...
}
```

## Testing Your Code with Config

```go
package myapp

import (
    "testing"
    "github.com/daveaugustus/vlock/config"
)

func TestMyFunction(t *testing.T) {
    // Create test configuration
    cfg := config.NewConfig()
    cfg.AppName = "TestApp"
    cfg.AppVersion = "1.0.0"
    cfg.AppEnv = "DEV"
    cfg.DEKSharedSecret = "test_secret"
    
    if err := cfg.Validate(); err != nil {
        t.Fatal(err)
    }
    
    // Use cfg in your tests...
}
```

## See Also

- Run tests: `go test -v` in the config directory
- Full test program: `test_config_main.go` in project root
- Unit tests: `config_test.go` (11 test functions, 95.2% coverage)
