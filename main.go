package main

import (
	"fmt"
	"io"
	"net/http"
)

const requestURL = "https://index.golang.org/index"

//type moduleInfo struct {
//	Path      string `json:"Path"`
//	Version   string `json:"Version"`
//	Timestamp string `json:"Timestamp"`
//}
//
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
	fmt.Printf("response body: %s\n", respBody)
}
