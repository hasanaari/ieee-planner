// backend/scraper/colly_scraper.go
package scraper

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
	"path"
	"sync"
	"time"
)

const CLASS_DESCRIPTIONS = "https://class-descriptions.northwestern.edu/"

type GenericScrapedObj struct {
	name string
	href string
	url  string
}

type ScrapedObj interface {
	GetName() string
	GetHref() string
	GetURL() string
}

func (q GenericScrapedObj) GetName() string {
	return q.name
}

func (q GenericScrapedObj) GetHref() string {
	return q.href
}

func (q GenericScrapedObj) GetURL() string {
	return q.url
}

func ScrapeGeneric[T any](
	urls []string,
	create func(name, href, url string) T,
	filter func(lowerText, href, url string) bool) map[string][]T {

	itemsByURL := make(map[string][]T)
	var mu sync.Mutex

	for _, url := range urls {
		itemsByURL[url] = []T{}
	}

	c := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(1),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 16, // nproc : 16
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		text := strings.TrimSpace(e.Text)
		lowerText := strings.ToLower(text)
		url := e.Request.URL.String()

		if filter(lowerText, href, url) {
			item := create(text, href, url)

			mu.Lock()
			itemsByURL[url] = append(itemsByURL[url], item)
			mu.Unlock()
		}
	})

	for _, url := range urls {
		c.Visit(url)
	}

	c.Wait()

	return itemsByURL
}

func PrintScrapeResult[T ScrapedObj](itemsByURL map[string][]T) {
	for url, items := range itemsByURL {
		fmt.Printf("URL: %s\n", url)
		for i, item := range items {
			fmt.Printf("  Item %d: %v\n", i, item)
		}
	}
}

// For Quarters
func ScrapeQuarters(urls []string) map[string][]GenericScrapedObj {
	return ScrapeGeneric(
		urls,
		func(name, href, url string) GenericScrapedObj { return GenericScrapedObj{name, href, url} },
		func(lowerText, href, _ string) bool { return startsWithFourDigitsRegex(lowerText) },
	)
}

const (
	WCAS = "WCAS"
	MEAS = "MEAS"
)

var SCHOOLS_WHITELIST = []string{WCAS, MEAS}

// For Schools
func ScrapeSchools(urls []string) map[string][]GenericScrapedObj {
	return ScrapeGeneric(
		urls,
		func(name, href, url string) GenericScrapedObj { return GenericScrapedObj{name, href, url} },
		func(lowerText, href, _ string) bool {
			return strings.Contains(lowerText, "school") || strings.Contains(lowerText, "college")
		},
	)
}

// For Subjects
func ScrapeSubjects(urls []string) map[string][]GenericScrapedObj {
	return ScrapeGeneric(
		urls,
		func(name, href, url string) GenericScrapedObj { return GenericScrapedObj{name, href, url} },
		func(lowerText, href, url string) bool {
			return strings.HasPrefix(href, path.Base(url))
		},
	)
}

// For Sections
func ScrapeSections(urls []string) map[string][]GenericScrapedObj {
	return ScrapeGeneric(
		urls,
		func(name, href, url string) GenericScrapedObj { return GenericScrapedObj{name, href, url} },
		func(lowerText, href, url string) bool {
			return startsWithNumberColon(lowerText)
		},
	)
}

func ScrapeCourseDescriptionHierarchy() {
	// TODO: build a full ass tree

	var nexturls []string
	nexturls = append(nexturls, CLASS_DESCRIPTIONS)
	quartersByURL := ScrapeQuarters(nexturls)

	nexturls = getNextURLSet(quartersByURL, nexturls) // TODO: this has to be a map
	schoolsByURL := ScrapeSchools(nexturls)

	nexturls = getNextURLSet(schoolsByURL, nexturls)
	subjectsByURL := ScrapeSubjects(nexturls)

	nexturls = getNextURLSet(subjectsByURL, nexturls)
	sectionsByURL := ScrapeSections(nexturls)

	nexturls = getNextURLSet(sectionsByURL, nexturls)
	coursesByURL := ScrapeNorthwesternCourses(
		nexturls,
		quartersByURL,
		schoolsByURL,
		subjectsByURL,
		sectionsByURL,
	)

	for _, course := range coursesByURL {
		PrintCourse(course)
		return
	}
}

