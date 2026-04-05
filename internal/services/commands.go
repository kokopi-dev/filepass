package services

import (
	"os/exec"
	"strings"
)

const defaultPort = "22"
const defaultStoragePath = "~/.filepass_storage"

// shellQuote wraps s in single quotes, escaping any single quotes within it.
// This is safe for use in remote shell commands passed over SSH.
func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

func serverPort(s Server) string {
	if s.Port == "" {
		return defaultPort
	}
	return s.Port
}

// SSHCmd returns an exec.Cmd for running a single command on the server.
func SSHCmd(s Server, remoteCmd string) *exec.Cmd {
	return exec.Command(
		"ssh",
		"-i", s.PrivateKey,
		"-p", serverPort(s),
		"-o", "StrictHostKeyChecking=no",
		"-o", "BatchMode=yes",
		s.User+"@"+s.Host,
		remoteCmd,
	)
}

// RsyncCmd returns an exec.Cmd for an rsync transfer.
// src and dst follow standard rsync syntax (local path or user@host:path).
func RsyncCmd(s Server, src, dst string) *exec.Cmd {
	sshFlag := "ssh -i " + s.PrivateKey + " -p " + serverPort(s) +
		" -o StrictHostKeyChecking=no -o BatchMode=yes"
	return exec.Command(
		"rsync",
		"-avz",
		"--partial",
		"-e", sshFlag,
		src,
		dst,
	)
}

// RemotePath returns the full remote rsync path for a filename inside storage.
// Only the filename is shell-quoted so ~ expands correctly on the remote shell.
func RemotePath(s Server, filename string) string {
	return s.User + "@" + s.Host + ":" + defaultStoragePath + "/" + shellQuote(filename)
}

// RemoteStorageRoot returns the remote storage root for rsync operations.
func RemoteStorageRoot(s Server) string {
	return s.User + "@" + s.Host + ":" + defaultStoragePath + "/"
}
