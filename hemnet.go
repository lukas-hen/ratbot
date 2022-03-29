package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/bwmarrin/discordgo"
)

func PostGigaVilla(s *discordgo.Session, m *discordgo.MessageCreate) {
	res, err := http.Get("https://www.hemnet.se/bostader?item_types%5B%5D=villa&rooms_min=6&living_area_min=200&price_min=10000000")
	if err != nil {
		log.Fatal("Hemnet GET request failed, ", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	// Find the review items
	villaUrls := []string{}
	doc.Find(".js-listing-card-link").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			villaUrls = append(villaUrls, href)
		}
	})

	rand.Seed(time.Now().UnixNano())
	min := 0
	max := len(villaUrls)
	url := villaUrls[rand.Intn(max-min+1)+min]

	s.ChannelMessageSend(m.ChannelID, url)
}
