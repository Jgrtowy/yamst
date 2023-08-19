package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func ExistsCheck() bool {
	_, err := os.Stat(GetCacheDirectory())
	content, err := readFileContents(GetCacheDirectory() + "/manifest.json")
	if err == nil {
		var releaseInfo ReleaseInfo
		err = json.Unmarshal(content, &releaseInfo)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return false
		}
		fmt.Printf("Manifest exists. Use '-u' to update. Latest release: %s\n", releaseInfo.Latest.Release)
		return true
	} else {
		fmt.Println("Manifest does not exist. Creating...")
		UpdateManifest()
		return true
	}
}

func readFileContents(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	contentBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return contentBytes, nil
}

func GetManifest(from string) ([]byte, error) {
	if from == "http" {
		response, err := http.Get("https://launchermeta.mojang.com/mc/game/version_manifest.json")
		if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return nil, err
		}
		return body, nil
	} else if from == "toolCache" {
		body, err := readFileContents(GetCacheDirectory() + "/manifest.json")
		if err != nil {
			fmt.Println("Error reading manifest:", err)
			return nil, err
		}
		return body, nil
	}
	return nil, nil
}
func UpdateManifest() {
	fmt.Println("Updating manifest...")
	response, err := GetManifest("http")
	var releaseInfo ReleaseInfo
	err = json.Unmarshal(response, &releaseInfo)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	os.WriteFile(GetCacheDirectory()+"/manifest.json", response, 0644)
	fmt.Println("Done! Latest release:", releaseInfo.Latest.Release)
}
