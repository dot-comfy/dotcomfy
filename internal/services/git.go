package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	Config "dotcomfy/internal/config"
	Log "dotcomfy/internal/logger"
)

func Clone(url, branch, commit_hash, path string) error {
	LOGGER := Log.GetLogger()

	// Default to clone with ssh
	use_ssh := true

	config := Config.GetConfig()
	auth := config.Auth

	// Debug: Log authentication configuration
	ssh_file_path, _ := auth.GetSSHKeyPath()
	LOGGER.Debugf("Authentication configuration loaded:")
	LOGGER.Debugf("  Username: %s", auth.GetUsername())
	LOGGER.Debugf("  Email: %s", auth.GetEmail())
	LOGGER.Debugf("  SSH file: %s", ssh_file_path)
	if auth.GetSSHKeyPassphrase() != "" {
		LOGGER.Debugf("  SSH key passphrase: (provided)")
	} else {
		LOGGER.Debugf("  SSH key passphrase: (not provided)")
	}

	ssh_key_path, err := auth.GetSSHKeyPath()
	if err != nil {
		LOGGER.Errorf("Error getting ssh key path: %v", err)
		LOGGER.Errorf("SSH authentication will be disabled, falling back to HTTPS")
		use_ssh = false
	} else {
		LOGGER.Debugf("Resolved SSH key path: %s", ssh_key_path)
	}

	// Validate SSH key file exists and is readable
	if use_ssh && ssh_key_path != "" {
		if _, err := os.Stat(ssh_key_path); os.IsNotExist(err) {
			LOGGER.Errorf("SSH key file does not exist: %s", ssh_key_path)
			LOGGER.Errorf("SSH authentication will be disabled, falling back to HTTPS")
			use_ssh = false
		} else if err != nil {
			LOGGER.Errorf("Error accessing SSH key file %s: %v", ssh_key_path, err)
			LOGGER.Errorf("SSH authentication will be disabled, falling back to HTTPS")
			use_ssh = false
		} else {
			LOGGER.Debugf("SSH key file exists and is accessible: %s", ssh_key_path)

			// Basic validation: check if key looks like PEM format
			ssh_key, err := os.ReadFile(ssh_key_path)
			if err != nil {
				LOGGER.Errorf("Error reading the ssh key: %v", err)
				use_ssh = false
			} else {
				keyContent := string(ssh_key)
				if !strings.Contains(keyContent, "-----BEGIN") || !strings.Contains(keyContent, "-----END") {
					LOGGER.Warn("SSH key may not be in PEM format (missing BEGIN/END markers)")
					use_ssh = false
				} else {
					LOGGER.Debugf("SSH key validation passed, size: %d bytes", len(ssh_key))
				}
			}
		}
	}

	// Prepare git clone command arguments
	var cloneArgs []string
	cloneArgs = append(cloneArgs, "clone")

	// Add branch option if specified
	if branch != "" && branch != "main" && branch != "master" {
		cloneArgs = append(cloneArgs, "--branch", branch)
		cloneArgs = append(cloneArgs, "--single-branch")
	}

	// Add URL and path
	cloneArgs = append(cloneArgs, url, path)

	// Prepare environment for git command
	env := os.Environ()
	clone_url := url

	if use_ssh && ssh_key_path != "" {
		// Configure SSH command for git
		sshCmd := fmt.Sprintf("ssh -i %s -o IdentitiesOnly=yes -o StrictHostKeyChecking=no", ssh_key_path)
		if auth.GetSSHKeyPassphrase() != "" {
			// For keys with passphrases, we'd need ssh-agent or expect script
			// For now, assume passphrase-less keys or keys loaded in ssh-agent
			LOGGER.Warn("SSH keys with passphrases are not fully supported in this implementation")
			LOGGER.Warn("Consider using ssh-agent or converting to passphrase-less keys")
		}
		env = append(env, fmt.Sprintf("GIT_SSH_COMMAND=%s", sshCmd))

		// Convert HTTPS URL to SSH format for SSH cloning
		if strings.HasPrefix(url, "https://github.com/") {
			// https://github.com/user/repo.git -> git@github.com:user/repo.git
			ssh_url := strings.Replace(url, "https://github.com/", "git@github.com:", 1)
			ssh_url = strings.Replace(ssh_url, ".git", ".git", 1) // keep .git suffix
			clone_url = ssh_url
			LOGGER.Infof("Converted HTTPS URL to SSH: %s", clone_url)
		} else if strings.HasPrefix(url, "https://gitlab.com/") {
			// https://gitlab.com/user/repo.git -> git@gitlab.com:user/repo.git
			ssh_url := strings.Replace(url, "https://gitlab.com/", "git@gitlab.com:", 1)
			clone_url = ssh_url
			LOGGER.Infof("Converted HTTPS URL to SSH: %s", clone_url)
		} else {
			LOGGER.Warn("SSH authentication configured but URL format not recognized for conversion")
		}
		LOGGER.Infof("Using SSH authentication with key: %s", ssh_key_path)
	} else {
		LOGGER.Info("Cloning with HTTPS (SSH not available)")
	}

	// Update the URL in clone args
	cloneArgs[len(cloneArgs)-2] = clone_url

	// Execute git clone command
	LOGGER.Infof("Executing: git %s", strings.Join(cloneArgs, " "))
	cmd := exec.Command("git", cloneArgs...)
	cmd.Env = env
	cmd.Dir = path

	// Get parent directory for command execution
	parentDir := filepath.Dir(path)
	if parentDir != "." {
		cmd.Dir = parentDir
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		LOGGER.Errorf("Git clone failed: %v", err)
		LOGGER.Errorf("Git output: %s", string(output))
		return fmt.Errorf("git clone failed: %v", err)
	}

	LOGGER.Infof("Git clone successful")
	LOGGER.Debugf("Git output: %s", string(output))

	// Handle commit hash checkout if specified
	if commit_hash != "" {
		LOGGER.Infof("Checking out commit: %s", commit_hash)

		checkoutArgs := []string{"-C", path, "checkout", commit_hash}
		checkoutCmd := exec.Command("git", checkoutArgs...)
		checkoutOutput, checkoutErr := checkoutCmd.CombinedOutput()

		if checkoutErr != nil {
			LOGGER.Errorf("Git checkout failed: %v", checkoutErr)
			LOGGER.Errorf("Git checkout output: %s", string(checkoutOutput))
			return fmt.Errorf("git checkout failed: %v", checkoutErr)
		}

		LOGGER.Infof("Git checkout successful")
	}

	return nil
}

