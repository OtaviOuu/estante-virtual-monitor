package scraper

import (
	"net/http"
	"os"
	"strings"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/PuerkitoBio/goquery"
)

type IScraper interface {
	Parse(url string) (*goquery.Document, error)
	GetBooksLinks(doc *goquery.Document) ([]string, error)
}

type Scraper struct {
	Name string
}

func NewScraper(name string) *Scraper {
	return &Scraper{
		Name: name,
	}
}

func (s *Scraper) Parse(url string) (*goquery.Document, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	applyDefaultHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (s *Scraper) GetBooksLinks(doc *goquery.Document) ([]string, error) {
	var books []string

	doc.Find("a.product-item__link.smarthint-tracking-card").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		link = strings.TrimSpace(link)

		books = append(books, link)
	})

	return books, nil
}

// ""bypass"" captcha
// Por algum motivo isso funciona ¯\_(ツ)_/¯
func applyDefaultHeaders(req *http.Request) *http.Request {
	req.Header.Set("User-Agent", browser.Random())
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("cookie", os.Getenv("COOKIES"))
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("accept-language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7,fr;q=0.6")
	req.Header.Set("cache-control", "max-age=0")
	req.Header.Set("priority", "u=0, i")
	req.Header.Set("sec-ch-ua", `"Chromium";v="130", "Google Chrome";v="130", "Not?A_Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Linux"`)
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "none")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")

	return req
}
