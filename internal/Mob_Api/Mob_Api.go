package mobapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"main.go/internal/utils"
)

func UploadFile(filePath, domain, apiKey string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	formFile, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %v", err)
	}

	_, err = io.Copy(formFile, file)
	if err != nil {
		return "", fmt.Errorf("failed to copy file: %v", err)
	}

	writer.Close()

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/upload", domain), &requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	request.Header.Set("Authorization", apiKey)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response code: %d", response.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response body: %v", err)
	}

	hash, ok := result["hash"].(string)
	if !ok {
		return "", fmt.Errorf("hash not found in response")
	}

	return hash, nil
}

func UploadFilesInFolder(folderPath, domain, apiKey string) ([]string, error) {
	fileInfos, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read folder: %v", err)
	}

	var filePaths []string
	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			filePath := filepath.Join(folderPath, fileInfo.Name())
			filePaths = append(filePaths, filePath)
		}
	}

	return UploadFiles(filePaths, domain, apiKey)
}

func UploadFiles(filePaths []string, domain, apiKey string) ([]string, error) {
	var hashes []string

	for _, filePath := range filePaths {
		hash, err := uploadSingleFile(filePath, domain, apiKey)
		if err != nil {
			return nil, err
		}
		hashes = append(hashes, hash)
	}

	return hashes, nil
}

func uploadSingleFile(filePath, domain, apiKey string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	formFile, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %v", err)
	}

	_, err = io.Copy(formFile, file)
	if err != nil {
		return "", fmt.Errorf("failed to copy file: %v", err)
	}

	writer.Close()

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/upload", domain), &requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	request.Header.Set("Authorization", apiKey)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response code: %d", response.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response body: %v", err)
	}

	hash, ok := result["hash"].(string)
	if !ok {
		return "", fmt.Errorf("hash not found in response")
	}

	return hash, nil
}

func Scanfile(Hash, domain, apiKey string) (string, error) {
	payload := []byte("hash=" + Hash)
	url := (domain + "/api/v1/scan")

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Erreur creating request :", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erreur creating request :", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erreur reading respond body :", err)
		return "", err
	}

	return string(body), nil
}

func GetRepJson(Hash, domain, apiKey string) (string, error) {
	payload := []byte("hash=" + Hash)
	url := (domain + "/api/v1/report_json")

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Erreur creating request :", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erreur sending request :", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erreur reading respond body :", err)
		return "", err
	}

	return string(body), nil
}

func ScordREP(Hash, domain, apiKey string) (string, error) {
	payload := []byte("hash=" + Hash)
	url := (domain + "/api/v1/scorecard")

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Erreur creating request :", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erreur sending request :", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erreur reading respond body :", err)
		return "", err
	}

	return string(body), nil
}

func UploadScanAndReport(filePath, domain, apiKey, folderPath string, filename string) error {
	log.SetFlags(0)

	hash, err := UploadFile(filePath, domain, apiKey)
	if err != nil {
		log.Printf("[-\x1b[48;5;196mERROR\x1b[1;0m-] - upload file: %v", err)
		return err
	}

	fmt.Println("[\x1b[48;5;141mRUNNING\x1b[1;0m] - upload")
	fmt.Println("[\x1b[48;5;141mRUNNING\x1b[1;0m] - hash - ", hash)

	_, err = Scanfile(hash, domain, apiKey)
	if err != nil {
		log.Printf("[\x1b[48;5;196mERROR\x1b[1;0m] - scan file: %v", err)
		return err
	}

	fmt.Println("[\x1b[48;5;141mRUNNING\x1b[1;0m] - Scan file")

	report, err := GetRepJson(hash, domain, apiKey)
	if err != nil {
		log.Printf("[\x1b[48;5;196mERROR\x1b[1;0m] - get report: %v", err)
	}

	utils.CheckAPIResponse(report)
	fmt.Println("[\x1b[48;5;141mRUNNING\x1b[1;0m] - get report")

	cleanedJSON, err := utils.FormatJSON(report)
	if err != nil {
		log.Printf("[\x1b[48;5;196mERROR\x1b[1;0m] - format JSON: %v", err)
		return err
	}

	fmt.Println("[\x1b[48;5;141mRUNNING\x1b[1;0m] - format json")

	err = utils.SaveStringToJSON(cleanedJSON, folderPath, filename)
	if err != nil {
		log.Printf("[\x1b[48;5;196mERROR\x1b[1;0m] - save JSON to file: %v", err)
		return err
	}

	fmt.Println("[\x1b[48;5;141mRUNNING\x1b[1;0m] - save json")
	return nil
}
