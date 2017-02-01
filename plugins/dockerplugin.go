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
	api "github.com/crunchydata/crunchy-watch/watchapi"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func DockerFailover() {
	api.Logger.Println("docker failover begins....")
	api.Logger.Println("creating the trigger file on " + api.EnvVars.PG_MASTER_SERVICE)
	//docker exec $PG_SLAVE_SERVICE touch /tmp/pg-failover-trigger

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		api.Logger.Printf("%s %s\n", container.ID[:10], container.Image)
	}
	api.Logger.Println("docker failover ends....")
}
