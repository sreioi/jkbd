package log

import (
	"github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
)

var Log *logrus.Logger
var once sync.Once

func init() {
	once.Do(func() {
		Log = logrus.New()
		file, err := os.OpenFile("jkbd.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			Log.Out = file
		} else {
			Log.Info("Failed to log to file, using default stderr")
		}

		Log.SetFormatter(&logrus.TextFormatter{
			DisableColors:   true,
			TimestampFormat: time.DateTime,
			FullTimestamp:   true,
		})
	})
}
