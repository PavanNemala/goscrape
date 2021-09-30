package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gorilla/mux"
)

type Product struct {
	Name         string `json:"name"`
	ImageURL     string `json:"imageURL"`
	Description  string `json:"description"`
	Price        string `json:"price"`
	TotalReviews int    `json:"totalReviews"`
}

type ApiBody struct {
	Url     string  `json:"url"`
	Product Product `json:"product"`
}

func scrape(w http.ResponseWriter, r *http.Request) {
	final := make(map[string]string)
	c := colly.NewCollector()

	c.OnHTML("title", func(e *colly.HTMLElement) {
		final["Title"] = e.Text
	})

	c.OnHTML("span[id=priceblock_ourprice]", func(e *colly.HTMLElement) {
		final["price"] = strings.TrimSpace(e.Text)
	})

	if final["price"] == "" {
		c.OnHTML("span[id=priceblock_dealprice]", func(e *colly.HTMLElement) {
			final["price"] = strings.TrimSpace(e.Text)
		})
	}

	c.OnHTML("div[data-hook=total-review-count]", func(e *colly.HTMLElement) {
		final["reviews"] = strings.TrimSpace(e.Text)
	})

	c.OnHTML("div[id=feature-bullets]", func(e *colly.HTMLElement) {
		final["description"] = strings.Replace(e.Text, "\n", "", -1)
	})

	c.Visit(string(r.URL.Query().Get("url")))
	json.NewEncoder(w).Encode(final)
}

func createdoc(w http.ResponseWriter, r *http.Request) {
	var ApiBodyin ApiBody
	_ = json.NewDecoder(r.Body).Decode(&ApiBodyin)
	url := ApiBodyin.Url
	product := ApiBodyin.Product

	file, _ := os.Create("Result.docx")
	file.WriteString("URL:" + url + "\n")
	file.WriteString("NAME:" + product.Name + "\n")
	file.WriteString("IMAGE:" + product.ImageURL + "\n")
	file.WriteString("PRICE:" + product.Price + "\n")
	file.WriteString("DESCRIPTION:" + strings.Replace(product.Description, "\n", "", -1) + "\n")
	file.WriteString("REVIEWS:" + fmt.Sprintln(product.TotalReviews) + "\n")
	file.WriteString(fmt.Sprintln(time.Now()))
	json.NewEncoder(w).Encode("success")
	defer file.Close()
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/amazonscrape", scrape)

	router.HandleFunc("/createdoc", createdoc)
	log.Fatal(http.ListenAndServe(":8000", router))
}
