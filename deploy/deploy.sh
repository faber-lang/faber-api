#!/bin/bash

set -euo pipefail

readonly KEYPAIR_NAME="coordws"
readonly CLUSTER_CONFIG="faber-api"
readonly INSTANCE_TYPE="t2.small"
readonly SIZE=1

readonly BASE_DIR="$(dirname "$BASH_SOURCE")"
readonly COMPOSE_CONFIG="$BASE_DIR/../docker-compose.prod.yml"

function cluster_up() {
    ecs-cli up \
        --keypair "$KEYPAIR_NAME" \
        --capabilty-iam \
        --lanuch-type EC2 \
        --size $SIZE \
        --port 443 \
        --instance-type $INSTANCE_TYPE \
        "$@"
}

function cluster_down() {
    ecs-cli down "$@"
}

function service_up() {
    ecs-cli compose -f "$COMPOSE_CONFIG" service up "$@"
}

function service_down() {
    ecs-cli compose -f "$COMPOSE_CONFIG" service down "$@"
}

function main() {
    local func="$1_$2"
    shift; shift
    $func "$@"
}

main "$@"
