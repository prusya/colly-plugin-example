package main

import (
	"github.com/gocolly/colly"
	"time"
	"regexp"
	"sync"
	"encoding/json"
	"net/http"
	"crypto/tls"
)

type Hash struct {
	Hash      string `json:"hash"`
	Name      string `json:"name"`
	ContentId string `json:"content_id"`
}

type Results struct {
	ID     string `json:"id"`
	Hashes []Hash `json:"hashes"`
}

var url string

func init() {
	url = "https://bitcq.com/search?q="
}

func Search(query string, id string) []byte {
	CheckSignature(query)
	results := crawl(query, id)
	encoded, _ := json.Marshal(results)
	//fmt.Printf("%+v", string(encoded))

	return encoded
}

func CheckSignature(query string) {
	c := colly.NewCollector()
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c.WithTransport(tr)

	tableCounter := 0

	c.OnHTML("body div.container table", func(e *colly.HTMLElement) {
		var mux sync.Mutex
		mux.Lock()
		tableCounter += 1
		mux.Unlock()
	})

	c.Visit(url + query)

	if tableCounter != 1 {
		panic("Must be only 1 'table' inside 'div.container'")
	}
}

func crawl(query string, id string) Results {
	results := Results{}
	results.ID = id

	c := colly.NewCollector()
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c.WithTransport(tr)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*bitcq.*",
		Parallelism: 2,
		Delay:       1 * time.Second,
	})

	c.OnHTML(".pagination a", func(e *colly.HTMLElement) {
		if e.Text == "›" {
			e.Request.Visit(e.Attr("href"))
		}
	})

	c.OnHTML("tr", func(e *colly.HTMLElement) {
		var mux sync.Mutex

		rawBtih := e.ChildAttr("td:nth-child(1) a", "href")
		re := regexp.MustCompile(".*?xt=urn:btih:([a-z0-9]+)&.*")
		btih := re.ReplaceAllString(rawBtih, "$1")
		title := e.ChildText("td:nth-child(2) a")
		if btih == "" {
			return
		}

		hash := Hash{
			Hash:      btih,
			Name:      title,
			ContentId: id,
		}
		mux.Lock()
		results.Hashes = append(results.Hashes, hash)
		mux.Unlock()
	})

	c.Visit(url + query)

	return results
}
