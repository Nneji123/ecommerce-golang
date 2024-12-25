package emailscraper

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lawzava/emailscraper"
)

type scrapeRequest struct {
	URLs []string `json:"urls"`
}

type scrapeResponse struct {
	Emails map[string][]string `json:"emails"`
	Error  string              `json:"error,omitempty"`
}

func ScrapeEmails(c echo.Context) error {
	req := new(scrapeRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, scrapeResponse{Error: err.Error()})
	}

	s := emailscraper.New(emailscraper.DefaultConfig())
	emailData := make(map[string][]string)

	for _, url := range req.URLs {
		extractedEmails, err := s.Scrape(url)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, scrapeResponse{Error: err.Error()})
		}
		emailData[url] = extractedEmails
	}

	return c.JSON(http.StatusOK, scrapeResponse{Emails: emailData})
}
