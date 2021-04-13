package config

import (
	"log"
	"os"
)

// Config holds data required to start the application
type Config struct {
	SecretKey []byte
	Realm     string
	Port      string
	CertFile  string
	KeyFile   string
	Smtp      SmtpData
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
	certFile := os.Getenv("CERT_FILE")
	if certFile == "" {
		//log.Fatal("No certification file provided")
	}
	keyFile := os.Getenv("KEY_FILE")
	if keyFile == "" {
		//log.Fatal("No key file provided")
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
		SecretKey: []byte("SuperSecretKey"),
		Realm:     "localhost",
		Port:      "8080",
		CertFile:  certFile,
		KeyFile:   keyFile,
		Smtp: SmtpData{
			Mail:     smtpMail,
			Password: smtpPassword,
			Host:     smtpHost,
			Port:     smtpPort,
		},
	}
}
