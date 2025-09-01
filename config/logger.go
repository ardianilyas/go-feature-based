package config

import (
	"encoding/json"
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

type OrderedJSONFormatter struct {}

func (f *OrderedJSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	log := map[string]interface{}{
		"level": entry.Level.String(),
		"message": entry.Message,
		"data": entry.Data,
		"time": entry.Time.Format("2006-01-02 15:04:05"),
	}

	b, err := json.Marshal(log)
	if err != nil {
		return nil, err
	}
	return append(b, '\n'), nil
}

func InitLogger() {
	Log = logrus.New()

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Log.Out = file
	} else {
		Log.Out = os.Stdout
	}

	Log.SetFormatter(&OrderedJSONFormatter{})
	Log.SetLevel(logrus.InfoLevel)
}