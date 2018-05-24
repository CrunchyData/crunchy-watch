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

package flags

import (
	"time"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type FlagInfo struct {
	Name        string
	Shorthand   string
	EnvVar      string
	Namespace   string
	Description string
}

func String(f *flag.FlagSet, flagInfo FlagInfo, defaultValue string) {
	f.String(flagInfo.Name, defaultValue, flagInfo.Description)
	viper.BindEnv(flagInfo.EnvVar)
	viper.BindPFlag(flagInfo.EnvVar, f.Lookup(flagInfo.Name))
}

func Int(f *flag.FlagSet, flagInfo FlagInfo, defaultValue int) {
	f.Int(flagInfo.Name, defaultValue, flagInfo.Description)
	viper.BindEnv(flagInfo.EnvVar)
	viper.BindPFlag(flagInfo.EnvVar, f.Lookup(flagInfo.Name))
}

func Duration(f *flag.FlagSet, flagInfo FlagInfo, defaultValue time.Duration) {
	f.Duration(flagInfo.Name, defaultValue, flagInfo.Description)
	viper.BindEnv(flagInfo.EnvVar)
	viper.BindPFlag(flagInfo.EnvVar, f.Lookup(flagInfo.Name))
}

func Bool(f *flag.FlagSet, flagInfo FlagInfo, defaultValue bool) {
	f.Bool(flagInfo.Name, defaultValue, flagInfo.Description)
	viper.BindEnv(flagInfo.EnvVar)
	viper.BindPFlag(flagInfo.EnvVar, f.Lookup(flagInfo.Name))
}
