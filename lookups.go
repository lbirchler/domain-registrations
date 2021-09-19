package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const baseURL = "https://websitebiography.com/new_domain_registrations"

var lastPageRegex = regexp.MustCompile(`([^\/]+$)`)

type Lookup struct {
	date  string
	page  int
	regex *regexp.Regexp
}

func (l *Lookup) GetDoc() (*goquery.Document, error) {
	url := fmt.Sprintf("%s/%s/%d", baseURL, l.date, l.page)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (l *Lookup) FindLastPage() (int, error) {
	var lastPage int

	doc, err := l.GetDoc()
	if err != nil {
		return lastPage, err
	}

	doc.Find("a.btn.btn-outline-info.btn-rounded.btn_backg.inblock").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "Last") {
			lastPagePath, _ := s.Attr("href")
			extractLastPage := lastPageRegex.FindString(lastPagePath)
			pageNo, err := strconv.Atoi(extractLastPage)
			if err != nil {
				return
			}
			lastPage = pageNo
		}
	})

	return lastPage, nil
}

func (l *Lookup) ScrapeDomains() ([]string, error) {
	var domains []string

	doc, err := l.GetDoc()
	if err != nil {
		return nil, err
	}

	doc.Find("div.main_content.in_block div.hold_domains.whitebg_p a").Each(func(i int, s *goquery.Selection) {
		domain, _ := s.Attr("id")
		if domain != "" {
			domains = append(domains, domain)
		}
	})

	return domains, nil
}
