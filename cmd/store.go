package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chetankale/wikipedia-extraction/internal/extractor"
	"github.com/chetankale/wikipedia-extraction/internal/storage"
	"github.com/spf13/cobra"
)

var storeCmd = &cobra.Command{
	Use:   "store [URL]",
	Short: "Extract and store structured data from a Wikipedia page",
	Long: `Extract structured information from a Wikipedia page URL and store it in the database.
The tool will parse infoboxes and extract quads in the form of:
(subject/entity, relationship, value, citation) and store them persistently.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		
		// Validate URL
		if !strings.Contains(url, "wikipedia.org") {
			log.Fatal("URL must be a Wikipedia page")
		}

		// Initialize storage
		dbPath := "quads.db"
		store, err := storage.NewSQLiteStorage(dbPath)
		if err != nil {
			log.Fatalf("Failed to initialize storage: %v", err)
		}
		defer store.Close()

		// Create extractor
		ext := extractor.NewExtractor()

		// Extract data
		quads, err := ext.ExtractFromURL(url)
		if err != nil {
			log.Fatalf("Failed to extract data: %v", err)
		}

		// Store data
		err = store.Store(quads, url, time.Now())
		if err != nil {
			log.Fatalf("Failed to store data: %v", err)
		}

		// Output results
		fmt.Printf("Extracted and stored %d quads from %s\n", len(quads), url)
		
		// Display first few quads as preview
		fmt.Println("\nPreview of extracted data:")
		for i, quad := range quads {
			if i >= 5 { // Show only first 5
				break
			}
			fmt.Printf("Quad %d: %s | %s | %s | %s\n", 
				i+1, quad.Subject, quad.Relationship, quad.Value, quad.Citation)
		}
	},
}

func init() {
	rootCmd.AddCommand(storeCmd)
} 