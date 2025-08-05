package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/chetankale/wikipedia-extraction/internal/extractor"
	_ "github.com/mattn/go-sqlite3"
)

// Storage interface defines methods for storing and retrieving quads
type Storage interface {
	// Store stores a collection of quads with metadata
	Store(quads []extractor.Quad, sourceURL string, extractedAt time.Time) error
	
	// GetBySubject retrieves all quads for a given subject
	GetBySubject(subject string) ([]extractor.Quad, error)
	
	// GetByRelationship retrieves all quads with a specific relationship
	GetByRelationship(relationship string) ([]extractor.Quad, error)
	
	// GetBySourceURL retrieves all quads from a specific source URL
	GetBySourceURL(sourceURL string) ([]extractor.Quad, error)
	
	// Search searches quads by text in any field
	Search(query string) ([]extractor.Quad, error)
	
	// GetStats returns storage statistics
	GetStats() (*Stats, error)
	
	// Close closes the storage connection
	Close() error
}

// Stats represents storage statistics
type Stats struct {
	TotalQuads     int    `json:"total_quads"`
	TotalSubjects  int    `json:"total_subjects"`
	TotalSources   int    `json:"total_sources"`
	LastExtraction string `json:"last_extraction"`
}

// QuadRecord represents a quad with metadata for storage
type QuadRecord struct {
	ID          int64     `json:"id"`
	Subject     string    `json:"subject"`
	Relationship string   `json:"relationship"`
	Value       string    `json:"value"`
	Citation    string    `json:"citation"`
	SourceURL   string    `json:"source_url"`
	ExtractedAt time.Time `json:"extracted_at"`
}

// SQLiteStorage implements Storage interface using SQLite
type SQLiteStorage struct {
	db *sql.DB
}

// NewSQLiteStorage creates a new SQLite storage instance
func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	
	// Create tables if they don't exist
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}
	
	return &SQLiteStorage{db: db}, nil
}

