#!/bin/bash
export PROJECT_TYPE=kube
export PG_MASTER_SERVICE=master
export PG_SLAVE_SERVICE=replica
export PG_MASTER_PORT=5432
export PG_MASTER_USER=master
export PG_DATABASE=postgres
export SLAVE_TO_TRIGGER_LABEL=
export SLEEP_TIME=10
export WAIT_TIME=50
export DOCKER_API_VERSION=1.12
go run watchserver.go
