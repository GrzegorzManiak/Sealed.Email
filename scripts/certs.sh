#!/bin/bash

clear
echo "!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!"
echo "!!! WARNING: DO NOT RUN WITH GOLAND / INTELIJ, IT WONT WORK !!!"
echo "!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!"

SERVICES=("notification" "domain" "api" "smtp")

CERTS_DIR="../.certs"
CA_DIR="$CERTS_DIR/ca"
CA_KEY="$CA_DIR/key.pem"
CA_CERT="$CA_DIR/ca.pem"
DAYS_VALID=365

# Check if CA certificate exists
if [[ ! -f "$CA_CERT" ]]; then
    echo "CA certificate not found at $CA_CERT. Please generate the CA certificate first."
    exit 1
fi

generate_service_cert() {
    # Generate certificate for a single service
    local service=$1
    echo "-------------[$service START]-------------"
    SERVICE_DIR="$CERTS_DIR/$service"
    echo "Generating certificate for $service..."
    mkdir -p "$SERVICE_DIR"
    pwd
    echo "$SERVICE_DIR"

    CRT_FILE="$SERVICE_DIR/$service.crt"
    CSR_FILE="$SERVICE_DIR/$service.csr"
    SUBJECT="/C=IE/ST=Dublin/L=Dublin/O=NOISE EMAIL V1.0.0/CN=$service.noise.email"


    # -- Generate CSR & Key
    openssl req -new -newkey ec -pkeyopt ec_paramgen_curve:prime256v1 -nodes -keyout "$SERVICE_DIR/$service.key" -out "$CSR_FILE" -subj "$SUBJECT"
    echo "CSR & Key for $service generated successfully."

    # -- Sign the CSR
    openssl x509 -req -in "$CSR_FILE" -CA "$CA_CERT" -CAkey "$CA_KEY" -CAcreateserial -out "$CRT_FILE" -days $DAYS_VALID
    echo "Certificate for $service generated successfully."
    echo "--------------[$service END]--------------"
}

for service in "${SERVICES[@]}"; do
    generate_service_cert "$service"
    echo
done

echo "All service certificates generated successfully."
