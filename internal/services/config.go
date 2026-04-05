package services

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Server struct {
	Host       string `json:"host"`
	User       string `json:"user"`
	PrivateKey string `json:"private_key"`
	Port       string `json:"port"`
}

type ConfigService struct {
	path    string
	servers map[string]Server
}

func NewConfigService() (*ConfigService, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(configDir, "filepass")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}

	path := filepath.Join(dir, "servers.json")

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		if err := os.WriteFile(path, []byte("{}"), 0o600); err != nil {
			return nil, err
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var servers map[string]Server
	if err := json.Unmarshal(data, &servers); err != nil {
		return nil, err
	}

	return &ConfigService{path: path, servers: servers}, nil
}

func (c *ConfigService) Servers() map[string]Server {
	return c.servers
}
