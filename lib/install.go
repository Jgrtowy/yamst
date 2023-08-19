package lib

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func download(version string) error {
	fmt.Printf("Selected version: %s \n", version)
	manifest, err := GetManifest("toolCache")
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

	fmt.Printf("Downloading server on version: %s \n", versionInfo.Id)
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

	abs, err := filepath.Abs(GetCacheDirectory() + "/" + version + ".jar")
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

func JarExists(version string) bool {
	abs, err := filepath.Abs(GetCacheDirectory() + "/" + version + ".jar")
	if err != nil {
		return false
	}
	_, err = os.Stat(abs)
	if err != nil {
		return false
	}
	return true
}

func InstallServer(version string, workingDir string) {
	// check if the version is downloaded already in cache
	if !JarExists(version) {
		fmt.Println("Server not found in cache. Downloading...")
		if version == "" {
			version = "latest"
		}
		err := download(version)
		if err != nil {
			fmt.Printf("Error downloading server: %s \n", err)
			return
		}
	} else {
		fmt.Println("Server found in cache. Copying...")
	}
	abs, err := filepath.Abs(GetCacheDirectory() + "/" + version + ".jar")
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

func ApplyDefaultSettings(workingDir string) error {
	fmt.Println("By using -d option you agree to the Minecraft EULA: https://account.mojang.com/documents/minecraft_eula")
	// wait for y or n
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you agree? (y/n): ")
	char, _, err := reader.ReadRune()
	if err != nil {
		return fmt.Errorf("Error reading input: %s \n", err)
	}
	switch char {
	case 'y':
		fmt.Println("Applying default settings...")
		abs, err := filepath.Abs(fmt.Sprintf("%s/eula.txt", workingDir))
		if err != nil {
			return fmt.Errorf("Error getting absolute path: %s \n", err)
		}
		_, err = os.Create(abs)
		if err != nil {
			return fmt.Errorf("Error creating eula.txt: %s \n", err)
		}
		file, err := os.OpenFile(abs, os.O_RDWR, 0644)
		if err != nil {
			return fmt.Errorf("Error opening eula.txt: %s \n", err)
		}
		_, err = file.WriteString("eula=true")
		if err != nil {
			return fmt.Errorf("Error writing to eula.txt: %s \n", err)
		}
		defer file.Close()

		_, err = os.Stat(fmt.Sprintf("%s/start.sh", workingDir))
		_, err = os.Stat(fmt.Sprintf("%s/start.bat", workingDir))
		if err != nil {
			_, err = os.Create(fmt.Sprintf("%s/start.sh", workingDir))
			_, err = os.Create(fmt.Sprintf("%s/start.bat", workingDir))

			if err != nil {
				return fmt.Errorf("Error creating start scripts: %s \n", err)
			}
		}

		err = os.WriteFile(fmt.Sprintf("%s/start.sh", workingDir), []byte("java -Xms1024M -Xmx1024M -jar server.jar nogui"), 0644)
		err = os.WriteFile(fmt.Sprintf("%s/start.bat", workingDir), []byte("java -Xms1024M -Xmx1024M -jar server.jar nogui"), 0644)

		if err != nil {

		}
		fmt.Println("Successfully applied default settings")

		return nil

	case 'n':
		fmt.Println("You must agree to the EULA to use the -d option")
		return nil

	default:
		fmt.Println("Invalid input")
		return nil
	}
}
