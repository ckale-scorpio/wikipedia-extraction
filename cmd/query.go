package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/chetankale/wikipedia-extraction/internal/extractor"
	"github.com/chetankale/wikipedia-extraction/internal/storage"
	"github.com/spf13/cobra"
)

var (
	querySubject     string
	queryRelationship string
	querySourceURL   string
	querySearch      string
	queryStats       bool
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query stored quads from the database",
	Long: `Query stored quads from the database using various filters.
You can search by subject, relationship, source URL, or use full-text search.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize storage
		dbPath := "quads.db"
		store, err := storage.NewSQLiteStorage(dbPath)
		if err != nil {
			log.Fatalf("Failed to initialize storage: %v", err)
		}
		defer store.Close()

		var quads []extractor.Quad
		var err2 error

		// Handle different query types
		switch {
		case queryStats:
			stats, err := store.GetStats()
			if err != nil {
				log.Fatalf("Failed to get stats: %v", err)
			}
			
			fmt.Printf("Database Statistics:\n")
			fmt.Printf("  Total Quads: %d\n", stats.TotalQuads)
			fmt.Printf("  Total Subjects: %d\n", stats.TotalSubjects)
			fmt.Printf("  Total Sources: %d\n", stats.TotalSources)
			fmt.Printf("  Last Extraction: %s\n", stats.LastExtraction)
			return

		case querySubject != "":
			quads, err2 = store.GetBySubject(querySubject)

		case queryRelationship != "":
			quads, err2 = store.GetByRelationship(queryRelationship)

		case querySourceURL != "":
			quads, err2 = store.GetBySourceURL(querySourceURL)

		case querySearch != "":
			quads, err2 = store.Search(querySearch)

		default:
			fmt.Println("Please specify a query type. Use --help for options.")
			return
		}

		if err2 != nil {
			log.Fatalf("Failed to query data: %v", err2)
		}

		// Output results
		if len(quads) == 0 {
			fmt.Println("No quads found matching the query.")
			return
		}

		fmt.Printf("Found %d quads:\n\n", len(quads))

		// Output in the specified format
		switch format {
		case "json":
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			encoder.Encode(quads)
		case "csv":
			// Simple CSV output
			fmt.Println("Subject,Relationship,Value,Citation")
			for _, quad := range quads {
				fmt.Printf("\"%s\",\"%s\",\"%s\",\"%s\"\n",
					quad.Subject, quad.Relationship, quad.Value, quad.Citation)
			}
		default:
			// Default table format
			for i, quad := range quads {
				fmt.Printf("Quad %d:\n", i+1)
				fmt.Printf("  Subject: %s\n", quad.Subject)
				fmt.Printf("  Relationship: %s\n", quad.Relationship)
				fmt.Printf("  Value: %s\n", quad.Value)
				fmt.Printf("  Citation: %s\n", quad.Citation)
				fmt.Println()
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)
	
	// Query flags
	queryCmd.Flags().StringVar(&querySubject, "subject", "", "Search by subject")
	queryCmd.Flags().StringVar(&queryRelationship, "relationship", "", "Search by relationship")
	queryCmd.Flags().StringVar(&querySourceURL, "source", "", "Search by source URL")
	queryCmd.Flags().StringVar(&querySearch, "search", "", "Full-text search")
	queryCmd.Flags().BoolVar(&queryStats, "stats", false, "Show database statistics")
} 