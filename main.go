package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

type Request struct {
	Site       []string
	SearchText string
}

type Response struct {
	FoundAtSite string
}

func main() {
	r := gin.Default()
	var request Request
	r.POST("/checkText", func(c *gin.Context) {

		if c.BindJSON(&request) == nil {
			if len(request.Site) >= 0 && request.SearchText != "" {
				for index, site := range request.Site {
					resp, err := http.Get(site)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err})
					} else {
						theBody, _ := ioutil.ReadAll(resp.Body)
						resp.Body.Close()
						regexStringCreation := fmt.Sprintf(".*%s.*", request.SearchText)
						theRegEx := regexp.MustCompile(regexStringCreation)
						s := string(theBody[:])
						if !theRegEx.MatchString(s) {
							request.Site = request.Site[:index+copy(request.Site[index:], request.Site[index+1:])]
						}
					}
				}
				var response Response
				for _, finishUp := range request.Site {
					if len(request.Site) > 1 {
						response.FoundAtSite += fmt.Sprintf("%s,", finishUp)
					} else {
						response.FoundAtSite = finishUp
					}
				}
				c.JSON(http.StatusOK, response)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"Site": "[]string", "SearchText": "string"})
			}
		} else {
			fmt.Printf("Posted body is not usable\n")
			c.JSON(http.StatusBadRequest, gin.H{"Site": "[]string", "SearchText": "string"})
		}

	})
	r.Run()
}