func Pull(repo_path string) error {
	LOGGER := Log.GetLogger()

	config := Config.GetConfig()
	auth := config.Auth

	// Get SSH key path for authentication
	ssh_key_path, err := auth.GetSSHKeyPath()
	if err != nil {
		LOGGER.Errorf("Error getting ssh key path: %v", err)
		return err
	}

	// Validate SSH key file exists
	if ssh_key_path != "" {
		if _, err := os.Stat(ssh_key_path); os.IsNotExist(err) {
			LOGGER.Errorf("SSH key file does not exist: %s", ssh_key_path)
			return fmt.Errorf("SSH key file not found: %s", ssh_key_path)
		}
	}

	// Prepare environment for git command with SSH authentication
	env := os.Environ()
	if ssh_key_path != "" {
		sshCmd := fmt.Sprintf("ssh -i %s -o IdentitiesOnly=yes -o StrictHostKeyChecking=no", ssh_key_path)
		env = append(env, fmt.Sprintf("GIT_SSH_COMMAND=%s", sshCmd))
		LOGGER.Infof("Using SSH authentication with key: %s", ssh_key_path)
	}

	// Execute git pull command
	LOGGER.Infof("Executing: git pull --ff-only")
	cmd := exec.Command("git", "pull", "--ff-only")
	cmd.Dir = repo_path
	cmd.Env = env

	output, err := cmd.CombinedOutput()
	if err != nil {
		LOGGER.Errorf("Git pull failed: %v", err)
		LOGGER.Errorf("Git output: %s", string(output))
		return fmt.Errorf("git pull failed: %v", err)
	}

	LOGGER.Infof("Git pull successful")
	LOGGER.Debugf("Git output: %s", string(output))

	return nil
}

