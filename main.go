package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func GetNewBooks() (books []Book) {
	mathBooksUrl := "https://www.estantevirtual.com.br/busca?categoria=ciencias-exatas&sort=new-releases"
	req, err := http.NewRequest("GET", mathBooksUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
	req.Header.Set("Cookie", "qualquer coisa aq vai funcionar")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	finds := []Book{}
	doc.Find(".product-item.product-list__item").Each(func(_ int, s *goquery.Selection) {
		price := strings.TrimSpace(s.Find(".product-item__text--darken").Text())
		publisher := strings.TrimSpace(s.Find(".product-item__publishing").Text())
		author := strings.TrimSpace(s.Find(".product-item__author").Text())
		book := Book{
			Title:     s.Find("a").AttrOr("title", "Sem título"),
			Image:     s.Find("img").AttrOr("data-src", "Sem imagem"),
			Author:    author,
			Link:      s.Find("a").AttrOr("href", "Sem link"),
			Publisher: publisher,
			Price:     price,
		}
		finds = append(finds, book)
		log.Printf("Title: %s\nImage: %s\nAuthor: %s\nLink: %s\nPublisher: %s\nPrice: %s\n\n", book.Title, book.Image, book.Author, book.Link, book.Publisher, book.Price)
	})
	return finds
}

func main() {
	godotenv.Load()

	books := GetNewBooks()
	botToken := os.Getenv("BOT_TOKEN")
	chatID := os.Getenv("CHANNEL_ID")
	for _, book := range books {
		link := fmt.Sprintf("estantevirtual.com.br%s", book.Link)
		caption := fmt.Sprintf("*%s*\n%s\n%s\n%s\n\n[Link](%s)", book.Title, book.Author, book.Publisher, book.Price, link)

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

		req, err := http.NewRequest("POST", fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto", botToken), bytes.NewBuffer(jsonData))
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
