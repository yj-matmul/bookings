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

func CustomLogger() *log.Logger {
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	logger := log.New(multiWriter, "[INFO] ", log.LstdFlags|log.Lshortfile)
	// Logger = log.New(logFile, "[INFO] ", log.LstdFlags|log.Lshortfile)
	logger.Print("End of Program")

	return logger
}

func TestWrite(msg string) {
	file, err := os.OpenFile("test.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString(msg)
}