func Push(repo_path string) error {
	LOGGER := Log.GetLogger()

	config := Config.GetConfig()
	auth := config.Auth

	// Get SSH key path for authentication
	ssh_key_path, err := auth.GetSSHKeyPath()
	if err != nil {
		LOGGER.Errorf("Error getting ssh key path: %v", err)
		return err
	}

	// Validate SSH key file exists
	if ssh_key_path != "" {
		if _, err := os.Stat(ssh_key_path); os.IsNotExist(err) {
			LOGGER.Errorf("SSH key file does not exist: %s", ssh_key_path)
			return fmt.Errorf("SSH key file not found: %s", ssh_key_path)
		}
	}

	// Prepare environment for git command with SSH authentication
	env := os.Environ()
	if ssh_key_path != "" {
		sshCmd := fmt.Sprintf("ssh -i %s -o IdentitiesOnly=yes -o StrictHostKeyChecking=no", ssh_key_path)
		env = append(env, fmt.Sprintf("GIT_SSH_COMMAND=%s", sshCmd))
		LOGGER.Infof("Using SSH authentication with key: %s", ssh_key_path)
	}

	// Execute git add to stage all changes
	LOGGER.Infof("Executing: git add -A")
	addCmd := exec.Command("git", "add", "-A")
	addCmd.Dir = repo_path
	addCmd.Env = env

	addOutput, addErr := addCmd.CombinedOutput()
	if addErr != nil {
		LOGGER.Errorf("Git add failed: %v", addErr)
		LOGGER.Errorf("Git add output: %s", string(addOutput))
		return fmt.Errorf("git add failed: %v", addErr)
	}

	// Check if there are changes to commit
	LOGGER.Infof("Executing: git status --porcelain")
	statusCmd := exec.Command("git", "status", "--porcelain")
	statusCmd.Dir = repo_path
	statusCmd.Env = env

	statusOutput, statusErr := statusCmd.CombinedOutput()
	if statusErr != nil {
		LOGGER.Errorf("Git status failed: %v", statusErr)
		LOGGER.Errorf("Git status output: %s", string(statusOutput))
		return fmt.Errorf("git status failed: %v", statusErr)
	}

	if len(statusOutput) == 0 {
		fmt.Println("No changes to commit")
		return nil
	}

	// Create commit message
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	username := auth.GetUsername()
	if username == "" {
		username = "unknown"
	}

	commit_message := fmt.Sprintf("Update from %s@%s at %s",
		username,
		hostname,
		time.Now().Format("2006-01-02 15:04:05 MST"),
	)

	// Execute git commit
	LOGGER.Infof("Executing: git commit -m '%s'", commit_message)
	commitCmd := exec.Command("git", "commit", "-m", commit_message)
	commitCmd.Dir = repo_path
	commitCmd.Env = env

	commitOutput, commitErr := commitCmd.CombinedOutput()
	if commitErr != nil {
		LOGGER.Errorf("Git commit failed: %v", commitErr)
		LOGGER.Errorf("Git commit output: %s", string(commitOutput))
		return fmt.Errorf("git commit failed: %v", commitErr)
	}

	// Execute git push
	LOGGER.Infof("Executing: git push origin HEAD")
	pushCmd := exec.Command("git", "push", "origin", "HEAD")
	pushCmd.Dir = repo_path
	pushCmd.Env = env

	pushOutput, pushErr := pushCmd.CombinedOutput()
	if pushErr != nil {
		LOGGER.Errorf("Git push failed: %v", pushErr)
		LOGGER.Errorf("Git push output: %s", string(pushOutput))
		return fmt.Errorf("git push failed: %v", pushErr)
	}

	LOGGER.Infof("Git push successful")
	LOGGER.Debugf("Git push output: %s", string(pushOutput))

	return nil
}
