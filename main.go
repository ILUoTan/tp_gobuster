package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}
	return lines, scanner.Err()
}

func sendRequest(target string, path string) (int, error) {
	url := fmt.Sprintf("http://%s/%s", target, path)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

func scan(target string, paths []string, workers int, quiet bool) {
	var wg sync.WaitGroup
	jobs := make(chan string)
	results := make(chan string)

	start := time.Now()

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range jobs {
				status, err := sendRequest(target, path)
				if err != nil {
					continue
				}
				if status == 200 || !quiet {
					results <- fmt.Sprintf("/%s\t%d", path, status)
				}
			}
		}()
	}

	go func() {
		for _, path := range paths {
			jobs <- path
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	fmt.Println("Starting scan...")
	for result := range results {
		fmt.Println(result)
	}

	elapsed := time.Since(start)
	fmt.Printf("Scan done in %.6fs\n", elapsed.Seconds())
}

func main() {
	dictionary := flag.String("d", "", "Path to dictionary file")
	quiet := flag.Bool("q", false, "Quiet mode, only show HTTP 200 results")
	target := flag.String("t", "", "Target to enumerate (including port)")
	workers := flag.Int("w", 1, "Number of workers to run")
	flag.Parse()

	if *dictionary == "" || *target == "" {
		fmt.Println("Usage of mygb:")
		flag.PrintDefaults()
		return
	}

	paths, err := readLines(*dictionary)
	if err != nil {
		fmt.Printf("Failed to read dictionary: %v\n", err)
		return
	}

	fmt.Println("Starting MyGB")
	fmt.Println("--")
	fmt.Printf("Target: http://%s\nList: %s\nWorkers: %d\n", *target, *dictionary, *workers)
	fmt.Println("--")

	scan(*target, paths, *workers, *quiet)
}
