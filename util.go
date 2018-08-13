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
	"os/exec"
	"path/filepath"
	"plugin"

	log "github.com/sirupsen/logrus"
)

func loadPlatformModule(platform string) FailoverHandler {
	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))

	if err != nil {
		fmt.Println(err.Error())
	}

	pluginPath := fmt.Sprintf("%s/plugins/%s.so", currentDir, platform)
	plug, err := plugin.Open(pluginPath)

	if err != nil {
		fmt.Println(err.Error())
	}

	sym, err := plug.Lookup("FailoverHandler")

	if err != nil {
		fmt.Println(err.Error())
	}

	handler, ok := sym.(FailoverHandler)

	if !ok {
		log.Errorf("Could not load platform module: %s", platform)
		log.Error("Unexpected type from module symbol")
		os.Exit(1)
	}

	return handler
}

func validPlatform(platform string) bool {
	for _, p := range platforms {
		if p == platform {
			return true
		}
	}
	return false
}

func execute(command string) error {
	if command == "" {
		return nil
	}

	cmd := exec.Command("/bin/sh", command)

	return cmd.Run()
}
