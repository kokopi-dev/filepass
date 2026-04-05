package services

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
