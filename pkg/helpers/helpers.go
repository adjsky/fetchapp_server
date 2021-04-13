package helpers

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/smtp"
	"os"
	"path/filepath"
	"server/config"

	"github.com/dchest/uniuri"
)

// SendEmail sends an email by using smtp protocol
func SendEmail(smtpData *config.SmtpData, to []string, message []byte) error {
	auth := smtp.PlainAuth("", smtpData.Mail, smtpData.Password, smtpData.Host)
	err := smtp.SendMail(smtpData.Host+":"+smtpData.Port, auth, smtpData.Mail, to, message)
	return err
}

// ParseBodyPartToJson parses a given multipart and unmarshalls its content
func ParseBodyPartToJson(part *multipart.Part, v interface{}) error {
	metadataBody, err := io.ReadAll(part)
	if err != nil {
		return err
	}
	err = json.Unmarshal(metadataBody, v)
	if err != nil {
		return err
	}
	return nil
}

// SaveToFile saves data to a given path
func SaveToFile(path string, data []byte) (filename string) {
	filename = uniuri.NewLen(32)
	_ = os.WriteFile(filepath.Join(path, filename), data, 0770)
	return
}
