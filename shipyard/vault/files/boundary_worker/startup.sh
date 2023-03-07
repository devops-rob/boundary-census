#!/bin/sh -e

# trap the interupt so that the worker deregisters on exit
trap "/boundary/deregister.sh" HUP INT QUIT TERM USR1

echo "[$(date +%T)] Generating controller led token for boundary worker"

# The name to use for the worker
worker_name="${worker_name}"

# The HCP cluster id, cluster id will be set in the system.d job as an environment var
cluster_id="${cluster_id}"

# Username and password used to obtain the worker registration token
username="${username}"
password="${password}"

# The auth id used for authentication
auth_method_id="${auth_method_id}"

# Base url for the HCP cluster
base_url="https://${cluster_id}.boundary.hashicorp.cloud/v1"
auth_url="${base_url}/auth-methods/${auth_method_id}:authenticate"
token_url="${base_url}/workers:create:controller-led"

# Authenticate with Boundary using the username and password and fetch the token
echo "[$(date +%T)] Authenticating with Boundary controller"
auth_request="{\"attributes\":{\"login_name\":\"${username}\",\"password\":\"${password}\"}}"
resp=$(curl ${auth_url} -s -d "${auth_request}")
token=$(echo ${resp} | sed 's/.*"token":"\([^"]*\)".*/\1/g')

# Generate the controller led token request
echo "[$(date +%T)] Calling boundary API to generate controller led token"
auth_request="{\"attributes\":{\"login_name\":\"${username}\",\"password\":\"${password}\"}}"
resp=$(curl ${token_url} -s -H "Authorization: Bearer ${token}" -d "{\"scope_id\":\"global\",\"name\":\"${worker_name}\"}")
controller_generated_activation_token=$(echo ${resp} | sed 's/.*"controller_generated_activation_token":"\([^"]*\)".*/\1/g')
worker_id=$(echo ${resp} | sed 's/{"id":"\([^"]*\)".*/\1/g')

# Write the worker id so we can use this to delete the worker on deallocation
echo "[$(date +%T)] Writing worker id file to ./worker_id"
echo ${worker_id} > ./worker_id

# Write the config
echo "[$(date +%T)] Writing config to ./worker_config.hcl"
cat <<-EOT > ./worker_config.hcl
  disable_mlock = true
  log_level = "debug"

  hcp_boundary_cluster_id = "${cluster_id}"

  listener "tcp" {
    address = "0.0.0.0:9202"
    purpose = "proxy"
  }

  worker {
    auth_storage_path="/boundary/auth_data"

    controller_generated_activation_token = "${controller_generated_activation_token}"
  
    tags {
      type   = ["vault"]
    }
  }
EOT

echo "[$(date +%T)] Generated worker config for worker: ${worker_id}"

boundary-worker server --config ./worker_config.hcl &
dpid=$!
wait $dpid