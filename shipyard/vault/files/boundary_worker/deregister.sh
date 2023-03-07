#!/bin/sh -e
echo "[$(date +%T)] Deregister boundary worker"

# Read the worker id from the file written on startup
worker_id=$(cat ./worker_id)

# Base url for the HCP cluster
base_url="https://${cluster_id}.boundary.hashicorp.cloud/v1"
auth_url="${base_url}/auth-methods/${auth_method_id}:authenticate"
dereg_url="${base_url}/workers/${worker_id}"

# Authenticate with Boundary using the username and password and fetch the token
echo "[$(date +%T)] Authenticating with Boundary controller"
auth_request="{\"attributes\":{\"login_name\":\"${username}\",\"password\":\"${password}\"}}"
resp=$(curl ${auth_url} -s -d "${auth_request}")
token=$(echo ${resp} | sed 's/.*"token":"\([^"]*\)".*/\1/g')

# Deregister the worker
echo "[$(date +%T)] Calling boundary API to delete the worker ${worker_id}"
curl ${dereg_url} -s -H "Authorization: Bearer ${token}" -X DELETE

echo "[$(date +%T)] Deregistered worker: ${worker_id}"

# Remove the auth folder
echo "[$(date +%T)] Remove auth folder"
rm -rf /boundary/auth_data