// createTables creates the necessary database tables
func createTables(db *sql.DB) error {
	quadsTable := `
	CREATE TABLE IF NOT EXISTS quads (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		subject TEXT NOT NULL,
		relationship TEXT NOT NULL,
		value TEXT NOT NULL,
		citation TEXT,
		source_url TEXT NOT NULL,
		extracted_at DATETIME NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	
	// Create indexes for better performance
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_quads_subject ON quads(subject);",
		"CREATE INDEX IF NOT EXISTS idx_quads_relationship ON quads(relationship);",
		"CREATE INDEX IF NOT EXISTS idx_quads_source_url ON quads(source_url);",
		"CREATE INDEX IF NOT EXISTS idx_quads_extracted_at ON quads(extracted_at);",
	}
	
	if _, err := db.Exec(quadsTable); err != nil {
		return err
	}
	
	for _, index := range indexes {
		if _, err := db.Exec(index); err != nil {
			return err
		}
	}
	
	return nil
}

// Store stores a collection of quads with metadata
func (s *SQLiteStorage) Store(quads []extractor.Quad, sourceURL string, extractedAt time.Time) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	
	stmt, err := tx.Prepare(`
		INSERT INTO quads (subject, relationship, value, citation, source_url, extracted_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()
	
	for _, quad := range quads {
		_, err := stmt.Exec(
			quad.Subject,
			quad.Relationship,
			quad.Value,
			quad.Citation,
			sourceURL,
			extractedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert quad: %w", err)
		}
	}
	
	return tx.Commit()
}

// GetBySubject retrieves all quads for a given subject
func (s *SQLiteStorage) GetBySubject(subject string) ([]extractor.Quad, error) {
	rows, err := s.db.Query(`
		SELECT subject, relationship, value, citation
		FROM quads
		WHERE subject LIKE ?
		ORDER BY extracted_at DESC
	`, "%"+subject+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to query quads: %w", err)
	}
	defer rows.Close()
	
	var quads []extractor.Quad
	for rows.Next() {
		var quad extractor.Quad
		err := rows.Scan(&quad.Subject, &quad.Relationship, &quad.Value, &quad.Citation)
		if err != nil {
			return nil, fmt.Errorf("failed to scan quad: %w", err)
		}
		quads = append(quads, quad)
	}
	
	return quads, nil
}

// GetByRelationship retrieves all quads with a specific relationship
func (s *SQLiteStorage) GetByRelationship(relationship string) ([]extractor.Quad, error) {
	rows, err := s.db.Query(`
		SELECT subject, relationship, value, citation
		FROM quads
		WHERE relationship LIKE ?
		ORDER BY extracted_at DESC
	`, "%"+relationship+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to query quads: %w", err)
	}
	defer rows.Close()
	
	var quads []extractor.Quad
	for rows.Next() {
		var quad extractor.Quad
		err := rows.Scan(&quad.Subject, &quad.Relationship, &quad.Value, &quad.Citation)
		if err != nil {
			return nil, fmt.Errorf("failed to scan quad: %w", err)
		}
		quads = append(quads, quad)
	}
	
	return quads, nil
}

// GetBySourceURL retrieves all quads from a specific source URL
func (s *SQLiteStorage) GetBySourceURL(sourceURL string) ([]extractor.Quad, error) {
	rows, err := s.db.Query(`
		SELECT subject, relationship, value, citation
		FROM quads
		WHERE source_url = ?
		ORDER BY extracted_at DESC
	`, sourceURL)
	if err != nil {
		return nil, fmt.Errorf("failed to query quads: %w", err)
	}
	defer rows.Close()
	
	var quads []extractor.Quad
	for rows.Next() {
		var quad extractor.Quad
		err := rows.Scan(&quad.Subject, &quad.Relationship, &quad.Value, &quad.Citation)
		if err != nil {
			return nil, fmt.Errorf("failed to scan quad: %w", err)
		}
		quads = append(quads, quad)
	}
	
	return quads, nil
}

// Search searches quads by text in any field
func (s *SQLiteStorage) Search(query string) ([]extractor.Quad, error) {
	rows, err := s.db.Query(`
		SELECT subject, relationship, value, citation
		FROM quads
		WHERE subject LIKE ? OR relationship LIKE ? OR value LIKE ? OR citation LIKE ?
		ORDER BY extracted_at DESC
	`, "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to query quads: %w", err)
	}
	defer rows.Close()
	
	var quads []extractor.Quad
	for rows.Next() {
		var quad extractor.Quad
		err := rows.Scan(&quad.Subject, &quad.Relationship, &quad.Value, &quad.Citation)
		if err != nil {
			return nil, fmt.Errorf("failed to scan quad: %w", err)
		}
		quads = append(quads, quad)
	}
	
	return quads, nil
}

// GetStats returns storage statistics
func (s *SQLiteStorage) GetStats() (*Stats, error) {
	var stats Stats
	
	// Get total quads
	err := s.db.QueryRow("SELECT COUNT(*) FROM quads").Scan(&stats.TotalQuads)
	if err != nil {
		return nil, fmt.Errorf("failed to get total quads: %w", err)
	}
	
	// Get total unique subjects
	err = s.db.QueryRow("SELECT COUNT(DISTINCT subject) FROM quads").Scan(&stats.TotalSubjects)
	if err != nil {
		return nil, fmt.Errorf("failed to get total subjects: %w", err)
	}
	
	// Get total unique sources
	err = s.db.QueryRow("SELECT COUNT(DISTINCT source_url) FROM quads").Scan(&stats.TotalSources)
	if err != nil {
		return nil, fmt.Errorf("failed to get total sources: %w", err)
	}
	
	// Get last extraction time
	err = s.db.QueryRow("SELECT MAX(extracted_at) FROM quads").Scan(&stats.LastExtraction)
	if err != nil {
		stats.LastExtraction = "Never"
	}
	
	return &stats, nil
}

// Close closes the storage connection
func (s *SQLiteStorage) Close() error {
	return s.db.Close()
} 