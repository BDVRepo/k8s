package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"
)

type YourJSONStruct struct {
	ID string `json:"id"`
}

func main() {
	path := flag.String("p", "./", "path to json files")
	flag.Parse()
	files, err := os.ReadDir(*path)
	if err != nil {
		log.Fatalf("Error reading directory: %s", err)
	}

	for {
		for _, file := range files {
			if filepath.Ext(file.Name()) == ".json" {
				filePath := filepath.Join(*path, file.Name())
				jsonFile, err := os.Open(filePath)
				if err != nil {
					log.Printf("Error opening file %s: %s", filePath, err)
					continue
				}
				defer jsonFile.Close()

				var data YourJSONStruct
				decoder := json.NewDecoder(jsonFile)
				if err := decoder.Decode(&data); err != nil {
					log.Printf("Error decoding JSON from file %s: %s", filePath, err)
					continue
				}

				if data.ID != "" {
					log.Printf("ID found in %s: %s", filePath, data.ID)
				} else {
					log.Printf("No ID found in %s", filePath)
				}
			}
		}
		time.Sleep(60 * time.Minute)
	}
}
