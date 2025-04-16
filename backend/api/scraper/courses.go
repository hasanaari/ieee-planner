// backend/scraper/northwestern.go
package scraper

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
)

// CourseData holds the scraped course information for Northwestern.
type CourseData struct {
	School string   `json:"school"`
	Links  []string `json:"links"`
}

// ScrapeNorthwesternCourses scrapes course data for a given quarter from Northwestern.
func ScrapeNorthwesternCourses(quarter string) ([]CourseData, error) {
	baseURL := "https://class-descriptions.northwestern.edu/"
	startURL := baseURL + quarter

	c := colly.NewCollector()

	schools := make([]string, 0)
	classes := make(map[string][]string)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Ensure that we're processing Northwestern-specific links
		if !contains(schools, link) {
			schools = append(schools, link)
			if strings.Contains(link, quarter) {
				classes[link] = []string{}
				// Visit links relative to Northwestern's course page.
				e.Request.Visit(e.Request.AbsoluteURL(link))
			}
		}
		// Attach the link to a school if it matches.
		for _, school := range schools {
			if strings.Contains(link, school) {
				classes[school] = append(classes[school], link)
			}
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error: %v\n", err)
	})

	if err := c.Visit(startURL); err != nil {
		return nil, err
	}

	var results []CourseData
	for school, links := range classes {
		results = append(results, CourseData{School: school, Links: links})
	}
	return results, nil
}

// contains checks if a slice contains a specific string.
func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
