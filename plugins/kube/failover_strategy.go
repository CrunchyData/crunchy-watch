package main

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
	config "github.com/spf13/viper"

	"github.com/crunchydata/crunchy-watch/util"
	"k8s.io/kubernetes/pkg/apis/core/pods"
)

/*
	return the name of the first pod with the name=CRUNCHY_WATCH_REPLICA label
 */
func defaultStrategy() (string, error) {


	selectors := map[string]string{"name":config.GetString("CRUNCHY_WATCH_REPLICA")}

	podList, err := getPods(config.GetString("CRUNCHY_WATCH_KUBE_NAMESPACE"), nil, selectors)

	if err != nil {
		log.Error("Error getting pods command")
		return "",nil
	}

	// If not found then return an error
	if len(podList.Items) == 0 {
		return "", errors.New("No replicas found")

	}

	pod := podList.Items[0]
	return pod.Name, nil

}

/*
	return the name of the first pod named with CRUNCHY_WATCH_REPLICA
	and replicatype trigger
	if nothing is found then use the default strategy
 */
func labelStrategy() (string, error) {

	selectors := map[string]string{"name":config.GetString("CRUNCHY_WATCH_REPLICA"), "replicatype":"trigger"}

	podList, err := getPods(config.GetString("CRUNCHY_WATCH_KUBE_NAMESPACE"), nil, selectors)

	if err != nil {
		log.Error("Error getting pods command")
		return "",nil
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
	selectors := map[string]string{"name":config.GetString("CRUNCHY_WATCH_REPLICA")}

	podList, err := getPods(config.GetString("CRUNCHY_WATCH_KUBE_NAMESPACE"), nil, selectors)

	if err != nil {
		log.Error("Error getting pods command")
		return "",nil
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
	var replicas [len(podList.Items)]ReplicaInfoName
	var i=0

	selectedReplica := replicas[0]
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

		replicas[i].ReplicationInfo, err = util.GetReplicationInfo(target)
		replicas[i].Name = p.Name

	}

	for _, replica := range replicas {
		if replica.ReceiveLocation > value {
			value = replica.ReceiveLocation
			selectedReplica = replica
		}
	}

	return selectedReplica.Name, nil
}
