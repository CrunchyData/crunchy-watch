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
)

func defaultStrategy() (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	replicas := make([]string, 0)

	cmd := exec.Command(ocCmd,
		"get",
		"pod",
		fmt.Sprintf("--namespace=%s", config.GetString(OSNamespace.EnvVar)),
		fmt.Sprintf("--selector=name=%s",
			config.GetString("CRUNCHY_WATCH_REPLICA")),
		"--no-headers",
	)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		return "", err
	}

	rows := strings.Split(stdout.String(), "\n")

	if len(rows) == 0 {
		return "", errors.New("No replicas found")
	}

	for _, row := range rows {
		if len(row) > 0 {
			pod := strings.Split(row, " ")
			replicas = append(replicas, pod[0])
		}
	}

	return replicas[0], nil
}

func labelStrategy() (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	replicas := make([]string, 0)

	// Get the list of replicas with the 'trigger' label.
	cmd := exec.Command(ocCmd,
		"get",
		"pod",
		fmt.Sprintf("--namespace=%s", config.GetString(OSNamespace.EnvVar)),
		fmt.Sprintf("--selector=name=%s,replicatype=trigger",
			config.GetString("CRUNCHY_WATCH_REPLICA")),
		"--no-headers",
	)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		log.Error("Error running 'oc' command")
		log.Error(stderr.String())
		return "", err
	}

	rows := strings.Split(stdout.String(), "\n")

	// If no 'trigger' replicas were found, then fall back to the default
	// strategy.
	if len(rows) == 0 {
		log.Info("No 'trigger' replicas were found, falling back to the 'default' failover strategy")
		return defaultStrategy()
	}

	for _, row := range rows {
		if len(row) > 0 {
			pod := strings.Split(row, " ")
			replicas = append(replicas, pod[0])
		}
	}

	return replicas[0], nil
}

func latestStrategy() (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	replicas := make([]util.Replica, 0)

	// Get the list of replicas.
	cmd := exec.Command(ocCmd,
		"get",
		"pod",
		fmt.Sprintf("--namespace=%s", config.GetString(OSNamespace.EnvVar)),
		fmt.Sprintf("--selector=name=%s",
			config.GetString("CRUNCHY_WATCH_REPLICA")),
		"--output=custom-columns=NAME:.metadata.name,IP:.status.podIP",
		"--no-headers",
	)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		log.Error("Error running 'oc' command")
		log.Error(stderr.String())
		return "", err
	}

	rows := strings.Split(stdout.String(), "\n")

	for _, row := range rows {
		if len(row) > 0 {
			pod := strings.Fields(row)
			replica := util.Replica{Name: pod[0], IP: pod[1]}
			replicas = append(replicas, replica)
		}
	}

	if len(replicas) == 0 {
		return "", errors.New("No replicas found")
	}

	// If only one replica exists then simply return it.
	if len(replicas) == 1 {
		return replicas[0].Name, nil
	}

	// Determine current replication status information for each replica
	for i, _ := range replicas {
		target := fmt.Sprintf(
			"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
			config.GetString("CRUNCHY_WATCH_USERNAME"),
			config.GetString("CRUNCHY_WATCH_PASSWORD"),
			replicas[i].IP,
			config.GetString("CRUNCHY_WATCH_REPLICA_PORT"),
			config.GetString("CRUNCHY_WATCH_DATABASE"),
		)

		replicas[i].Status, err = util.GetReplicationInfo(target)

		if err != nil {
			log.Error("Could not determine replication status for replica")
			return "", err
		}
	}

	var value uint64 = 0
	selectedReplica := replicas[0]

	for _, replica := range replicas {
		if replica.Status.ReceiveLocation > value {
			value = replica.Status.ReceiveLocation
			selectedReplica = replica
		}
	}

	return selectedReplica.Name, nil
}
