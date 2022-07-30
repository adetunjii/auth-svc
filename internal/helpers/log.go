package helpers

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

func InitializeLogDir() {
	logDir := os.Getenv("LOG_DIR")
	logFile := os.Getenv("LOG_FILE")
	if logDir == "" {
		logDir = "logs"
	}

	_ = os.Mkdir(logDir, os.ModePerm)
	f, err := os.OpenFile(logDir+logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalf("error opening file:%v", err)
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetFlags(0)
	log.SetOutput(f)
}

func LogEvent(level string, message interface{}) {
	data, err := json.Marshal(struct {
		TimeStamp string      `json:"time_stamp"`
		Level     string      `json:"level"`
		Message   interface{} `json:"message"`
	}{
		TimeStamp: time.Now().Format(time.RFC3339),
		Level:     level,
		Message:   message,
	})

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s\n", data)
}
