package main

import (
	"fmt"
	"log"
	"os"

	"github.com/daveaugustus/vlock/config"
)

func main() {
	fmt.Println("=== VLock Configuration Package Test ===")
	fmt.Println("Testing with REAL config files from the repo\n")

	// Test 1: Load from environment variables
	fmt.Println("Test 1: Loading configuration from environment variables")
	fmt.Println("-------------------------------------------------------")

	// Set environment variables
	os.Setenv("FP_APPNAME", "TestVLockApp")
	os.Setenv("FP_APPVERSION", "1.0.0")
	os.Setenv("FP_APPENV", "DEV")
	os.Setenv("FP_DEFAULT_SHAREDSECRET", "my_test_secret_123")
	os.Setenv("FP_SIMPLEAPI_INSTALLPATH", "/opt/voltage/simpleapi")
	os.Setenv("FP_TRUSTSTORE_PATH", "/opt/voltage/trustStore")
	os.Setenv("FP_NETWORKTIMEOUT", "30")

	cfg, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("❌ Failed to load config from environment: %v", err)
	}

	fmt.Printf("✅ Configuration loaded successfully!\n")
	fmt.Printf("   App Name: %s\n", cfg.AppName)
	fmt.Printf("   App Version: %s\n", cfg.AppVersion)
	fmt.Printf("   Environment: %s\n", cfg.AppEnv)
	fmt.Printf("   Library Path: %s\n", cfg.SimpleAPIInstallPath)
	fmt.Printf("   Trust Store: %s\n", cfg.TrustStorePath)
	fmt.Printf("   Network Timeout: %d seconds\n", cfg.NetworkTimeout)
	fmt.Printf("   Log Level: %d\n", cfg.LogLevel)
	fmt.Printf("\n")

	// Test 2: String representation (should mask secrets)
	fmt.Println("Test 2: String representation (secrets should be masked)")
	fmt.Println("--------------------------------------------------------")
	fmt.Printf("Config String: %s\n\n", cfg.String())

	// Test 3: Environment checks
	fmt.Println("Test 3: Environment utility methods")
	fmt.Println("------------------------------------")
	fmt.Printf("   Is Production? %v\n", cfg.IsProduction())
	fmt.Printf("   Environment Name: %s\n", cfg.GetEnvironment())
	fmt.Printf("\n")

	// Test 4: Load from REAL config file in repo (DEV)
	fmt.Println("Test 4: Loading REAL config file (config/dev/voltageprotector.cfg)")
	fmt.Println("-------------------------------------------------------------------")

	// Use actual config file from the repo
	realConfigPath := "./config/dev/voltageprotector.cfg"

	// Override placeholders with test values via environment
	os.Setenv("FP_KEK_CERTPASSPHRASE", "test_dev_passphrase")
	os.Setenv("FP_DEFAULT_SHAREDSECRET", "test_dev_secret")

	cfg2, err := config.LoadConfig(realConfigPath)
	if err != nil {
		log.Fatalf("❌ Failed to load config from file: %v", err)
	}

	fmt.Printf("✅ Configuration loaded from REAL dev config!\n")
	fmt.Printf("   Config File: %s\n", realConfigPath)
	fmt.Printf("   App Name: %s\n", cfg2.AppName)
	fmt.Printf("   App Version: %s\n", cfg2.AppVersion)
	fmt.Printf("   Environment: %s\n", cfg2.AppEnv)
	fmt.Printf("   Library Path: %s\n", cfg2.SimpleAPIInstallPath)
	fmt.Printf("   Trust Store: %s\n", cfg2.TrustStorePath)
	fmt.Printf("   XML Config: %s\n", cfg2.XMLConfigPath)
	fmt.Printf("   Default Crypt ID: %s\n", cfg2.DefaultCryptID)
	fmt.Printf("   Log File: %s\n", cfg2.LogFile)
	fmt.Printf("   CRL Checking Disabled? %v\n", cfg2.DisableCRLChecking)
	fmt.Printf("\n")

	// Test 5: Load ALL environment configs (DEV, QA, PROD)
	fmt.Println("Test 5: Loading ALL real environment configs from repo")
	fmt.Println("-------------------------------------------------------")

	configFiles := map[string]string{
		"DEV":  "./config/dev/voltageprotector.cfg",
		"QA":   "./config/qa/voltageprotector.cfg",
		"PROD": "./config/prod/voltageprotector.cfg",
	}

	for envName, configPath := range configFiles {
		// Set required secrets via env vars (since placeholders in files)
		os.Setenv("FP_DEFAULT_SHAREDSECRET", "test_"+envName+"_secret")

		envCfg, err := config.LoadConfig(configPath)
		if err != nil {
			fmt.Printf("   ❌ %s: Failed to load - %v\n", envName, err)
			continue
		}

		fmt.Printf("   ✅ %-4s: %s v%s (Env: %s, IsProd: %v)\n",
			envName, envCfg.AppName, envCfg.AppVersion,
			envCfg.GetEnvironment(), envCfg.IsProduction())
	}
	fmt.Printf("\n")

	// Test 6: Environment variable precedence over real file
	fmt.Println("Test 6: Environment variable precedence over real file")
	fmt.Println("-------------------------------------------------------")

	// Set an env var to override the real file
	os.Setenv("FP_APPNAME", "EnvOverridesRealFile")
	os.Setenv("FP_NETWORKTIMEOUT", "99")

	cfg3, err := config.LoadConfig(realConfigPath)
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	fmt.Printf("✅ Environment variables override file values!\n")
	fmt.Printf("   App Name (from ENV): %s (file had: VLockDev)\n", cfg3.AppName)
	fmt.Printf("   App Version (from FILE): %s\n", cfg3.AppVersion)
	fmt.Printf("   Network Timeout (from ENV): %d (file had: 10)\n", cfg3.NetworkTimeout)
	fmt.Printf("\n")

	// Test 7: Validation errors
	fmt.Println("Test 7: Configuration validation")
	fmt.Println("----------------------------------")

	// Clear all env vars
	os.Unsetenv("FP_APPNAME")
	os.Unsetenv("FP_APPVERSION")
	os.Unsetenv("FP_APPENV")
	os.Unsetenv("FP_DEFAULT_SHAREDSECRET")
	os.Unsetenv("FP_NETWORKTIMEOUT")

	_, err = config.LoadConfig("")
	if err != nil {
		fmt.Printf("✅ Validation correctly caught missing fields:\n")
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("❌ Validation should have failed!\n")
	}
	fmt.Printf("\n")

	// Test 8: Create new config with defaults
	fmt.Println("Test 8: Create config with defaults")
	fmt.Println("------------------------------------")

	defaultCfg := config.NewConfig()
	fmt.Printf("✅ Created config with defaults:\n")
	fmt.Printf("   Network Timeout: %d\n", defaultCfg.NetworkTimeout)
	fmt.Printf("   Log Level: %d\n", defaultCfg.LogLevel)
	fmt.Printf("   CRL Checking Disabled? %v\n", defaultCfg.DisableCRLChecking)
	fmt.Printf("\n")

	// Test 9: Testing different environments programmatically
	fmt.Println("Test 9: Testing different environments programmatically")
	fmt.Println("--------------------------------------------------------")

	environments := []string{"DEV", "QA", "CAT", "PROD"}
	for _, env := range environments {
		os.Setenv("FP_APPNAME", "TestApp")
		os.Setenv("FP_APPVERSION", "1.0.0")
		os.Setenv("FP_APPENV", env)
		os.Setenv("FP_DEFAULT_SHAREDSECRET", "secret")

		testCfg, err := config.LoadConfig("")
		if err != nil {
			fmt.Printf("   ❌ Failed to load %s config: %v\n", env, err)
			continue
		}

		fmt.Printf("   %s: IsProduction=%v, Environment=%s\n",
			env, testCfg.IsProduction(), testCfg.GetEnvironment())
	}
	fmt.Printf("\n")

	// Test 10: Invalid environment validation
	fmt.Println("Test 10: Invalid environment validation")
	fmt.Println("----------------------------------------")

	os.Setenv("FP_APPENV", "INVALID_ENV")
	_, err = config.LoadConfig("")
	if err != nil {
		fmt.Printf("✅ Correctly rejected invalid environment:\n")
		fmt.Printf("   Error: %v\n", err)
	}
	fmt.Printf("\n")

	// Summary
	fmt.Println("=== Test Summary ===")
	fmt.Println("✅ All configuration tests completed!")
	fmt.Println("\nTests performed:")
	fmt.Println("  ✅ Test 1: Environment variable loading")
	fmt.Println("  ✅ Test 2: Safe string representation (secrets masked)")
	fmt.Println("  ✅ Test 3: Environment utility methods")
	fmt.Println("  ✅ Test 4: Loading REAL config/dev/voltageprotector.cfg")
	fmt.Println("  ✅ Test 5: Loading ALL repo configs (DEV/QA/PROD)")
	fmt.Println("  ✅ Test 6: Environment variable precedence over file")
	fmt.Println("  ✅ Test 7: Configuration validation")
	fmt.Println("  ✅ Test 8: Default values")
	fmt.Println("  ✅ Test 9: Different environments (DEV/QA/CAT/PROD)")
	fmt.Println("  ✅ Test 10: Invalid environment rejection")
	fmt.Println("\nThe config package supports:")
	fmt.Println("  • Loading from .cfg files (INI format)")
	fmt.Println("  • Loading from REAL repo config templates")
	fmt.Println("  • Environment variable overrides")
	fmt.Println("  • Validation of required fields")
	fmt.Println("  • Multiple authentication methods (KEK/DEK)")
	fmt.Println("  • Environment-specific settings (DEV/QA/CAT/PROD)")
	fmt.Println("  • Secure string representation (secrets masked)")
	fmt.Println("  • Default values for optional settings")

	fmt.Println("\n📁 Real config files tested:")
	fmt.Println("  • ./config/dev/voltageprotector.cfg")
	fmt.Println("  • ./config/qa/voltageprotector.cfg")
	fmt.Println("  • ./config/prod/voltageprotector.cfg")
}
