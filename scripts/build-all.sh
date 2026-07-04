#!/usr/bin/env sh
set -eu
wails build -platform darwin/universal
wails build -platform windows/amd64
wails build -platform linux/amd64
