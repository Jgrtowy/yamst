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
		return fmt.Errorf("Error getting manifest: %s \n", err)
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
		return fmt.Errorf("Version not found \n")
	}
	fmt.Printf("Downloading server\n")
	var response *http.Response
	switch serverType {
	case "vanilla":
		var PackageInfo PackageInfo
		response, err := http.Get(versionInfo.Url)

		if err != nil {
			return fmt.Errorf("Error downloading package info: %s \n", err)
		}

		err = json.NewDecoder(response.Body).Decode(&PackageInfo)
		if err != nil {
			return fmt.Errorf("Error decoding package info: %s \n", err)
		}

		response, err = http.Get(PackageInfo.Downloads.Server.Url)
		defer response.Body.Close()
		if err != nil {
			return fmt.Errorf("Error downloading server: %s \n", err)
		}
		break
	case "paper":
		latestBuild, err := GetLatestPaperBuild(version)
		if err != nil {
			return fmt.Errorf("Error getting latest build: %s \n", err)
		}
		url := fmt.Sprintf("https://papermc.io/api/v2/projects/paper/versions/%s/builds/%d", version, latestBuild)
		response, err = http.Get(url)
		var paperBuildInfo PaperBuildInfo
		err = json.NewDecoder(response.Body).Decode(&paperBuildInfo)
		if err != nil {
			return fmt.Errorf("Error decoding package info: %s \n", err)
		}
		url = url + "/downloads/" + paperBuildInfo.Downloads.Application.Name
		fmt.Println(url)
		response, err = http.Get(url)
		if err != nil {
			return fmt.Errorf("Error downloading server: %s \n", err)
		}
		defer response.Body.Close()

		break
	}

	abs, err := filepath.Abs(GetCacheDirectory() + "/" + version + "_" + serverType + ".jar")
	file, err := os.Create(abs)
	if err != nil {
		return fmt.Errorf("Error creating server file: %s \n", err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return fmt.Errorf("Error saving file: %s", err)
	}

	return nil
}

func JarExists(version string, serverType string) bool {
	abs, err := filepath.Abs(GetCacheDirectory() + "/" + version + "_" + serverType + ".jar")
	if err != nil {
		return false
	}
	_, err = os.Stat(abs)
	if err != nil {
		return false
	}
	return true
}

func InstallServer(version string, workingDir string, serverType string) {
	if !JarExists(version, serverType) {
		fmt.Println("Server not found in cache. Downloading...")
		if version == "" {
			var err error
			version, err = GetLatest()
			if err != nil {
				fmt.Printf("Error getting latest version: %s \n", err)
				return
			}
		}
		err := download(version, serverType)
		if err != nil {
			fmt.Printf("Error downloading server: %s \n", err)
			return
		}
	} else {
		fmt.Println("Server found in cache. Copying...")
	}
	abs, err := filepath.Abs(GetCacheDirectory() + "/" + version + "_" + serverType + ".jar")
	err = copyFile(fmt.Sprintf(abs), fmt.Sprintf("%s/server.jar", workingDir))

	if err != nil {
		fmt.Printf("Error copying server: %s \n", err)
		return
	}

	fmt.Printf("Successfully copied server to this directory on version: %s \n", version)
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
