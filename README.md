# faber-api

[![Go Report Card](https://goreportcard.com/badge/github.com/coord-e/faber-api)](https://goreportcard.com/report/github.com/coord-e/faber-api)
![Docker Cloud Automated build](https://img.shields.io/docker/cloud/automated/coorde/faber-api.svg)
![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/coorde/faber-api.svg)
[![MicroBadger](https://images.microbadger.com/badges/image/coorde/faber-api.svg)](https://microbadger.com/images/coorde/faber-api)

an online compilation service for [coord-e/faber](https://github.com/coord-e/faber)

## deploy your own

[ecs-cli](https://github.com/aws/amazon-ecs-cli) must be configured properly before the deployment.

```shell
cp .env.sample .env
# edit .env
./deploy/deploy.sh cluster up
./deploy/deploy.sh service up
```
