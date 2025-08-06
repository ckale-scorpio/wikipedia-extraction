package cmd

import (
	"fmt"
	"os"
	"log"
	"strings"

	"github.com/chetankale/wikipedia-extraction/internal/extractor"
	"github.com/chetankale/wikipedia-extraction/internal/output"
	"github.com/spf13/cobra"
)

var extractCmd = &cobra.Command{
	Use:   "extract [URL]",
	Short: "Extract structured data from a Wikipedia page",
	Long: `Extract structured information from a Wikipedia page URL.
The tool will parse infoboxes and extract quads in the form of:
(subject/entity, relationship, value, citation)`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		
		// Validate URL
		if !strings.Contains(url, "wikipedia.org") {
			log.Fatal("URL must be a Wikipedia page")
		}

		// Create extractor
		ext := extractor.NewExtractor()

		// Extract data
		quads, err := ext.ExtractFromURL(url)
		if err != nil {
			log.Fatalf("Failed to extract data: %v", err)
		}

		// Output results
		fmt.Printf("Extracted %d quads from %s\n", len(quads), url)
		
		fileWriter, err := os.Create(outputFile)
		if err != nil {
			fmt.Errorf("failed to create output file: %w", err)
			return
		}
		defer fileWriter.Close()

		// Save to file
		formatter := output.NewFormatter()
		if err := formatter.WriteQuads(quads, fileWriter, format); err != nil {
			log.Fatalf("Failed to write output: %v", err)
		}
		
		fmt.Printf("Results saved to %s in %s format\n", outputFile, format)
		
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
	rootCmd.AddCommand(extractCmd)
} 