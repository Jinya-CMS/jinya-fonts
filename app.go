package main

import (
	"flag"
	"jinya-fonts/config"
	"jinya-fonts/fontsync"
	"os"
)

func ContainsString(slice []string, search string) bool {
	for _, item := range slice {
		if item == search {
			return true
		}
	}

	return false
}

func main() {
	configFileFlag := flag.String("config-file", "./config.yaml", "The config file, check the sample for the structure")
	flag.Parse()

	configuration, err := config.LoadConfiguration(*configFileFlag)
	if err != nil {
		panic(err)
	}

	if ContainsString(os.Args, "sync") {
		err = fontsync.Sync(configuration)
		if err != nil {
			panic(err)
		}
	}
}
