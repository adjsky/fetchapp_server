package helpers

import (
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/adjsky/fetchapp_server/config"
	"github.com/dchest/uniuri"
)

// SendEmail sends an email via the smtp protocol
func SendEmail(smtpData *config.SMTPData, to []string, message []byte) error {
	auth := smtp.PlainAuth("", smtpData.Mail, smtpData.Password, smtpData.Host)
	err := smtp.SendMail(smtpData.Host+":"+smtpData.Port, auth, smtpData.Mail, to, message)
	return err
}

// ParseBodyPartToJSON parses a given multipart and unmarshalls its content
func ParseBodyPartToJSON(part *multipart.Part, v interface{}) error {
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

// GetBoundary returns a boundary provided in a Content-Type header or an empty string if there's no boundary
func GetBoundary(header string) string {
	contentType, params, err := mime.ParseMediaType(header)
	if err != nil {
		return ""
	}
	if strings.HasPrefix(contentType, "multipart") {
		boundary, ok := params["boundary"]
		if !ok {
			return ""
		}
		return boundary
	}
	return ""
}

// ResponseInvalidBody responses with 400 status code
func ResponseInvalidBody(c *gin.Context) {
	code := http.StatusBadRequest
	c.JSON(code, gin.H{
		"code":    code,
		"message": "invalid request body",
	})
}
