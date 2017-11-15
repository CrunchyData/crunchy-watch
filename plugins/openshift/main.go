package main

import (
	"fmt"

	flag "github.com/spf13/pflag"

	"github.com/crunchydata/crunchy-watch/flags"
)

type failoverHandler struct{}

var (
	Namespace = flags.FlagInfo{
		Name:        "openshift-namespace",
		EnvVar:      "CRUNCHY_WATCH_OPENSHIFT_NAMESPACE",
		Description: "openshift namespace",
	}
)

func (h failoverHandler) Hello() {
	fmt.Println("Hello From Openshift!")
}

func (h failoverHandler) SetFlags(f *flag.FlagSet) {
	flags.String(f, Namespace, "")
}

var FailoverHandler failoverHandler
