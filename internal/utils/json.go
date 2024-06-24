package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
)

// just clean the json
func FormatJSON(jsonString string) (string, error) {
	var data interface{}
	err := json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	formattedJSON, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", fmt.Errorf("failed to format JSON: %v", err)
	}

	return string(formattedJSON), nil
}

func SaveStringToJSON(input string, folderPath string, filename string) error {

	fileName := fmt.Sprintf("Scan_report_%s.json", filename)
	filePath := filepath.Join(folderPath, fileName)
	err := ioutil.WriteFile(filePath, []byte(input), 0644)
	if err != nil {
		return fmt.Errorf("failed to write JSON to file: %v", err)
	}

	fmt.Printf("[-\x1b[48;5;75mJSON\x1b[1;0m--] - SAVED : %s\n", filePath)
	return nil
}

func CheckAPIResponse(response string) {
	var apiError map[string]string
	if jsonErr := json.Unmarshal([]byte(response), &apiError); jsonErr == nil {
		// Check specific error conditions
		if apiError["report"] == "Report not Found" {
			log.Println("[-\x1b[48;5;196mERROR\x1b[1;0m-] - Report not Found")
		} else if apiError["error"] == "Invalid Hash" {
			log.Println("[-\x1b[48;5;196mERROR\x1b[1;0m-] - Invalid Hash")
		}
	}
}
