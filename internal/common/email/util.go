package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"path/filepath"

	"github.com/Boostport/mjml-go"
	"jaytaylor.com/html2text"
)

func RenderTemplate(templatePath string, data map[string]interface{}) (string, string, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse template: %w", err)
	}

	var mjmlBuf bytes.Buffer
	if err := tmpl.Execute(&mjmlBuf, data); err != nil {
		return "", "", fmt.Errorf("failed to execute template: %w", err)
	}

	ctx := context.Background()

	html, err := mjml.ToHTML(ctx, mjmlBuf.String(), mjml.WithMinify(true))
	if err != nil {
		return "", "", fmt.Errorf("failed to convert MJML to HTML: %w", err)
	}

	plainText := StripHTML(html)
	return html, plainText, nil
}

func StripHTML(htmlContent string) string {
	text, err := html2text.FromString(htmlContent, html2text.Options{
		PrettyTables: false,
	})
	if err != nil {
		log.Printf("Error stripping HTML: %v", err)
		return htmlContent
	}
	return text
}

func GetMimeType(filename string) string {
	ext := filepath.Ext(filename)
	switch ext {
	case ".pdf":
		return "application/pdf"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	default:
		return "application/octet-stream"
	}
}
