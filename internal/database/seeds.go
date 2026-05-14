package database

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func RunSeeds(ctx context.Context, db *sql.DB, seedsPath string) error {
	entries, err := os.ReadDir(seedsPath)
	if err != nil {
		return err
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	for _, name := range names {
		body, err := os.ReadFile(filepath.Join(seedsPath, name))
		if err != nil {
			return err
		}
		if _, err := db.ExecContext(ctx, string(body)); err != nil {
			return err
		}
	}
	return nil
}
