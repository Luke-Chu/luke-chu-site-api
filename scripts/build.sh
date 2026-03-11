#!/usr/bin/env bash
set -euo pipefail

go mod tidy
mkdir -p bin
go build -o bin/luke-chu-site-api ./cmd/server

