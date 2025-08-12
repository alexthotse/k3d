#!/bin/bash

CURR_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
[ -d "$CURR_DIR" ] || { echo "FATAL: no current dir (maybe running in zsh?)";  exit 1; }

# shellcheck source=./common.sh
source "$CURR_DIR/common.sh"

export CURRENT_STAGE="local | remote_docker"


info "Starting dind with TLS (sleeping for 10s to give it time to get ready)"
$RUNTIME_CMD run -d -p 3376:2376 -e DOCKER_TLS_CERTDIR=/certs -v /tmp/dockercerts:/certs --privileged --rm --name k3dlocaltestdindsec docker:20.10-dind # Using docker image here is intentional for dind
sleep 10

info "Setting Docker Context (Skipped if RUNTIME_CMD is not 'docker')"
if [ "$RUNTIME_CMD" = "docker" ]; then
  $RUNTIME_CMD context create k3dlocaltestdindsec --description "dind local secure" --docker "host=tcp://127.0.0.1:3376,ca=/tmp/dockercerts/client/ca.pem,cert=/tmp/dockercerts/client/cert.pem,key=/tmp/dockercerts/client/key.pem"
  $RUNTIME_CMD context use k3dlocaltestdindsec
  $RUNTIME_CMD context list
else
  info "Skipping Docker context setup for $RUNTIME_CMD"
fi

info "Running k3d"
k3d_test_cmd cluster create test1
k3d_test_cmd cluster list

if [ "$RUNTIME_CMD" = "docker" ]; then
  info "Switching to default context"
  $RUNTIME_CMD context list
  $RUNTIME_CMD ps
  $RUNTIME_CMD context use default
  $RUNTIME_CMD ps
else
  info "Skipping Docker context switch for $RUNTIME_CMD. Current $RUNTIME_CMD ps output:"
  $RUNTIME_CMD ps
fi

info "Checking DOCKER_TLS env var based setting (Relevant mostly for Docker)"
export DOCKER_HOST=tcp://127.0.0.1:3376 # This will point k3d to the dind instance
export DOCKER_TLS_VERIFY=1
export DOCKER_CERT_PATH=/tmp/dockercerts/client

if [ "$RUNTIME_CMD" = "docker" ]; then
  $RUNTIME_CMD context list # Should show current context is overridden by DOCKER_HOST
fi
$RUNTIME_CMD ps
k3d_test_cmd cluster create test2 # k3d_test_cmd will add --runtime, DOCKER_HOST will be used by the selected runtime client
k3d_test_cmd cluster list
$RUNTIME_CMD ps

info "Cleaning up"
unset DOCKER_HOST
unset DOCKER_TLS_VERIFY
unset DOCKER_CERT_PATH
k3d_test_cmd cluster rm -a
if [ "$RUNTIME_CMD" = "docker" ]; then
  $RUNTIME_CMD context use default
  $RUNTIME_CMD rm -f k3dlocaltestdindsec
  $RUNTIME_CMD context rm k3dlocaltestdindsec
else
  # For podman, if k3dlocaltestdindsec container was started with podman, it needs podman rm
  $RUNTIME_CMD rm -f k3dlocaltestdindsec || true # true to prevent failure if not started with podman
  info "Skipped Docker context cleanup for $RUNTIME_CMD"
fi


info ">>> DONE <<<"
