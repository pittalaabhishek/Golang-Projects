package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	rootDir := "sample"
	fmt.Printf("Renaming files in: %s\n", rootDir)

	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		fmt.Printf("Directory '%s' does not exist. Creating it...\n", rootDir)
		err := os.MkdirAll(rootDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create directory: %v", err)
		}
		
		createSampleFiles(rootDir)
		fmt.Println("Created sample files for testing.")
	}

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", path, err)
			return err
		}
		
		if info.IsDir() {
			return nil
		}

		fileName := info.Name()
		dir := filepath.Dir(path)
		var newName string

		// birthday_001.txt -> birthday_renamed_001.txt
		if strings.HasPrefix(fileName, "birthday_") && strings.HasSuffix(fileName, ".txt") {
			newName = strings.Replace(fileName, "birthday_", "birthday_renamed_", 1)
		}

		// christmas files -> christmas_renamed files
		if strings.Contains(fileName, "christmas") && strings.HasSuffix(fileName, ".txt") {
			parts := strings.Split(fileName, "(")
			if len(parts) > 1 {
				newName = "christmas_renamed_" + parts[1]
				newName = strings.Replace(newName, " of ", "_of_", -1)
				newName = strings.Replace(newName, ")", "", -1)
			} else {
				// Handle christmas files without parentheses
				newName = strings.Replace(fileName, "christmas", "christmas_renamed", 1)
			}
		}

		// n_xxx.txt -> nested_xxx.txt
		if strings.HasPrefix(fileName, "n_") && strings.Contains(path, "nested") {
			newName = strings.Replace(fileName, "n_", "nested_", 1)
		}

		if newName != "" {
			newPath := filepath.Join(dir, newName)
			fmt.Printf("Renaming: %s -> %s\n", fileName, newName)
			err := os.Rename(path, newPath)
			if err != nil {
				fmt.Printf("Error renaming file: %v\n", err)
			}
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done!")
}

func createSampleFiles(rootDir string) {
	testFiles := []string{
		"birthday_001.txt",
		"birthday_002.txt",
		"christmas_special.txt",
		"christmas (day of joy).txt",
		"regular_file.txt",
	}

	for _, file := range testFiles {
		filePath := filepath.Join(rootDir, file)
		f, err := os.Create(filePath)
		if err != nil {
			fmt.Printf("Warning: Could not create test file %s: %v\n", file, err)
			continue
		}
		f.Close()
	}

	nestedDir := filepath.Join(rootDir, "nested")
	os.MkdirAll(nestedDir, 0755)
	
	nestedFiles := []string{
		"n_test.txt",
		"testing.txt",
	}

	for _, file := range nestedFiles {
		filePath := filepath.Join(nestedDir, file)
		f, err := os.Create(filePath)
		if err != nil {
			fmt.Printf("Warning: Could not create nested test file %s: %v\n", file, err)
			continue
		}
		f.Close()
	}
}