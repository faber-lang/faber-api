#!/bin/bash

set -euo pipefail

# these variables will be set in .env
# readonly KEYPAIR_NAME="some_keypair"
# readonly INSTANCE_TYPE="t2.small"
# readonly SIZE=1
# readonly AWS_REGION=ap-northeast-1

readonly BASE_DIR="$(dirname "$BASH_SOURCE")/.."
readonly COMPOSE_CONFIG="$BASE_DIR/docker-compose.prod.yml"

source "$BASE_DIR/.env"

function cluster_up() {
    ecs-cli up \
        --keypair "$KEYPAIR_NAME" \
        --capability-iam \
        --launch-type EC2 \
        --size $SIZE \
        --port 443 \
        --instance-type $INSTANCE_TYPE \
        --region $AWS_REGION
        "$@"
}

function cluster_down() {
    ecs-cli down "$@"
}

function service_up() {
    ecs-cli compose -f "$COMPOSE_CONFIG" service up --create-log-groups "$@"
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
