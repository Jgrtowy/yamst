package main

import (
	"flag"
	"github.com/Jgrtowy/ServerSetup/lib"
)

func main() {
	updateManifest := flag.Bool("u", false, "Update the manifest")
	flag.Parse()
	if !lib.ExistsCheck() || *updateManifest {
		lib.UpdateManifest()
	}
}
