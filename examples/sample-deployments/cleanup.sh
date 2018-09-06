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

source $WATCH_ROOT/examples/common.sh
echo_info "Cleaning up.."

${WATCH_CLI?} delete --namespace=${WATCH_NAMESPACE?} deployment watchprimary watchreplica
${WATCH_CLI?} delete --namespace=${WATCH_NAMESPACE?} pod --selector=name=watchprimary
${WATCH_CLI?} delete --namespace=${WATCH_NAMESPACE?} pod --selector=name=watchreplica
${WATCH_CLI?} delete --namespace=${WATCH_NAMESPACE?} secret watchprimary-secret
${WATCH_CLI?} delete --namespace=${WATCH_NAMESPACE?} service watchprimary watchreplica
${WATCH_CLI?} delete --namespace=${WATCH_NAMESPACE?} pvc --selector=name=watchprimary
${WATCH_CLI?} delete --namespace=${WATCH_NAMESPACE?} pvc --selector=name=watchreplica

if [ -z "$CCP_STORAGE_CLASS" ]
then
  ${WATCH_CLI?} delete --namespace=${WATCH_NAMESPACE?} pv watchprimary-pgdata watchreplica-pgdata 
fi

dir_check_rm "watchprimary"
dir_check_rm "watchreplica"
