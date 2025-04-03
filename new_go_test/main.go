package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
	// "database/sql"
	// "github.com/lib/pq"
)

func getCourseListingLink(quarter string) string {
	return "https://class-descriptions.northwestern.edu/" + quarter
}

func main() {
	// var conninfo string = ""
	// db, err := sql.Open("postgres", conninfo)
	// if err != nil { // go style of error checking
	// 	panic(err)
	// }
	// err := db.QueryRow("SELECT get_work()").Scan(&work)

	c := colly.NewCollector()

	schools := make([]string, 0)
	classes := make(map[string][]string, 0)
	quarter := "4980" // 2025 spring quarter
	// +10 for each quarter

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link

		schools = append(schools, link) // school search
		if strings.Contains(link, quarter) {
			fmt.Printf("school link found: %q -> %s\n", e.Text, link)
			classes[link] = make([]string, 0)
			c.Visit(e.Request.AbsoluteURL(link)) // recurse
		}

		// TODO: improve this can be O(1)
		for _, school := range schools { // course search
			if strings.Contains(link, school) {
				fmt.Printf("class link found: %q -> %s\n", e.Text, link)
			}
			classes[school] = append(classes[school], link)
		}

		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		// c.Visit(e.Request.AbsoluteURL(link))
	})

	c.Visit(getCourseListingLink(quarter))
}

// //
// // You can see the program in action by defining a function similar to
// // the following:
// //
// // CREATE OR REPLACE FUNCTION public.get_work()
// //   RETURNS bigint
// //   LANGUAGE sql
// //   AS $$
// //     SELECT CASE WHEN random() >= 0.2 THEN int8 '1' END
// //   $$
// // ;

// package test

// import (
// 	"database/sql"
// 	"fmt"
// 	"time"

// 	"github.com/lib/pq"
// )

// func doWork(db *sql.DB, work int64) {
// 	// work here
// }

// func getWork(db *sql.DB) {
// 	for {
// 		// get work from the database here
// 		var work sql.NullInt64
// 		err := db.QueryRow("SELECT get_work()").Scan(&work)
// 		if err != nil {
// 			fmt.Println("call to get_work() failed: ", err)
// 			time.Sleep(10 * time.Second)
// 			continue
// 		}
// 		if !work.Valid {
// 			// no more work to do
// 			fmt.Println("ran out of work")
// 			return
// 		}

// 		fmt.Println("starting work on ", work.Int64)
// 		go doWork(db, work.Int64)
// 	}
// }

// func waitForNotification(l *pq.Listener) {
// 	select {
// 	case <-l.Notify:
// 		fmt.Println("received notification, new work available")
// 	case <-time.After(90 * time.Second):
// 		go l.Ping()
// 		// Check if there's more work available, just in case it takes
// 		// a while for the Listener to notice connection loss and
// 		// reconnect.
// 		fmt.Println("received no work for 90 seconds, checking for new work")
// 	}
// }

// func main() {
// 	var conninfo string = ""

// 	db, err := sql.Open("postgres", conninfo)
// 	if err != nil {
// 		panic(err)
// 	}

// 	reportProblem := func(ev pq.ListenerEventType, err error) {
// 		if err != nil {
// 			fmt.Println(err.Error())
// 		}
// 	}

// 	minReconn := 10 * time.Second
// 	maxReconn := time.Minute
// 	listener := pq.NewListener(conninfo, minReconn, maxReconn, reportProblem)
// 	err = listener.Listen("getwork")
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("entering main loop")
// 	for {
// 		// process all available work before waiting for notifications
// 		getWork(db)
// 		waitForNotification(listener)
// 	}
// }
