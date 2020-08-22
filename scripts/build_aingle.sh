#!/bin/bash
#
# Builds AINGLE with the latest git tag and commit hash (short)
# E.g.: ./aingle -v --> AINGLE 0.3.0-75316fe

latest_tag=$(git describe --tags $(git rev-list --tags --max-count=1))
commit_hash=$(git rev-parse --short HEAD)

go build -ldflags="-s -w -X github.com/Ariwonto/aingle-alpha/plugins/cli.AppVersion=${latest_tag:1}-$commit_hash"
