package main

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	config "github.com/spf13/viper"

	"github.com/crunchydata/crunchy-watch/flags"
	"github.com/crunchydata/crunchy-watch/util"
)

// Valid supported platforms.
var platforms = []string{
	"docker",
	"kube",
	"openshift",
}

var flagSet *flag.FlagSet

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
)

const (
	DefaultHealthCheckInterval = 10 * time.Second
	DefaultFailoverWait        = 50 * time.Second
)

type FailoverHandler interface {
	Failover() error
	SetFlags(*flag.FlagSet)
}

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
	flags.Int(flagSet, Timeout, 10)
	flags.Int(flagSet, MaxFailures, 0)
	flags.Duration(flagSet, HealthcheckInterval, DefaultHealthCheckInterval)
	flags.Duration(flagSet, FailoverWait, DefaultFailoverWait)
}

func main() {
	platform := os.Args[1]
	validPlatform := checkPlatform(platform)

	// Check that platform is valid.
	if !validPlatform {
		log.Error("Usage: crunchy-watch <platform> [flags]")
		log.Errorf("Error: '%s' is not a valid platform\n\n", platform)
		log.Error("Valid platform options are:")

		for _, p := range platforms {
			log.Errorf(" - %s", p)
		}

		os.Exit(1)
	}

	// Load platform module
	log.Infof("Loading Platform Module: %s", platform)
	handler := loadPlatformModule(platform)

	// Allow platform module to add it's command-line flags
	handler.SetFlags(flagSet)

	// Parse all command-line flags
	err := flagSet.Parse(os.Args[2:])

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	// Check that required flags/envs were set
	if config.GetString(Primary.EnvVar) == "" {
		log.Error("Must specify a primary PostgreSQL instance.")
		log.Errorf("This value can be provided by either the '--%s' flag or the '%s' environment variable",
			Primary.Name, Primary.EnvVar)
		os.Exit(1)
	}

	if config.GetString(Replica.EnvVar) == "" {
		log.Error("Must specify a replica PostgreSQL instance for failover.")
		log.Errorf("This value can be provided by either the '--%s' flag or the '%s' environment variable",
			Replica.Name, Replica.EnvVar)
		os.Exit(1)
	}

	// Construct connection string to primary
	target := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable&connect_timeout=%d",
		config.GetString(Username.EnvVar),
		config.GetString(Password.EnvVar),
		config.GetString(Primary.EnvVar),
		config.GetInt(PrimaryPort.EnvVar),
		config.GetString(Database.EnvVar),
		config.GetInt(Timeout.EnvVar),
	)

	// Watch
	failures := 0

	for {
		log.Infof("Health Checking: '%s'", config.GetString(Primary.EnvVar))
		err := util.HealthCheck(target)

		if err == nil {
			log.Infof("Successfully reached '%s'", config.GetString(Primary.EnvVar))
		} else {
			log.Errorf("Could not reach '%s'", config.GetString(Primary.EnvVar))
			log.Debug(err.Error())

			failures += 1

			// If max failure has been exceeded then process failover
			if failures > config.GetInt(MaxFailures.EnvVar) {
				err := handler.Failover()

				if err != nil {
					log.Error(err.Error())
				}

				// reset retry count.
				failures = 0
			}
		}

		time.Sleep(config.GetDuration(HealthcheckInterval.EnvVar))
	}
}
