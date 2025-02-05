![dotcomfy Logo](logo.jpg)

[![Go Report Card](https://goreportcard.com/badge/github.com/your-username/dotcomfy)](https://goreportcard.com/report/github.com/your-username/dotcomfy)

**dotcomfy** is a CLI tool designed to simplify the management of configuration files for developer tools like Neovim, Tmux, Alacritty, and more. With dotcomfy, you can install, switch, and manage your config files with ease, automating the setup of package dependencies along the way.

## Features

- **One-command installation** of config sets for various developer tools.
- **Config switching** between different setups or environments.
- **Automated dependency management** for packages required by your configs.
- **Customizable installation scripts** for tools without standard package management.

### Note
The dependency management feature is still in development and may not work as expected.

## Installation

### Prerequisites

- Go
- Git
- PackageKit

### Building from Source

```sh
git clone https://github.com/dot-comfy/dotcomfy.git
cd dotcomfy
make build
sudo make install
```

This will build the binare and install it to `/usr/local/bin/`.

## Usage

### Installation
`dotcomfy install [REPO] --branch [BRANCH]`
- REPO: can be either a GitHub username or a repository URL.
  - If you're using a GitHub username, dotcomfy will attempt to clone the `dotfiles` repository under that user.
- BRANCH: the branch of the repository to install. If not specified, the `main` branch will be used.

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
```toml
[dependencies]
# Version can be specified for a package being installed from a package manager
fzf = { version = "0.57.0" }
nvim = { 
    # Custom installations can be specified step by step
    steps = [
        "curl -LO https://github.com/neovim/neovim/releases/latest/download/nvim-linux-x86_64.tar.gz",
        "sudo rm -rf /opt/nvim",
        "sudo tar -C /opt -xzf nvim-linux-x86_64.tar.gz",
        "echo 'export PATH=\"$PATH:/opt/nvim-linux-x86_64/bin\"' >> .zshrc"
    ]
}
# Empty configs will default to installing the latest version
# of the package found in the package manager
tmux = {}
# Commands needed after package installation can also be specified
zsh = { 
    post_install_steps = [
        "chsh -s $(which zsh)"
    ]
}
oh-my-zsh = {
    # Custom installation scripts can be specified. These should be located in the same directory as the config file.
    script = "oh-my-zsh.sh"
}
```
