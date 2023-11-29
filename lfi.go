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

	var (
		url       string
		urlsFile  string
		payloadsFile string
		outputFile string
		verbosity bool
		threads   int
	)

	flag.StringVar(&url, "u", "", "Single target URL")
	flag.StringVar(&urlsFile, "f", "", "File containing target URLs")
	flag.StringVar(&payloadsFile, "p", "", "File containing payloads")
	flag.StringVar(&outputFile, "o", "", "Output file for vulnerable URLs")
	flag.BoolVar(&verbosity, "v", false, "Enable verbosity for all requests")
	flag.IntVar(&threads, "t", 15, "Number of threads")
	flag.Parse()

	// Set the verbosity level based on the flag
	verbose = verbosity

	// Generate the banner
	banner := `
______  ___                      ______ _____ __        ______ _____________ 
___   |/  /_____________________ __|__ \__  // /        ___  / ___  __/___(_)
__  /|_/ / __  ___/_  ___/_  __ \____/ /_  // /___________  /  __  /_  __  / 
_  /  / /  _  /    / /__  / /_/ /_  __/ /__  __/_/_____/_  /____  __/  _  /  
/_/  /_/   /_/     \___/  \____/ /____/   /_/           /_____//_/     /_/     

`
	fmt.Print(CYAN, banner, NC)

	if (url == "" && urlsFile == "") || payloadsFile == "" || outputFile == "" {
		fmt.Println("Invalid input. Please provide the required flags.")
		return
	}

	var urls []string

	if url != "" {
		urls = append(urls, url)
	} else {
		_, err := os.Stat(urlsFile)
		if err != nil {
			fmt.Println("URLs file not found. Make sure the file exists.")
			return
		}

		urlsData, readURLsErr := ioutil.ReadFile(urlsFile)
		if readURLsErr != nil {
			fmt.Println("Error reading URLs file:", readURLsErr)
			return
		}
		urls = strings.Split(string(urlsData), "\n")
	}

	payloadsData, readPayloadsErr := ioutil.ReadFile(payloadsFile)
	if readPayloadsErr != nil {
		fmt.Println("Error reading payloads file:", readPayloadsErr)
		return
	}
	payloads := strings.Split(string(payloadsData), "\n")

	output, createOutputErr := os.Create(outputFile)
	if createOutputErr != nil {
		fmt.Println("Error creating output file:", createOutputErr)
		return
	}
	defer output.Close()

	var wg sync.WaitGroup
	threadCh := make(chan struct{}, threads)

	for _, u := range urls {
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

					fmt.Printf("%sRequest:%s %s\n", GREEN, NC, fullURL)

					body, readErr := ioutil.ReadAll(resp.Body)
					if readErr == nil && strings.Contains(string(body), "root:") {
						fmt.Printf("%sVulnerable:%s %s\n", RED, NC, fullURL)
						output.WriteString(fmt.Sprintf("Request URL: %s\nVulnerable URL: %s\n", fullURL, fullURL))
					}
				}
			}(u, payload)
		}
	}

	wg.Wait()
}

// isValidURL checks if a URL is valid
func isValidURL(url string) bool {
	_, err := http.Get(url)
	return err == nil
}
