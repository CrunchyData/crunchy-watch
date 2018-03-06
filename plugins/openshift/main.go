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
	OSProject = flags.FlagInfo{
		Namespace:   "openshift",
		Name:        "openshift-project",
		EnvVar:      "CRUNCHY_WATCH_OPENSHIFT_PROJECT",
		Description: "the openshift project",
	}

	OSFailoverStrategy = flags.FlagInfo{
		Namespace:   "openshift",
		Name:        "openshift-failover-strategy",
		EnvVar:      "CRUNCHY_WATCH_OPENSHIFT_FAILOVER_STRATEGY",
		Description: "the openshift failover strategy",
	}
)

var failoverStrategy = []string{
	"default",
	"label",
	"latest",
}

const (
	ocCmd string = "/opt/cpm/bin/oc"
)

func getReplica() (string, error) {
	switch config.GetString(OSFailoverStrategy.EnvVar) {
	case "default":
		return defaultStrategy()
	case "label":
		return labelStrategy()
	case "latest":
		return latestStrategy()
	default:
		return "", errors.New("invalid openshift failover strategy")
	}
}

func relabelReplica(replica string) error {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(ocCmd,
		"label",
		"pod",
		"--overwrite=true",
		fmt.Sprintf("--namespace=%s", config.GetString(OSProject.EnvVar)),
		replica,
		fmt.Sprintf("name=%s", config.GetString("CRUNCHY_WATCH_PRIMARY")),
	)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		log.Error(stderr.String())
	}

	log.Info(stdout.String())
	log.Info(stderr.String())

	return err
}

func deletePrimaryPod() error {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(ocCmd,
		"delete",
		"pod",
		fmt.Sprintf("--namespace=%s", config.GetString("CRUNCHY_WATCH_KUBE_NAMESPACE")),
		fmt.Sprintf("name=%s", config.GetString("CRUNCHY_WATCH_PRIMARY")),
	)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		log.Error(stderr.String())
	}

	log.Info(stdout.String())
	log.Info(stderr.String())

	return err
}

func promoteReplica(replica string) error {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(ocCmd,
		"exec",
		fmt.Sprintf("--namespace=%s", config.GetString(OSProject.EnvVar)),
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

	log.Info(stdout.String())
	log.Info(stderr.String())

	return err
}

func (h failoverHandler) Failover() error {
	log.Infof("Processing Failover: Strategy - %s",
		config.GetString(OSFailoverStrategy.EnvVar))

	// shoot the old primary in the head
	log.Info("Deleting existing primary...")
	err := deletePrimaryPod()

	if err != nil {
		log.Error(err)
		log.Error("An error occurred while deleting the old primary")
	}
	log.Info("Deleted old primary ")

	log.Info("Choosing failover replica...")
	replica, err := getReplica()

	if err != nil {
		log.Error("An error occurred while choosing the failover replica")
		return err
	}

	log.Infof("Chose failover target (%s)\n", replica)

	log.Info("Promoting failover replica...")
	err = promoteReplica(replica)

	if err != nil {
		log.Errorf("Cannot promote replica: %s", replica)
		return err
	}

	log.Info("Relabeling failover replica...")
	err = relabelReplica(replica)

	if err != nil {
		log.Errorf("Cannot relabel replica: %s", replica)
		return err
	}

	return nil
}

func (h failoverHandler) SetFlags(f *flag.FlagSet) {
	flags.String(f, OSProject, "default")
	flags.String(f, OSFailoverStrategy, "default")
}

func (h failoverHandler) Initialize() error {
	return nil
}

var FailoverHandler failoverHandler
