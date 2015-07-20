package docker

import "github.com/Sirupsen/logrus"

var (
	log = logrus.WithFields(logrus.Fields{
		"package": "docker",
	})
)
