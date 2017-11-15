package main

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
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
		fmt.Println("Unexpected type from module symbol")
		os.Exit(1)
	}

	return handler
}

func checkPlatform(platform string) bool {
	for _, p := range platforms {
		if p == platform {
			return true
		}
	}
	return false
}
