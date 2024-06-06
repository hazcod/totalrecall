package recall

import (
	"database/sql"
	"fmt"
)

type TextGrab struct {
	ProgramName string
	Contents    string
}

func (r *Recall) GrabText() ([]TextGrab, error) {
	conn, err := sql.Open("sqlite3", r.dbPath)
	if err != nil {
		return nil, fmt.Errorf("could not open database connection: %w", err)
	}
	defer conn.Close()

	query := `SELECT c1, c2 AS ProgramName AS URL FROM WindowCaptureTextIndex_content;`
	rows, err := conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	defer rows.Close()

	var results []TextGrab

	for rows.Next() {
		var windowName string
		var contents string

		err := rows.Scan(&windowName, &contents)
		if err != nil {
			return nil, fmt.Errorf("could not scan row: %w", err)
		}

		r.logger.WithField("window", windowName).Debug("processing Recall OCR text")

		results = append(results, TextGrab{
			ProgramName: windowName,
			Contents:    contents,
		})
	}

	return results, nil
}
