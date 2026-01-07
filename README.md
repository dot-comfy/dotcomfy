![dotcomfy Logo](logo.jpg)

[![Go Report Card](https://goreportcard.com/badge/github.com/dot-comfy/dotcomfy)](https://goreportcard.com/report/github.com/dot-comfy/dotcomfy)

**dotcomfy** is a CLI tool written in Go that simplifies the management of configuration files (dotfiles) for developer tools like Neovim, Tmux, Alacritty, and more. With dotcomfy, you can install, switch, and manage your config files with ease, automating the setup of package dependencies along the way.

Whether you're SSHing into brand new cloud servers, bouncing between different operating systems, or just wanting to try out different Linux rices, dotcomfy has you covered!

## Features

- One-command installation of config sets from Git repositories
- Config switching between different setups/environments
- Automated dependency management with custom installation scripts and package manager prioritization
- Support for both public and private repositories with SSH authentication
- Containerized testing across multiple Linux distributions (Arch, Fedora, Ubuntu)

### Note
The dependency management feature, including package manager prioritization, is still in development and may not work as expected.

## Prerequisites

- Go 1.23.0+
- Git
- Package manager (pacman, yum, apt, brew, etc.) - auto-detected with user prioritization

## Installation

### WARNING

**Please make sure you back up your `.config` directory before using dotcomfy. It is currently a WIP, and stability is not guaranteed.**

### Building from Source

```sh
git clone https://github.com/dot-comfy/dotcomfy.git
cd dotcomfy
make build
make install
```

This will build the binary and install it to `/usr/local/bin/`. If you want to specify a different directory for installation, you can specify it in the `make install` command:
```bash
make install INSTALL_DIR=~/bin
```

## Usage

### Installation
`dotcomfy install [REPO] --branch [BRANCH] --package-manager [PM] --skip-dependencies`
- REPO: can be either a GitHub username or a repository URL.
  - If you're using a GitHub username, dotcomfy will attempt to clone the `dotfiles` repository under that user.
- BRANCH: the branch of the repository to install. If not specified, the `main` branch will be used.
- `--package-manager` or `-pm`: Specify preferred package manager for dependency installation (apt, dnf, yum, yay, pacman, zypper, brew). Falls back to auto-detection if not available.
- `--skip-dependencies` skips the dependency installation step.

### Switch
`dotcomfy switch --repo [REPO] --branch [BRANCH]`
- One or both of `--repo` and `--branch` must be specified.
- If only `--branch` is specified, the current installation will switch to that branch of the current repository.

### Uninstall
`dotcomfy uninstall --yes`
- Uninstalls the currently installed config set.
- `--yes` autoconfirms the uninstallation process.

### Pull
`dotcomfy pull`
- Pulls the latest changes from the current branch of the dotcomfy installation. Please note that any files locally changed that conflict with changes being pulled in **will automatically be overwritten**.

### Push
`dotcomfy push`
- Stages all changes, commits them with an auto-generated message including username, hostname, and timestamp, then pushes to the remote origin branch.
- Requires write permissions on the remote repository.
- Uses SSH authentication from config.

### Global Flags
- `--config`: Specify a custom config file path (default: `$HOME/.config/dotcomfy/config.toml`)
- `-v`: Increase logging verbosity (can be used multiple times, e.g.: `-vvvv`)

## Configuration

dotcomfy's config file supports YAML. It lives at `$HOME/.config/dotcomfy/config.yaml`.

### Dependencies
Define packages or tools required for your dotfiles. Each dependency supports:
- **version**: Package manager installation with specific version (e.g., "0.57.0", "latest")
- **script**: Path to custom shell script for installation (relative to config file)
- **steps**: Array of shell commands for installation
- **needs**: Array of other dependencies that must be installed first
- **post_install_steps**: Array of commands to run after package installation (requires version)
- **post_install_script**: Path to script to run after installation (requires version)

**Validation rules:**
- Must specify at least `version`, `script`, or `steps`
- Only one of `version`, `steps` or `script` can be specified
- Cannot mix `post_install_steps` with `post_install_script`
- No self-dependencies or cycles

Example:
```yaml
dependencies:
  fzf:
    version: "0.57.0"
  nvim:
    script: "nvim.sh"
    needs:
      - zsh
  tmux:
    version: "latest"
   zsh:
     post_install_steps:
       - chsh -s $(which zsh)
 ```

#### Package Manager
Specify a preferred package manager for dependency installation. If not set or unavailable, dotcomfy will auto-detect from supported managers.
- **preferred_package_manager**: One of: apt, dnf, yum, yay, pacman, zypper, brew

Example:
```yaml
preferred_package_manager: pacman
```

#### Authentication
Required for pushing to private repositories or SSH-based operations.
- **username**: Your Git username
- **email**: Your Git email
- **ssh_file**: Path to SSH private key
- **ssh_key_passphrase**: Passphrase for SSH key (optional)

Example:
```yaml
authentication:
  username: "your_username"
  email: "your_email@example.com"
  ssh_file: "~/.ssh/id_rsa"
  ssh_key_passphrase: "optional_passphrase"
```

## Development

### Project Structure
```
dotcomfy/
├── main.go                   # Entry point with sudo protection
├── go.mod/go.sum             # Go module dependencies
├── Makefile                  # Build, test, and installation targets
├── README.md                 # User documentation
├── bin/                      # Built binary output
├── cmd/dotcomfy/cobra/       # CLI command implementations
├── internal/
│   ├── config/               # Configuration structures and validation
│   ├── services/             # Business logic (Git ops, dependencies, etc.)
│   └── logger/               # Logging configuration
├── tests/scripts/            # Test scripts for container testing
├── docs/                     # Documentation and references
└── Containerfile*            # Container definitions for different distros
```

### Building
```bash
make build          # Build binary to bin/dotcomfy
make install        # Install to /usr/local/bin (or custom path)
make references     # Build docs/REFERENCES.md
```

### Testing
```bash
make test-<script>  # Run specific test in container (e.g., make test-install)
make test          # Run all tests
make container     # Interactive container for testing
```

Containerized tests support multiple Linux distributions with Podman/Docker.

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## License

The MIT License (MIT)

Copyright © 2025 Ethan Harmon, Stephen Reaves

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
