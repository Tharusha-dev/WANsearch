package handlers

import (
	"WANsearchAPI/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetVideos(c *gin.Context) {

	var videos []utils.Video

	query := c.Query("q")

	query = strings.TrimSuffix(query, " ")
	query = strings.TrimPrefix(query, " ")

	if isQuotedSearch(query) {
		videos = utils.QuotedVideos(query)
	} else {
		videos = utils.Videos(query)
	}

	c.JSON(200, videos)
}

func isQuotedSearch(query string) bool {

	query = strings.ToLower(strings.TrimSuffix(query, " "))
	characters := strings.Split(query, "")

	if characters[0] == `"` && characters[len(characters)-1] == `"` {
		return true

	} else {
		return false
	}

}
