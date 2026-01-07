package services

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	GitConfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/utils/merkletrie"

	Config "dotcomfy/internal/config"
	Log "dotcomfy/internal/logger"
)

func Clone(url, branch, commit_hash, path string) error {
	LOGGER = Log.GetLogger()

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

	// NOTE: This uses a raw Git command to clone the repo with SSH instead of
	//		 `git-go` because that module has quirks with using SSH keys for
	//		 auth.

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
	LOGGER = Log.GetLogger()
	user, err := user.Current()
	if err != nil {
		os.Exit(1)
	}

	dotcomfy_dir := user.HomeDir + "/.dotcomfy"
	old_dotfiles_dir := user.HomeDir + "/.config"

	repo, err := git.PlainOpen(repo_path)
	if err != nil {
		LOGGER.Errorf("Error opening the local repo in %s: %v", repo_path, err)
		return err
	}

	head, err := repo.Head()
	if err != nil {
		LOGGER.Errorf("Error getting HEAD: %v", err)
		return err
	}

	branch_name := string(head.Name())
	if strings.HasPrefix(branch_name, "refs/heads/") {
		branch_name = strings.TrimPrefix(branch_name, "refs/heads/")
	}

	err = repo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		Force:      true,
		RefSpecs: []GitConfig.RefSpec{
			GitConfig.RefSpec("+refs/heads/" + branch_name + ":refs/remotes/origin/" + branch_name),
		},
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		LOGGER.Errorf("Error fetching: %v", err)
		return err
	}

	remote_ref_name := plumbing.NewRemoteReferenceName("origin", branch_name)
	origin_ref, err := repo.Reference(remote_ref_name, true)
	if err != nil {
		LOGGER.Errorf("Error getting origin reference: %v", err)
		return err
	}

	remote_commit, err := repo.CommitObject(origin_ref.Hash())
	if err != nil {
		LOGGER.Errorf("Error getting origin commit hash: %v", err)
		return err
	}

	remote_tree, err := remote_commit.Tree()
	if err != nil {
		LOGGER.Errorf("Error getting origin commit tree: %v", err)
		return err
	}

	local_ref_name := plumbing.NewBranchReferenceName(branch_name)
	local_ref, err := repo.Reference(local_ref_name, true)
	if err != nil {
		LOGGER.Errorf("Error getting local reference: %v", err)
		return err
	}

	local_commit, err := repo.CommitObject(local_ref.Hash())
	if err != nil {
		LOGGER.Errorf("Error getting local commit hash: %v", err)
		return err
	}

	local_tree, err := local_commit.Tree()
	if err != nil {
		LOGGER.Errorf("Error getting origin commit tree: %v", err)
		return err
	}

	changes, err := object.DiffTree(local_tree, remote_tree)
	if err != nil {
		LOGGER.Errorf("Error getting changes: %v", err)
		return err
	}

	LOGGER.Infof("Origin ref after fetch: %s", origin_ref.Hash())

	branch := plumbing.NewBranchReferenceName(branch_name)
	// Bypass dirty worktree checks and just "fast forward" to the latest commit
	err = repo.Storer.SetReference(plumbing.NewHashReference(branch, origin_ref.Hash()))
	if err != nil {
		LOGGER.Errorf("Error switching local reference to latest from origin: %v", err)
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		LOGGER.Errorf("Error getting the worktree: %v", err)
		return err
	}

	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: branch,
		Force:  true,
	})
	if err != nil {
		LOGGER.Errorf("Error checking out branch: %v", err)
		return err
	}

	err = worktree.Pull(&git.PullOptions{
		RemoteName:    "origin",
		ReferenceName: branch,
		SingleBranch:  true,
		Force:         true,
		Progress:      os.Stdout, // May omit this, we'll see how it looks
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		LOGGER.Errorf("Error pulling: %v", err)
		return err
	}

	head, err = repo.Head()
	if err != nil {
		LOGGER.Errorf("Error getting HEAD: %v", err)
		return err
	}

	LOGGER.Infof("HEAD is now at %s\n", head.Hash())
	LOGGER.Infof("Changes from local to remote HEAD:")

	for _, change := range changes {
		action, err := change.Action()
		if err != nil {
			return err
		}

		patch, err := change.Patch()
		if err != nil {
			return err
		}

		patch_string := patch.String()

		lines := strings.Split(patch_string, "\n")
		if len(lines) == 0 {
			continue
		}

		file_path := strings.TrimPrefix(lines[0], "diff --git a/")
		file_path = "/" + strings.Split(file_path, " ")[0]

		if action == merkletrie.Insert {
			file_path, err = RenameSymlinkUnix(old_dotfiles_dir, dotcomfy_dir, file_path)
		}
	}

	return nil
}

