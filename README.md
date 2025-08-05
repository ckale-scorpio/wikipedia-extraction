# Wikipedia Extraction Tool

A Go-based tool to extract structured information from Wikipedia pages. Pulls, parses and stores data from infoboxes in the form of quads (subject or entity, relationship, value and citation).

## Features

- Extract structured data from Wikipedia infoboxes
- Parse Wikipedia tables for additional data
- Output in multiple formats (JSON, CSV, XML)
- Persistent storage with SQLite database
- Advanced querying capabilities
- Command-line interface with configurable options
- Respectful web scraping with proper user agent

## Installation

### Prerequisites

- Go 1.21 or later

### Build from source

```bash
# Clone the repository
git clone https://github.com/chetankale/wikipedia-extraction.git
cd wikipedia-extraction

# Install dependencies
make install

# Build the application
make build
```

## Usage

### Basic usage

```bash
# Extract data from a Wikipedia page (output to file)
./bin/wikipedia-extraction extract "https://en.wikipedia.org/wiki/Go_(programming_language)"

# Specify output file and format
./bin/wikipedia-extraction extract "https://en.wikipedia.org/wiki/Python_(programming_language)" \
  --output python_data.json --format json

# Store data in database
./bin/wikipedia-extraction store "https://en.wikipedia.org/wiki/Go_(programming_language)"

# Query stored data
./bin/wikipedia-extraction query --subject "Go"
./bin/wikipedia-extraction query --relationship "Designed"
./bin/wikipedia-extraction query --search "Robert"
./bin/wikipedia-extraction query --stats
```

### Command options

#### Extract command
- `--output`: Output file path (default: output.json)
- `--format`: Output format - json, csv, or xml (default: json)
- `--config`: Configuration file path

#### Query command
- `--subject`: Search by subject name
- `--relationship`: Search by relationship type
- `--source`: Search by source URL
- `--search`: Full-text search across all fields
- `--stats`: Show database statistics

### Example output

The tool extracts quads in the format:
```
Subject | Relationship | Value | Citation
```

For example:
```
Go (programming language) | Designed by | Robert Griesemer, Rob Pike, Ken Thompson | infobox
Go (programming language) | First appeared | November 10, 2009 | infobox
Go (programming language) | Paradigm | Multi-paradigm: concurrent, functional, imperative, object-oriented | infobox
```

## Development

### Project structure

```
wikipedia-extraction/
├── cmd/                    # Command-line interface
│   ├── root.go            # Root command setup
│   ├── extract.go         # Extract command
│   ├── store.go           # Store command
│   └── query.go           # Query command
├── internal/              # Internal packages
│   ├── extractor/         # Wikipedia extraction logic
│   ├── output/           # Output formatting
│   └── storage/          # Database storage layer
├── main.go               # Application entry point
├── go.mod               # Go module definition
├── Makefile             # Build and development tasks
└── README.md           # This file
```

### Available make targets

```bash
make build      # Build the application
make run        # Build and run the application
make install    # Install dependencies
make test       # Run tests
make clean      # Clean build artifacts
make example    # Run with example Wikipedia page
make build-all  # Build for multiple platforms
make help       # Show all available targets
```

### Database

The application uses SQLite for persistent storage. The database file (`quads.db`) is created automatically when you first use the `store` command. The database includes:

- **Quads table**: Stores all extracted quads with metadata
- **Indexes**: Optimized for fast querying by subject, relationship, and source
- **Statistics**: Track total quads, subjects, and sources

### Running tests

```bash
go test ./...
```

## Output Formats

### JSON
```json
[
  {
    "subject": "Go (programming language)",
    "relationship": "Designed by",
    "value": "Robert Griesemer, Rob Pike, Ken Thompson",
    "citation": "infobox"
  }
]
```

### CSV
```csv
Subject,Relationship,Value,Citation
Go (programming language),Designed by,Robert Griesemer Rob Pike Ken Thompson,infobox
```

### XML
```xml
<quads>
  <quad>
    <subject>Go (programming language)</subject>
    <relationship>Designed by</relationship>
    <value>Robert Griesemer, Rob Pike, Ken Thompson</value>
    <citation>infobox</citation>
  </quad>
</quads>
```

## License

This project is licensed under the CC0 1.0 Universal license - see the [LICENSE](LICENSE) file for details. 
