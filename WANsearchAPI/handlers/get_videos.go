package handlers

import (
	"WANsearchAPI/utils"
	"fmt"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetVideos(c *gin.Context) {

	fmt.Println("hh")

	var response utils.Response

	query := c.Query("q")

	query = strings.TrimSuffix(query, " ")
	query = strings.TrimPrefix(query, " ")

	if isQuotedSearch(query) {

		response = utils.QuotedVideos(query)
	} else {
		response = utils.Videos(query)
	}

	c.JSON(200, response)
}

func isQuotedSearch(query string) bool {

	query = strings.ToLower(strings.TrimSuffix(query, " "))

	match, _ := regexp.MatchString(`^["“][^"]*["”]$`, query)

	return match

}
