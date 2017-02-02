/*
 Copyright 2017 Crunchy Data Solutions, Inc.
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package watchapi

import (
	"log"
	"os"
	"strconv"
)

type Env struct {
	//required for kube and openshift
	PG_MASTER_SERVICE string
	//required for kube and openshift
	PG_SLAVE_SERVICE string
	//required, defaults to 5432
	PG_MASTER_PORT string
	//required, defaults to master
	PG_MASTER_USER string
	//required
	PG_MASTER_PASSWORD string
	//defaults to postgres
	PG_DATABASE string
	//optional, used for kube and openshift
	SLAVE_TO_TRIGGER_LABEL string
	//required, defaults to 10 seconds
	SLEEP_TIME int
	//required, defaults to 50 seconds
	WAIT_TIME    int
	PROJECT_TYPE string
}

const DOCKER_PROJECT = "docker"
const OSE_PROJECT = "ose"
const KUBE_PROJECT = "kube"

var EnvVars Env

var Logger *log.Logger

func init() {
	log.Println("initializing env...")
	EnvVars = Env{}
	Logger = log.New(os.Stdout, "logger: ", log.Lshortfile|log.Ldate|log.Ltime)
	EnvVars.SLEEP_TIME = 10
	EnvVars.WAIT_TIME = 50

}

func GetEnv() {

	var str string
	var err error

	str = os.Getenv("PROJECT_TYPE")
	if str != "" {
		EnvVars.PROJECT_TYPE = str
	}

	//for backward compat
	str = os.Getenv("KUBE_PROJECT")
	if str != "" {
		EnvVars.PROJECT_TYPE = KUBE_PROJECT
	} else {
		str = os.Getenv("OSE_PROJECT")
		if str != "" {
			EnvVars.PROJECT_TYPE = OSE_PROJECT
		} else {
			str = os.Getenv("DOCKER_PROJECT")
			if str != "" {
				EnvVars.PROJECT_TYPE = DOCKER_PROJECT
			}
		}
	}
	Logger.Println("PROJECT_TYPE set to " + EnvVars.PROJECT_TYPE)

	EnvVars.PG_MASTER_SERVICE = os.Getenv("PG_MASTER_SERVICE")
	if EnvVars.PG_MASTER_SERVICE == "" {
		log.Println("PG_MASTER_SERVICE is not supplied and is required")
		os.Exit(2)
	}
	Logger.Printf("PG_MASTER_SERVICE is %s\n", EnvVars.PG_MASTER_SERVICE)
	EnvVars.PG_SLAVE_SERVICE = os.Getenv("PG_SLAVE_SERVICE")
	if EnvVars.PG_SLAVE_SERVICE == "" {
		log.Println("PG_SLAVE_SERVICE is not supplied and is required")
		os.Exit(2)
	}
	Logger.Printf("PG_SLAVE_SERVICE is %s\n", EnvVars.PG_SLAVE_SERVICE)
	str = os.Getenv("PG_MASTER_PORT")
	if str != "" {
		_, err = strconv.Atoi(str)
		if err != nil {
			log.Println("PG_MASTER_PORT is not a valid integer")
			log.Println(err)
			os.Exit(2)
		}
		EnvVars.PG_MASTER_PORT = str
	}
	Logger.Printf("PG_MASTER_PORT is %s\n", EnvVars.PG_MASTER_PORT)
	str = os.Getenv("PG_MASTER_USER")
	if str != "" {
		EnvVars.PG_MASTER_USER = str
	}
	Logger.Printf("PG_MASTER_USER is %s\n", EnvVars.PG_MASTER_USER)
	str = os.Getenv("PG_DATABASE")
	if str != "" {
		EnvVars.PG_DATABASE = str
	}
	Logger.Printf("PG_DATABASE is %s\n", EnvVars.PG_DATABASE)
	str = os.Getenv("SLAVE_TO_TRIGGER_LABEL")
	if str != "" {
		EnvVars.SLAVE_TO_TRIGGER_LABEL = str
	}
	Logger.Printf("SLAVE_TO_TRIGGER_LABEL is %s\n", EnvVars.SLAVE_TO_TRIGGER_LABEL)
	str = os.Getenv("SLEEP_TIME")
	if str != "" {
		EnvVars.SLEEP_TIME, err = strconv.Atoi(str)
		if err != nil {
			log.Println("SLEEP_TIME is not a valid integer")
			log.Println(err)
			os.Exit(2)
		}
	}
	Logger.Printf("SLEEP_TIME is %d\n", EnvVars.SLEEP_TIME)
	str = os.Getenv("WAIT_TIME")
	if str != "" {
		EnvVars.WAIT_TIME, err = strconv.Atoi(str)
		if err != nil {
			log.Println("WAIT_TIME is not a valid integer")
			log.Println(err)
			os.Exit(2)
		}
	}
	Logger.Printf("WAIT_TIME is %d\n", EnvVars.WAIT_TIME)
	str = os.Getenv("PG_MASTER_PASSWORD")
	if str != "" {
		EnvVars.PG_MASTER_PASSWORD = str
	} else {
		log.Println("PG_MASTER_PASSWORD is required")
		os.Exit(2)
	}

}
