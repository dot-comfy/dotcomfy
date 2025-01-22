#!/usr/bin/env bash
#
# Description: Basic test emulating installing a repo

set -xeuv

mkdir .config/

./bin/dotcomfy install ethangamma24
ls -al .dotcomfy
ls -al .config
pwd
ls -al

./bin/dotcomfy uninstall --yes
ls -al .dotcomfy
ls -al .config
