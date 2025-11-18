package vlock

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()

	if config.NetworkTimeout != 10 {
		t.Errorf("Expected NetworkTimeout to be 10, got %d", config.NetworkTimeout)
	}

	if config.DisableCRLChecking != false {
		t.Errorf("Expected DisableCRLChecking to be false, got %v", config.DisableCRLChecking)
	}

	if config.LogLevel != 2 {
		t.Errorf("Expected LogLevel to be 2, got %d", config.LogLevel)
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	// Create temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.cfg")

	configContent := `# Test Configuration
[ProtectorConfig]
fp_appName=TestApp
fp_appVersion=1.0.0
fp_appEnv=DEV
fp_simpleAPI_installPath=/opt/voltage
fp_trustStore_path=/opt/truststore
XMLConfig=./vsconfig.xml
fp_default_sharedSecret=test_secret
fp_networkTimeout=15
fp_disableCRLChecking=true
DefaultCryptId=TEST_ID
LogLevel=3
LogFile=/tmp/test.log
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify all fields were loaded correctly
	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{"AppName", config.AppName, "TestApp"},
		{"AppVersion", config.AppVersion, "1.0.0"},
		{"AppEnv", config.AppEnv, "DEV"},
		{"SimpleAPIInstallPath", config.SimpleAPIInstallPath, "/opt/voltage"},
		{"TrustStorePath", config.TrustStorePath, "/opt/truststore"},
		{"XMLConfigPath", config.XMLConfigPath, "./vsconfig.xml"},
		{"DEKSharedSecret", config.DEKSharedSecret, "test_secret"},
		{"NetworkTimeout", config.NetworkTimeout, 15},
		{"DisableCRLChecking", config.DisableCRLChecking, true},
		{"DefaultCryptID", config.DefaultCryptID, "TEST_ID"},
		{"LogLevel", config.LogLevel, 3},
		{"LogFile", config.LogFile, "/tmp/test.log"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("Expected %s to be %v, got %v", tt.name, tt.expected, tt.got)
			}
		})
	}
}

func TestLoadConfigFromEnv(t *testing.T) {
	// Set environment variables
	envVars := map[string]string{
		"FP_APPNAME":               "EnvApp",
		"FP_APPVERSION":            "2.0.0",
		"FP_APPENV":                "PROD",
		"FP_SIMPLEAPI_INSTALLPATH": "/env/voltage",
		"FP_TRUSTSTORE_PATH":       "/env/truststore",
		"FP_KEK_CERTPATH":          "/env/cert.pfx",
		"FP_KEK_CERTPASSPHRASE":    "env_passphrase",
		"FP_KEK_SHAREDSECRET":      "env_kek_secret",
		"FP_DEFAULT_SHAREDSECRET":  "env_dek_secret",
		"FP_DEFAULT_USERNAME":      "env_user",
		"FP_DEFAULT_PASSWORD":      "env_pass",
		"FP_NETWORKTIMEOUT":        "25",
		"FP_DISABLECRLCHECKING":    "true",
	}

	// Set environment variables
	for key, value := range envVars {
		os.Setenv(key, value)
	}

	// Clean up after test
	defer func() {
		for key := range envVars {
			os.Unsetenv(key)
		}
	}()

	// Load config without a file
	config, err := LoadConfig("")
	if err != nil {
		t.Fatalf("Failed to load config from env: %v", err)
	}

	// Verify environment variables were loaded
	if config.AppName != "EnvApp" {
		t.Errorf("Expected AppName to be EnvApp, got %s", config.AppName)
	}
	if config.AppVersion != "2.0.0" {
		t.Errorf("Expected AppVersion to be 2.0.0, got %s", config.AppVersion)
	}
	if config.AppEnv != "PROD" {
		t.Errorf("Expected AppEnv to be PROD, got %s", config.AppEnv)
	}
	if config.SimpleAPIInstallPath != "/env/voltage" {
		t.Errorf("Expected SimpleAPIInstallPath to be /env/voltage, got %s", config.SimpleAPIInstallPath)
	}
	if config.KEKCertPath != "/env/cert.pfx" {
		t.Errorf("Expected KEKCertPath to be /env/cert.pfx, got %s", config.KEKCertPath)
	}
	if config.NetworkTimeout != 25 {
		t.Errorf("Expected NetworkTimeout to be 25, got %d", config.NetworkTimeout)
	}
}

func TestEnvVarPrecedence(t *testing.T) {
	// Create temporary config file with one value
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.cfg")

	configContent := `
fp_appName=FileApp
fp_appVersion=1.0.0
fp_appEnv=DEV
fp_default_sharedSecret=file_secret
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Set environment variable to override
	os.Setenv("FP_APPNAME", "EnvApp")
	os.Setenv("FP_DEFAULT_SHAREDSECRET", "env_secret")
	defer func() {
		os.Unsetenv("FP_APPNAME")
		os.Unsetenv("FP_DEFAULT_SHAREDSECRET")
	}()

	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Environment variable should override file value
	if config.AppName != "EnvApp" {
		t.Errorf("Expected AppName from env (EnvApp), got %s from file", config.AppName)
	}
	if config.DEKSharedSecret != "env_secret" {
		t.Errorf("Expected DEKSharedSecret from env (env_secret), got %s from file", config.DEKSharedSecret)
	}
	// File value should be used for non-overridden values
	if config.AppVersion != "1.0.0" {
		t.Errorf("Expected AppVersion from file (1.0.0), got %s", config.AppVersion)
	}
}

