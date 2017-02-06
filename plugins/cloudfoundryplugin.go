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
	consulapi "github.com/hashicorp/consul/api"
)

func CloudFoundryFailover() {
	api.Logger.Println("cloud foundry failover begins....")

	config := consulapi.DefaultConfig()
	config.Address = "192.168.0.7:8500"
	consul, err := consulapi.NewClient(config)
	if err == nil {
		api.Logger.Println("error getting CF client..." + err.Error())
		return
	}
	kv := consul.KV()
	if kv == nil {
		api.Logger.Println("error getting CF kv...")
		return
	}

	//get slaves
	//trigger failover on a slave
	//update proxy
	//delete old slaves
	//delete old master

	api.Logger.Println("cloud foundry failover ends....")
}
