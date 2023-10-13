package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func ExistsCheck(serverType string, version string) (bool, error) {
	_, err := os.Stat(GetCacheDirectory())
	if err != nil {
		return false, err
	}
	switch serverType {
	case "paper":
		_, err := os.Stat(GetCacheDirectory() + fmt.Sprintf("/manifest_%s_%s.json", serverType, version))
		if err != nil {

			return false, err
		}
	case "vanilla":
		_, err := os.Stat(GetCacheDirectory() + fmt.Sprintf("/manifest_%s.json", serverType))
		if err != nil {
			return false, err
		}
	}
	return true, nil
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

func GetManifest(from string, serverType string, version string) ([]byte, error) {
	if from == "http" {
		var response *http.Response
		var err error

		switch serverType {
		case "vanilla":
			response, err = http.Get("https://launchermeta.mojang.com/mc/game/version_manifest.json")

		case "paper":
			if version == "" {
				version, err = GetLatest()
				if err != nil {
					return nil, err
				}
			}
			response, err = http.Get(fmt.Sprintf("https://papermc.io/api/v2/projects/paper/versions/%s", version))
			if err != nil {
				return nil, err
			}

		default:
			return nil, fmt.Errorf("invalid server type: %s", serverType)
		}
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
		exists, err := ExistsCheck(serverType, version)
		if err != nil {
			return nil, err
		}
		if !exists {
			err = UpdateManifest(serverType, version)
			if err != nil {
				return nil, err
			}
		}
		switch serverType {
		case "paper":
			body, err := readFileContents(GetCacheDirectory() + "/manifest_" + serverType + "_" + version + ".json")
			if err != nil {
				return nil, err
			}
			return body, nil
		case "vanilla":
			body, err := readFileContents(GetCacheDirectory() + "/manifest_" + serverType + ".json")
			if err != nil {
				return nil, err
			}
			return body, nil
		}
	}
	return nil, nil
}
func UpdateManifest(serverType string, version string) error {
	fmt.Println("Updating manifest...")
	if serverType == "" {
		serverType = "vanilla"
	}
	response, err := GetManifest("http", serverType, version)
	if err != nil {
		return fmt.Errorf("error getting manifest: %s", err)
	}
	var releaseInfo ReleaseInfo
	err = json.Unmarshal(response, &releaseInfo)
	if err != nil {
		return fmt.Errorf("error decoding JSON: %s", err)
	}
	latest, err := GetLatest()
	if err != nil {
		return fmt.Errorf("error getting latest version: %s", err)
	}
	switch serverType {
	case "paper":
		err = os.WriteFile(GetCacheDirectory()+fmt.Sprintf("/manifest_%s_%s.json", serverType, version), response, 0644)
		if err != nil {
			return fmt.Errorf("error writing to file: %s", err)
		}
		fmt.Println("Done! Using release:", version)
	case "vanilla":
		err = os.WriteFile(GetCacheDirectory()+"/manifest_vanilla.json", response, 0644)
		if err != nil {
			return fmt.Errorf("error writing to file: %s", err)
		}
		fmt.Println("Done! Latest release:", latest)
	default:
		return fmt.Errorf("invalid server type: %s", serverType)
	}

	return nil
}

func GetLatest() (string, error) {
	response, err := GetManifest("http", "vanilla", "")
	if err != nil {
		return "", err
	}
	var releaseInfo ReleaseInfo
	err = json.Unmarshal(response, &releaseInfo)
	if err != nil {
		return "", err
	}
	return releaseInfo.Latest.Release, nil
}

func GetLatestPaperBuild(version string) (int32, error) {
	exists, err := ExistsCheck("paper", version)
	if err != nil {
		return 0, err
	}
	if !exists {
		err = UpdateManifest("paper", version)
		if err != nil {
			return 0, err
		}
	}
	response, err := GetManifest("toolCache", "paper", version)
	if err != nil {
		return 0, err
	}
	var buildsList PaperBuilds
	err = json.Unmarshal(response, &buildsList)
	if err != nil {
		return 0, err
	}
	return buildsList.Builds[len(buildsList.Builds)-1], nil
}
