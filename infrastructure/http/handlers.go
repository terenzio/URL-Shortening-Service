package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/terenzio/URL-Shortening-Service/application"
	urlModel "github.com/terenzio/URL-Shortening-Service/domain"
)

type Handler struct {
	service *application.URLService
}

// NewHandler creates a new instance of Handler
func NewHandler(service *application.URLService) *Handler {
	return &Handler{service: service}
}

// HandleHomePage displays the list of all shortened URLs mapped to their original ones.
// @Summary Displays the list of all shortened URLs mapped to their original ones in JSON format.
// @Description Displays the list of all shortened URLs mapped to their original ones in JSON format.
// @Tags URL
// @Produce json
// @Success 200 {object} urlModel.URLMapping "URL Mappings"
// @Router /url/display [get]
func (h *Handler) HandleHomePage(c *gin.Context) {
	urls, err := h.service.FetchAllURLs(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to fetch URLs: %v", err)
		return
	}

	//var linksListHtml string
	var urlMappings []urlModel.URLMapping
	for _, url := range urls {
		//linksListHtml += fmt.Sprintf("<div>Shortened: <a href=\"/short/%s\">/short/%s</a> Original: %s</div>", url.ShortCode, url.ShortCode, url.OriginalURL)
		urlMappings = append(urlMappings, urlModel.URLMapping{ShortCode: url.ShortCode, OriginalURL: url.OriginalURL, Expiry: url.Expiry})
	}
	//c.Header("Content-Type", "text/html")
	//c.String(http.StatusOK, "<h2>Go URL Shortener</h2><p>Welcome! Here's a list of all shortened URLs:</p>%s", linksListHtml)
	c.IndentedJSON(http.StatusOK, urlMappings)
}

// HandleAddLink creates a shortened link for the given original URL.
// @Summary Creates a shortened link for the given original URL.
// @Description NOTE 1: In the JSON body, the "original_url" should contain proper formatting with either http or https. Example: https://www.google.com.
// @Description NOTE 2: In the JSON body, the "expiry" date is optional, with the default expiration set to 30 days from now. The expiry time can be customized like this example: 2024-04-02T00:00:00Z
// @Tags URL
// @Accept json
// @Param original_url body urlModel.AddURLRequest true "Original URL and Expiry Time (optional)"
// @Produce json
// @Success 200 {object} urlModel.AddSuccessResponse "Shortened URL"
// @Router /url/add [post]
func (h *Handler) HandleAddLink(c *gin.Context) {

	var newUrl = urlModel.URL{}
	if err := c.BindJSON(&newUrl); err != nil || newUrl.OriginalURL == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad request - original_url is required"})
		return
	}
	originalURL := newUrl.OriginalURL
	if !isValidUrl(originalURL) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request - Missing http or https - example: https://www.google.com"})
		return
	}

	// Get the expiry time
	customExpiryTime := newUrl.Expiry
	fmt.Println("customExpiryTime: ", customExpiryTime)

	now := time.Now()
	var defaultExpiryDays time.Duration = 30
	defaultExpiryTime := now.Add(defaultExpiryDays * 24 * time.Hour)
	var adjustedExpiryTime time.Time

	// Check if the expiry time is not set
	// If not, set the expiry time to the defaultExpiryDays
	if customExpiryTime == (time.Time{}) {
		adjustedExpiryTime = defaultExpiryTime
	} else {
		// Check if the futureTime is after the current time
		if customExpiryTime.After(now) {
			fmt.Println("The given time is in the future.")
			adjustedExpiryTime = customExpiryTime
		} else {
			fmt.Println("The given time is not in the future.")
			adjustedExpiryTime = defaultExpiryTime
		}
	}
	fmt.Println("adjustedExpiryTime: ", adjustedExpiryTime)

	shortCode, err := h.service.ShortenURL(c, originalURL, adjustedExpiryTime)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error shortening URL: %v", err)
		return
	}

	shortenedURL := fmt.Sprintf("http://localhost:9000/api/v1/redirect/%s", shortCode)
	c.IndentedJSON(http.StatusOK, gin.H{"shortened_url": shortenedURL})
}

// HandleRedirectToOriginalLink redirects the user to the original URL based on the short code.
// @Summary Redirects the user to the original URL based on the input short code.
// @Description NOTE: Copy the full url including the short code to the browser to be redirected. Do not use the Swagger UI here as it does not support redirection.
// @Tags REDIRECT
// @Param shortcode path string true "Short Code"
// @Produce plain
// @Success 307 {string} string "Redirected to original url - example: http://localhost:9000/api/v1/redirect/2v5ompxD"
// @Failure 400  {string}  string "Parameter missing - enter the short code in the URL path"
// @Failure 404  {string}  string "No original URL exists for the given short code"
// @Router /redirect/{shortcode} [get]
func (h *Handler) HandleRedirectToOriginalLink(c *gin.Context) {
	shortCode := c.Param("shortcode")
	if shortCode == "" {
		c.String(http.StatusBadRequest, "Parameter missing - enter the short code in the URL path")
		return
	}

	originalURL, err := h.service.GetOriginalURL(c, shortCode)
	if err != nil {
		c.String(http.StatusNotFound, "No original URL exists for the given short code: %v", err)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, originalURL)
}
