#!/bin/bash

cd /home/greg/
IP=$(hostname -I | awk '{print $1}')
etcd --listen-client-urls http://$IP:2379 \
--advertise-client-urls http://$IP:2379 \
--listen-peer-urls http://$IP:2380 \
--initial-advertise-peer-urls http://$IP:2380