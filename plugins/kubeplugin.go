/*
 Copyright 2017 Crunchy Data Solutions, Inc.
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package plugins

import (
	"bytes"
	api "github.com/crunchydata/crunchy-watch/watchapi"
	"os/exec"
)

func KubeFailover() {
	api.Logger.Println("kube failover begins....")
	api.Logger.Println("creating the trigger file on " + api.EnvVars.PG_MASTER_SERVICE)
	//get slaves
	//TRIGGERSLAVES=`kubectl get pod --selector=name=$PG_SLAVE_SERVICE --selector=slavetype=trigger --no-headers | cut -f1 -d' '`
	var cmd *exec.Cmd
	cmd = exec.Command("kubectl", "get", "pod", "--selector=name="+api.EnvVars.PG_SLAVE_SERVICE, "--no-headers")

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		api.Logger.Println("error running kubectl get pod " + err.Error())
		api.Logger.Println(out.String() + stderr.String())
		return
	}
	api.Logger.Println("trigger slaves found are ..." + out.String())

	targetSlave := api.EnvVars.SLAVE_TO_TRIGGER_LABEL
	if api.EnvVars.SLAVE_TO_TRIGGER_LABEL != "" {
		api.Logger.Println("trigger to specific replica..using SLAVE_TO_TRIGGER_LABEL env var " + api.EnvVars.SLAVE_TO_TRIGGER_LABEL)
	} else {
		api.Logger.Println("trigger to first replica ")
	}
	api.Logger.Println("targetSlave is " + targetSlave)

	api.Logger.Println("kube failover ends....")
}
