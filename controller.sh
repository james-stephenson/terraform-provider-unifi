#!/usr/bin/env bash

set -eou pipefail

default_tag="stable-6"

# Local Administrator
# Username: tfacctest
# Password: tfacctest1234
# Email: tfacctest@example.com
DOCKER_HTTP_PORT="${DOCKER_HTTP_PORT:-8080}"
DOCKER_HTTPS_PORT="${DOCKER_HTTPS_PORT:-8443}"
DOCKER_STUN_PORT="${DOCKER_STUN_PORT:-3478}"
DOCKER_AIRCONTROL_PORT="${DOCKER_AIRCONTROL_PORT:-10001}"

if test $# -eq 0; then
  echo "please specify either 'start' or 'test'"
  exit 1
fi

case "$1" in
  "start")
    if [[ ! -d testata/out ]]; then
      cp -r testdata/{unifi,out}
    fi

    if [[ $(( $(docker ps --filter=name=unifi -q | wc -l) )) == 0 ]]; then
      docker run --rm --init -d \
        -p ${DOCKER_HTTP_PORT}:8080 \
        -p ${DOCKER_HTTPS_PORT}:8443 \
        -p ${DOCKER_STUN_PORT}:3478/udp \
        -p ${DOCKER_AIRCONTROL_PORT}:10001/udp \
        -e TZ='America/New_York' \
        -v $(pwd)/testdata/out:/unifi \
        --name unifi \
        jacobalberty/unifi:${2:-$default_tag}
    fi

    echo "Waiting for login page..."
    timeout 300 bash -c 'while [[ "$(curl --insecure -s -o /dev/null -w "%{http_code}" '"https://localhost:${DOCKER_HTTPS_PORT}/manage/account/login"')" != "200" ]]; do sleep 5; done'
    echo "Controller running."
    ;;
  "test")
    TF_ACC=1 \
    UNIFI_USERNAME=tfacctest \
    UNIFI_PASSWORD=tfacctest1234 \
    UNIFI_API="https://localhost:${DOCKER_HTTPS_PORT}/" \
    UNIFI_ACC_WLAN_CONCURRENCY="4" \
    UNIFI_INSECURE="true" \
    go test -v -cover -count 1 ./internal/provider
    ;;
  "stop")
    docker stop unifi
    ;;
  "update")
    docker pull jacobalberty/unifi:${2:-$default_tag}
    ;;
  "reset")
    rm -rf testdata/out
    ;;
  *)
    echo "unrecognized command"
    exit 1
    ;;
esac
