package main

import (
	flag "github.com/spf13/pflag"
)

type failoverHandler struct{}

func init() {}

func (h failoverHandler) Failover() error {
	return nil
}

func (h failoverHandler) SetFlags(f *flag.FlagSet) {
	// No docker specific flags
}
func (h failoverHandler) Initialize() error {
	return nil
}

var FailoverHandler failoverHandler
