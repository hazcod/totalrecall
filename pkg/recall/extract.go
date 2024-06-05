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

func getCurrentDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %w", err)
	}

	return dir, nil
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

func (r *Recall) ExtractImagesForCurrentUser() ([]ExtractResult, error) {
	r.logger.Debug("retrieving current user")
	username, err := getUserName()
	if err != nil {
		return nil, err
	}

	return r.ExtractImages(username)
}

func (r *Recall) GetRecallPathsForCurrentUser() (string, string, error) {
	r.logger.Debug("retrieving current user")
	username, err := getUserName()
	if err != nil {
		return "", "", err
	}

	return r.GetRecallPaths(username)
}

func (r *Recall) GetRecallPaths(username string) (string, string, error) {
	basePath := fmt.Sprintf("C:\\Users\\%s\\AppData\\Local\\CoreAIPlatform.00\\UKP", username)

	r.logger.WithField("basepath", basePath).Debug("finding Recall GUID folder")
	guidFolder, err := findGuidFolder(basePath)
	if err != nil {
		// RecallNotEnabledErr
		return "", "", err
	}

	r.logger.Infof("📁 Recall folder found: %s\n", guidFolder)

	dbPath := filepath.Join(guidFolder, "ukg.db")
	imageStorePath := filepath.Join(guidFolder, "ImageStore")

	r.logger.Debugf("checking image store path: %s\n", imageStorePath)

	if !fileExists(dbPath) || !fileExists(imageStorePath) {
		return "", "", NotEnabledError
	}

	return dbPath, imageStorePath, nil
}

func (r *Recall) IsRecallEnabled(username string) (bool, error) {
	if username == "" {
		r.logger.Debug("retrieving current user")
		var err error
		username, err = getUserName()
		if err != nil {
			return false, err
		}
	}

	_, _, err := r.GetRecallPaths(username)
	return errors.Is(err, NotEnabledError), nil
}

func (r *Recall) ExtractImages(userName string) ([]ExtractResult, error) {
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