func getNextURLSet[T ScrapedObj](itemsByURL map[string][]T, urls []string) []string {
	for _, url := range urls {
		nexturls := getUrls(itemsByURL[url])
		var lhs []string
		lhs = append(lhs, url)
		return strCartesian(lhs, nexturls)
	}
	return make([]string, 0)
}

func getUrls[T ScrapedObj](objs []T) []string {
	var res []string
	for _, obj := range objs {
		res = append(res, obj.GetHref())
	}
	return res
}

// get cartesian product of two url arrays
func strCartesian(lhs []string, rhs []string) []string {
	var res []string
	for _, l := range lhs {
		for _, r := range rhs {
			joined := ""

			lSegments := strings.Split(strings.TrimRight(l, "/"), "/")
			lLastSegment := ""
			if len(lSegments) > 0 {
				lLastSegment = lSegments[len(lSegments)-1]
			}

			rSegments := strings.Split(strings.TrimLeft(r, "/"), "/")
			rFirstSegment := ""
			if len(rSegments) > 0 {
				rFirstSegment = rSegments[0]
			}

			if lLastSegment != "" && rFirstSegment != "" && lLastSegment == rFirstSegment {
				if strings.HasSuffix(l, "/") {
					joined = l + strings.Join(rSegments[1:], "/")
				} else {
					joined = l + "/" + strings.Join(rSegments[1:], "/")
				}
			} else {
				if strings.HasSuffix(l, "/") && strings.HasPrefix(r, "/") {
					joined = l + r[1:]
				} else if !strings.HasSuffix(l, "/") && !strings.HasPrefix(r, "/") {
					joined = l + "/" + r
				} else {
					joined = l + r
				}
			}

			res = append(res, joined)
		}
	}
	return res
}

func startsWithFourDigitsRegex(s string) bool {
	re := regexp.MustCompile("^[0-9]{4}")
	return re.MatchString(s)
}

func startsWithNumberColon(s string) bool {
	re := regexp.MustCompile(`^[0-9]+:`)
	return re.MatchString(s)
}

type Instructor struct {
	Name        string
	Phone       string
	Email       string
	OfficeHours string
	Address     string
}

type MeetingTime struct {
	Location  string
	Days      []string
	StartTime time.Time
	EndTime   time.Time
	TimeRange string
}

type Course struct {
	Title  string
	Number string

	Topic        string
	Instructors  []Instructor
	MeetingTimes []MeetingTime
	Overview     string
	URL          string

	Section string
	Subject string
	School  string
	Quarter string
}

// standard days
func ParseDay(day string) string {
	dayMap := map[string]string{
		"mon":   "Monday",
		"tues":  "Tuesday",
		"tue":   "Tuesday",
		"wed":   "Wednesday",
		"thurs": "Thursday",
		"thu":   "Thursday",
		"fri":   "Friday",
		"sa":    "Saturday",
		"su":    "Sunday",
	}

	day = strings.ToLower(strings.TrimSpace(day))
	if standardDay, ok := dayMap[day]; ok {
		return standardDay
	}
	return day
}

// structured times
func ParseTimeRange(timeStr string) (time.Time, time.Time, string) {
	parts := strings.Split(timeStr, "-")
	if len(parts) != 2 {
		return time.Time{}, time.Time{}, timeStr
	}

	startStr := strings.TrimSpace(parts[0])
	endStr := strings.TrimSpace(parts[1])

	refDate := "2006-01-02 " // ignore

	startTime, err := time.Parse("2006-01-02 3:04PM", refDate+startStr)
	if err != nil {
		return time.Time{}, time.Time{}, timeStr
	}

	endTime, err := time.Parse("2006-01-02 3:04PM", refDate+endStr)
	if err != nil {
		return time.Time{}, time.Time{}, timeStr
	}

	return startTime, endTime, timeStr
}

