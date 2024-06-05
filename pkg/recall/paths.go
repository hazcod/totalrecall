package recall

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

func getRecallPaths(logger *logrus.Logger, username string) (string, string, error) {
	escapedUsername := filepath.Clean(username)
	basePath := fmt.Sprintf("C:\\Users\\%s\\AppData\\Local\\CoreAIPlatform.00\\UKP", escapedUsername)

	logger.WithField("basepath", basePath).Debug("finding Recall GUID folder")
	guidFolder, err := findGuidFolder(basePath)
	if err != nil {
		// RecallNotEnabledErr
		return "", "", err
	}

	logger.Debugf("📁 Recall folder found: %s\n", guidFolder)

	dbPath := filepath.Join(guidFolder, "ukg.db")
	imageStorePath := filepath.Join(guidFolder, "ImageStore")

	logger.Debugf("checking image store path: %s\n", imageStorePath)

	if !fileExists(dbPath) || !fileExists(imageStorePath) {
		return "", "", NotEnabledError
	}

	logger.WithField("db_path", dbPath).WithField("image_path", imageStorePath).
		Debug("Recall feature found enabled")

	return dbPath, imageStorePath, nil
}
