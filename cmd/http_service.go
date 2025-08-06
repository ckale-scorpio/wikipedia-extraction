package cmd

import (
	"log"
	"net/http"

	"github.com/chetankale/wikipedia-extraction/internal/extractor"
	"github.com/chetankale/wikipedia-extraction/internal/output"
	"github.com/spf13/cobra"
)

var httpServiceCmd = &cobra.Command{
	Use:   "http-service",
	Short: "Start a HTTP service that extracts structured data from Wikipedia pages",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		StartHTTPServer()
	},
}

func init() {
	rootCmd.AddCommand(httpServiceCmd)
}

func StartHTTPServer() {

	
	http.HandleFunc("/extract", func(w http.ResponseWriter, r *http.Request) {
		src := r.URL.Query().Get("src")
		if src == "" {
			log.Println("No source URL provided")
			http.Error(w, "No source URL provided", http.StatusBadRequest)
			return
		}
		// Create extractor
		ext := extractor.NewExtractor()

		quads, err := ext.ExtractFromURL(src)
		if err != nil {
			log.Println("Error: %v", err)
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		formatter := output.NewFormatter()
		if err := formatter.WriteQuads(quads, w, format); err != nil {
			log.Println("Failed to write output: %v", err)
			http.Error(w, "Failed to write output: "+err.Error(), http.StatusInternalServerError)
			return
		}
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}