package migrations

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Migration struct {
	Version int
	Name    string
	Up      string
	Down    string
}

func RunMigrations(db *sql.DB) error {
	// Create migrations table if not exists
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			version INTEGER PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %v", err)
	}

	// Read migration files
	files, err := os.ReadDir("infrastructure/repository/migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %v", err)
	}

	// Sort files by name
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	// Execute migrations
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		// Check if migration has already been applied
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM migrations WHERE name = $1", file.Name()).Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %v", err)
		}

		if count > 0 {
			continue
		}

		// Read migration file content
		content, err := os.ReadFile(filepath.Join("infrastructure/repository/migrations", file.Name()))
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %v", file.Name(), err)
		}

		// Execute migration
		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %v", file.Name(), err)
		}

		// Record migration in migrations table
		_, err = db.Exec("INSERT INTO migrations (version, name) VALUES ($1, $2)", 1, file.Name())
		if err != nil {
			return fmt.Errorf("failed to record migration %s: %v", file.Name(), err)
		}

		fmt.Printf("Migration %s applied successfully\n", file.Name())
	}

	return nil
} 