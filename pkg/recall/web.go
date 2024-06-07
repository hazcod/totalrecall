package recall

import (
	"fmt"
	"time"
	"zombiezen.com/go/sqlite"
)

type WebResult struct {
	Timestamp time.Time
	Domain    string
	URL       string
}

func (r *Recall) ExtractWeb() ([]WebResult, error) {
	conn, err := sqlite.OpenConn(r.dbPath, sqlite.OpenReadOnly)
	if err != nil {
		return nil, fmt.Errorf("could not open database connection: %w", err)
	}
	defer conn.Close()

	query := `SELECT Timestamp, Domain, Uri AS URL FROM Web;`
	stmt, err := conn.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}

	var results []WebResult

	for {
		hasRow, err := stmt.Step()
		if err != nil {
			return nil, fmt.Errorf("could not execute query: %w", err)
		}
		if !hasRow {
			break
		}

		timestamp := stmt.GetInt64("Timestamp")
		domain := stmt.GetText("Domain")
		url := stmt.GetText("Uri")

		r.logger.WithField("domain", domain).Debug("processing Recall web")

		results = append(results, WebResult{
			Timestamp: time.Unix(timestamp/1000, 0),
			Domain:    domain,
			URL:       url,
		})
	}

	return results, nil
}
