#!/usr/bin/env bash
#
# Description: Basic test emulating installing a repo

set -xeuv

./bin/dotcomfy install ethangamma24 --branch hyprland --at-commit de02cedcfa7eb6d5186fdb81e13290841081d6f8 --skip-dependencies
ls -al .dotcomfy/ghostty/
sleep 5

./bin/dotcomfy sync
ls -al .dotcomfy/ghostty/
sleep 5
