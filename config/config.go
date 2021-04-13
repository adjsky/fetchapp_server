package config

import (
	"log"
	"os"

	"github.com/dchest/uniuri"
)

// Config holds data required to start the application
type Config struct {
	SecretKey        []byte
	Port             string
	DatabaseUrl      string
	CertFile         string
	KeyFile          string
	PythonScriptPath string
	TempDir          string
	Smtp             SmtpData
}

// SmtpData struct provides data required to send emails
type SmtpData struct {
	Mail     string
	Password string
	Host     string
	Port     string
}

// Get config instance filled with the required data to start the application
func Get() *Config {
	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatal("No secret provided")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("No port provided")
	}
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		log.Fatal("No database url provided")
	}
	certFile := os.Getenv("CERT_FILE")
	if certFile == "" {
		//log.Fatal("No certification file provided")
	}
	keyFile := os.Getenv("KEY_FILE")
	if keyFile == "" {
		//log.Fatal("No key file provided")
	}
	pythonScriptPath := os.Getenv("PYTHON_SCRIPT_PATH")
	if pythonScriptPath == "" {
		log.Fatal("No python script path provided")
	}
	tempDir, err := os.MkdirTemp("", uniuri.New())
	if err != nil {
		log.Fatal(err)
	}
	smtpMail := os.Getenv("SMTP_MAIL")
	if smtpMail == "" {
		log.Fatal("No smtp mail provided")
	}
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	if smtpPassword == "" {
		log.Fatal("No smtp password provided")
	}
	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		log.Fatal("No smtp host provided")
	}
	smtpPort := os.Getenv("SMTP_PORT")
	if smtpPort == "" {
		log.Fatal("No smtp port provided")
	}

	return &Config{
		SecretKey:        []byte(secret),
		Port:             port,
		DatabaseUrl:      databaseUrl,
		CertFile:         certFile,
		KeyFile:          keyFile,
		PythonScriptPath: pythonScriptPath,
		TempDir:          tempDir,
		Smtp: SmtpData{
			Mail:     smtpMail,
			Password: smtpPassword,
			Host:     smtpHost,
			Port:     smtpPort,
		},
	}
}
