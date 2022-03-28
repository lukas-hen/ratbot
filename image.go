package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

type UnsplashResponse struct {
	Id   string `json:"id"`
	Urls struct {
		Small string `json:"small"`
	} `json:"urls"`
}

func PostRatImage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bot's own messages
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!råtta" {
		client := &http.Client{}

		req, err := http.NewRequest("GET", "https://api.unsplash.com/photos/random?query=rat", nil)
		if err != nil {
			log.Fatal("Could not create request, ", err)
			return

		}

		req.Header.Add("Authorization", fmt.Sprintf("Client-ID %s", os.Getenv("UNSPLASH_KEY")))

		res, err := client.Do(req)
		if err != nil {
			log.Fatal("Failed to fetch rat image")
			return
		}

		body, err := ioutil.ReadAll(res.Body)
		fmt.Println(string(body))
		if err != nil {
			log.Fatal("Unable to read response body, ", err)
		}

		var data UnsplashResponse

		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Fatal("Unmarshal of json failed, ", err)
		}
		fmt.Println("URL: ", data.Urls.Small)
		embed := &discordgo.MessageEmbed{
			Author:      &discordgo.MessageEmbedAuthor{},
			Description: "Here is your rat!",
			Image: &discordgo.MessageEmbedImage{
				URL: string(data.Urls.Small),
			},
			Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
			Title:     "Rat image",
		}

		s.ChannelMessageSendEmbed(m.ChannelID, embed)
	}
}