func TestValidateRequiredFields(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		shouldError bool
		errorMsg    string
	}{
		{
			name: "Valid configuration",
			config: &Config{
				AppName:         "TestApp",
				AppVersion:      "1.0.0",
				AppEnv:          "DEV",
				DEKSharedSecret: "secret",
			},
			shouldError: false,
		},
		{
			name: "Missing AppName",
			config: &Config{
				AppVersion:      "1.0.0",
				AppEnv:          "DEV",
				DEKSharedSecret: "secret",
			},
			shouldError: true,
			errorMsg:    "AppName is required",
		},
		{
			name: "Missing AppVersion",
			config: &Config{
				AppName:         "TestApp",
				AppEnv:          "DEV",
				DEKSharedSecret: "secret",
			},
			shouldError: true,
			errorMsg:    "AppVersion is required",
		},
		{
			name: "Missing AppEnv",
			config: &Config{
				AppName:         "TestApp",
				AppVersion:      "1.0.0",
				DEKSharedSecret: "secret",
			},
			shouldError: true,
			errorMsg:    "AppEnv is required",
		},
		{
			name: "Invalid AppEnv",
			config: &Config{
				AppName:         "TestApp",
				AppVersion:      "1.0.0",
				AppEnv:          "INVALID",
				DEKSharedSecret: "secret",
			},
			shouldError: true,
			errorMsg:    "AppEnv must be one of",
		},
		{
			name: "Missing authentication",
			config: &Config{
				AppName:    "TestApp",
				AppVersion: "1.0.0",
				AppEnv:     "DEV",
			},
			shouldError: true,
			errorMsg:    "authentication method",
		},
		{
			name: "Valid with KEK cert",
			config: &Config{
				AppName:     "TestApp",
				AppVersion:  "1.0.0",
				AppEnv:      "PROD",
				KEKCertPath: "/path/to/cert.pfx",
			},
			shouldError: false,
		},
		{
			name: "Valid with KEK shared secret",
			config: &Config{
				AppName:         "TestApp",
				AppVersion:      "1.0.0",
				AppEnv:          "QA",
				KEKSharedSecret: "kek_secret",
			},
			shouldError: false,
		},
		{
			name: "Valid with DEK username/password",
			config: &Config{
				AppName:     "TestApp",
				AppVersion:  "1.0.0",
				AppEnv:      "CAT",
				DEKUsername: "user",
				DEKPassword: "pass",
			},
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.shouldError {
				if err == nil {
					t.Error("Expected validation error, got nil")
				} else if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
			}
		})
	}
}

