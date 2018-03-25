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
	flag "github.com/spf13/pflag"
)

type failoverHandler struct{}

func init() {}

func (h failoverHandler) Failover() error {
	return nil
}

func (h failoverHandler) SetFlags(f *flag.FlagSet) error {
	// No docker specific flags
	return nil
}
func (h failoverHandler) Initialize() error {
	return nil
}

var FailoverHandler failoverHandler
