#!/usr/bin/env bash
#
# Description: Basic test emulating installing a repo

set -xeuv

echo ".config"
echo "------------------"
ls -al .config
./bin/dotcomfy install https://gitlab.com/reavessm/dot-files
echo ".config"
echo "------------------"
ls -al .config
echo ".dotcomfy"
echo "------------------"
ls -al .dotcomfy
