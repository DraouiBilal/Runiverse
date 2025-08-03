#!/bin/bash

image=$1

if [[ -z image ]]; then 
    echo "You need to pass the image to build"
fi

podman build -f services/$image/Containerfile -t $image:test services/$image
