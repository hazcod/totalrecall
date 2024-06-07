package main

import (
	"errors"
	"flag"
	"github.com/hazcod/totalrecall/pkg/recall"
	"github.com/sirupsen/logrus"
	"os"
	"time"
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

	if *username == "" {
		usr, err := recall.GetUserName()
		if err != nil {
			logger.WithError(err).Fatal("failed to detect current user, please hardcode with -username")
		}
		username = &usr
	}

	hasRecall, err := recall.IsRecallEnabled(nil, *username)
	if err != nil {
		logger.WithError(err).Error("could not determine if Recall is enabled")
	}

	if !hasRecall {
		logger.Fatalf("user %s does not have recall enabled.", *username)
	}

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
		logger.Infof("%d - %s - %s - %s", i+1, extract.Timestamp.Format(time.DateTime), extract.WindowTitle, extract.WindowToken)
	}

	logger.WithField("total", len(extracts)).Info("extracted all Recall images")
}
