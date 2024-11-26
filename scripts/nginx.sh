#!/bin/bash

# -- Kill all NGINX processes
echo "Killing all NGINX processes"
sudo killall nginx
sudo kill -9 $(sudo lsof -t -i:2095 -c nginx) 2> /dev/null

# -- Start NGINX
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
CONFIG="/../dev/nginx.conf"
echo "Starting NGINX in DIR: $DIR$CONFIG"
sudo nginx -c $DIR$CONFIG
