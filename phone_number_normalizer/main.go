package main

import (
	"fmt"
	"log"
	"regexp"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PhoneNumber struct {
	gorm.Model
	Number string `gorm:"unique;not null"`
}

var initialNumbers = []string{
	"1234567890",
	"123 456 7891",
	"(123) 456 7892",
	"(123) 456-7893",
	"123-456-7894",
	"123-456-7890",
	"1234567892",
	"(123)456-7892",
}

func normalizePhoneNumber(number string) string {
	re := regexp.MustCompile(`\D`)
	return re.ReplaceAllString(number, "")
}

func main() {
	dbPath := "./numbers.db"

	if err := os.Remove(dbPath); err != nil && !os.IsNotExist(err) {
		log.Fatalf("Error removing existing DB file: %v", err)
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&PhoneNumber{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate database schema: %v", err)
	}
	fmt.Println("Database table created/migrated.")

	fmt.Println("\n--- Inserting Initial Numbers ---")
	for _, numStr := range initialNumbers {
		phoneNumber := PhoneNumber{Number: numStr}
		result := db.Create(&phoneNumber)
		if result.Error != nil {
			fmt.Printf("Warning: Failed to insert '%s' (might be a duplicate): %v\n", numStr, result.Error)
		} else {
			fmt.Printf("Inserted: %s\n", numStr)
		}
	}

	fmt.Println("\n--- Numbers Before Normalization ---")
	printNumbers(db)

	fmt.Println("\n--- Normalizing and Deduplicating Numbers ---")

	var numbersInDB []PhoneNumber
	if err := db.Find(&numbersInDB).Error; err != nil {
		log.Fatalf("Failed to retrieve numbers for normalization: %v", err)
	}

	uniqueNormalized := make(map[string]struct{})
	for _, p := range numbersInDB {
		uniqueNormalized[normalizePhoneNumber(p.Number)] = struct{}{}
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM phone_numbers").Error; err != nil {
			return fmt.Errorf("failed to clear table: %w", err)
		}
		fmt.Println("Existing numbers cleared from database.")

		for normalizedNum := range uniqueNormalized {
			newPhoneNumber := PhoneNumber{Number: normalizedNum}
			if err := tx.Create(&newPhoneNumber).Error; err != nil {
				return fmt.Errorf("failed to insert normalized number '%s': %w", normalizedNum, err)
			}
			fmt.Printf("Inserted normalized: %s\n", normalizedNum)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error during normalization and deduplication transaction: %v", err)
	}

	fmt.Println("\n--- Numbers After Normalization and Deduplication ---")
	printNumbers(db)

	fmt.Println("\nProgram finished successfully!")
}

func printNumbers(db *gorm.DB) {
	var numbers []PhoneNumber
	db.Order("number").Find(&numbers)

	if len(numbers) == 0 {
		fmt.Println("  (No numbers found)")
		return
	}

	for _, p := range numbers {
		fmt.Printf("  ID: %d, Number: %s\n", p.ID, p.Number)
	}
}