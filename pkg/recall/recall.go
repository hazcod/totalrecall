package recall

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type Recall struct {
	logger    *logrus.Logger
	dbPath    string
	imagePath string
}

func New(logger *logrus.Logger, userName string) (*Recall, error) {
	if logger == nil {
		logger = logrus.New()
	}

	if userName == "" {
		var err error
		userName, err = getUserName()
		if err != nil {
			return nil, fmt.Errorf("could not extract current username: %w", err)
		}
	}

	dbPath, imagePath, err := getRecallPaths(logger, userName)
	if err != nil {
		return nil, err
	}

	if dbPath == "" || imagePath == "" {
		return nil, fmt.Errorf("could not extract current recall paths, empty")
	}

	recall := Recall{
		logger:    logger,
		dbPath:    dbPath,
		imagePath: imagePath,
	}

	return &recall, nil
}
