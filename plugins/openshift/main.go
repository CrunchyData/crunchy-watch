package main

import (

	"errors"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	config "github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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

var client *kubernetes.Clientset
var restConfig *rest.Config

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



func (h failoverHandler) Failover() error {
	log.Infof("Processing Failover: Strategy - %s",
		config.GetString(OSFailoverStrategy.EnvVar))

	// shoot the old primary in the head
	log.Info("Deleting existing primary...")
	err := deletePrimaryPod(config.GetString(OSProject.EnvVar), config.GetString("CRUNCHY_WATCH_PRIMARY"))

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
	err = promoteReplica(config.GetString(OSProject.EnvVar), replica)

	if err != nil {
		log.Errorf("Cannot promote replica: %s", replica)
		return err
	}

	log.Info("Relabeling failover replica...")
	err = relabelReplica(config.GetString(OSProject.EnvVar), replica, config.GetString("CRUNCHY_WATCH_PRIMARY"))

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
