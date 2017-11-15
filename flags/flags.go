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
