#!/bin/bash

PIDFILE=/tmp/crunchy-watch.pid

export PATH=$PATH:/opt/cpm/crunchy-watch

function trap_sigterm() {
	echo "doing trap logic..."

	echo "Clean shutdown crunchy-watch..."
	kill -SIGINT $(head -1 $PIDFILE)
}

CMD=/opt/cpm/bin/crunchy-watch/crunchy-watch

ENVS=$(env | grep CRUNCHY_WATCH)

echo "Starting crunchy watch"
echo ${ENVS} | ${CMD} ${CRUNCHY_WATCH_PLATFORM} &
echo $! > ${PIDFILE}

echo "Started..."

wait
