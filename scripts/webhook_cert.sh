#!/bin/bash

set -e

service="webhook-service"
namespace="system"

fullServiceDomain="${service}.${namespace}.svc"

if [ ${#fullServiceDomain} -gt 64 ] ; then
  echo "

  Common name exceeds the 64 character limit: ${fullServiceDomain}

  "
  exit 1
fi

if [ ! -x "$(command -v openssl)" ]; then
  echo "

  You need to install openssl to continue

  "
  exit 1
fi


dir="/var/folders/lc/g7hl9x1j2q308_wqs9msn0080000gn/T/k8s-webhook-server/serving-certs"

echo "

  Creating certs in dir ${dir}

"

mkdir -p ${dir}

prefix="${dir}/tls"

#openssl rand -base64 48 > passphrase.txt

# Generate a Private Key
openssl genrsa -aes128 -passout file:${prefix}.passphrase.txt -out ${prefix}.server.key 2048

# Generate a CSR (Certificate Signing Request)
openssl req -new -passin file:${prefix}.passphrase.txt -key ${prefix}.server.key -out ${prefix}.server.csr \
    -subj "/C=FR/O=krkr/OU=Domain Control Validated/CN=*.krkr.io"

# Remove Passphrase from Key
cp ${prefix}.server.key ${prefix}.server.key.org
openssl rsa -in ${prefix}.server.key.org -passin file:${prefix}.passphrase.txt -out ${prefix}.server.key

# Generating a Self-Signed Certificate for 100 years
openssl x509 -req -days 36500 -in ${prefix}.server.csr -signkey ${prefix}.server.key -out ${prefix}.server.crt

mv ${prefix}.server.crt ${prefix}.crt
mv ${prefix}.server.key ${prefix}.key