/*************************************
 * TODO: *
 *************************************/
func Push(repo_path string) error {
	LOGGER = Log.GetLogger()
	// var repo_url string
	var branch string
	// var failed_files []string

	config := Config.GetConfig()
	auth := config.Auth

	ssh_key_path, err := auth.GetSSHKeyPath()
	if err != nil {
		LOGGER.Errorf("Error getting ssh key path: %v", err)
		return err
	}
	LOGGER.Debugf("Push resolved SSH key path: %s", ssh_key_path)

	// Validate SSH key file exists before attempting to read
	if _, err := os.Stat(ssh_key_path); os.IsNotExist(err) {
		LOGGER.Errorf("SSH key file does not exist: %s", ssh_key_path)
		return fmt.Errorf("SSH key file not found: %s", ssh_key_path)
	} else if err != nil {
		LOGGER.Errorf("Error accessing SSH key file %s: %v", ssh_key_path, err)
		return fmt.Errorf("cannot access SSH key file: %v", err)
	}

	ssh_key, err := os.ReadFile(ssh_key_path)
	if err != nil {
		LOGGER.Errorf("Error reading the ssh key: %v", err)
		return err
	}
	LOGGER.Debugf("SSH key read successfully for push, size: %d bytes", len(ssh_key))

	// Basic validation: check if key looks like PEM format
	keyContent := string(ssh_key)
	if !strings.Contains(keyContent, "-----BEGIN") || !strings.Contains(keyContent, "-----END") {
		LOGGER.Warn("SSH key may not be in PEM format (missing BEGIN/END markers)")
	}

	passphrase := auth.GetSSHKeyPassphrase()
	if passphrase != "" {
		LOGGER.Debugf("SSH key is encrypted (passphrase provided)")
	} else {
		LOGGER.Debugf("SSH key is not encrypted (no passphrase)")
	}

	ssh_auth, err := ssh.NewPublicKeys("git", ssh_key, auth.GetSSHKeyPassphrase())
	if err != nil {
		LOGGER.Errorf("Error creating the SSH authenticatior: %v", err)
		return err
	}
	LOGGER.Infof("SSH authentication method created successfully for push")

	// Grab URL to be pushed to? Do I need this?
	/*
		urls := remote.Config().URLs
		if len(urls) > 0 {
			repo_url = urls[0]
		} else {
			LOGGER.Fatal("No URL found for the remote 'origin'")
		}
	*/

	worktree, err := repo.Worktree()
	if err != nil {
		LOGGER.Errorf("Error getting the worktree: %v", err)
		return err
	}

	// Simple error-tolerant approach
	status, err := worktree.Status()
	if err != nil {
		LOGGER.Errorf("Error getting status: %v", err)
		return err
	}

	hasChanges := false
	for file := range status {
		_, err = worktree.Add(file)
		if err != nil {
			// Log warning but continue - empty directories can't be added anyway
			LOGGER.Warnf("Could not add %s (likely empty directory): %v", file, err)
			continue
		}
		hasChanges = true
	}

	if !hasChanges {
		fmt.Println("No changes to commit")
		return nil
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	username := auth.GetUsername()
	if err != nil {
		username = "unknown"
	}

	commit_message := fmt.Sprintf("Update from %s@%s at %s",
		username,
		hostname,
		time.Now().Format("12:30:00 CST 1963-11-22"),
	)

	_, err = worktree.Commit(commit_message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  username,
			Email: auth.GetEmail(),
			When:  time.Now(),
		},
	})
	if err != nil {
		LOGGER.Fatalf("Error committing: %v", err)
	}

	remote, err := repo.Remote("origin")
	if err != nil {
		LOGGER.Error("Error getting remote from the origin:", err)
		return err
	}

	// TODO:
	// Set up auth for private repos and/or username/password auth at runtime
	// var auth *http.BasicAuth
	// if auth != nil {

	// }

	err = remote.Push(&git.PushOptions{
		Auth:       ssh_auth,
		RemoteName: "origin",
		Progress:   os.Stdout,
		Force:      false,
		RefSpecs: []GitConfig.RefSpec{
			GitConfig.RefSpec("+refs/heads/" + branch_name + ":refs/remotes/origin/" + branch),
		},
	})
	if err != nil {
		LOGGER.Fatalf("Error pushing: %v", err)
	}

	return nil
}
