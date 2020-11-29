#!/bin/bash
CONTAINER_NAME=$1

if [ -x "$(command -v docker)" ]; then
    echo "Updating docker"
    docker update
else
    echo "Installing docker"

    if [[ "$OSTYPE" == "darwin"* ]]; then
        # Install docker with brew if on Mac OS X
        brew install docker
    elif [[ "$OSTYPE" == "linux"* || "$OSTYPE" == "bsd"* ]] then
        # Instacll docker with apt if on Linux or BSD
        apt install docker.io
    else
        echo "Cannot install docker, please install it manually"
    fi
fi

# Retrieve the Alpine Linux image distro
docker pull alpine:latest
# Get the overlayfs location
OVERLAY_ID = $(docker image inspect alpine:latest -f {{.GraphDriver.Data.UpperDir}} | sed "s/\'//g")

# Copy the overlayfs to the container location
cp -r $OVERLAY_ID /root/$CONTAINER_NAME

# Create an identifier in the container
CONTAINER_IDENTIFIER = "testContainer"
mkdir -p /root/$CONTAINER_NAME/$CONTAINER_IDENTIFIER