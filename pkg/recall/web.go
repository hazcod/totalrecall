package recall

import (
	"database/sql"
	"fmt"
)

type WebResult struct {
	Domain string
	URL    string
}

func (r *Recall) ExtractWeb() ([]WebResult, error) {
	conn, err := sql.Open("sqlite3", r.dbPath)
	if err != nil {
		return nil, fmt.Errorf("could not open database connection: %w", err)
	}
	defer conn.Close()

	query := `SELECT Domain, Uri AS URL FROM Web;`
	rows, err := conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	defer rows.Close()

	var results []WebResult

	for rows.Next() {
		var domain string
		var url string

		err := rows.Scan(&domain, &url)
		if err != nil {
			return nil, fmt.Errorf("could not scan row: %w", err)
		}

		r.logger.WithField("domain", domain).Debug("processing Recall web")

		results = append(results, WebResult{
			Domain: domain,
			URL:    url,
		})
	}

	return results, nil
}