func ParseMeetingInfo(meetingStr string) []MeetingTime {
	var meetings []MeetingTime

	meetingLines := strings.Split(meetingStr, "\n")

	for _, line := range meetingLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		meeting := MeetingTime{}

		if colonIndex := strings.Index(line, ":"); colonIndex != -1 {
			meeting.Location = strings.TrimSpace(line[:colonIndex])

			dayTimeStr := strings.TrimSpace(line[colonIndex+1:])

			for _, day := range []string{"Mon", "Tues", "Wed", "Thurs", "Fri", "Sa", "Su"} {
				dayTimeStr = strings.Replace(dayTimeStr, day+",", day+" ", -1)
				dayTimeStr = strings.Replace(dayTimeStr, day, "|"+day, -1)
			}

			dayTimeParts := strings.Split(dayTimeStr, "|")

			if len(dayTimeParts) >= 2 {
				for i := 1; i < len(dayTimeParts); i++ {
					dayPart := strings.TrimSpace(dayTimeParts[i])
					if spaceIndex := strings.Index(dayPart, " "); spaceIndex != -1 {
						day := dayPart[:spaceIndex]
						meeting.Days = append(meeting.Days, ParseDay(day))
					} else if len(dayPart) > 0 && !strings.Contains(dayPart, ":") {
						meeting.Days = append(meeting.Days, ParseDay(dayPart))
					}
				}

				timeStr := strings.TrimSpace(dayTimeParts[len(dayTimeParts)-1])
				timeRegex := regexp.MustCompile(`(\d+:\d+[AP]M)\s*-\s*(\d+:\d+[AP]M)`)
				matches := timeRegex.FindStringSubmatch(timeStr)

				if len(matches) >= 3 {
					startTimeStr := matches[1]
					endTimeStr := matches[2]

					refDate := "2006-01-02 "
					startTime, err1 := time.Parse("2006-01-02 3:04PM", refDate+startTimeStr)
					endTime, err2 := time.Parse("2006-01-02 3:04PM", refDate+endTimeStr)

					if err1 == nil && err2 == nil {
						meeting.StartTime = startTime
						meeting.EndTime = endTime
						meeting.TimeRange = startTimeStr + " - " + endTimeStr
					} else {
						meeting.TimeRange = timeStr
					}
				} else {
					singleTimeRegex := regexp.MustCompile(`(\d+:\d+[AP]M)`)
					timeMatches := singleTimeRegex.FindStringSubmatch(timeStr)

					if len(timeMatches) >= 2 {
						timeValue := timeMatches[1]
						refDate := "2006-01-02 "
						parsedTime, err := time.Parse("2006-01-02 3:04PM", refDate+timeValue)

						if err == nil {
							meeting.EndTime = parsedTime
							meeting.TimeRange = "End: " + timeValue
						} else {
							meeting.TimeRange = timeStr
						}
					} else {
						meeting.TimeRange = timeStr
					}
				}
			} else {
				meeting.TimeRange = dayTimeStr
			}
		} else {
			meeting.TimeRange = line
		}
		meetings = append(meetings, meeting)
	}

	return meetings
}

