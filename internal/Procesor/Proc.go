package Procesor

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
	mobapi "main.go/internal/Mob_Api"
	"main.go/internal/config"
)

var (
	processingChan   = make(chan struct{}, 20)
	processedFiles   = make(map[string]bool)
	mutex            sync.Mutex
	wg               sync.WaitGroup
	maxConcurrency   = semaphore.NewWeighted(11)
	processedCounter int
	counterMutex     sync.Mutex
	startTime        time.Time
)

func ScanFolder() {
	log.SetFlags(0)
	cfg, err := config.LoadConfig("\\internal\\config\\config.json")
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
		return
	}

	fmt.Printf("--------------RUNNING V°%s clean----------------\n", cfg.Version)
	startTime = time.Now()

	go printProcessedCount()

	for {
		scanAndProcessFiles(cfg)
		time.Sleep(2 * time.Second)
	}
}

func processFile(ctx context.Context, cfg *config.Config, filePath, fileName string) {
	defer wg.Done()

	processingChan <- struct{}{}
	defer func() { <-processingChan }()

	defer func() {
		if r := recover(); r != nil {
			log.Printf("[-\x1b[48;5;196mERROR\x1b[1;0m-] Panic while processing file %s: %v", fileName, r)
		}
	}()

	if err := maxConcurrency.Acquire(ctx, 1); err != nil {
		log.Printf("[-\x1b[48;5;196mERROR\x1b[1;0m-] Acquiring semaphore for file %s: %v", fileName, err)
		return
	}
	defer maxConcurrency.Release(1)

	fmt.Printf("[--\x1b[48;5;40mNEW\x1b[1;0m--] - %s\n", fileName)
	if err := mobapi.UploadScanAndReport(filePath, cfg.Domain, cfg.ApiKey, cfg.Dump, fileName); err != nil {
		log.Printf("[-\x1b[48;5;196mERROR\x1b[1;0m-] - UploadScanAndReport - %s", err)
	} else {
		fmt.Printf("[--\x1b[48;5;226mRMV\x1b[1;0m--] - %s\n", fileName)
		if err := os.Remove(filePath); err != nil {
			log.Printf("[\x1b[48;5;196mERROR\x1b[1;0m] Removing %s: %v", filePath, err)
		}

		mutex.Lock()
		processedFiles[fileName] = true
		mutex.Unlock()

		counterMutex.Lock()
		processedCounter++
		counterMutex.Unlock()
	}
}

func scanAndProcessFiles(cfg *config.Config) {
	files, err := os.ReadDir(cfg.InFolder)
	if err != nil {
		log.Fatalf("[\x1b[48;5;196mERROR\x1b[1;0m] - read dir: %v", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()
			filePath := filepath.Join(cfg.InFolder, fileName)

			mutex.Lock()
			processed := processedFiles[fileName]
			mutex.Unlock()

			if !processed && hasValidExtension(fileName, cfg.Etx_list) {
				wg.Add(1)
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
				defer cancel()

				go processFile(ctx, cfg, filePath, fileName)
				mutex.Lock()
				processedFiles[fileName] = true
				mutex.Unlock()
			} else if !processed {
				fmt.Printf("[\x1b[48;5;215mINVALID\x1b[1;0m] - %s: file ext not right\n", fileName)
				fmt.Printf("[--\x1b[48;5;226mRMV\x1b[1;0m--] - %s\n", fileName)
				if err := os.Remove(filePath); err != nil {
					log.Printf("[-\x1b[48;5;196mERROR\x1b[1;0m-] Removing %s: %v", filePath, err)
				}
			}
		}
	}

	wg.Wait()
}

func printProcessedCount() {
	for {
		time.Sleep(30 * time.Second)
		counterMutex.Lock()
		elapsed := time.Since(startTime)
		elapsedMinutes := int(elapsed.Minutes())
		elapsedSeconds := int(elapsed.Seconds()) % 60
		fmt.Printf("[\x1b[48;5;83mPROCESS\x1b[1;0m] - Processed_file: %d | Time taken %dm %ds \n", processedCounter, elapsedMinutes, elapsedSeconds)
		counterMutex.Unlock()
	}
}

func hasValidExtension(fileName string, allowedExts []string) bool {
	for _, ext := range allowedExts {
		if strings.HasSuffix(fileName, ext) {
			return true
		}
	}
	return false
}
