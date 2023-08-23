package main

import (
	"flag"
	"fmt"
	"github.com/jgrtowy/yamst/lib"
	"os"
)

func main() {
	updateManifest := flag.Bool("u", false, "Update the manifest")
	installServer := flag.Bool("i", false, "Install the server")
	version := flag.String("v", "", "Version to install")
	clearCache := flag.Bool("c", false, "Clear the cache")
	defaultSettings := flag.Bool("d", false, "Apply default settings")
	help := flag.Bool("h", false, "Show help")
	serverType := flag.String("t", "", "Server type")
	flag.Parse()

	if *serverType == "" {
		*serverType = "vanilla"
	}

	if *version == "" {
		var err error
		*version, err = lib.GetLatest()
		if err != nil {
			fmt.Printf("Error getting latest version: %s \n", err)
		}
	}

	workingDir, _ := os.Getwd()

	if *help {
		fmt.Println("Usage: yamst [options]")
		flag.PrintDefaults()
		return
	}

	if *clearCache {
		fmt.Println("Clearing cache...")
		err := os.RemoveAll(lib.GetCacheDirectory())
		if err != nil {
			fmt.Errorf("Error clearing cache: %s \n", err)
		}
		fmt.Println("Successfully cleared cache")
	}

	if *updateManifest {
		err := lib.UpdateManifest(*serverType, *version)
		if err != nil {
			fmt.Errorf("Error updating manifest: %s \n", err)
			return
		}
	}

	if *installServer {
		lib.InstallServer(*version, workingDir, *serverType)
	}

	if *defaultSettings {
		err := lib.ApplyDefaultSettings(workingDir)
		if err != nil {
			fmt.Printf("Error applying default settings: %s \n", err)
		}
	}
}
