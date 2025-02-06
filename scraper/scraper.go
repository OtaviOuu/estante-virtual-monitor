package scraper

import "github.com/PuerkitoBio/goquery"

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
