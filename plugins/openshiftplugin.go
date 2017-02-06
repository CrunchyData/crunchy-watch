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

func OpenshiftFailover() {
	api.Logger.Println("openshift failover begins....")
	//get slaves
	var cmd, cmd2, cmd3, cmd4, cmd5 *exec.Cmd
	var out, out2, out3, out4, out5 bytes.Buffer
	var stderr, stderr2, stderr3, stderr4, stderr5 bytes.Buffer
	var err error

	api.Logger.Println("setting project to " + api.EnvVars.NAMESPACE)
	cmd = exec.Command("/opt/cpm/bin/oc", "project", api.EnvVars.NAMESPACE)

	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		api.Logger.Println("error running oc project " + api.EnvVars.NAMESPACE + err.Error())
		api.Logger.Println(out.String() + stderr.String())
		return
	}
	api.Logger.Println("oc project output..." + out.String())

	cmd2 = exec.Command("/opt/cpm/bin/oc", "get", "pod", "--selector=name="+api.EnvVars.PG_SLAVE_SERVICE, "--no-headers")

	cmd2.Stdout = &out2
	cmd2.Stderr = &stderr2
	err = cmd2.Run()
	if err != nil {
		api.Logger.Println("error running oc get pod " + err.Error())
		api.Logger.Println(out2.String() + stderr2.String())
		return
	}
	api.Logger.Println("trigger slaves found are ..." + out2.String())
	rows := strings.Split(out2.String(), "\n")
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

	cmd3 = exec.Command("/opt/cpm/bin/oc", "exec", failoverTarget, "touch", "/tmp/pg-failover-trigger")

	cmd3.Stdout = &out3
	cmd3.Stderr = &stderr3
	err = cmd3.Run()
	if err != nil {
		api.Logger.Println("error running oc exec touch " + err.Error())
		api.Logger.Println(out3.String() + stderr3.String())
		return
	}
	api.Logger.Println("trigger output ..." + out3.String())
	api.Logger.Printf("sleeping %d %s\n", api.EnvVars.WAIT_TIME, " seconds to let failover process")

	time.Sleep(time.Duration(api.EnvVars.WAIT_TIME) * time.Second)
	api.Logger.Println("changing label of slave to " + api.EnvVars.PG_MASTER_SERVICE)
	//oc label --overwrite=true pod $i name=$PG_MASTER_SERVICE
	cmd4 = exec.Command("/opt/cpm/bin/oc", "label", "--overwrite=true", "pod", failoverTarget, "name="+api.EnvVars.PG_MASTER_SERVICE)

	cmd4.Stdout = &out4
	cmd4.Stderr = &stderr4
	err = cmd4.Run()
	if err != nil {
		api.Logger.Println("error running oc label " + err.Error())
		api.Logger.Println(out4.String() + stderr4.String())
		return
	}
	api.Logger.Println("oc label output ..." + out4.String())
	api.Logger.Println("deleting all other replica pods...")
	for _, pod := range allPods {
		if pod != failoverTarget {
			cmd5 = exec.Command("/opt/cpm/bin/oc", "delete", "pod", pod)
			cmd5.Stdout = &out5
			cmd5.Stderr = &stderr5
			err = cmd5.Run()
			if err != nil {
				api.Logger.Println("error running oc delete pod  on " + pod + err.Error())
				api.Logger.Println(out5.String() + stderr5.String())
			} else {
				api.Logger.Println("deleted pod " + pod)
			}
		}
	}

	api.Logger.Println("openshift failover ends....")
}
