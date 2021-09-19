package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	dateFlag            string
	domainNameRegexFlag string
	domainNameRegex     *regexp.Regexp
	csvOutputFlag       string
)

func init() {
	flag.StringVar(&dateFlag, "date", "", "domain registration date e.g. 2021-01-01 or date range e.g. 2021-01-01,2021-01-15")
	flag.StringVar(&domainNameRegexFlag, "regex", "", `only return domains that match provided regex e.g. "^[a-zA-Z]\-[a-zA-Z0-9]{2,3}\.(xyz|club|shop|online)"`)
	flag.StringVar(&csvOutputFlag, "out", "domains.csv", "csv output path")
}

func fetchDomains(lookups []Lookup, file *CsvWriter) error {
	errChan := make(chan error, len(lookups))
	sem := make(chan int, 10)

	var wg sync.WaitGroup
	wg.Add(len(lookups))

	for _, lookup := range lookups {
		go worker(lookup, sem, &wg, errChan, file)
	}
	wg.Wait()

	close(errChan)
	return <-errChan
}

func worker(lookup Lookup, sem chan int, wg *sync.WaitGroup, errChan chan error, file *CsvWriter) {
	defer wg.Done()
	sem <- 1

	doms, err := lookup.ScrapeDomains()
	if err != nil {
		errChan <- err
	}

	for _, d := range doms {
		switch domainNameRegexFlag {
		case "":
			record := []string{lookup.date, d}
			file.Write(record)
			break
		default:
			if lookup.regex.MatchString(d) {
				record := []string{lookup.date, d}
				file.Write(record)
			}
		}
	}

	<-sem
}

func main() {

	flag.Parse()

	file, err := NewCsvWriter(csvOutputFlag)
	if err != nil {
		log.Fatalf("error creating csv file: %s\n", err)
	}

	// check if required flags were provided
	if dateFlag == "" {
		fmt.Println("date flag required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// compile domainNameRegex
	domainNameRegex = regexp.MustCompile(domainNameRegexFlag)

	// parse date flag
	dates := getDateRage(dateFlag)

	// find number of pages that need to be scraped
	datesPages := make(map[string]int)
	for _, date := range dates {
		l := Lookup{date: date, page: 1}
		lp, err := l.FindLastPage()
		if err == nil {
			datesPages[date] = lp
		}
	}

	// create slice of lookups
	var lookups []Lookup
	for d, p := range datesPages {
		for i := 1; i <= p; i++ {
			lookups = append(lookups, Lookup{
				date:  d,
				page:  i,
				regex: domainNameRegex,
			})
		}
	}

	fetchDomains(lookups, file)

}

func getDateRage(dateFlag string) []string {
	var dates []string
	// single date
	if !strings.Contains(dateFlag, ",") {
		dates = append(dates, dateFlag)
		return dates
	}
	// date range
	dateFlagSplit := strings.Split(dateFlag, ",")
	fromDt, _ := time.Parse("2006-01-02", dateFlagSplit[0])
	toDt, _ := time.Parse("2006-01-02", dateFlagSplit[1])
	for dt := fromDt; !dt.After(toDt); dt = dt.AddDate(0, 0, 1) {
		dates = append(dates, dt.Format("2006-01-02"))
	}
	return dates
}
