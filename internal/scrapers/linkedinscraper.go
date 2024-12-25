package scrapers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/nneji123/ecommerce-golang/config"
)

// Define a model for scraped data
type ScrapeData struct {
	ID                 uint   `gorm:"primaryKey"`
	LinkedinProfileURL string `json:"linkedin_profile_url"`
	Extra              string `json:"extra"`
}

// scrapeLinkedInProfile is the handler for scraping LinkedIn profile
func scrapeLinkedInProfile(c echo.Context) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err)
	}
	// Get the LinkedIn profile URL from the request
	linkedinURL := c.QueryParam("linkedinurl")

	// Prepare the request to the ProxyCrawl person profile endpoint
	proxyCrawlURL := "https://nubela.co/proxycurl/api/v2/linkedin"
	reqData := url.Values{}
	reqData.Set("linkedin_profile_url", linkedinURL)
	reqData.Set("extra", "include")
	reqData.Set("github_profile_id", "include")
	reqData.Set("facebook_profile_id", "include")
	reqData.Set("twitter_profile_id", "include")
	reqData.Set("personal_contact_number", "include")
	reqData.Set("personal_email", "include")
	reqData.Set("use_cache", "if-present")
	reqData.Set("fallback_to_cache", "on-error")

	// Create HTTP client
	client := http.Client{}

	// Prepare the request
	req, err := http.NewRequest("GET", proxyCrawlURL, nil)
	if err != nil {
		return err
	}

	// Set authorization header using API key from config
	req.Header.Set("Authorization", "Bearer "+cfg.ProxyCurlAPIKey)

	// Encode request data and append it to the URL query
	req.URL.RawQuery = reqData.Encode()

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Unmarshal the response into a ScrapeData struct
	var scrapeData ScrapeData
	if err := json.Unmarshal(responseBody, &scrapeData); err != nil {
		return err
	}

	// Return success response
	return c.String(http.StatusOK, "Scraped data saved to database successfully")
}
