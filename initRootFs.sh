#!/bin/bash
CONTAINER_NAME=$1

echo "[Docker Install] Checking installation of Docker..."
if [ -x "$(command -v docker)" ]; then
    echo "[Docker Install] Docker installed, skipping installation"
else
    echo "[Docker Install] Installing Docker..."

    if [[ "$OSTYPE" == "darwin"* ]]; then
        # Install docker with brew if on Mac OS X
        brew install docker
        echo "[Docker Install] Installed Docker"
    elif [[ "$OSTYPE" == "linux"* || "$OSTYPE" == "bsd"* ]]; then
        # Instacll docker with apt if on Linux or BSD
        apt install docker.io
        echo "[Docker Install] Installed Docker"
    else
        echo "[Docker Install] Cannot install Docker, please install it manually"
    fi
fi

# Retrieve the Alpine Linux image distro
docker pull alpine:latest
# Get the overlayfs location
OVERLAY_ID=$(docker image inspect alpine:latest -f '{{.GraphDriver.Data.UpperDir}}' | sed "s/\'//g")

# Copy the overlayfs to the container location
cp -r $OVERLAY_ID /root/$CONTAINER_NAME
echo "[FS Init] Copied overlay filesystem from Alpine to [/root/$CONTAINER_NAME]"

# Create an identifier in the container
CONTAINER_IDENTIFIER="testContainer"
mkdir -p /root/$CONTAINER_NAME/$CONTAINER_IDENTIFIER
echo "[FS Init] Created identifier in container [/root/$CONTAINER_NAME/$CONTAINER_IDENTIFIER]"
echo "[FS Init] Filesystem creation for [$CONTAINER_NAME] complete"