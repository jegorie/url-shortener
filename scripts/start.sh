#!/bin/sh
go build ./cmd/url-shortener
CONFIG_PATH=./config/local.yaml ./url-shortener
