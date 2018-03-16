/*
* Copyright 2016-2018 Crunchy Data Solutions, Inc.
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package main

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	config "github.com/spf13/viper"

	"github.com/crunchydata/crunchy-watch/util"
)

func defaultStrategy() (string, error) {

	selectors := map[string]string{"name": config.GetString("CRUNCHY_WATCH_REPLICA")}

	podList, err := getPods(config.GetString(OSProject.EnvVar), nil, selectors)

	if err != nil {
		log.Error("Error getting pods command")
		return "", nil
	}

	// If not found then return an error
	if len(podList.Items) == 0 {
		return "", errors.New("No replicas found")

	}

	pod := podList.Items[0]
	return pod.Name, nil
}

func labelStrategy() (string, error) {

	selectors := map[string]string{"name": config.GetString("CRUNCHY_WATCH_REPLICA"), "replicatype": "trigger"}

	podList, err := getPods(config.GetString(OSProject.EnvVar), nil, selectors)

	if err != nil {
		log.Error("Error getting pods command")
		return "", nil
	}

	// If no 'trigger' replicas were found, then fall back to the default strategy.
	if len(podList.Items) == 0 {
		log.Info("No 'trigger' replicas were found, falling back to 'default' failover strategy")
		return defaultStrategy()
	}

	pod := podList.Items[0]
	return pod.Name, nil

}

func latestStrategy() (string, error) {
	selectors := map[string]string{"name": config.GetString("CRUNCHY_WATCH_REPLICA")}

	podList, err := getPods(config.GetString(OSProject.EnvVar), nil, selectors)

	if err != nil {
		log.Error("Error getting pods command")
		return "", nil
	}

	// If not found then return an error
	if len(podList.Items) == 0 {
		return "", errors.New("No replicas found")

	}

	// If only one replica exists then simply return it.
	if len(podList.Items) == 1 {
		pod := podList.Items[0]
		return pod.Name, nil
	}

	type ReplicaInfoName struct {
		*util.ReplicationInfo
		Name string
	}

	var value uint64 = 0
	var replicas []ReplicaInfoName
	var replicaInfoName ReplicaInfoName
	var replicaInfo *util.ReplicationInfo

	// Determine current replication status information for each replica
	for _, p := range podList.Items {

		target := fmt.Sprintf(
			"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
			config.GetString("CRUNCHY_WATCH_USERNAME"),
			config.GetString("CRUNCHY_WATCH_PASSWORD"),
			p.Status.PodIP,
			config.GetString("CRUNCHY_WATCH_REPLICA_PORT"),
			config.GetString("CRUNCHY_WATCH_DATABASE"),
		)

		replicaInfo, err = util.GetReplicationInfo(target)
		replicaInfoName = ReplicaInfoName{replicaInfo, p.Name}
		replicaInfoName.Name = p.Name
		replicas = append(replicas, replicaInfoName)

	}

	selectedReplica := replicas[0]
	for _, replica := range replicas {
		if replica.ReceiveLocation > value {
			value = replica.ReceiveLocation
			selectedReplica = replica
		}
	}

	return selectedReplica.Name, nil

}
