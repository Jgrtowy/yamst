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
	flag.Parse()

	workingDir, _ := os.Getwd()
	if *help {
		fmt.Println("Usage: ServerSetup [options]")
		flag.PrintDefaults()
		return
	}

	if *updateManifest && *installServer {
		panic("Cannot update manifest and install server at the same time")
	}

	if *clearCache {
		err := os.RemoveAll(lib.GetCacheDirectory())
		if err != nil {
			panic(err)
		}
	}

	if *updateManifest {
		lib.UpdateManifest()
	}

	if *installServer {
		lib.InstallServer(*version, workingDir)
	}

	if *defaultSettings {
		lib.ApplyDefaultSettings(workingDir)
	}
}
