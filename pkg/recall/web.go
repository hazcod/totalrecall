package recall

import (
	"database/sql"
	"fmt"
)

type WebResult struct {
	Domain string
	URL    string
}

func (r *Recall) ExtractWeb(userName string) ([]WebResult, error) {
	if userName == "" {
		return nil, fmt.Errorf("no username provided, provide one or use ExtractImagesForCurrentUser")
	}

	dbPath, imagePath, err := r.GetRecallPaths(userName)
	if err != nil {
		return nil, err
	}

	r.logger.WithField("db_path", dbPath).WithField("image_path", imagePath).
		Debug("Recall feature found enabled")

	conn, err := sql.Open("sqlite3", dbPath)
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
