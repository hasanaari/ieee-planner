// backend/scraper/colly_scraper.go
package scraper

import (
    "fmt"
    "strings"
    "time"

    "github.com/gocolly/colly/v2"
    "github.com/nynniaw12/ieee-planner/api/models"
)

var (
    CachedCourses []models.Course
    ScrapeError   error
    onRequest func(*colly.Request)
)

func SetLogger(fn func(*colly.Request)) {
    onRequest = fn
}

const baseURL = "https://class-descriptions.northwestern.edu"

func ScrapeCourses(quarter string) ([]models.Course, error) {
    c := colly.NewCollector(
        colly.AllowedDomains("class-descriptions.northwestern.edu"),
        colly.Async(true),
        colly.MaxDepth(4),
    )

    // throttle & parallelism: at most 4 concurrent requests
    c.Limit(&colly.LimitRule{
        DomainGlob:  "*.northwestern.edu",
        Parallelism: 4,
        Delay:       200 * time.Millisecond,
    })

    if onRequest != nil {
        c.OnRequest(onRequest)
    }

    // 1) Log every request
    c.OnRequest(func(r *colly.Request) {
        fmt.Println("→ REQUEST:", r.URL)
    })

    // 2) Log every response
    c.OnResponse(func(r *colly.Response) {
        fmt.Println("← RESPONSE:", r.Request.URL, "Status:", r.StatusCode)
    })

    // 3) Log any errors
    c.OnError(func(r *colly.Response, err error) {
        fmt.Printf("!! ERROR fetching %s: %v (status %d)\n",
            r.Request.URL, err, r.StatusCode)
    })

    var courses []models.Course


    // 1) Quarter page: follow each school link
    c.OnHTML("#main-content ul li a[href]", func(e *colly.HTMLElement) {
        href := e.Attr("href")       // e.g. "4980", "4980/WCAS", "4980/WCAS/CHEM", or "4980/WCAS/CHEM/35064"
        parts := strings.Split(href, "/")
    
        switch len(parts) {
        case 1:
            // on /4980, this a school link
            e.Request.Visit(e.Request.AbsoluteURL(href))
        case 2:
            // on /4980/SCHOOL, this is a subject link
            e.Request.Visit(e.Request.AbsoluteURL(href))
        case 3:
            // on /4980/SCHOOL/SUBJ, this is a section link
            e.Request.Visit(e.Request.AbsoluteURL(href))
        // case 4:
            // could even handle detail pages here instead of using OnHTML(".expander…")
        default:
            // ignore everything else (footer, header, breadcrumbs, etc.)
            return
        }
    })
    

    // 2) School page: follow each subject link
    c.OnHTML("#main-content ul li a[href]", func(e *colly.HTMLElement) {
        href := e.Attr("href")
        // links like "4980/WCAS/CHEM"
        if strings.Count(href, "/") == 2 {
            e.Request.Visit(e.Request.AbsoluteURL(href))
        }
    })

    // 3) Subject page: follow each course‐section link
    c.OnHTML(".expander ul li a[href]", func(e *colly.HTMLElement) {
        // this is the "/4980/WCAS/CHEM/35064" link
        e.Request.Visit(e.Request.AbsoluteURL(e.Attr("href")))
    })

    // 4) Detail page: scrape title, desc, prereqs
    c.OnHTML("main .content", func(e *colly.HTMLElement) {
        title := strings.TrimSpace(e.ChildText("h1"))
        desc := strings.TrimSpace(e.ChildText(".overview-of-class"))
        preText := e.ChildText("p:contains('Prerequisites')")
        prereqs := []string{}
        if parts := strings.SplitN(preText, ":", 2); len(parts) == 2 {
            for _, p := range strings.Split(parts[1], " and ") {
                prereqs = append(prereqs, strings.TrimSpace(p))
            }
        }

        courses = append(courses, models.Course{
            Department:    "",
            ID:            0,
            Name:          title,
            Description:   desc,
            Prerequisites: prereqs,
        })
        fmt.Println("  + parsed course:", title)
    })

    start := fmt.Sprintf("%s/%s", baseURL, quarter)
    if err := c.Visit(start); err != nil {
        return nil, err
    }

    // block until all queued visits have completed
    c.Wait()
    fmt.Printf("DEBUG: c.Wait() returned, total courses=%d\n", len(courses))
    return courses, nil
}

