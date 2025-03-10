#!/bin/sh
#
# Use this script to run your program LOCALLY.
#
# Note: Changing this script WILL NOT affect how CodeCrafters runs your program.
#
# Learn more: https://codecrafters.io/program-interface

set -e # Exit early if any commands fail

# Compile the server if needed
go build -o bin/server app/cmd/server/main.go

# Run the server with the provided arguments
exec bin/server "$@"
