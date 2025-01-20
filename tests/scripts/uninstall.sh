#!/usr/bin/env bash
#
# Description: Basic test emulating installing a repo

set -xeuv

./bin/dotcomfy install ethangamma24
ls -al .dotcomfy
ls -al .config

./bin/dotcomfy uninstall --yes
ls -al .dotcomfy
ls -al .config
