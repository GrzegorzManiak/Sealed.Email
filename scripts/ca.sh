#!/bin/bash

clear

echo "!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!"
echo "!!! WARNING: DO NOT RUN WITH GOLAND / INTELIJ, IT WONT WORK !!!"
echo "!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!"

CA_DIR="../dev/.certs/ca"
mkdir -p "$CA_DIR"
CA_KEY="$CA_DIR/key.pem"
CA_CERT="$CA_DIR/ca.pem"
DAYS_VALID=365

echo "Generating CA certificate..."
openssl ecparam -name prime256v1 -genkey -noout -out $CA_KEY
openssl req -new -x509 -key $CA_KEY -out $CA_CERT -days $DAYS_VALID -subj "/C=IE/ST=Dublin/L=Dublin/O=NOISE EMAIL V1.0.0/CN=noise.email"
openssl x509 -in $CA_CERT -text -noout
echo "CA certificate generated successfully."