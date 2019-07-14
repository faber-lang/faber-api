# faber-api

[![Go Report Card](https://goreportcard.com/badge/github.com/faber-lang/faber-api)](https://goreportcard.com/report/github.com/faber-lang/faber-api)
![Docker Cloud Automated build](https://img.shields.io/docker/cloud/automated/faberlang/faber-api.svg)
![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/faberlang/faber-api.svg)
[![MicroBadger](https://images.microbadger.com/badges/image/faberlang/faber-api.svg)](https://microbadger.com/images/faberlang/faber-api)

an online compilation service for [faber-lang/faber](https://github.com/faber-lang/faber)

## deploy your own

[ecs-cli](https://github.com/aws/amazon-ecs-cli) must be configured properly before the deployment.

```shell
cp .env.sample .env
# edit .env
./deploy/deploy.sh cluster up
./deploy/deploy.sh service up
```
