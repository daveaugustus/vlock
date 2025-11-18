package vlock

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config represents the complete configuration for the Voltage Protector library
type Config struct {
	// Application settings (Required)
	AppName    string
	AppVersion string
	AppEnv     string // DEV, QA, CAT, PROD

	// Paths (Required for most deployments)
	SimpleAPIInstallPath string
	TrustStorePath       string
	XMLConfigPath        string

	// KEK (Key Encryption Key) settings
	KEKCertPath       string
	KEKCertPassphrase string
	KEKSharedSecret   string

	// DEK (Data Encryption Key) settings
	DEKSharedSecret string
	DEKUsername     string
	DEKPassword     string

	// Optional settings
	NetworkTimeout     int
	DisableCRLChecking bool
	DefaultCryptID     string
	LogLevel           int
	LogFile            string

	// Internal
	ConfigFilePath string
}

// Environment variable names mapped to configuration fields
const (
	EnvAppName              = "FP_APPNAME"
	EnvAppVersion           = "FP_APPVERSION"
	EnvAppEnv               = "FP_APPENV"
	EnvNetworkTimeout       = "FP_NETWORKTIMEOUT"
	EnvDisableCRLChecking   = "FP_DISABLECRLCHECKING"
	EnvSimpleAPIInstallPath = "FP_SIMPLEAPI_INSTALLPATH"
	EnvTrustStorePath       = "FP_TRUSTSTORE_PATH"
	EnvKEKCertPath          = "FP_KEK_CERTPATH"
	EnvKEKCertPassphrase    = "FP_KEK_CERTPASSPHRASE"
	EnvKEKSharedSecret      = "FP_KEK_SHAREDSECRET"
	EnvDEKSharedSecret      = "FP_DEFAULT_SHAREDSECRET"
	EnvDEKUsername          = "FP_DEFAULT_USERNAME"
	EnvDEKPassword          = "FP_DEFAULT_PASSWORD"
)

// ConfigError represents a configuration-related error
type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("configuration error [%s]: %s", e.Field, e.Message)
	}
	return fmt.Sprintf("configuration error: %s", e.Message)
}

// NewConfig creates a new configuration with default values
func NewConfig() *Config {
	return &Config{
		NetworkTimeout:     10,
		DisableCRLChecking: false,
		LogLevel:           2,
	}
}

// LoadConfig loads configuration from a file and applies environment variable overrides
func LoadConfig(configPath string) (*Config, error) {
	config := NewConfig()
	config.ConfigFilePath = configPath

	// Load from file if path is provided and file exists
	if configPath != "" {
		if _, err := os.Stat(configPath); err == nil {
			if err := config.loadFromFile(configPath); err != nil {
				return nil, fmt.Errorf("failed to load config file: %w", err)
			}
		} else if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to access config file: %w", err)
		}
		// If file doesn't exist, continue with defaults and env vars
	}

	// Apply environment variable overrides
	config.loadFromEnv()

	// Validate required fields
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// loadFromFile reads configuration from a .cfg file (INI format)
func (c *Config) loadFromFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		// Skip section headers like [ProtectorConfig]
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			continue
		}

		// Parse key=value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Map configuration keys to struct fields
		c.setConfigValue(key, value)
	}

	return nil
}

// setConfigValue maps a configuration key to the appropriate struct field
func (c *Config) setConfigValue(key, value string) {
	switch key {
	case "fp_appName":
		c.AppName = value
	case "fp_appVersion":
		c.AppVersion = value
	case "fp_appEnv":
		c.AppEnv = strings.ToUpper(value)
	case "fp_simpleAPI_installPath":
		c.SimpleAPIInstallPath = value
	case "fp_trustStore_path":
		c.TrustStorePath = value
	case "XMLConfig":
		c.XMLConfigPath = value
	case "fp_kek_certPath":
		c.KEKCertPath = value
	case "fp_kek_certPassphrase":
		c.KEKCertPassphrase = value
	case "fp_kek_sharedSecret":
		c.KEKSharedSecret = value
	case "fp_default_sharedSecret":
		c.DEKSharedSecret = value
	case "fp_default_userName":
		c.DEKUsername = value
	case "fp_default_password":
		c.DEKPassword = value
	case "DefaultCryptId":
		c.DefaultCryptID = value
	case "LogLevel":
		if val, err := strconv.Atoi(value); err == nil {
			c.LogLevel = val
		}
	case "LogFile":
		c.LogFile = value
	case "fp_networkTimeout":
		if val, err := strconv.Atoi(value); err == nil {
			c.NetworkTimeout = val
		}
	case "fp_disableCRLChecking":
		c.DisableCRLChecking = strings.ToLower(value) == "true"
	}
}

