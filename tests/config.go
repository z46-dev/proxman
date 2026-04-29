package tests

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
)

type configuration struct {
	EnableProxmoxTesting bool `toml:"enable_proxmox_testing" default:"false"` // Enable Proxmox VE API testing (requires valid Proxmox VE API credentials)
	ProxmoxTesting       struct {
		Hostname string `toml:"hostname" default:""` // Proxmox VE server hostname or IP address (e.g. "proxmox.cyber.lab")
		Port     string `toml:"port" default:""`     // Proxmox VE API port (usually "8006")
		TokenID  string `toml:"token_id" default:""` // Proxmox VE API token ID (e.g. "api-token-id")
		Secret   string `toml:"secret" default:""`   // Proxmox VE API token secret
	} `toml:"proxmox_testing"` // Proxmox VE API testing configuration
}

var config configuration

func loadConfig(path string) (err error) {
	// Apply struct defaults BEFORE loading TOML (so TOML overrides)
	if err = defaults.Set(&config); err != nil {
		err = fmt.Errorf("set defaults: %w", err)
		return
	}

	// Decode TOML file into struct
	if _, err = toml.DecodeFile(path, &config); err != nil {
		err = fmt.Errorf("decode toml: %w", err)
		return
	}

	// Validate required fields
	if err = validator.New(validator.WithRequiredStructEnabled()).Struct(config); err != nil {
		err = fmt.Errorf("validate config: %w", err)
		return
	}

	return
}

// generateDefaultConfig writes a config.toml with all default values filled in.
// It will overwrite any existing file at path.
func generateDefaultConfig(path string) (err error) {
	var cfg configuration

	// 1. Apply struct defaults
	if err = defaults.Set(&cfg); err != nil {
		err = fmt.Errorf("set defaults: %w", err)
		return
	}

	// NOTE: Do NOT validate here.
	// The default config is allowed to be "invalid" from a required-fields POV;
	// it's just a template for the user to fill in.
	// Validation happens in LoadConfig() when we actually load the file.

	// 2. Create / truncate the file
	var file *os.File
	if file, err = os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644); err != nil {
		err = fmt.Errorf("create config file: %w", err)
		return
	}

	defer file.Close()

	// 3. Encode as TOML
	var encoder *toml.Encoder = toml.NewEncoder(file)
	encoder.Indent = "    "
	if err = encoder.Encode(cfg); err != nil {
		err = fmt.Errorf("encode toml: %w", err)
	}

	return
}

func initConfig(path string) (err error) {
	if !filepath.IsAbs(path) {
		if path, err = filepath.Abs(path); err != nil {
			return err
		}
	}

	if _, err = os.Stat(path); err != nil {
		if err = generateDefaultConfig(path); err != nil {
			return
		}

		err = fmt.Errorf("no config file found, created a default config at %s. Please fill in the required values and try again", path)
		return
	}

	if err = loadConfig(path); err != nil {
		return err
	}

	return nil
}
