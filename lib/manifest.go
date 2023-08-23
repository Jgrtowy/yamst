package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func ExistsCheck(serverType string, version string) error {
	_, err := os.Stat(GetCacheDirectory())
	if err != nil {
		return err
	}
	switch serverType {
	case "paper":
		_, err := os.Stat(GetCacheDirectory() + fmt.Sprintf("/manifest_%s_%s.json", serverType, version))
		if err != nil {
			return err
		}
		break
	case "vanilla":
		_, err := os.Stat(GetCacheDirectory() + fmt.Sprintf("/manifest_%s.json", serverType))
		if err != nil {
			return err
		}
		break
	}
	return nil
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

		if err != nil {
			return nil, err
		}
		switch serverType {
		case "vanilla":
			response, err = http.Get("https://launchermeta.mojang.com/mc/game/version_manifest.json")
			break
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
			break
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
		if err := ExistsCheck(serverType, version); err != nil {
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
	var releaseInfo ReleaseInfo
	err = json.Unmarshal(response, &releaseInfo)
	if err != nil {
		return fmt.Errorf("Error decoding JSON: %s \n", err)
	}
	latest, err := GetLatest()
	if err != nil {
		return fmt.Errorf("Error getting latest version: %s \n", err)
	}
	switch serverType {
	case "paper":
		err = os.WriteFile(GetCacheDirectory()+fmt.Sprintf("/manifest_%s_%s.json", serverType, version), response, 0644)
		if err != nil {
			return fmt.Errorf("Error writing to file: %s \n", err)
		}
		fmt.Println("Done! Using release:", version)
		break
	case "vanilla":
		err = os.WriteFile(GetCacheDirectory()+"/manifest_vanilla.json", response, 0644)
		if err != nil {
			return fmt.Errorf("Error writing to file: %s \n", err)
		}
		fmt.Println("Done! Latest release:", latest)
		break
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
	if err := ExistsCheck("paper", version); err != nil {
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
