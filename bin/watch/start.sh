#!/bin/bash  -x

# Copyright 2017 Crunchy Data Solutions, Inc.
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

#
# this script is useful for development of the watch container
# or if you need to add some setup logic but as of right now, 1.2.8,
# the watchserver is run as the ENTRYPOINT directly in the container
#

function trap_sigterm() {
	echo "doing trap logic..."  >> /tmp/trap.out
}

trap 'trap_sigterm' SIGINT SIGTERM

export PATH=$PATH:/opt/cpm/bin

watchserver &

# this loop is for debugging only
while true; do 
	sleep 10
done

wait
