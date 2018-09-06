#!/bin/bash
# Copyright 2017 - 2018 Crunchy Data Solutions, Inc.
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

source ${WATCH_ROOT}/examples/common.sh
echo_info "Cleaning up.."

${WATCH_CLI?} delete --namespace=${WATCH_NAMESPACE?} pod pr-replica
${WATCH_CLI?} delete --namespace=${WATCH_NAMESPACE?} pod pr-replica-2
sleep  2
${WATCH_CLI?} delete --namespace=${WATCH_NAMESPACE?} service pr-replica
${WATCH_CLI?} delete --namespace=${WATCH_NAMESPACE?} service pr-primary
${WATCH_CLI?} delete --namespace=${WATCH_NAMESPACE?} pod pr-primary
$WATCH_ROOT/examples/waitforterm.sh pr-primary ${WATCH_CLI?}
$WATCH_ROOT/examples/waitforterm.sh pr-replica ${WATCH_CLI?}
$WATCH_ROOT/examples/waitforterm.sh pr-replica-2 ${WATCH_CLI?}
