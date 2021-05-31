package config

import (
	"errors"
	"os"
	"strconv"
)

// Config holds data required to start the application
type Config struct {
	SecretKey        []byte
	Port             string
	DatabaseURL      string
	CertFile         string
	KeyFile          string
	PythonScriptPath string
	FrontendPath     string
	TempDir          string
	SMTP             SMTPData
	OnlyAPI          bool
}

// SMTPData struct provides data required to send emails
type SMTPData struct {
	Mail     string
	Password string
	Host     string
	Port     string
}

// Get config instance filled with the required data to start the application
func Get() (*Config, error) {
	secret := os.Getenv("SECRET")
	if secret == "" {
		return nil, errors.New("no secret provided")
	}
	port := os.Getenv("PORT")
	if port == "" {
		return nil, errors.New("no port provided")
	}
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("no database url provided")
	}
	certFile := os.Getenv("CERT_FILE")
	if certFile == "" {
		// return nil, errors.New("no certification file provided")
	}
	keyFile := os.Getenv("KEY_FILE")
	if keyFile == "" {
		// return nil, errors.New("no key file provided")
	}
	pythonScriptPath := os.Getenv("PYTHON_SCRIPT_PATH")
	if pythonScriptPath == "" {
		return nil, errors.New("no python script path provided")
	}
	frontendPath := os.Getenv("FRONTEND_PATH")
	if frontendPath == "" {
		return nil, errors.New("no frontend path provided")
	}
	tempDir := os.Getenv("TEMP_DIR_PATH")
	if tempDir == "" {
		return nil, errors.New("no temporary dir path provided")
	}
	_ = os.MkdirAll(tempDir, 0770) // create if not exists
	smtpMail := os.Getenv("SMTP_MAIL")
	if smtpMail == "" {
		return nil, errors.New("no smtp mail provided")
	}
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	if smtpPassword == "" {
		return nil, errors.New("no smtp password provided")
	}
	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		return nil, errors.New("no smtp host provided")
	}
	smtpPort := os.Getenv("SMTP_PORT")
	if smtpPort == "" {
		return nil, errors.New("no smtp port provided")
	}
	onlyAPI, _ := strconv.ParseBool(os.Getenv("ONLY_API"))

	return &Config{
		SecretKey:        []byte(secret),
		Port:             port,
		DatabaseURL:      databaseURL,
		CertFile:         certFile,
		KeyFile:          keyFile,
		PythonScriptPath: pythonScriptPath,
		FrontendPath:     frontendPath,
		TempDir:          tempDir,
		SMTP: SMTPData{
			Mail:     smtpMail,
			Password: smtpPassword,
			Host:     smtpHost,
			Port:     smtpPort,
		},
		OnlyAPI: onlyAPI,
	}, nil
}
