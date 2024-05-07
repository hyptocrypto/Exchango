#!/bin/bash

# Function to handle SIGINT
cleanup() {
    echo "Caught SIGINT. Cleaning up and exiting."
    pkill exchangego
    exit 1
}

# Register the cleanup function to be called on SIGINT
trap cleanup INT

# Kill any running server process
pkill exchangego

# Build and start the server
go run .
# Watch for changes in .go files
fswatch . | while read f; do 
    extension="${f##*.}"
    if [ "$extension" = "go" ]; then
        echo "Changes detected. Restarting server."
        pkill exchangego
        go run .
    fi
done
