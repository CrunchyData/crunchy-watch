#!/bin/bash

# Copyright 2016-2018 Crunchy Data Solutions, Inc.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export CONTAINER_NAME=watch
$DIR/cleanup.sh

docker run \
	--privileged=true \
	--link primary:primary \
	--link replica:replica \
	-e CRUNCHY_WATCH_PLATFORM="docker" \
	-e CRUNCHY_WATCH_PRIMARY="primary:5432" \
	-e CRUNCHY_WATCH_REPLICA="replica:5432" \
	-e CRUNCHY_WATCH_USERNAME="primaryuser" \
	-e CRUNCHY_WATCH_PASSWORD="password" \
	-e CRUNCHY_WATCH_DATABASE="postgres" \
	-e CRUNCHY_WATCH_HEALTHCHECK_INTERVAL="30s" \
	-e CRUNCHY_WATCH_FAILOVER_WAIT="10s" \
	--name=$CONTAINER_NAME \
	--hostname=$CONTAINER_NAME \
	-d crunchydata/crunchy-watch:$CCP_IMAGE_TAG

# -e PG_FAILOVER_WAIT="10s" \
