package main

import (
	"errors"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	config "github.com/spf13/viper"

	"github.com/crunchydata/crunchy-watch/flags"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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
	kubeConfig = flags.FlagInfo{
		Namespace:   "kubernetes",
		Name:        "kube-config-file",
		EnvVar:      "CRUNCHY_KUBE_CONFIG_FILE",
		Description: "kubernetes client config file",
	}
)

var client *kubernetes.Clientset
var restConfig *rest.Config

var failoverStrategies = []string{
	"default",
	"label",
	"latest",
}


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

func (h failoverHandler) Failover( ) error {
	log.Infof("Processing Failover: Strategy - %s",
		config.GetString(KubeFailoverStrategy.EnvVar))

	// shoot the old primary in the head
	log.Info("Deleting existing primary...")
	err := deletePrimaryPod(config.GetString("CRUNCHY_WATCH_KUBE_NAMESPACE"), config.GetString("CRUNCHY_WATCH_PRIMARY"))

	if err != nil {
		log.Error(err)
		log.Error("An error occurred while deleting the old primary")
	}
	log.Info("Deleted old primary ")

	// Get the list of replicas available
	log.Info("Choosing failover replica...")
	replica, err := getReplica()

	if err != nil {
		log.Error("An error occurred while choosing the failover replica")
		return err
	}
	log.Infof("Chose failover target (%s)\n", replica)

	// Promote replica to be new primary.
	log.Info("Promoting failover replica...")
	err = promoteReplica(config.GetString("CRUNCHY_WATCH_KUBE_NAMESPACE"), replica)

	if err != nil {

		log.Error("An error occurred while promoting the failover replica")
		return err
	}

	// Change labels so that the replica becomes the new primary.
	log.Info("Relabeling failover replica...")
	err = relabelReplica(config.GetString("CRUNCHY_WATCH_KUBE_NAMESPACE"), replica, config.GetString("CRUNCHY_WATCH_PRIMARY"))

	if err != nil {
		log.Error("An error occurred while relabeling the failover replica")
		return err
	}

	return nil
}

func (h failoverHandler) SetFlags(f *flag.FlagSet) {
	flags.String(f, KubeNamespace, "default")
	flags.String(f, KubeFailoverStrategy, "default")
	flags.String(f, kubeConfig, "")
}

func (h failoverHandler) Initialize() error {
	cfg, err := buildConfig(config.GetString(kubeConfig.EnvVar))
	if err != nil {
		log.Error("An error occurred initializing the client")
		return err
	}
	restConfig = cfg
	c, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Error("An error occurred initializing the client")
		return err
	}
	client = c
	return nil
}

func buildConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		log.Debugf("building config from flags %s\n", kubeConfig)
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	log.Debug("building config in cluster\n")
	return rest.InClusterConfig()
}

var FailoverHandler failoverHandler
