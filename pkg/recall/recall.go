package recall

import "github.com/sirupsen/logrus"

type Recall struct {
	logger *logrus.Logger
}

func New(logger *logrus.Logger) (*Recall, error) {
	if logger == nil {
		logger = logrus.New()
	}

	recall := Recall{
		logger: logger,
	}

	return &recall, nil
}
