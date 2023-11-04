package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

// Define ANSI color codes
const (
	RED   = "\033[0;31m"
	CYAN  = "\033[0;36m"
	GREEN = "\033[0;32m"
	NC    = "\033[0m" // No Color
)

var (
	verbose = false
)

func main() {
	// Command-line flags
	urlsFile := flag.String("u", "urls.txt", "File containing target URLs")
	payloadsFile := flag.String("p", "payloads.txt", "File containing payloads")
	outputFile := flag.String("o", "vulnerable_urls.txt", "Output file for vulnerable URLs")
	verbosity := flag.Bool("v", false, "Enable verbosity for all requests")
	flag.Parse()

	// Set the verbosity level based on the flag
	verbose = *verbosity

	// Generate the banner
	banner := `
___  ___                    _____    ___          _     ______  _____ 
|  \/  |                   / __  \  /   |        | |    |  ___||_   _|
| .  . | _ __   ___   ___  `' / /' / /| | ______ | |    | |_     | |  
| |\/| || '__| / __| / _ \   / /  / /_| ||______|| |    |  _|    | |  
| |  | || |   | (__ | (_) |./ /___\___  |        | |____| |     _| |_ 
\_|  |_/|_|    \___| \___/ \_____/    |_/        \_____/\_|     \___/                                                                
                                                
`
	fmt.Print(CYAN, banner, NC)

	// Check if the URLs and payloads files exist
	_, urlsErr := os.Stat(*urlsFile)
	_, payloadsErr := os.Stat(*payloadsFile)

	if urlsErr != nil || payloadsErr != nil {
		fmt.Println("URLs file or payloads file not found. Make sure both files exist.")
		return
	}

	// Read the list of target URLs from the URLs file
	urlsData, readURLsErr := ioutil.ReadFile(*urlsFile)
	if readURLsErr != nil {
		fmt.Println("Error reading URLs file:", readURLsErr)
		return
	}
	urls := strings.Split(string(urlsData), "\n")

	// Read the list of payloads from the payloads file
	payloadsData, readPayloadsErr := ioutil.ReadFile(*payloadsFile)
	if readPayloadsErr != nil {
		fmt.Println("Error reading payloads file:", readPayloadsErr)
		return
	}
	payloads := strings.Split(string(payloadsData), "\n")

	// Initialize the output file
	output, createOutputErr := os.Create(*outputFile)
	if createOutputErr != nil {
		fmt.Println("Error creating output file:", createOutputErr)
		return
	}
	defer output.Close()

	// Create a wait group for synchronization
	var wg sync.WaitGroup

	// Iterate through URLs and payloads concurrently
	for _, url := range urls {
		for _, payload := range payloads {
			wg.Add(1)
			go func(url, payload string) {
				defer wg.Done()
				fullURL := url + payload

				// Check if the URL is valid before making the request
				if isValidURL(fullURL) {
					// Send a GET request to the URL and store the response
					resp, getErr := http.Get(fullURL)
					if getErr != nil {
						if verbose {
							fmt.Printf("%sError:%s %s\n", RED, NC, fullURL)
						}
						return
					}
					defer resp.Body.Close()

					if verbose {
						fmt.Printf("%sRequest:%s %s\n", GREEN, NC, fullURL)
					}

					// Check if the response contains "root:" to identify vulnerabilities
					body, readErr := ioutil.ReadAll(resp.Body)
					if readErr == nil && strings.Contains(string(body), "root:") {
						fmt.Printf("%sVulnerable:%s %s\n", RED, NC, fullURL)
						output.WriteString(fullURL + "\n")
					}
				}
			}(url, payload)
		}
	}

	// Wait for all tasks to complete
	wg.Wait()
}

// isValidURL checks if a URL is valid
func isValidURL(url string) bool {
	_, err := http.Get(url)
	return err == nil
}
