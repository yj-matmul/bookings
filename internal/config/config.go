package config

import (
	"html/template"
	"io"
	"log"
	"os"

	"github.com/alexedwards/scs/v2"
	"github.com/yj-matmul/bookings/internal/models"
)

// AppConfig holds the application configuration
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	MailChan      chan models.MailData
}

// CustomLogger wirtes log to txt file and os standard out
func CustomLogger() (*log.Logger, *os.File) {
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	logger := log.New(multiWriter, "[INFO] ", log.LstdFlags|log.Lshortfile)

	return logger, logFile
}
