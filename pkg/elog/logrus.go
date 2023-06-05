package elog

import (
	"github.com/sirupsen/logrus"
	"os"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		DisableHTMLEscape: true,
		PrettyPrint: false,
	})
	logrus.SetOutput(os.Stdout)
	logrus.Infof("log module initialization is complete")
}
