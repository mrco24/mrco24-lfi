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
	urlsFile := flag.String("u", "urls.txt", "File containing target URLs")
	payloadsFile := flag.String("p", "payloads.txt", "File containing payloads")
	outputFile := flag.String("o", "vulnerable_urls.txt", "Output file for vulnerable URLs")
	verbosity := flag.Bool("v", false, "Enable verbosity for all requests")
	threads := flag.Int("t", 10, "Number of threads")
	flag.Parse()

	// Set the verbosity level based on the flag
	verbose = *verbosity

	// Generate the banner
	banner := `
______  ___                      ______ _____ __        ______ _____________ 
___   |/  /_____________________ __|__ \__  // /        ___  / ___  __/___(_)
__  /|_/ / __  ___/_  ___/_  __ \____/ /_  // /___________  /  __  /_  __  / 
_  /  / /  _  /    / /__  / /_/ /_  __/ /__  __/_/_____/_  /____  __/  _  /  
/_/  /_/   /_/     \___/  \____/ /____/   /_/           /_____//_/     /_/     
`
	fmt.Print(CYAN, banner, NC)

	_, urlsErr := os.Stat(*urlsFile)
	_, payloadsErr := os.Stat(*payloadsFile)

	if urlsErr != nil || payloadsErr != nil {
		fmt.Println("URLs file or payloads file not found. Make sure both files exist.")
		return
	}

	urlsData, readURLsErr := ioutil.ReadFile(*urlsFile)
	if readURLsErr != nil {
		fmt.Println("Error reading URLs file:", readURLsErr)
		return
	}
	urls := strings.Split(string(urlsData), "\n")

	payloadsData, readPayloadsErr := ioutil.ReadFile(*payloadsFile)
	if readPayloadsErr != nil {
		fmt.Println("Error reading payloads file:", readPayloadsErr)
		return
	}
	payloads := strings.Split(string(payloadsData), "\n")

	output, createOutputErr := os.Create(*outputFile)
	if createOutputErr != nil {
		fmt.Println("Error creating output file:", createOutputErr)
		return
	}
	defer output.Close()

	var wg sync.WaitGroup

	threadCh := make(chan struct{}, *threads)

	for _, url := range urls {
		for _, payload := range payloads {
			threadCh <- struct{}{}

			wg.Add(1)
			go func(url, payload string) {
				defer wg.Done()
				defer func() {
					<-threadCh
				}()

				fullURL := url + payload

				if isValidURL(fullURL) {
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

					body, readErr := ioutil.ReadAll(resp.Body)
					if readErr == nil && strings.Contains(string(body), "root:") {
						fmt.Printf("%sVulnerable:%s %s\n", RED, NC, fullURL)
						output.WriteString(fullURL + "\n")
					}
				}
			}(url, payload)
		}
	}

	wg.Wait()
}

// isValidURL checks if a URL is valid
func isValidURL(url string) bool {
	_, err := http.Get(url)
	return err == nil
}
