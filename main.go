package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
)

func parseBookCategory(page int) {
	c := colly.NewCollector(colly.AllowURLRevisit(), colly.Async(true))

	var books []string

	c.OnHTML("h2.product-item__title.product-item__title--mt.product-item__name", func(e *colly.HTMLElement) {
		title := strings.Trim(e.Text, " \n")
		books = append(books, title)
		fmt.Println(title)
	})

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 100,
		Delay:       0 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("cookie", os.Getenv("COOKIES"))
		r.Headers.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		r.Headers.Set("accept-language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7,fr;q=0.6")
		r.Headers.Set("cache-control", "max-age=0")
		r.Headers.Set("priority", "u=0, i")
		r.Headers.Set("sec-ch-ua", `"Chromium";v="130", "Google Chrome";v="130", "Not?A_Brand";v="99"`)
		r.Headers.Set("sec-ch-ua-mobile", "?0")
		r.Headers.Set("sec-ch-ua-platform", `"Linux"`)
		r.Headers.Set("sec-fetch-dest", "document")
		r.Headers.Set("sec-fetch-mode", "navigate")
		r.Headers.Set("sec-fetch-site", "none")
		r.Headers.Set("sec-fetch-user", "?1")
		r.Headers.Set("upgrade-insecure-requests", "1")
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	for i := 1; i <= page; i++ {
		c.Visit("https://www.estantevirtual.com.br/busca?categoria=ciencias-exatas&page=" + fmt.Sprint(i))
	}
	c.Wait()

}

func main() {
	loadEnv()
	parseBookCategory(100)
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
