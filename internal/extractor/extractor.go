package extractor

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

// Quad represents a structured data point extracted from Wikipedia
type Quad struct {
	Subject     string `json:"subject"`
	Relationship string `json:"relationship"`
	Value       string `json:"value"`
	Citation    string `json:"citation"`
}

// Extractor handles Wikipedia page extraction
type Extractor struct {
	colly *colly.Collector
}

// NewExtractor creates a new Wikipedia extractor
func NewExtractor() *Extractor {
	c := colly.NewCollector(
		colly.UserAgent("Wikipedia-Extraction/1.0"),
	)

	return &Extractor{
		colly: c,
	}
}

// ExtractFromURL extracts structured data from a Wikipedia URL
func (e *Extractor) ExtractFromURL(url string) ([]Quad, error) {
	var quads []Quad
	var references map[string]string

	e.colly.OnHTML("body", func(h *colly.HTMLElement) {
		doc := h.DOM

		// Extract page title
		title := doc.Find("h1#firstHeading").Text()
		if title == "" {
			title = doc.Find("title").Text()
		}

		// First, extract all references from the references section
		references = e.extractReferences(h.DOM)

		// Find and parse infoboxes
		doc.Find(".infobox").Each(func(i int, s *goquery.Selection) {
			infoboxQuads := e.parseInfobox(s, title, references)
			quads = append(quads, infoboxQuads...)
		})

		// Find and parse other structured data tables
		doc.Find("table.wikitable").Each(func(i int, s *goquery.Selection) {
			tableQuads := e.parseTable(s, title, references)
			quads = append(quads, tableQuads...)
		})
	})

	err := e.colly.Visit(url)
	if err != nil {
		return nil, fmt.Errorf("failed to visit URL: %w", err)
	}

	return quads, nil
}

// parseInfobox extracts quads from a Wikipedia infobox
func (e *Extractor) parseInfobox(infobox *goquery.Selection, subject string, references map[string]string) []Quad {
	var quads []Quad

	infobox.Find("tr").Each(func(i int, s *goquery.Selection) {
		// Skip header rows
		if s.HasClass("infobox-header") || s.HasClass("infobox-subheader") {
			return
		}

		// Extract label and value
		label := strings.TrimSpace(s.Find("th").Text())
		valueCell := s.Find("td")
		value := strings.TrimSpace(valueCell.Text())

		if label != "" && value != "" {
			// Extract citations from the value cell
			citations := e.extractCitations(valueCell, references)
			
			quad := Quad{
				Subject:     subject,
				Relationship: label,
				Value:       value,
				Citation:    citations,
			}
			quads = append(quads, quad)
		}
	})

	return quads
}

// parseTable extracts quads from a Wikipedia table
func (e *Extractor) parseTable(table *goquery.Selection, subject string, references map[string]string) []Quad {
	var quads []Quad

	table.Find("tr").Each(func(i int, s *goquery.Selection) {
		cells := s.Find("td, th")
		if cells.Length() >= 2 {
			label := strings.TrimSpace(cells.Eq(0).Text())
			valueCell := cells.Eq(1)
			value := strings.TrimSpace(valueCell.Text())

			if label != "" && value != "" {
				// Extract citations from the value cell
				citations := e.extractCitations(valueCell, references)
				
				quad := Quad{
					Subject:     subject,
					Relationship: label,
					Value:       value,
					Citation:    citations,
				}
				quads = append(quads, quad)
			}
		}
	})

	return quads
}

// extractCitations extracts citation links by following named anchors to the references section
func (e *Extractor) extractCitations(cell *goquery.Selection, references map[string]string) string {
	var citations []string
	citationMap := make(map[string]bool)
	
	// Find all citation links in the cell
	cell.Find("a[href*='#cite_note']").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			// Extract the citation ID from the href
			if strings.Contains(href, "#cite_note-") {
				citationID := strings.TrimPrefix(href, "#cite_note-")
				// Look up the actual citation from the references map
				referenceKey := "cite_note-" + citationID
				if actualCitation, exists := references[referenceKey]; exists {
					if !citationMap[actualCitation] {
						citationMap[actualCitation] = true
						citations = append(citations, actualCitation)
					}
				}
			}
		}
	})
	
	// Also look for superscript citation links
	cell.Find("sup a").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			// Extract the citation ID from the href
			if strings.Contains(href, "#cite_note-") {
				citationID := strings.TrimPrefix(href, "#cite_note-")
				// Look up the actual citation from the references map
				referenceKey := "cite_note-" + citationID
				if actualCitation, exists := references[referenceKey]; exists {
					if !citationMap[actualCitation] {
						citationMap[actualCitation] = true
						citations = append(citations, actualCitation)
					}
				}
			}
		}
	})
	
	// If no citations found, return "no citation"
	if len(citations) == 0 {
		return "no citation"
	}
	
	return strings.Join(citations, "; ")
}

// extractReferences extracts all references from the references section
func (e *Extractor) extractReferences(doc *goquery.Selection) map[string]string {
	references := make(map[string]string)
	
	// Find the references section - Wikipedia uses various selectors
	doc.Find("#References, #references, .reflist, .references").Each(func(i int, s *goquery.Selection) {
		// Find all reference list items
		s.Find("li").Each(func(j int, li *goquery.Selection) {
			// Extract the reference ID
			if id, exists := li.Attr("id"); exists {
				// Look for external links in the reference
				li.Find("a[href^='http']").Each(func(k int, a *goquery.Selection) {
					if href, exists := a.Attr("href"); exists {
						references[id] = href
					}
				})
			}
		})
	})
	
	// Also look for cite_note references
	doc.Find("ol.references li").Each(func(i int, li *goquery.Selection) {
		if id, exists := li.Attr("id"); exists {
			// Look for external links in the reference
			li.Find("a[href^='http']").Each(func(k int, a *goquery.Selection) {
				if href, exists := a.Attr("href"); exists {
					references[id] = href
				}
			})
		}
	})
	
	return references
} 