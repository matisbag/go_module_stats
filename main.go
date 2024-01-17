package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

const requestURL = "https://index.golang.org/index"

// Forge represents a single forge with its path, version, and timestamp.
type Forge struct {
	Path      string `json:"Path"`
	Version   string `json:"Version"`
	Timestamp string `json:"Timestamp"`
}

// ForgeStats holds the statistics for a forge including its name (forge), number of modules, and versions.
type ForgeStats struct {
	Forge    string
	Modules  int
	Versions int
}

func main() {
	forges, err := getForges()
	if err != nil {
		fmt.Println("error getting forges:", err)
		return
	}

	stats := calculateStats(forges)

	sortedStatsVersions := sortStatsVersions(stats)
	printTable(sortedStatsVersions)

	sortedStatsModules := sortStatsModules(stats)
	printTable(sortedStatsModules)

	sortedStatsForge := sortStatsForge(stats)
	printTable(sortedStatsForge)
}

func getForges() ([]Forge, error) {
	// Create a new HTTP client
	client := &http.Client{}

	// Create a new HTTP request
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	// Set a custom header for the request
	req.Header.Add("Disable-Module-Fetch", "true")
	req.Header.Add("Content-Type", "application/json")

	// Send the HTTP request and get the response
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making http request: %w", err)
	}
	// Ensure the response body is closed after function return
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}

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
			return nil, fmt.Errorf("error decoding JSON: %w", err)
		}

		// Add the forge to the slice
		forges = append(forges, mod)
	}

	return forges, nil
}

func calculateStats(forges []Forge) []ForgeStats {
	// Use a map to store the statistics of each forge
	forgeStats := make(map[string]*ForgeStats)

	// Use a map to store the checked modules
	modulesChecked := make(map[string]bool)

	for _, element := range forges {
		// Split the path once and store the result
		forgeName := strings.Split(element.Path, "/")[0]

		// Check if the module has already been checked
		_, foundModule := modulesChecked[element.Path]
		if !foundModule {
			modulesChecked[element.Path] = true
		}

		// Get the statistics of the forge
		stat, foundForge := forgeStats[forgeName]
		if !foundForge {
			// If the forge does not exist, create a new statistic
			stat = &ForgeStats{
				Forge:    forgeName,
				Modules:  0,
				Versions: 0,
			}
			forgeStats[forgeName] = stat
		}

		// Update the statistics
		stat.Versions++
		if !foundModule {
			stat.Modules++
		}
	}

	// Convert the map into a slice to sort the results
	var stats []ForgeStats
	for _, stat := range forgeStats {
		stats = append(stats, *stat)
	}

	return stats
}

func sortStatsVersions(stats []ForgeStats) []ForgeStats {
	// Sort the slice according to the requested criteria
	sort.Slice(stats, func(i, j int) bool {
		// Sort: 1. Versions DESC, 2. Modules DESC, 3. Forge ASC
		if stats[i].Versions == stats[j].Versions {
			if stats[i].Modules == stats[j].Modules {
				return stats[i].Forge < stats[j].Forge
			}
			return stats[i].Modules > stats[j].Modules
		}
		return stats[i].Versions > stats[j].Versions
	})

	return stats
}

func sortStatsModules(stats []ForgeStats) []ForgeStats {
	// Sort the slice according to the requested criteria
	sort.Slice(stats, func(i, j int) bool {
		// Sort: 1. Modules DESC, 2. Versions DESC, 3. Forge ASC
		if stats[i].Modules == stats[j].Modules {
			if stats[i].Versions == stats[j].Versions {
				return stats[i].Forge < stats[j].Forge
			}
			return stats[i].Versions > stats[j].Versions
		}
		return stats[i].Modules > stats[j].Modules
	})

	return stats
}

func sortStatsForge(stats []ForgeStats) []ForgeStats {
	// Sort the slice according to the requested criteria
	sort.Slice(stats, func(i, j int) bool {
		// Sort: 1. Forge ASC, 2. Versions DESC, 3. Modules DESC
		if stats[i].Forge == stats[j].Forge {
			if stats[i].Versions == stats[j].Versions {
				return stats[i].Modules > stats[j].Modules
			}
			return stats[i].Versions > stats[j].Versions
		}
		return stats[i].Forge < stats[j].Forge
	})

	return stats
}

func printTable(stats []ForgeStats) {
	// Use tab writer for a clean formatting of the table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	defer w.Flush()

	// Table headers
	fmt.Fprintln(w, "Forge\tModules\tVersions")

	// Data lines
	for _, stat := range stats {
		fmt.Fprintf(w, "%s\t%d\t%d\n", stat.Forge, stat.Modules, stat.Versions)
	}
}
