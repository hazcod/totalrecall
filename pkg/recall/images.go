package recall

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"zombiezen.com/go/sqlite"
)

var (
	NotEnabledError = errors.New("recall not enabled")
)

type ExtractResult struct {
	WindowTitle string
	WindowToken string
	Timestamp   time.Time
	Screenshot  []byte
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
	conn, err := sqlite.OpenConn(r.dbPath, sqlite.OpenReadOnly)
	if err != nil {
		return nil, fmt.Errorf("could not open database connection: %w", err)
	}
	defer conn.Close()

	query := `SELECT WindowTitle, TimeStamp, ImageToken FROM WindowCapture WHERE (WindowTitle IS NOT NULL OR ImageToken IS NOT NULL);`
	stmt, err := conn.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}

	var results []ExtractResult

	for {
		hasRow, err := stmt.Step()
		if err != nil {
			return nil, fmt.Errorf("could not execute query: %w", err)
		}
		if !hasRow {
			break
		}

		windowTitle := stmt.GetText("WindowTitle")
		timestamp := stmt.GetInt64("TimeStamp")
		imageToken := stmt.GetText("ImageToken")

		readableTimestamp := time.Unix(0, timestamp*int64(time.Millisecond)).Format("2006-01-02 15:04:05")
		r.logger.WithField("timestamp", readableTimestamp).WithField("window", windowTitle).
			Debug("processing Recall image")

		var imageBytes []byte

		if imageToken == "" {
			r.logger.WithField("timestamp", readableTimestamp).WithField("window", windowTitle).
				Debug("image token is empty")
		} else {
			screenshotPath := filepath.Join(r.imagePath, imageToken)

			imageBytes, err = os.ReadFile(screenshotPath)
			if err != nil {
				r.logger.WithError(err).WithField("path", screenshotPath).Fatal("could not load image")
			}

			r.logger.WithField("size", len(imageBytes)).WithField("program", windowTitle).
				Debug("loaded screenshot")
		}

		results = append(results, ExtractResult{
			WindowTitle: windowTitle,
			WindowToken: imageToken,
			Screenshot:  imageBytes,
			Timestamp:   time.Unix(timestamp/1000, 0),
		})
	}

	return results, nil
}
