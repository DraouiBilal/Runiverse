#!/bin/bash

arg=$1

if [[ -z arg ]]; then
    
    echo "Command needed: start, stop"

fi

case "$arg" in
    start)
        podman compose -f podman-compose.yml up
        ;;
    stop)
        podman compose -f podman-compose.yml down
        ;;
    *)
        echo "Command not supported, use: start or stop"
        ;;
esac
