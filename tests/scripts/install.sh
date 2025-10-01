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
mkdir ~/.config
touch ~/.config/.viminfo
touch ~/.config/.vimrc
mkdir ~/.config/nvim
touch ~/.config/nvim/init.lua
ls ~/.config
./bin/dotcomfy install ethangamma24 --branch macOS --skip-dependencies -vv
ls -al .dotcomfy
sleep 5
ls -al .config
sleep 5

./bin/dotcomfy uninstall --yes
ls -al .dotcomfy
sleep 5
ls -al .config
ls -al .config/nvim

# ./bin/dotcomfy uninstall --yes
# ls -al .dotcomfy
# sleep 5

: <<'END_COMMENT'
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
# nvim --help
which nvim || true
sleep 5
END_COMMENT
