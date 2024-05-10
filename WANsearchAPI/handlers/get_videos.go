package handlers

import (
	"WANsearchAPI/utils"
	"regexp"
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

	match, _ := regexp.MatchString(`^["“][^"]*["”]$`, query)

	return match

}