// loadFromEnv applies environment variable overrides to the configuration
// Environment variables take precedence over file-based configuration
func (c *Config) loadFromEnv() {
	if val := os.Getenv(EnvAppName); val != "" {
		c.AppName = val
	}
	if val := os.Getenv(EnvAppVersion); val != "" {
		c.AppVersion = val
	}
	if val := os.Getenv(EnvAppEnv); val != "" {
		c.AppEnv = strings.ToUpper(val)
	}
	if val := os.Getenv(EnvSimpleAPIInstallPath); val != "" {
		c.SimpleAPIInstallPath = val
	}
	if val := os.Getenv(EnvTrustStorePath); val != "" {
		c.TrustStorePath = val
	}
	if val := os.Getenv(EnvKEKCertPath); val != "" {
		c.KEKCertPath = val
	}
	if val := os.Getenv(EnvKEKCertPassphrase); val != "" {
		c.KEKCertPassphrase = val
	}
	if val := os.Getenv(EnvKEKSharedSecret); val != "" {
		c.KEKSharedSecret = val
	}
	if val := os.Getenv(EnvDEKSharedSecret); val != "" {
		c.DEKSharedSecret = val
	}
	if val := os.Getenv(EnvDEKUsername); val != "" {
		c.DEKUsername = val
	}
	if val := os.Getenv(EnvDEKPassword); val != "" {
		c.DEKPassword = val
	}
	if val := os.Getenv(EnvNetworkTimeout); val != "" {
		if timeout, err := strconv.Atoi(val); err == nil {
			c.NetworkTimeout = timeout
		}
	}
	if val := os.Getenv(EnvDisableCRLChecking); val != "" {
		c.DisableCRLChecking = strings.ToLower(val) == "true"
	}
}

// Validate ensures all required configuration parameters are present
func (c *Config) Validate() error {
	var errors []string

	// Required fields
	if c.AppName == "" {
		errors = append(errors, "AppName is required (set fp_appName or FP_APPNAME)")
	}
	if c.AppVersion == "" {
		errors = append(errors, "AppVersion is required (set fp_appVersion or FP_APPVERSION)")
	}
	if c.AppEnv == "" {
		errors = append(errors, "AppEnv is required (set fp_appEnv or FP_APPENV)")
	}

	// Validate AppEnv value
	validEnvs := []string{"DEV", "QA", "CAT", "PROD"}
	if c.AppEnv != "" {
		valid := false
		for _, env := range validEnvs {
			if c.AppEnv == env {
				valid = true
				break
			}
		}
		if !valid {
			errors = append(errors, fmt.Sprintf("AppEnv must be one of: %s (got: %s)", strings.Join(validEnvs, ", "), c.AppEnv))
		}
	}

	// At least one of the authentication methods must be configured
	hasKEK := c.KEKCertPath != "" || c.KEKSharedSecret != ""
	hasDEK := c.DEKSharedSecret != "" || (c.DEKUsername != "" && c.DEKPassword != "")

	if !hasKEK && !hasDEK {
		errors = append(errors, "At least one authentication method must be configured (KEK or DEK credentials)")
	}

	if len(errors) > 0 {
		return &ConfigError{
			Message: strings.Join(errors, "; "),
		}
	}

	return nil
}

// String returns a string representation of the configuration (with sensitive data masked)
func (c *Config) String() string {
	return fmt.Sprintf(
		"Config{AppName: %s, AppVersion: %s, AppEnv: %s, SimpleAPIInstallPath: %s, TrustStorePath: %s, XMLConfigPath: %s, NetworkTimeout: %d}",
		c.AppName,
		c.AppVersion,
		c.AppEnv,
		c.SimpleAPIInstallPath,
		c.TrustStorePath,
		c.XMLConfigPath,
		c.NetworkTimeout,
	)
}

// IsProduction returns true if the configuration is for a production environment
func (c *Config) IsProduction() bool {
	return c.AppEnv == "PROD"
}

// GetEnvironment returns the current environment name
func (c *Config) GetEnvironment() string {
	return c.AppEnv
}