func TestConfigString(t *testing.T) {
	config := &Config{
		AppName:              "TestApp",
		AppVersion:           "1.0.0",
		AppEnv:               "DEV",
		SimpleAPIInstallPath: "/opt/voltage",
		TrustStorePath:       "/opt/trust",
		XMLConfigPath:        "./config.xml",
		NetworkTimeout:       15,
		KEKCertPassphrase:    "secret123", // Should not appear in string
		DEKSharedSecret:      "secret456", // Should not appear in string
	}

	str := config.String()

	// Verify important fields are present
	if !contains(str, "TestApp") {
		t.Error("AppName should be in string representation")
	}
	if !contains(str, "1.0.0") {
		t.Error("AppVersion should be in string representation")
	}
	if !contains(str, "DEV") {
		t.Error("AppEnv should be in string representation")
	}

	// Verify secrets are not in string representation
	if contains(str, "secret123") {
		t.Error("KEKCertPassphrase should not be in string representation")
	}
	if contains(str, "secret456") {
		t.Error("DEKSharedSecret should not be in string representation")
	}
}

func TestIsProduction(t *testing.T) {
	tests := []struct {
		env      string
		expected bool
	}{
		{"PROD", true},
		{"DEV", false},
		{"QA", false},
		{"CAT", false},
	}

	for _, tt := range tests {
		t.Run(tt.env, func(t *testing.T) {
			config := &Config{AppEnv: tt.env}
			if config.IsProduction() != tt.expected {
				t.Errorf("Expected IsProduction() to be %v for %s", tt.expected, tt.env)
			}
		})
	}
}

func TestGetEnvironment(t *testing.T) {
	config := &Config{AppEnv: "QA"}
	if config.GetEnvironment() != "QA" {
		t.Errorf("Expected GetEnvironment() to return QA, got %s", config.GetEnvironment())
	}
}

func TestConfigError(t *testing.T) {
	tests := []struct {
		name     string
		err      *ConfigError
		expected string
	}{
		{
			name:     "Error with field",
			err:      &ConfigError{Field: "AppName", Message: "is required"},
			expected: "configuration error [AppName]: is required",
		},
		{
			name:     "Error without field",
			err:      &ConfigError{Message: "general error"},
			expected: "configuration error: general error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.expected {
				t.Errorf("Expected error message '%s', got '%s'", tt.expected, tt.err.Error())
			}
		})
	}
}

func TestLoadConfigNonExistentFile(t *testing.T) {
	// Try to load config from non-existent file with required env vars
	os.Setenv("FP_APPNAME", "TestApp")
	os.Setenv("FP_APPVERSION", "1.0.0")
	os.Setenv("FP_APPENV", "DEV")
	os.Setenv("FP_DEFAULT_SHAREDSECRET", "secret")
	defer func() {
		os.Unsetenv("FP_APPNAME")
		os.Unsetenv("FP_APPVERSION")
		os.Unsetenv("FP_APPENV")
		os.Unsetenv("FP_DEFAULT_SHAREDSECRET")
	}()

	// Should succeed with env vars even if file doesn't exist
	config, err := LoadConfig("/nonexistent/path/config.cfg")
	if err != nil {
		t.Fatalf("Should succeed with env vars even if file doesn't exist: %v", err)
	}

	if config.AppName != "TestApp" {
		t.Errorf("Expected AppName from env, got %s", config.AppName)
	}
}

func TestAppEnvCaseInsensitive(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.cfg")

	configContent := `
fp_appName=TestApp
fp_appVersion=1.0.0
fp_appEnv=dev
fp_default_sharedSecret=secret
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Should be normalized to uppercase
	if config.AppEnv != "DEV" {
		t.Errorf("Expected AppEnv to be normalized to DEV, got %s", config.AppEnv)
	}
}

func TestCommentsParsing(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.cfg")

	configContent := `# This is a comment
fp_appName=TestApp
; This is also a comment
fp_appVersion=1.0.0
# Another comment
fp_appEnv=DEV
fp_default_sharedSecret=secret
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if config.AppName != "TestApp" {
		t.Error("Failed to parse config with comments")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
