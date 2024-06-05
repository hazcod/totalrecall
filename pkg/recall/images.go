package recall

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

var (
	NotEnabledError = errors.New("recall not enabled")
)

type ExtractResult struct {
	WindowTitle string
	WindowToken string
	Timestamp   int64
}

func getUserName() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("error getting current user: %w", err)
	}

	return usr.Username, nil
}

func findGuidFolder(basePath string) (string, error) {
	files, err := os.ReadDir(basePath)
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			return filepath.Join(basePath, file.Name()), nil
		}
	}

	return "", NotEnabledError
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (r *Recall) ExtractImages() ([]ExtractResult, error) {
	conn, err := sql.Open("sqlite3", r.dbPath)
	if err != nil {
		return nil, fmt.Errorf("could not open database connection: %w", err)
	}
	defer conn.Close()

	query := `SELECT WindowTitle, TimeStamp, ImageToken FROM WindowCapture WHERE (WindowTitle IS NOT NULL OR ImageToken IS NOT NULL);`
	rows, err := conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	defer rows.Close()

	var results []ExtractResult

	for rows.Next() {
		var windowTitle string
		var timestamp int64
		var imageToken string

		err := rows.Scan(&windowTitle, &timestamp, &imageToken)
		if err != nil {
			return nil, fmt.Errorf("could not scan row: %w", err)
		}

		readableTimestamp := time.Unix(0, timestamp*int64(time.Millisecond)).Format("2006-01-02 15:04:05")
		r.logger.WithField("timestamp", readableTimestamp).WithField("window", windowTitle).
			Debug("processing Recall image")

		if imageToken == "" {
			r.logger.WithField("timestamp", readableTimestamp).WithField("window", windowTitle).
				Warn("image token is empty")
		}

		results = append(results, ExtractResult{
			WindowTitle: windowTitle,
			WindowToken: imageToken,
			Timestamp:   timestamp,
		})
	}

	return results, nil
}
