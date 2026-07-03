package config

import (
	"os"
	"path/filepath"
)

const (
	AppName       = "psw"
	ConfigFile    = "config.json"
	ConfigPerm    = 0700
	FilePerm      = 0600
)

// GetConfigDir returns the application config directory
func GetConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, AppName), nil
}

// GetConfigPath returns the full path to the config file
func GetConfigPath() (string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ConfigFile), nil
}

// EnsureConfigDir creates the config directory if it doesn't exist
func EnsureConfigDir() error {
	dir, err := GetConfigDir()
	if err != nil {
		return err
	}
	return os.MkdirAll(dir, ConfigPerm)
}
