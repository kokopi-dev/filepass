package services

import (
	"fmt"
	"strings"
)

// StorageService executes file operations against a single server's storage.
type StorageService struct {
	server Server
}

func NewStorageService(s Server) *StorageService {
	return &StorageService{server: s}
}

// Check returns the list of files currently in the remote storage directory.
func (s *StorageService) Check() ([]string, error) {
	cmd := SSHCmd(s.server,
		"find "+defaultStoragePath+" -type f -printf '%f\n' 2>/dev/null",
	)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("check failed: %w", err)
	}
	raw := strings.TrimSpace(string(out))
	if raw == "" {
		return []string{}, nil
	}
	return strings.Split(raw, "\n"), nil
}

// Send transfers one or more local files to the remote storage.
// Multiple files are archived into a temp tarball first.
func (s *StorageService) Send(localPaths []string) error {
	// TODO: implement
	return fmt.Errorf("send: not yet implemented")
}

// Get downloads one or more files from remote storage to destDir.
// Multiple files are archived server-side, transferred, then extracted.
func (s *StorageService) Get(remoteFiles []string, destDir string) error {
	// TODO: implement
	return fmt.Errorf("get: not yet implemented")
}

// Clean removes specific files from remote storage.
// Pass a nil or empty slice to remove all files.
func (s *StorageService) Clean(remoteFiles []string) error {
	// TODO: implement
	return fmt.Errorf("clean: not yet implemented")
}
