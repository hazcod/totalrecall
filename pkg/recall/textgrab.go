package recall

import (
	"fmt"
	"zombiezen.com/go/sqlite"
)

type TextGrab struct {
	ProgramName string
	Contents    string
}

func (r *Recall) GrabText() ([]TextGrab, error) {
	conn, err := sqlite.OpenConn(r.dbPath, sqlite.OpenReadOnly)
	if err != nil {
		return nil, fmt.Errorf("could not open database connection: %w", err)
	}
	defer conn.Close()

	query := `SELECT c1, c2 AS ProgramName AS URL FROM WindowCaptureTextIndex_content;`
	stmt, err := conn.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}

	var results []TextGrab

	for {
		hasRow, err := stmt.Step()
		if err != nil {
			return nil, fmt.Errorf("could not execute query: %w", err)
		}
		if !hasRow {
			break
		}

		windowName := stmt.GetText("c1")
		contents := stmt.GetText("c2")

		r.logger.WithField("window", windowName).Debug("processing Recall OCR text")

		results = append(results, TextGrab{
			ProgramName: windowName,
			Contents:    contents,
		})
	}

	return results, nil
}
