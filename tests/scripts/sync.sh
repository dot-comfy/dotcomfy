#!/usr/bin/env bash
#
# Description: Basic test emulating installing a repo

set -xeuv

./bin/dotcomfy install ethangamma24 --branch hyprland --skip-dependencies -vvvv
ls -al .dotcomfy/ghostty/
sleep 5

./bin/dotcomfy sync
ls -al .dotcomfy/ghostty/
sleep 5
