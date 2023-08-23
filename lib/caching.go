package lib

import (
	"os"
	"path/filepath"
	"runtime"
)

func GetCacheDirectory() string {
	var cacheDir string

	switch runtime.GOOS {
	case "windows":
		cacheDir = os.Getenv("LocalAppData")
	case "darwin":
		cacheDir = filepath.Join(os.Getenv("HOME"), "Library", "Caches")
	default:
		cacheDir = os.Getenv("XDG_CACHE_HOME")
		if cacheDir == "" {
			cacheDir = filepath.Join(os.Getenv("HOME"), ".cache")
		}
	}

	if _, err := os.Stat(filepath.Join(cacheDir, "YAMST")); os.IsNotExist(err) {
		err = os.Mkdir(filepath.Join(cacheDir, "YAMST"), 0755)
		if err != nil {
			panic(err)
		}
	}

	return filepath.Join(cacheDir, "YAMST")
}
