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
	"fmt"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	config "github.com/spf13/viper"

	"github.com/crunchydata/crunchy-watch/util"
)

// Valid supported platforms.
var platforms = []string{
	"docker",
	"kube",
	"openshift",
}

type FailoverHandler interface {
	Failover() error
	SetFlags(*flag.FlagSet)
	Initialize() error
}

func main() {
	var pause bool

	if config.GetBool(Debug.EnvVar) {
		log.SetLevel(log.DebugLevel)
	}

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGUSR1,
		)

		log.Info("Waiting for signal...")
		s := <-ch
		log.Infof("Received signal: %s", s)
		if s == syscall.SIGUSR1 {
			pause = true
			failover()
			pause = false
		} else {
			os.Exit(0)
		}
	}()

	if len(os.Args) < 2 {
		errorExit()
	}

	log.SetLevel(log.DebugLevel)

	platform := os.Args[1]
	// Check that platform is valid.
	if !validPlatform(platform) {
		log.Errorf("Error: '%s' is not a valid platform\n\n", platform)
		errorExit()
	}

	// Load platform module
	log.Infof("Loading Platform Module: %s", platform)
	handler = loadPlatformModule(platform)
	// Allow platform module to add it's command-line flags
	handler.SetFlags(flagSet)

	// initialize the handler
	err := handler.Initialize()

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	// Parse all command-line flags
	err = flagSet.Parse(os.Args[2:])

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

	timeout := config.GetDuration(Timeout.EnvVar)

	// Construct connection string to primary
	target := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable&connect_timeout=%d",
		config.GetString(Username.EnvVar),
		config.GetString(Password.EnvVar),
		config.GetString(Primary.EnvVar),
		config.GetInt(PrimaryPort.EnvVar),
		config.GetString(Database.EnvVar),
		int(timeout.Seconds()),
	)

	// Watch
	failures := 0

	for {

		if pause == false {
			log.Infof("Health Checking: '%s'", config.GetString(Primary.EnvVar))
			err := util.HealthCheck(target)

			if err == nil {
				log.Infof("Successfully reached '%s'", config.GetString(Primary.EnvVar))
			} else {
				failures += 1

				log.Errorf("Could not reach '%s' (Attempt: %d)",
					config.GetString(Primary.EnvVar),
					failures,
				)

				// If max failure has been exceeded then process failover
				if failures > config.GetInt(MaxFailures.EnvVar) {
					failover()
					// reset retry count.
					failures = 0
				}
			}
		} else {
			log.Info("Health Checking paused")
		}
		time.Sleep(config.GetDuration(HealthcheckInterval.EnvVar))

	}
}

var inFailOver int32 = 0

func failover() {

	if atomic.CompareAndSwapInt32(&inFailOver, 0, 1) == false {
		return
	}

	// process failover pre-hook
	preHook := config.GetString(PreHook.EnvVar)
	if preHook != "" {
		log.Infof("Executing pre-hook: %s", preHook)
		err := execute(preHook)
		if err != nil {
			log.Error("Could not execute pre-hook")
			log.Error(err.Error())
		}
	}

	if handler != nil {

		// Process failover
		err := handler.Failover()

		if err != nil {
			log.Errorf("Failover process failed: %s", err.Error())
		}
	} else {
		log.Error("Failover process failed handler not initialized yet")
	}

	// Process failover post-hook
	postHook := config.GetString(PostHook.EnvVar)
	if postHook != "" {
		log.Infof("Executing post-hook: %s", postHook)
		err := execute(postHook)

		if err != nil {
			log.Error("Could not execute post-hook")
			log.Error(err.Error())

		}
	}
}
func errorExit() {
	log.Error("Usage: crunchy-watch <platform> [flags]")
	log.Error("Valid platform options are:")

	for _, p := range platforms {
		log.Errorf(" - %s", p)
	}

	os.Exit(1)
}
