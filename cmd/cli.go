package main

import (
	"errors"
	"flag"
	"github.com/hazcod/totalrecall/pkg/recall"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	defaultLogLevel = logrus.InfoLevel
)

func main() {
	logger := logrus.New()
	logger.SetLevel(defaultLogLevel)

	logLevel := flag.String("log", defaultLogLevel.String(), "The log level to use.")
	username := flag.String("username", "", "The username to find Recall with.")
	flag.Parse()

	logrusLevel, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		logger.WithError(err).Error("invalid log level provided")
		logrusLevel = logrus.InfoLevel
	}
	logger.SetLevel(logrusLevel)

	// ---

	recallPkg, err := recall.New(logger, *username) // current user
	if err != nil {
		logger.WithError(err).Fatal("could not create init recall")
	}

	extracts, err := recallPkg.ExtractImages()
	if errors.Is(err, recall.NotEnabledError) {
		logger.Info("Recall is not enabled on this machine")
		os.Exit(2)
	}

	if err != nil {
		logger.WithError(err).Fatal("could not extract Recall Images")
	}

	for i, extract := range extracts {
		logger.Info("%d - %s - %s - %s", i+1, extract.Timestamp, extract.WindowTitle, extract.WindowToken)
	}

	logger.WithField("total", len(extracts)).Info("extracted all Recall images")
}
