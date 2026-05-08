package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

func main() {

	url := flag.String("url", "http://localhost:8080/orders", "Ziel-URL (Monolith: 8080, Microservice: 8082)")
	concurrency := flag.Int("users", 100, "Anzahl gleichzeitiger Nutzer")
	requestsPerUser := flag.Int("reqs", 10, "Bestellungen pro Nutzer")
	flag.Parse()

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 1000
	t.MaxConnsPerHost = 1000
	t.MaxIdleConnsPerHost = 1000

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}

	var (
		successCount uint64
		failCount    uint64
		totalLatency int64
	)

	payload := []byte(`{"productId": 1, "amount": 1, "userId": 1}`)

	fmt.Printf("--- START LASTTEST ---\n")
	fmt.Printf("Ziel: %s\n", *url)
	fmt.Printf("User: %d | Reqs/User: %d | Total Reqs: %d\n", *concurrency, *requestsPerUser, *concurrency**requestsPerUser)

	start := time.Now()

	var wg sync.WaitGroup

	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func(userId int) {
			defer wg.Done()

			for j := 0; j < *requestsPerUser; j++ {
				reqStart := time.Now()

				resp, err := client.Post(*url, "application/json", bytes.NewReader(payload))

				latency := time.Since(reqStart).Milliseconds()
				atomic.AddInt64(&totalLatency, latency)

				if err != nil {
					atomic.AddUint64(&failCount, 1)
					continue
				}

				io.Copy(io.Discard, resp.Body)
				resp.Body.Close()

				if resp.StatusCode >= 200 && resp.StatusCode < 300 {
					atomic.AddUint64(&successCount, 1)
				} else {
					atomic.AddUint64(&failCount, 1)
				}
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)

	totalReqs := successCount + failCount
	avgLatency := float64(totalLatency) / float64(totalReqs)
	throughput := float64(totalReqs) / duration.Seconds()

	fmt.Printf("\n--- ERGEBNISSE ---\n")
	fmt.Printf("Zeit gesamt:    %.2fs\n", duration.Seconds())
	fmt.Printf("Erfolgreich:    %d\n", successCount)
	fmt.Printf("Fehlgeschlagen: %d\n", failCount)
	fmt.Printf("Throughput:     %.2f req/s\n", throughput)
	fmt.Printf("Avg Latency:    %.2f ms\n", avgLatency)

	if failCount > 0 {
		fmt.Println("WARNUNG: Fehler aufgetreten! System überlastet oder leer?")
	}

}
