FROM hashicorp/boundary-worker-hcp

RUN apk add curl
RUN curl https://releases.hashicorp.com/vault/1.12.4/vault_1.12.4_linux_amd64.zip -o /tmp/vault.zip && \
  unzip /tmp/vault.zip && \
  mv ./vault /usr/bin

