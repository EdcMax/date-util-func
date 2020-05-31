BIN := goreman
VERSION := $$(make -s show-version)
CURRENT_REVISION := $(shell git rev-parse --short HEAD)
BU