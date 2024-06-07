package recall

import (
	"fmt"
	"zombiezen.com/go/sqlite"
)

type AppResult struct {
	Name string
	Path string
}

func (r *Recall) ExtractApps() ([]AppResult, error) {
	conn, err := sqlite.OpenConn(r.dbPath, sqlite.OpenReadOnly)
	if err != nil {
		return nil, fmt.Errorf("could not open database connection: %w", err)
	}
	defer conn.Close()

	query := `SELECT Name, Path AS URL FROM App;`
	stmt, err := conn.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}

	var results []AppResult

	for {
		hasRow, err := stmt.Step()
		if err != nil {
			return nil, fmt.Errorf("could not execute query: %w", err)
		}
		if !hasRow {
			break
		}

		appName := stmt.GetText("Name")
		appPath := stmt.GetText("Path")

		r.logger.WithField("app", appName).Debug("processing Recall web")

		results = append(results, AppResult{
			Name: appName,
			Path: appPath,
		})
	}

	return results, nil
}
