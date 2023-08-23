package lib

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

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
