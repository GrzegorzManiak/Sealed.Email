#!/bin/bash

GOPROJ_ROOT="$(pwd)/../server"
PROTO_DIR="${GOPROJ_ROOT}/proto"
PROTOS=("notification" "domain" "smtp")

for proto in "${PROTOS[@]}"; do
    protoc --go_out="${PROTO_DIR}" --go_opt=paths=source_relative \
        --go-grpc_out="${PROTO_DIR}" --go-grpc_opt=paths=source_relative \
        --proto_path="${PROTO_DIR}" "${PROTO_DIR}"/"${proto}"/"${proto}".proto
done