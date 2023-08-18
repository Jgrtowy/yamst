package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func ExistsCheck() bool {
	_, err := os.Stat("./cache/manifest.json")
	content, err := readFileContents("./cache/manifest.json")
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

	contentBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return contentBytes, nil
}

func getManifest() ([]byte, error) {
	response, err := http.Get("https://launchermeta.mojang.com/mc/game/version_manifest.json")
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer response.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, err
	}

	return body, nil
}
func UpdateManifest() {
	fmt.Println("Updating manifest...")
	response, err := getManifest()
	var releaseInfo ReleaseInfo
	err = json.Unmarshal(response, &releaseInfo)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	os.Mkdir("cache", 0755)
	os.WriteFile("./cache/manifest.json", response, 0644)
	// create directory
	fmt.Println("Done! Latest release:", releaseInfo.Latest.Release)
}
