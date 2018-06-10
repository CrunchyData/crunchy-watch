package main

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	config "github.com/spf13/viper"

	"github.com/crunchydata/crunchy-watch/flags"
)

var flagSet *flag.FlagSet
var handler FailoverHandler

// Define common flag information
var (
	Primary = flags.FlagInfo{
		Namespace:   "general",
		Name:        "primary",
		EnvVar:      "CRUNCHY_WATCH_PRIMARY",
		Description: "host of the primary PostgreSQL instance",
	}

	PrimaryPort = flags.FlagInfo{
		Namespace:   "general",
		Name:        "primary-port",
		EnvVar:      "CRUNCHY_WATCH_PRIMARY_PORT",
		Description: "port of the primary PostgreSQL instance",
	}

	Replica = flags.FlagInfo{
		Namespace:   "general",
		Name:        "replica",
		EnvVar:      "CRUNCHY_WATCH_REPLICA",
		Description: "host of the replica PostgreSQL instance",
	}

	ReplicaPort = flags.FlagInfo{
		Namespace:   "general",
		Name:        "replica-port",
		EnvVar:      "CRUNCHY_WATCH_REPLICA_PORT",
		Description: "port of the replica PostgreSQL instance",
	}

	Username = flags.FlagInfo{
		Namespace:   "general",
		Name:        "username",
		EnvVar:      "CRUNCHY_WATCH_USERNAME",
		Description: "login user to connect to PostgreSQL",
	}

	Password = flags.FlagInfo{
		Namespace:   "general",
		Name:        "password",
		EnvVar:      "CRUNCHY_WATCH_PASSWORD",
		Description: "login user's password to connect to PostgreSQL",
	}

	Database = flags.FlagInfo{
		Namespace:   "general",
		Name:        "database",
		EnvVar:      "CRUNCHY_WATCH_DATABASE",
		Description: "database connect to",
	}

	Timeout = flags.FlagInfo{
		Namespace:   "general",
		Name:        "timeout",
		EnvVar:      "CRUNCHY_WATCH_TIMEOUT",
		Description: "connection timeout",
	}

	MaxFailures = flags.FlagInfo{
		Namespace:   "general",
		Name:        "max-failures",
		EnvVar:      "CRUNCHY_WATCH_MAX_FAILURES",
		Description: "maximum number of failures before performing failover",
	}

	HealthcheckInterval = flags.FlagInfo{
		Namespace:   "general",
		Name:        "healthcheck-interval",
		EnvVar:      "CRUNCHY_WATCH_HEALTHCHECK_INTERVAL",
		Description: "interval between healthchecks",
	}

	FailoverWait = flags.FlagInfo{
		Namespace:   "general",
		Name:        "failover-wait",
		EnvVar:      "CRUNCHY_WATCH_FAILOVER_WAIT",
		Description: "time to wait for failover to process",
	}

	PreHook = flags.FlagInfo{
		Namespace:   "general",
		Name:        "pre-hook",
		EnvVar:      "CRUNCHY_WATCH_PRE_HOOK",
		Description: "failover pre-hook to execute before processing failover",
	}

	PostHook = flags.FlagInfo{
		Namespace:   "general",
		Name:        "post-hook",
		EnvVar:      "CRUNCHY_WATCH_POST_HOOK",
		Description: "failover post-hook to execute after processing failover",
	}

	Debug = flags.FlagInfo{
		Namespace:   "general",
		Name:        "debug",
		EnvVar:      "CRUNCHY_DEBUG",
		Description: "when set to true, debug output is enabled",
	}
)

const (
	DefaultHealthCheckInterval = 10 * time.Second
	DefaultFailoverWait        = 50 * time.Second
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	flagSet = flag.NewFlagSet("main", flag.ContinueOnError)

	flags.String(flagSet, Primary, "")
	flags.Int(flagSet, PrimaryPort, 5432)
	flags.String(flagSet, Replica, "")
	flags.Int(flagSet, ReplicaPort, 5432)
	flags.String(flagSet, Database, "postgres")
	flags.String(flagSet, Username, "postgres")
	flags.String(flagSet, Password, "")
	flags.Duration(flagSet, Timeout, 10*time.Second)
	flags.Int(flagSet, MaxFailures, 0)
	flags.Duration(flagSet, HealthcheckInterval, DefaultHealthCheckInterval)
	flags.Duration(flagSet, FailoverWait, DefaultFailoverWait)
	flags.String(flagSet, PreHook, "")
	flags.String(flagSet, PostHook, "")
	flags.Bool(flagSet, Debug, false)
}

func findConfigOrFail(fi flags.FlagInfo, errorMessage string, exitCode int) {
	if config.GetString(fi.EnvVar) != "" {
		return
	}

	log.Error(errorMessage)
	log.Errorf("This value can be provided by either the '--%s' flag or the '%s' environment variable",
		fi.Name, fi.EnvVar)
	os.Exit(exitCode)
}
