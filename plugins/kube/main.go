package main

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	config "github.com/spf13/viper"

	"github.com/crunchydata/crunchy-watch/flags"
)

type failoverHandler struct{}

var (
	KubeNamespace = flags.FlagInfo{
		Namespace:   "kubernetes",
		Name:        "kube-namespace",
		EnvVar:      "CRUNCHY_WATCH_KUBE_NAMESPACE",
		Description: "the kubernetes namespace",
	}

	KubeFailoverStrategy = flags.FlagInfo{
		Namespace:   "kubernetes",
		Name:        "kube-failover-strategy",
		EnvVar:      "CRUNCHY_WATCH_KUBE_FAILOVER_STRATEGY",
		Description: "the kubernetes failover strategy",
	}
)

var failoverStrategies = []string{
	"default",
	"label",
	"latest",
}

const (
	kubectlCmd string = "/opt/cpm/bin/kubectl"
)

func getReplica() (string, error) {
	switch config.GetString(KubeFailoverStrategy.EnvVar) {
	case "default":
		return defaultStrategy()
	case "label":
		return labelStrategy()
	case "latest":
		return latestStrategy()
	default:
		return "", errors.New("invalid kubernetes failover strategy")
	}
}

func relabelReplica(replica string) error {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(kubectlCmd,
		"label",
		"pod",
		"--overwrite=true",
		fmt.Sprintf("--namespace=%s", config.GetString("CRUNCHY_WATCH_KUBE_NAMESPACE")),
		replica,
		fmt.Sprintf("name=%s", config.GetString("CRUNCHY_WATCH_PRIMARY")),
	)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		log.Error(stderr.String())
	}

	return err
}

func promoteReplica(replica string) error {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(kubectlCmd,
		"exec",
		fmt.Sprintf("--namespace=%s", config.GetString("CRUNCHY_WATCH_KUBE_NAMESPACE")),
		replica,
		"touch",
		"/tmp/pg-failover-trigger",
	)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		log.Error(stderr.String())
	}

	return err
}

func (h failoverHandler) Failover() error {
	log.Infof("Processing Failover: Strategy - %s",
		config.GetString(KubeFailoverStrategy.EnvVar))

	// Get the list of replicas available
	log.Info("Choosing failover replica...")
	replica, err := getReplica()

	if err != nil {
		log.Error("An error occurred while choosing the failover replica")
		return err
	}

	// Promote replica to be new primary.
	log.Info("Promoting failover replica...")
	err = promoteReplica(replica)

	if err != nil {
		log.Error("An error occurred while promoting the failover replica")
		return err
	}

	// Change labels so that the replica becomes the new primary.
	log.Info("Relabeling failover replica...")
	err = relabelReplica(replica)

	if err != nil {
		log.Error("An error occurred while relabeling the failover replica")
		return err
	}

	return nil
}

func (h failoverHandler) SetFlags(f *flag.FlagSet) {
	flags.String(f, KubeNamespace, "default")
	flags.String(f, KubeFailoverStrategy, "default")
}

var FailoverHandler failoverHandler
