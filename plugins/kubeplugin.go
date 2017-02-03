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
	"strings"
	"time"
)

func KubeFailover() {
	api.Logger.Println("kube failover begins....")
	//get slaves
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
	rows := strings.Split(out.String(), "\n")
	var failoverTarget string
	var failoverFound = false
	allPods := []string{}
	for _, row := range rows {
		if len(row) > 0 {
			api.Logger.Println("row=[" + row + "]")
			words := strings.Split(row, " ")
			api.Logger.Println("name=[" + words[0] + "]")
			if failoverFound == false {
				failoverTarget = words[0]
				failoverFound = true
			}
			allPods = append(allPods, words[0])
		}
	}

	if api.EnvVars.SLAVE_TO_TRIGGER_LABEL != "" {
		api.Logger.Println("trigger to specific replica..using SLAVE_TO_TRIGGER_LABEL env var " + api.EnvVars.SLAVE_TO_TRIGGER_LABEL)
		failoverTarget = api.EnvVars.SLAVE_TO_TRIGGER_LABEL
	} else {
		api.Logger.Println("trigger to first replica which is  " + failoverTarget)
	}
	api.Logger.Println("creating the trigger file on " + failoverTarget)
	api.Logger.Println("targetSlave is " + failoverTarget)

	cmd = exec.Command("kubectl", "exec", failoverTarget, "touch", "/tmp/pg-failover-trigger")

	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		api.Logger.Println("error running kubectl exec touch " + err.Error())
		api.Logger.Println(out.String() + stderr.String())
		return
	}
	api.Logger.Println("trigger output ..." + out.String())
	api.Logger.Printf("sleeping %d %s\n", api.EnvVars.WAIT_TIME, " seconds to let failover process")

	time.Sleep(time.Duration(api.EnvVars.WAIT_TIME) * time.Second)
	api.Logger.Println("changing label of slave to " + api.EnvVars.PG_MASTER_SERVICE)
	//kubectl label --overwrite=true pod $i name=$PG_MASTER_SERVICE
	cmd = exec.Command("kubectl", "label", "--overwrite=true", "pod", failoverTarget, "name="+api.EnvVars.PG_MASTER_SERVICE)

	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		api.Logger.Println("error running kubectl label " + err.Error())
		api.Logger.Println(out.String() + stderr.String())
		return
	}
	api.Logger.Println("kubectl label output ..." + out.String())
	api.Logger.Println("deleting all other replica pods...")
	for _, pod := range allPods {
		if pod != failoverTarget {
			cmd = exec.Command("kubectl", "delete", "pod", pod)
			cmd.Stdout = &out
			cmd.Stderr = &stderr
			err = cmd.Run()
			if err != nil {
				api.Logger.Println("error running kubectl delete pod  on " + pod + err.Error())
				api.Logger.Println(out.String() + stderr.String())
			} else {
				api.Logger.Println("deleted pod " + pod)
			}
		}
	}

	api.Logger.Println("kube failover ends....")
}
