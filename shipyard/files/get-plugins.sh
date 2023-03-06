#!/bin/bash

if [ ! -f /files/vault_plugins/boundary ]; then
  mkdir /files/vault_plugins
  curl -L -o boundary.zip https://github.com/hashicorp-dev-advocates/vault-plugin-boundary-secrets-engine/releases/download/v1.0.0/vault-plugin-boundary-secrets-engine_v1.0.0_linux_amd64.zip

  unzip boundary.zip -d .

  mv vault-plugin-boundary-secrets-engine_v1.0.0 /files/vault_plugins/boundary
fi
