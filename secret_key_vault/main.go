package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"secret_key_vault/secret"
)

func getSecretsFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Cannot find user home directory")
	}
	return filepath.Join(home, ".secrets.enc")
}

func main() {
	var encodingKey string
	flag.StringVar(&encodingKey, "k", "", "Encryption key")
	flag.Parse()

	args := flag.Args() // non-flag arguments
	
	if len(args) < 2 {
		fmt.Println("Usage: secret [get|set] key [value] -k encryption-key")
		return
	}

	operation := args[0]
	key := args[1]
	value := ""
	
	if len(args) > 2 && operation == "set" {
		value = args[2]
	}

	if encodingKey == "" {
		log.Fatal("Encryption key required. Use -k flag.")
	}

	vault := secret.NewFileVault(encodingKey, getSecretsFilePath())

	switch operation {
	case "set":
		err := vault.Set(key, value)
		if err != nil {
			log.Fatalf("Failed to set key: %v", err)
		}
		fmt.Println("Value set!")
	case "get":
		val, err := vault.Get(key)
		if err != nil {
			log.Fatalf("Failed to get key: %v", err)
		}
		fmt.Printf("%s\n", val)
	default:
		fmt.Println("Unknown command. Use get or set.")
	}
}