package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// debugLog appends a timestamped line to /tmp/filepass-debug.log.
// Safe to call from any goroutine; errors are silently ignored.
func debugLog(format string, args ...any) {
	f, err := os.OpenFile("/tmp/filepass-debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Fprintf(f, "[%s] %s\n", time.Now().Format("15:04:05.000"), fmt.Sprintf(format, args...))
}

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

// Get downloads a single file from remote storage into destDir.
func (s *StorageService) Get(filename, destDir string) error {
	src := RemotePath(s.server, filename)
	dst := filepath.Join(destDir, filename)
	cmd := RsyncCmd(s.server, src, dst)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("get failed: %w\n%s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// Delete removes a single file from remote storage.
func (s *StorageService) Delete(filename string) error {
	remoteCmd := "rm -f " + defaultStoragePath + "/" + shellQuote(filename)
	cmd := SSHCmd(s.server, remoteCmd)
	debugLog("Delete | args: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	debugLog("Delete | exit_err: %v | output: %q", err, strings.TrimSpace(string(out)))
	if err != nil {
		return fmt.Errorf("delete failed: %w\n%s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// Send transfers one or more local files to the remote storage.
// Multiple files are archived into a temp tarball first.
func (s *StorageService) Send(localPaths []string) error {
	// TODO: implement
	return fmt.Errorf("send: not yet implemented")
}

// CleanAll removes all files from remote storage.
func (s *StorageService) CleanAll() error {
	cmd := SSHCmd(s.server,
		"rm -f "+defaultStoragePath+"/*",
	)
	debugLog("CleanAll | args: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	debugLog("CleanAll | exit_err: %v | output: %q", err, strings.TrimSpace(string(out)))
	if err != nil {
		return fmt.Errorf("clean all failed: %w\n%s", err, strings.TrimSpace(string(out)))
	}
	return nil
}
