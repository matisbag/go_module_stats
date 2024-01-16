package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const requestURL = "https://index.golang.org/index"

type Forge struct {
	Path      string `json:"Path"`
	Version   string `json:"Version"`
	Timestamp string `json:"Timestamp"`
}

//type forgeStats struct {
//	Forge    string
//	Modules  int
//	Versions int
//}

func main() {
	// Create a new HTTP client
	client := &http.Client{}

	// Create a new HTTP request
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		fmt.Println("error creating HTTP request:", err)
		return
	}

	// Set a custom header for the request
	req.Header.Set("Disable-Module-Fetch", "true")
	req.Header.Add("Content-Type", "application/json")

	// Send the HTTP request and get the response
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error making http request:", err)
		return
	}
	// Ensure the response body is closed after function return
	defer resp.Body.Close()

	fmt.Println("client: got response!")
	fmt.Printf("client: status code: %d\n", resp.StatusCode)

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("could not read response body:", err)
		return
	}
	// Print the response body
	// fmt.Printf("response body: %s\n", respBody)

	// Convert the response body to a reader
	reader := strings.NewReader(string(respBody))

	// Create a new JSON decoder
	decoder := json.NewDecoder(reader)

	// Create a slice to hold the forges
	var forges []Forge

	// Loop and read each object from the stream
	for {
		var mod Forge

		// Decode the next JSON object
		err := decoder.Decode(&mod)

		// If an error occurred, break the loop
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("error decoding JSON:", err)
			return
		}

		// Add the forge to the slice
		forges = append(forges, mod)
	}

	// Print the forges
	for _, forge := range forges {
		fmt.Printf("Path: %s, Version: %s, Timestamp: %s\n", forge.Path, forge.Version, forge.Timestamp)
	}
}
