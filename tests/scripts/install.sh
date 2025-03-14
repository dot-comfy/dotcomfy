#!/usr/bin/env bash
#
# Description: Basic test emulating installing a repo

set -xeuv

# ./bin/dotcomfy install https://gitlab.com/reavessm/dot-files --skip-dependencies
# ls -al .dotcomfy
# sleep 5
# 
# ./bin/dotcomfy uninstall --yes
# ls -al .dotcomfy
# sleep 5
# 
# ./bin/dotcomfy install ethangamma24 --branch macOS --skip-dependencies
# ls -al .dotcomfy
# sleep 5
# 
# ./bin/dotcomfy uninstall --yes
# ls -al .dotcomfy
# sleep 5


which fzf || true
which tmux || true
which zig || true
which nvm || true
which nvim || true
./bin/dotcomfy install ethangamma24 --branch hyprland -vvvv
source ~/.zshrc
fzf --help
# tmux --help
zig --help
nvm --help
nvim --help
sleep 5