// ParseInstructors parses instructor information, where:
// - First line is always the name (if provided)
// - Second line is phone number (if provided)
// - Third line is email (if provided)
// - Fourth line is address (if provided)
// - Fifth line is office hours (if provided)
// - Instructors are separated by empty lines
func ParseInstructor(instructorText string) Instructor {
	instructor := Instructor{}

	// Common patterns for validation
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	phoneRegex := regexp.MustCompile(`(\d{3}[-/\s]?\d{3}[-/\s]?\d{4})`)
	officeHoursRegex := regexp.MustCompile(`(?i)office\s+hours:(.+)`)

	// Split text into lines
	lines := strings.Split(instructorText, "\n")

	// Track which fields we've already found
	foundName := false
	foundPhone := false
	foundEmail := false
	foundAddress := false
	foundOfficeHours := false

	// Process each line
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Try to match patterns in priority order

		// If we haven't found a name yet, assume the first line is the name
		if !foundName {
			instructor.Name = line
			foundName = true
			continue
		}

		// Try to match phone
		if !foundPhone && phoneRegex.MatchString(line) {
			matches := phoneRegex.FindStringSubmatch(line)
			if len(matches) > 1 {
				instructor.Phone = matches[1]
				foundPhone = true
				continue
			}
		}

		// Try to match email
		if !foundEmail && emailRegex.MatchString(line) {
			instructor.Email = emailRegex.FindString(line)
			foundEmail = true
			continue
		}

		// Try to match office hours
		if !foundOfficeHours && officeHoursRegex.MatchString(line) {
			matches := officeHoursRegex.FindStringSubmatch(line)
			if len(matches) > 1 {
				instructor.OfficeHours = strings.TrimSpace(matches[1])
			} else {
				instructor.OfficeHours = line
			}
			foundOfficeHours = true
			continue
		}

		// If we haven't found an address yet and none of the above patterns matched,
		// assume this is the address
		if !foundAddress {
			instructor.Address = line
			foundAddress = true
			continue
		}

		// If none of the above patterns matched, append to office hours or set it
		if instructor.OfficeHours != "" {
			instructor.OfficeHours += " " + line
		} else {
			instructor.OfficeHours = line
		}
	}

	return instructor
}

// scrapes course information from a Northwestern course page
func ScrapeNorthwesternCourses(
	urls []string,
	quartersByURL map[string][]GenericScrapedObj,
	schoolsByURL map[string][]GenericScrapedObj,
	subjectsByURL map[string][]GenericScrapedObj,
	sectionsByURL map[string][]GenericScrapedObj,
) map[string]*Course {
	var mutex sync.Mutex

	c := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(1),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 16, // nproc : 16
	})

	coursesByURL := make(map[string]*Course)

	for _, url := range urls {
		coursesByURL[url] = &Course{
			URL: url,
		}
	}

	c.OnHTML("h1", func(e *colly.HTMLElement) {
		url := e.Request.URL.String()
		fullTitle := e.Text

		mutex.Lock()
		defer mutex.Unlock()

		course, exists := coursesByURL[url]
		if !exists {
			course = &Course{URL: url}
			coursesByURL[url] = course
		}

		parts := strings.Split(fullTitle, "(")
		if len(parts) >= 2 {
			course.Title = strings.TrimSpace(parts[0])
			numPart := strings.TrimSpace(parts[1])
			numPart = strings.TrimSuffix(numPart, ")")
			course.Number = numPart
		} else {
			course.Title = fullTitle
		}

		sectionURL := url
		subjectURL := PopLastURLPart(sectionURL)
		schoolURL := PopLastURLPart(subjectURL)
		quarterURL := PopLastURLPart(schoolURL)

		sectionID := GetLastURLPart(sectionURL)
		subjectID := GetLastURLPart(subjectURL)
		schoolID := GetLastURLPart(schoolURL)
		quarterID := GetLastURLPart(quarterURL)

		course.Section = sectionID
		course.Subject = subjectID
		course.School = schoolID
		course.Quarter = quarterID
	})

	// topic
	c.OnHTML("h2:contains('Topic') + p", func(e *colly.HTMLElement) {
		url := e.Request.URL.String()

		mutex.Lock()
		defer mutex.Unlock()

		course, exists := coursesByURL[url]
		if !exists {
			course = &Course{URL: url}
			coursesByURL[url] = course
		}

		course.Topic = strings.TrimSpace(e.Text)
	})

	// instructors
	c.OnHTML("h2:contains('Instructors') + p", func(e *colly.HTMLElement) {
		url := e.Request.URL.String()

		mutex.Lock()
		defer mutex.Unlock()

		course, exists := coursesByURL[url]
		if !exists {
			course = &Course{URL: url}
			coursesByURL[url] = course
		}

		course.Instructors = append(course.Instructors, ParseInstructor(e.Text))
	})

	// instructors
	c.OnHTML("h2:contains('Instructors') + p + p", func(e *colly.HTMLElement) {
		url := e.Request.URL.String()

		mutex.Lock()
		defer mutex.Unlock()

		course, exists := coursesByURL[url]
		if !exists {
			course = &Course{URL: url}
			coursesByURL[url] = course
		}

		course.Instructors = append(course.Instructors, ParseInstructor(e.Text))
	})

	// instructors
	c.OnHTML("h2:contains('Instructors') + p + p + p", func(e *colly.HTMLElement) {
		url := e.Request.URL.String()

		mutex.Lock()
		defer mutex.Unlock()

		course, exists := coursesByURL[url]
		if !exists {
			course = &Course{URL: url}
			coursesByURL[url] = course
		}

		course.Instructors = append(course.Instructors, ParseInstructor(e.Text))
	})

	// meeting times
	c.OnHTML("h2:contains('Meeting Info') + p", func(e *colly.HTMLElement) {
		url := e.Request.URL.String()
		meetingText := e.Text

		mutex.Lock()
		defer mutex.Unlock()

		course, exists := coursesByURL[url]
		if !exists {
			course = &Course{URL: url}
			coursesByURL[url] = course
		}

		course.MeetingTimes = ParseMeetingInfo(meetingText)
	})

	// overview
	c.OnHTML("h2:contains('Overview of class') + p", func(e *colly.HTMLElement) {
		url := e.Request.URL.String()

		mutex.Lock()
		defer mutex.Unlock()

		course, exists := coursesByURL[url]
		if !exists {
			course = &Course{URL: url}
			coursesByURL[url] = course
		}

		course.Overview = strings.TrimSpace(e.Text)
	})

	for _, url := range urls {
		err := c.Visit(url)
		if err != nil {
			fmt.Printf("Error visiting %s: %v\n", url, err)
		}
	}

	c.Wait()

	return coursesByURL
}

