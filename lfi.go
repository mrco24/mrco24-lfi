package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "strings"
    "sync"
)

// Define ANSI color codes
const (
    RED  = "\033[0;31m"
    CYAN = "\033[0;36m"
    NC   = "\033[0m" // No Color
)

func main() {
    // Generate the banner
    banner := `
  ____ _   _ ____ ____ ____ _  _ 
   __  |   | |__/ |    |__| |  | 
   __) |___| |  \ |___ |  | |__| 
`
    fmt.Print(CYAN, banner, NC)

    // Define the target URLs and payloads files
    urlsFile := "urls.txt"
    payloadsFile := "payloads.txt"
    outputFile := "vulnerable_urls.txt"

    // Check if the URLs and payloads files exist
    _, urlsErr := os.Stat(urlsFile)
    _, payloadsErr := os.Stat(payloadsFile)

    if urlsErr != nil || payloadsErr != nil {
        fmt.Println("URLs file or payloads file not found. Make sure both files exist.")
        return
    }

    // Read the list of target URLs from the URLs file
    urlsData, readURLsErr := ioutil.ReadFile(urlsFile)
    if readURLsErr != nil {
        fmt.Println("Error reading URLs file:", readURLsErr)
        return
    }
    urls := strings.Split(string(urlsData), "\n")

    // Read the list of payloads from the payloads file
    payloadsData, readPayloadsErr := ioutil.ReadFile(payloadsFile)
    if readPayloadsErr != nil {
        fmt.Println("Error reading payloads file:", readPayloadsErr)
        return
    }
    payloads := strings.Split(string(payloadsData), "\n")

    // Initialize the output file
    output, createOutputErr := os.Create(outputFile)
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

                // Send a GET request to the URL and store the response
                resp, getErr := http.Get(fullURL)
                if getErr != nil {
                    return
                }
                defer resp.Body.Close()

                // Check if the response contains "root:" to identify vulnerabilities
                body, readErr := ioutil.ReadAll(resp.Body)
                if readErr == nil && strings.Contains(string(body), "root:") {
                    fmt.Printf("%sVulnerable:%s %s\n", RED, NC, fullURL)
                    output.WriteString(fullURL + "\n")
                }
            }(url, payload)
        }
    }

    // Wait for all tasks to complete
    wg.Wait()
}
