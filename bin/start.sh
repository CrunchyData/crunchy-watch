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


PIDFILE=/tmp/crunchy-watch.pid

export PATH=$PATH:/opt/cpm/crunchy-watch

CMD=/opt/cpm/bin/crunchy-watch/crunchy-watch

function trap_sigterm() {
	echo "Shutting down crunchy-watch..."
	kill -SIGINT $(head -1 $PIDFILE)
}

trap trap_sigterm SIGINT SIGTERM

function check_deprecated() {
	if [ -v ${1} ]; then
		echo -e "${RED}WARNING:${NOCOLOR} ${1} is deprecated and will be removed in a future release, use ${2} instead."

		# if the value of $1 is a time/duration based value, then append
		# appropriate time unit to the value exported to $2
		if [ ${1} == "SLEEP_TIME" ] || [ ${1} == "WAIT_TIME" ]; then
			export ${2}=$(printf "%ss" $(printenv ${1}))
		else
			export ${2}=$(printenv ${1})
		fi
	fi
}

# Check for deprecated environment variables
check_deprecated PG_PRIMARY_SERVICE CRUNCHY_WATCH_PRIMARY
check_deprecated PG_PRIMARY_PORT CRUNCHY_WATCH_PRIMARY_PORT
check_deprecated PG_REPLICA_SERVICE CRUNCHY_WATCH_REPLICA
check_deprecated PG_PRIMARY_USER CRUNCHY_WATCH_USERNAME
check_deprecated PG_PASSWORD CRUNCHY_WATCH_PASSWORD
check_deprecated PG_DATABASE CRUNCHY_WATCH_DATABASE
check_deprecated WATCH_PRE_HOOK CRUNCHY_WATCH_PRE_HOOK
check_deprecated WATCH_POST_HOOK CRUNCHY_WATCH_POST_HOOK
check_deprecated SLEEP_TIME CRUNCHY_WATCH_HEALTHCHECK_INTERVAL
check_deprecated WAIT_TIME CRUNCHY_WATCH_FAILOVER_WAIT
check_deprecated MAX_FAILURES CRUNCHY_WATCH_MAX_FAILURES

if [ -v KUBE_PROJECT ]; then
    export CRUNCHY_WATCH_PLATFORM=kube
elif [ -v OSE_PROJECT ]; then
    export CRUNCHY_WATCH_PLATFORM=openshift
fi

ENVS=$(env | grep CRUNCHY_WATCH)

echo "Starting crunchy watch"
echo ${ENVS} | ${CMD} ${CRUNCHY_WATCH_PLATFORM} &
echo $! > ${PIDFILE}

echo "Started..."

wait
