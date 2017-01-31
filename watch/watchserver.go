/*
 Copyright 2016 Crunchy Data Solutions, Inc.
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

package main

import (
	"github.com/crunchydata/crunchy-watch/watchapi"
	"log"
	"os"
	"time"
)

var POLL_INT = int64(3)

var logger *log.Logger

func main() {
	logger = log.New(os.Stdout, "logger: ", log.Lshortfile|log.Ldate|log.Ltime)
	var VERSION = os.Getenv("VERSION")

	logger.Println("watchserver " + VERSION + ": starting")

	for true {
		watchapi.DoSomething()
		time.Sleep(time.Duration(POLL_INT) * time.Minute)
	}

}
