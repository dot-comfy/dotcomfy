#!/usr/bin/env bash
#
# Description: Basic test emulating installing a repo

set -xeuv

./bin/dotcomfy install https://gitlab.com/reavessm/dot-files
ls -al .dotcomfy
sleep 5

./bin/dotcomfy switch --repo ethangamma24 --branch macOS
ls -al .dotcomfy
sleep 5

./bin/dotcomfy switch --branch hyprland
ls -al .dotcomfy
sleep 5

./bin/dotcomfy switch --repo https://github.com/ethangamma24/dotfiles
ls -al .dotcomfy
sleep 5
