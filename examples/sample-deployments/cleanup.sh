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

${CCP_CLI?} delete --namespace=${CCP_NAMESPACE?} deployment watchprimary watchreplica
${CCP_CLI?} delete --namespace=${CCP_NAMESPACE?} pod --selector=name=watchprimary
${CCP_CLI?} delete --namespace=${CCP_NAMESPACE?} pod --selector=name=watchreplica
${CCP_CLI?} delete --namespace=${CCP_NAMESPACE?} secret watchprimary-secret
${CCP_CLI?} delete --namespace=${CCP_NAMESPACE?} service watchprimary watchreplica
${CCP_CLI?} delete --namespace=${CCP_NAMESPACE?} pvc --selector=name=watchprimary
${CCP_CLI?} delete --namespace=${CCP_NAMESPACE?} pvc --selector=name=watchreplica

if [ -z "$CCP_STORAGE_CLASS" ]
then
  ${CCP_CLI?} delete --namespace=${CCP_NAMESPACE?} pv watchprimary-pgdata watchreplica-pgdata 
fi

dir_check_rm "watchprimary"
dir_check_rm "watchreplica"