func PopLastURLPart(url string) string {
	lastSlashIndex := strings.LastIndex(url, "/")
	if lastSlashIndex == -1 {
		return url
	}

	return url[:lastSlashIndex]
}

// pretty print course information
func PrintCourse(course *Course) {
	fmt.Printf("Course: %s (%s)\n", course.Title, course.Number)
	fmt.Printf("Section: %s\n", course.Section)
	fmt.Printf("Subject: %s\n", course.Subject)
	fmt.Printf("School: %s\n", course.School)
	fmt.Printf("Quarter: %s\n", course.Quarter)
	fmt.Printf("Topic: %s\n", course.Topic)

	fmt.Println("\nInstructors:")
	for _, instructor := range course.Instructors {
		fmt.Printf("(name %s) (phone %s) (email %s) (office hours %s)\n", instructor.Name, instructor.Phone, instructor.Email, instructor.OfficeHours)
	}

	fmt.Println("\nMeeting Times:")
	for _, meeting := range course.MeetingTimes {
		days := strings.Join(meeting.Days, ", ")
		timeStr := meeting.TimeRange

		if !meeting.StartTime.IsZero() && !meeting.EndTime.IsZero() {
			timeStr = meeting.StartTime.Format("3:04PM") + " - " + meeting.EndTime.Format("3:04PM")
		}
		fmt.Printf("- (location %s) (days %s) (time %s)\n", meeting.Location, days, timeStr)
	}

	fmt.Println("\nOverview:")
	fmt.Println(course.Overview)
}

func GetLastURLPart(url string) string {
	lastSlashIndex := strings.LastIndex(url, "/")
	if lastSlashIndex == -1 || lastSlashIndex == len(url)-1 {
		return ""
	}

	return url[lastSlashIndex+1:]
}
