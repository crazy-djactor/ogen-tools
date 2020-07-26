#!/bin/bash

## Source: https://gist.github.com/brianclements/f72b2de8e307c7b56689
## This script will remove all docker images and containers

docker rm -vf $(docker ps -a -q) 2>/dev/null || echo "No more containers to remove."
docker rmi $(docker images -q) 2>/dev/null || echo "No more images to remove."

exit 0