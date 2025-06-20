![dotcomfy Logo](logo.jpg)

[![Go Report Card](https://goreportcard.com/badge/github.com/dot-comfy/dotcomfy)](https://goreportcard.com/report/github.com/dot-comfy/dotcomfy)

**dotcomfy** is a CLI tool designed to simplify the management of configuration files for developer tools like Neovim, Tmux, Alacritty, and more. With dotcomfy, you can install, switch, and manage your config files with ease, automating the setup of package dependencies along the way.

## Features

- **One-command installation** of config sets for various developer tools.
- **Config switching** between different setups or environments.
- **Automated dependency management** for packages required by your configs.
- **Customizable installation scripts** for tools without standard package management.

### Note
The dependency management feature is still in development and may not work as expected.

## Installation

### WARNING

**Please make sure you back up your `.config` directory before using dotcomfy. It is currently a WIP, and stability is not guaranteed.**

### Prerequisites

- Go
- Git

### Building from Source

```sh
git clone https://github.com/dot-comfy/dotcomfy.git
cd dotcomfy
make build
sudo make install
```

This will build the binary and install it to `/usr/local/bin/`.

## Usage

### Installation
`dotcomfy install [REPO] --branch [BRANCH] --skip-dependencies`
- REPO: can be either a GitHub username or a repository URL.
  - If you're using a GitHub username, dotcomfy will attempt to clone the `dotfiles` repository under that user.
- BRANCH: the branch of the repository to install. If not specified, the `main` branch will be used.
- `--skip-dependencies` skips the dependency installation step.

### Switch
`dotcomfy switch --repo [REPO] --branch [BRANCH]`
- One or both of `--repo` and `--branch` must be specified.
- If only `--branch` is specified, the current installation will switch to that branch of the current repository.

### Uninstall
`dotcomfy uninstall --yes`
- Uninstalls the currently installed config set.
- `--yes` autoconfirms the uninstallation process.

## Configuration

dotcomfy's config file lives at `$HOME/.config/dotcomfy/config.toml`.

### Dependencies

You can specify packages that need to be installed in order for the config set to function properly:
```toml filename="config.toml"
[dependencies]
# Version can be specified for a package being installed from a package manager
fzf = { version = "0.57.0" }
# Custom shell scripts for dependency installation can be specified.
# The `needs` field can be used to specify dependencies that need to be installed before this dependency.
nvim = { script = "nvim.sh", needs = ["zsh"] }
# Custom installations can be specified step by step
nvm = { steps = [ "curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.1/install.sh | bash", "source ~/.zshrc" ], needs = ["zsh"] }
tmux = { version = "latest" }
# Commands needed after package installation can also be specified
zsh = { post_install_steps = [ "chsh -s $(which zsh)" ] }
```
#### Custom Installation Scripts

Scripts should start with `#!/bin/sh` and should be located in the same directory as the config file.
```sh filename="nvim.sh"
#!/bin/bash

curl -LO https://github.com/neovim/neovim/releases/latest/download/nvim-linux-x86_64.tar.gz
sudo rm -rf /opt/nvim
sudo tar -C /opt -xzf nvim-linux-x86_64.tar.gz
if [ -n "$ZSH_VERSION" ]; then
    echo 'export PATH="$PATH:/opt/nvim-linux-x86_64/bin"' >> "$HOME/.zshrc"
    echo 'Adding path to .zshrc'
elif [ -n "$BASH_VERSION" ]; then
    echo 'export PATH="$PATH:/opt/nvim-linux-x86_64/bin"' >> "$HOME/.bashrc"
    echo 'Adding path to .bashrc'
else
    echo 'No idea what shell is on this system'
fi
export PATH="$PATH:/opt/nvim-linux-x86_64/bin"
```
