package services

import (
	"os/exec"
	"strings"
)

const defaultPort = "22"
const defaultStoragePath = ".filepass_storage"

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
// --protect-args (-s) prevents rsync from shell-expanding the remote path,
// which correctly handles filenames with spaces and special characters.
func RsyncCmd(s Server, src, dst string) *exec.Cmd {
	sshFlag := "ssh -i " + s.PrivateKey + " -p " + serverPort(s) +
		" -o StrictHostKeyChecking=no -o BatchMode=yes"
	return exec.Command(
		"rsync",
		"-avz",
		"--partial",
		"--protect-args",
		"-e", sshFlag,
		src,
		dst,
	)
}

// RemotePath returns the full remote rsync path for a filename inside storage.
// No shell quoting needed — --protect-args in RsyncCmd handles special characters.
func RemotePath(s Server, filename string) string {
	return s.User + "@" + s.Host + ":" + defaultStoragePath + "/" + filename
}

// RemoteStorageRoot returns the remote storage root for rsync operations.
func RemoteStorageRoot(s Server) string {
	return s.User + "@" + s.Host + ":" + defaultStoragePath + "/"
}
