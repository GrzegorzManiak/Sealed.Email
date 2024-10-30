#!/bin/bash

export ETCDCTL_API=3
ETCD_ENDPOINTS="localhost:2379"
ROOT_USERNAME="root"
ROOT_PASSWORD="root"

export ETCDCTL_ENDPOINTS=$ETCD_ENDPOINTS
etcdctl role add root
etcdctl role grant-permission root --prefix=true readwrite ""
etcdctl user add $ROOT_USERNAME <<EOF
$ROOT_PASSWORD
$ROOT_PASSWORD
EOF
etcdctl user grant-role $ROOT_USERNAME root
etcdctl auth enable