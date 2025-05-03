package cache

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func IsCacheValid (db *sql.DB, table_name string) bool {
	query := "SELECT ttl, updated_at FROM cacheentries WHERE table_name = $1"

	row := db.QueryRow(query, table_name)

	var ttlseconds int64
	var updated_at time.Time

	err := row.Scan(&ttlseconds, &updated_at)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Failed to retrieve rows: %w", err)
			return false
		}
		fmt.Println("Failed to scan row: %w", err)
		return false
	}

	ttl := time.Duration(ttlseconds) * time.Second

	return time.Since(updated_at) < ttl
}


func UpdateCacheTimestamp (db *sql.DB, table_name string) (bool, error) {
	query := "UPDATE updated_at SET updated_at = NOW() WHERE table_name = $1" 

	result, err := db.Exec(query, table_name)

	if err != nil {
		return false, fmt.Errorf("failed to update timestamp: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to update rows: %w", err)
	}

	return rowsAffected > 0, nil
}
