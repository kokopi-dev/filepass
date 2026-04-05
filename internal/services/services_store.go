package services

import "fmt"

type ServicesStore struct {
	Config *ConfigService
}

func NewServicesStore() (*ServicesStore, error) {
	cfg, err := NewConfigService()
	if err != nil {
		return nil, err
	}
	return &ServicesStore{Config: cfg}, nil
}

func (s *ServicesStore) NewStorageService(serverName string) (*StorageService, error) {
	srv, ok := s.Config.servers[serverName]
	if !ok {
		return nil, fmt.Errorf("server %q not found", serverName)
	}
	return NewStorageService(srv), nil
}
