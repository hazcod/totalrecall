package recall

import (
	"database/sql"
	"fmt"
)

type WebResult struct {
	Timestamp int64
	Domain    string
	URL       string
}

func (r *Recall) ExtractWeb() ([]WebResult, error) {
	conn, err := sql.Open("sqlite3", r.dbPath)
	if err != nil {
		return nil, fmt.Errorf("could not open database connection: %w", err)
	}
	defer conn.Close()

	query := `SELECT Timestamp, Domain, Uri AS URL FROM Web;`
	rows, err := conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	defer rows.Close()

	var results []WebResult

	for rows.Next() {
		var timestamp int64
		var domain string
		var url string

		err := rows.Scan(&timestamp, &domain, &url)
		if err != nil {
			return nil, fmt.Errorf("could not scan row: %w", err)
		}

		r.logger.WithField("domain", domain).Debug("processing Recall web")

		results = append(results, WebResult{
			Timestamp: timestamp,
			Domain:    domain,
			URL:       url,
		})
	}

	return results, nil
}
