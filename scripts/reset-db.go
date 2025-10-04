package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	// Delete the database file
	dbFile := "dev.db"
	
	if _, err := os.Stat(dbFile); err == nil {
		err := os.Remove(dbFile)
		if err != nil {
			log.Printf("Warning: Could not delete %s: %v", dbFile, err)
			log.Printf("Please manually delete the file or restart your system")
		} else {
			fmt.Printf("Successfully deleted %s\n", dbFile)
		}
	} else {
		fmt.Printf("Database file %s does not exist\n", dbFile)
	}
	
	fmt.Println("Database reset complete. You can now start the server.")
}