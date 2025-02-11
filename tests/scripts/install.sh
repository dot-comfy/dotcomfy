#!/usr/bin/env bash
#
# Description: Basic test emulating installing a repo

set -xeuv

# ./bin/dotcomfy install https://gitlab.com/reavessm/dot-files
# ls -al .dotcomfy
# sleep 5
# 
# ./bin/dotcomfy uninstall --yes
# ls -al .dotcomfy
# sleep 5
# 
# ./bin/dotcomfy install ethangamma24 --branch macOS
# ls -al .dotcomfy
# sleep 5
# 
# ./bin/dotcomfy uninstall --yes
# ls -al .dotcomfy
# sleep 5

./bin/dotcomfy install ethangamma24 --branch hyprland
ls -al .config/dotcomfy
ls -al .dotcomfy
sleep 5
