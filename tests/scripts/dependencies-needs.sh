#!/usr/bin/env bash
8
# Description: Tests handling misconfigured dependencies section of config file

set -xeuv

which fzf || true
which tmux || true
which zig || true
which nvm || true
which nvim || true
./bin/dotcomfy install ethangamma24 --branch dotcomfy-dependency-needs -vvvv
which fzf || true
which tmux || true
which zig || true
which nvm || true
which nvim || true
