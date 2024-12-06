package migration_test

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Configuration
var (
	migrations     = 200
	sessionSizes   = []int{1, 256, 512, 1024, 2048, 3072, 4096, 5120} // Different session sizes in MB
	ravennaNodeURL = "http://localhost:8080/migrate?size="
	outputFile     = "migration_results.csv"
)

// MigrationResult holds the results of a single migration
type MigrationResult struct {
	Size     int
	Duration time.Duration
}

func performMigration(sessionSize int) (time.Duration, error) {
	// Call the migration function
	resp, err := http.Post(ravennaNodeURL+strconv.Itoa(sessionSize), "application/json", nil)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	duration, err := time.ParseDuration(string(body))
	if err != nil {
		return 0, err
	}

	return duration, nil
}

func main() {
	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatalln("Failed to create file:", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Session Size (KB)", "Total Migration Time (ms)"})

	for _, size := range sessionSizes {
		for i := 0; i < migrations; i++ {
			duration, err := performMigration(size)
			if err != nil {
				log.Fatal("Migration failed:", err)
			}
			writer.Write([]string{strconv.Itoa(size), fmt.Sprintf("%v", duration.Milliseconds())})
		}
	}

	fmt.Println("Migration tests completed and results saved to", outputFile)
}
