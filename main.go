package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jgrtowy/yamst/lib"
)

func main() {
	updateManifest := flag.Bool("u", false, "Update the manifest")
	installServer := flag.Bool("i", false, "Install the server")
	version := flag.String("v", "", "Version to install")
	clearCache := flag.Bool("c", false, "Clear the cache")
	defaultSettings := flag.Bool("d", false, "Apply default settings")
	help := flag.Bool("h", false, "Show help")
	ngrok := flag.Bool("n", false, "Run ngrok tunnel")
	serverType := flag.String("t", "", "Server type")
	port := flag.Int("p", 25565, "Port to run ngrok tunnel on")
	flag.Parse()

	if *serverType == "" {
		*serverType = "vanilla"
	}
	if *ngrok {
		err := lib.RunTunnel(port)
		if err != nil {
			fmt.Printf("error running ngrok tunnel: %s", err)
		}
	}
	if *version == "" {
		var err error
		*version, err = lib.GetLatest()
		if err != nil {
			fmt.Printf("error getting latest version: %s", err)
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
			fmt.Printf("error clearing cache: %s", err)
			return
		}
		fmt.Println("Successfully cleared cache")
	}

	if *updateManifest {
		err := lib.UpdateManifest(*serverType, *version)
		if err != nil {
			fmt.Printf("error updating manifest: %s", err)
			return
		}
	}

	if *installServer {
		lib.InstallServer(*version, workingDir, *serverType)
	}

	if *defaultSettings {
		err := lib.ApplyDefaultSettings(workingDir)
		if err != nil {
			fmt.Printf("error applying default settings: %s", err)
			return
		}
	}
}
