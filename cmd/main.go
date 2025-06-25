package main

import (
	"Lead-Automation-Pipeline/cmd/utils"
	"fmt"
	"slices"
	"sync"
)

type SafeSeenLeadsWebsites struct {
	mu sync.Mutex
	websites []string
}

func main() {
	// The websites of all the leads whos been seen so far
	var seenLeadsWebsites = SafeSeenLeadsWebsites{}

	// Runs the google maps scraper

	// Reads the email, phone number, and website from the results
	var leads = utils.ReadCSV("output.csv")

	// Checks whether any of the clients have already been scraped
	for i := 0; i < len(leads); i++ {
		if slices.Contains(seenLeadsWebsites.websites, leads[i].Website) {
			leads = append(leads[:i], leads[i+1:]...)
			i--
		}
	}

	var wg sync.WaitGroup
	var numOfGoroutines = 5

	// Splits the leads into chunks, one chunk per goroutine
	var leadChunks = splitLeadsUp(leads, numOfGoroutines)

	// Starts a goroutine for each chunk
	for _, chunk := range leadChunks {
		wg.Add(1)

		// This anonymous function stops you from needing
		// to loop over the chunk inside the processLead() function
		go func(c []utils.Lead) {
			defer wg.Done()

			for _, lead := range chunk {
				processLead(lead, &seenLeadsWebsites)
			}
		}(chunk)
	}

	wg.Wait()

	// Saves the final results in the output file specified by the user
}

func processLead(lead utils.Lead, seenLeadsWebsites *SafeSeenLeadsWebsites) {
	// Scrapes the website for emails, and copies the HTML

	// Converts the HTML to Markdown, and limits it to 10,000 characters

	// Calls OpenAI to write a summarisation of the website's content

	// Calls OpenAI to write an email icebreaker

	// Save the email to the array used to check if a lead has been scraped before
	seenLeadsWebsites.mu.Lock()
	seenLeadsWebsites.websites = append(seenLeadsWebsites.websites, lead.Website)
	seenLeadsWebsites.mu.Unlock()
}

// Splits the lead array up into multiple chunks
func splitLeadsUp(leads []utils.Lead, numOfChunks int) [][]utils.Lead {
    var leadChunks [][]utils.Lead
    total := len(leads)

    if numOfChunks <= 0 {
        return nil // or return the whole slice as one chunk
    }
    if numOfChunks > total {
        numOfChunks = total
    }

    chunkSize := total / numOfChunks
    remainder := total % numOfChunks

    start := 0
    for i := 0; i < numOfChunks; i++ {
        // Add 1 to the chunk size for the first 'remainder' chunks to spread leftovers evenly
        size := chunkSize
        if i < remainder {
            size++
        }
        end := start + size
        leadChunks = append(leadChunks, leads[start:end])
        start = end
    }

    return leadChunks
}

