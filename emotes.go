package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func downloadFullEmoteSet() EmoteSet {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.twitch.tv/kraken/chat/emoticon_images", nil)
	checkError(err)

	clientId, clientIdAvailiable := os.LookupEnv("TEN_CLIENT_ID")
	if !clientIdAvailiable {
		log.Fatal("Environment Variable \"TEN_CLIENT_ID\" not availiable.")
	}
	req.Header.Add("Client-ID", clientId)
	req.Header.Add("Accept", "application/vnd.twitchtv.v5+json")

	response, err := client.Do(req)
	checkError(err)

	defer response.Body.Close()

	responseBytes, err := ioutil.ReadAll(response.Body)
	checkError(err)

	var emoteSet EmoteSet
	err = json.Unmarshal(responseBytes, &emoteSet)
	checkError(err)

	return emoteSet
}

func compareEmoteSets(oldSet EmoteSet, newSet EmoteSet) (newEmotes []Emote, removedEmotes []Emote) {
	log.Println("Searching for new Emotes")
	for _, emote := range newSet.Emotes {
		if !containsEmote(oldSet, emote) {
			newEmotes = append(newEmotes, emote)
		}
	}

	log.Println("Searching for removed Emotes")
	for _, emote := range oldSet.Emotes {
		if !containsEmote(newSet, emote) {
			removedEmotes = append(removedEmotes, emote)
		}
	}

	return newEmotes, removedEmotes
}

func containsEmote(set EmoteSet, emote Emote) bool {
	for _, emoteFromSet := range set.Emotes {
		if emoteFromSet.Id == emote.Id {
			return true
		}
	}
	return false
}
