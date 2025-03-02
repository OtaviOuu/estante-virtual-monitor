package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
)

type Book struct {
	Title     string
	Image     string
	Author    string
	Link      string
	Publisher string
	Price     string // 😁
}

func GetPage(url string) (*goquery.Document, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Cookie", "qualquer coisa aq vai funcionar")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return doc, nil
}

func GetNewBooks() (books []Book) {
	baseUrl := "https://www.estantevirtual.com.br"
	category := "ciencias-exatas"

	// f"https://www.estantevirtual.com.br/busca?q={category}&sort=new-releases" 😢
	mathBooksUrl := baseUrl + "/busca?q=" + category + "&sort=new-releases"
	doc, err := GetPage(mathBooksUrl)
	if err != nil {
		log.Fatal(err)
	}

	finds := []Book{}
	doc.Find(".product-item.product-list__item").Each(func(_ int, s *goquery.Selection) {
		author := strings.TrimSpace(s.Find(".product-item__author").Text())
		publisher := strings.TrimSpace(s.Find(".product-item__publishing").Text())
		price := strings.TrimSpace(s.Find(".product-item__text--darken").Text())

		book := Book{
			Author:    author,
			Publisher: publisher,
			Price:     price,
			Title:     s.Find("a").AttrOr("title", "Sem título"),
			Image:     s.Find("img").AttrOr("data-src", "Sem imagem"),
			Link:      s.Find("a").AttrOr("href", "Sem link"),
		}
		finds = append(finds, book)
		log.Println(book.Title + " | " + book.Price)
	})
	return finds
}

func main() {
	api_url := "https://api.telegram.org"

	godotenv.Load()

	books := GetNewBooks()
	botToken := os.Getenv("BOT_TOKEN")
	chatID := os.Getenv("CHANNEL_ID")
	for _, book := range books {
		link := "estantevirtual.com.br" + book.Link

		caption := `
			*` + book.Title + `*
			` + book.Author + `
			` + book.Publisher + `
			` + book.Price + `
			[Link](` + link + `)
		`

		data := map[string]interface{}{
			"chat_id":    chatID,
			"caption":    caption,
			"photo":      book.Image,
			"parse_mode": "Markdown",
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Fatal("Error marshalling JSON: ", err)
		}

		sendPhotoURL := api_url + "/bot" + botToken + "/sendPhoto"
		req, err := http.NewRequest("POST", sendPhotoURL, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatal("Error creating request: ", err)
		}

		req.Header.Set("Content-Type", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal("Error sending request: ", err)
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			body, err := io.ReadAll(res.Body)
			if err != nil {
				log.Fatal("Error reading response body: ", err)
			}
			log.Fatalf("status code error: %d %s, Response body: %s", res.StatusCode, res.Status, string(body))
		}

		log.Println("Photo sent successfully!")
	}
}
