package test_offloading

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// Simulate request to cloud node
func request(cloudURL string) time.Duration {
	start := time.Now()
	resp, err := http.Get(cloudURL)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	return time.Since(start)
}

func runBenchmarkEdgeOnly(n int, edgeURL string) time.Duration {
	var wg sync.WaitGroup
	totalDuration := time.Duration(0)
	var mu sync.Mutex

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			duration := request(edgeURL)
			mu.Lock()
			totalDuration += duration
			mu.Unlock()
		}()
	}

	wg.Wait()
	return totalDuration / time.Duration(n)
}

func runBenchmarkMixed(n int, edgeURL, cloudURL string) time.Duration {
	var wg sync.WaitGroup
	totalDuration := time.Duration(0)
	var mu sync.Mutex

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var duration time.Duration
			if i < 42 {
				duration = request(edgeURL)
			} else {
				duration = request(cloudURL)
			}
			mu.Lock()
			totalDuration += duration
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	return totalDuration / time.Duration(n)
}

func writeCSV(filename string, data [][]string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range data {
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	edgeURL := "http://<ip-address-edge>/api/empty-function"
	cloudURL := "http://<ip-address-cloud>/api/empty-function"

	var edgeData [][]string
	edgeData = append(edgeData, []string{"Requests", "AvgResponseTime"})

	var mixedData [][]string
	mixedData = append(mixedData, []string{"Requests", "AvgResponseTime"})

	for n := 1; n <= 100; n++ {
		avgTimeEdge := runBenchmarkEdgeOnly(n, edgeURL)
		avgTimeMixed := runBenchmarkMixed(n, edgeURL, cloudURL)

		edgeData = append(edgeData, []string{strconv.Itoa(n), fmt.Sprintf("%d", avgTimeEdge.Milliseconds())})
		mixedData = append(mixedData, []string{strconv.Itoa(n), fmt.Sprintf("%d", avgTimeMixed.Milliseconds())})
	}

	if err := writeCSV("./edge_only.csv", edgeData); err != nil {
		fmt.Println("Error writing edge_only.csv:", err)
	}

	if err := writeCSV("./mixed.csv", mixedData); err != nil {
		fmt.Println("Error writing mixed.csv:", err)
	}
}
