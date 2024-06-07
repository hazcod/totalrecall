package recall

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os/user"
	"strings"
)

type Recall struct {
	logger    *logrus.Logger
	dbPath    string
	imagePath string
}

func GetUserName() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("error getting current user: %w", err)
	}

	return strings.Split(usr.Username, "\\")[1], nil
}

func New(logger *logrus.Logger, userName string) (*Recall, error) {
	if logger == nil {
		logger = logrus.New()
	}

	if userName == "" {
		var err error
		userName, err = GetUserName()
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
