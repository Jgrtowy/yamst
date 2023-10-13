package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func download(version string, serverType string) error {
	fmt.Printf("Selected version: %s %s \n", version, serverType)
	manifest, err := GetManifest("toolCache", "vanilla", version)
	if err != nil {
		return fmt.Errorf("error getting manifest: %s", err)
	}

	var releaseInfo ReleaseInfo
	err = json.Unmarshal(manifest, &releaseInfo)
	if err != nil {
		return err
	}

	var versionInfo VersionInfo
	for _, v := range releaseInfo.Versions {
		if v.Id == version {
			versionInfo = v
		}
	}

	if versionInfo.Id == "" {
		return fmt.Errorf("version not found: %s", version)
	}

	fmt.Printf("Downloading server\n")
	var response *http.Response

	switch serverType {
	case "vanilla":
		var PackageInfo PackageInfo
		response, err = http.Get(versionInfo.Url)

		if err != nil {
			return fmt.Errorf("error downloading package info: %s", err)
		}

		err = json.NewDecoder(response.Body).Decode(&PackageInfo)
		if err != nil {
			return fmt.Errorf("error decoding package info: %s", err)
		}

		response, err = http.Get(PackageInfo.Downloads.Server.Url)
		defer response.Body.Close()
		if err != nil {
			return fmt.Errorf("error downloading server: %s", err)
		}
	case "paper":
		latestBuild, err := GetLatestPaperBuild(version)
		if err != nil {
			return fmt.Errorf("error getting latest build: %s", err)
		}
		url := fmt.Sprintf("https://papermc.io/api/v2/projects/paper/versions/%s/builds/%d", version, latestBuild)
		response, err = http.Get(url)
		if err != nil {
			return fmt.Errorf("error downloading package info: %s", err)
		}

		var paperBuildInfo PaperBuildInfo
		err = json.NewDecoder(response.Body).Decode(&paperBuildInfo)
		if err != nil {
			return fmt.Errorf("error decoding package info: %s", err)
		}
		url = url + "/downloads/" + paperBuildInfo.Downloads.Application.Name
		fmt.Println(url)
		response, err = http.Get(url)
		if err != nil {
			return fmt.Errorf("error downloading server: %s", err)
		}
		defer response.Body.Close()
	default:
		return fmt.Errorf("server type not found: %s", serverType)
	}

	abs, err := filepath.Abs(GetCacheDirectory() + "/" + version + "_" + serverType + ".jar")
	if err != nil {
		return fmt.Errorf("error getting absolute path: %s", err)
	}
	file, err := os.Create(abs)
	if err != nil {
		return fmt.Errorf("error creating server file: %s", err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return fmt.Errorf("error saving file: %s", err)
	}

	return nil
}

func JarExists(version string, serverType string) bool {
	abs, err := filepath.Abs(GetCacheDirectory() + "/" + version + "_" + serverType + ".jar")
	if err != nil {
		return false
	}
	_, err = os.Stat(abs)
	return err == nil
}

func InstallServer(version string, workingDir string, serverType string) error {
	if !JarExists(version, serverType) {
		fmt.Println("Server not found in cache. Downloading...")
		if version == "" {
			var err error
			version, err = GetLatest()
			if err != nil {
				return fmt.Errorf("error getting latest version: %s", err)
			}
		}

		err := download(version, serverType)
		if err != nil {
			return fmt.Errorf("error downloading server: %s", err)
		}
	} else {
		fmt.Println("Server found in cache. Copying...")
	}

	abs, err := filepath.Abs(GetCacheDirectory() + "/" + version + "_" + serverType + ".jar")
	if err != nil {
		return fmt.Errorf("error getting absolute path: %s", err)
	}

	err = copyFile(abs, workingDir+"/server.jar")
	if err != nil {
		return fmt.Errorf("error copying server: %s", err)
	}

	fmt.Printf("Successfully copied server to this directory on version: %s \n", version)
	return nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}
