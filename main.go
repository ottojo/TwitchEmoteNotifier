package main

import (
	"encoding/base64"
	"encoding/json"
	"github.com/ChimeraCoder/anaconda"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"fmt"
)

func main() {
	//Load old Emotes from file specified as command line argument
	fmt.Println("Loading old Emote Set...")
	emoteFile, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	oldEmotesFile, err := ioutil.ReadFile(emoteFile)
	checkError(err)
	var oldEmoteSet EmoteSet
	json.Unmarshal(oldEmotesFile, &oldEmoteSet)
	fmt.Println("done")

	//Download new Emotes
	fmt.Println("Downloading new Emote Set...")
	newEmoteSet := downloadFullEmoteSet()
	fmt.Println("done")

	//Compare new Emotes against old Emotes
	fmt.Println("Comparing Emotes...")
	newEmotes, removedEmotes := compareEmoteSets(oldEmoteSet, newEmoteSet)
	fmt.Println("done")

	fmt.Println("New Emotes:")
	fmt.Println(newEmotes)
	fmt.Println("Removed Emotes:")
	fmt.Println(removedEmotes)

	//Save new Emotes
	emoteSetBytes, err := json.Marshal(newEmoteSet)
	checkError(err)
	ioutil.WriteFile(emoteFile, emoteSetBytes, 0644)

	//Set up twitter
	fmt.Println("Setting up twitter...")
	twitterConsumerKey, twitterConsumerKeyAvailiable := os.LookupEnv("TEN_CONSUMER_KEY")
	if !twitterConsumerKeyAvailiable {
		log.Fatal("Environment Variable \"TEN_CONSUMER_KEY\" not availiable.")
	}
	anaconda.SetConsumerKey(twitterConsumerKey)

	twitterConsumerSecret, twitterConsumerSecretAvailiable := os.LookupEnv("TEN_CONSUMER_SECRET")
	if !twitterConsumerSecretAvailiable {
		log.Fatal("Environment Variable \"TEN_CONSUMER_SECRET\" not availiable.")
	}
	anaconda.SetConsumerSecret(twitterConsumerSecret)

	twitterAccessToken, twitterAccessTokenAvailiable := os.LookupEnv("TEN_ACCESS_TOKEN")
	if !twitterAccessTokenAvailiable {
		log.Fatal("Environment Variable \"TEN_ACCESS_TOKEN\" not availiable.")
	}
	twitterAccessTokenSecret, twitterAccessTokenSecretAvailiable := os.LookupEnv("TEN_ACCESS_TOKEN_SECRET")
	if !twitterAccessTokenSecretAvailiable {
		log.Fatal("Environment Variable \"TEN_ACCESS_TOKEN_SECRET\" not availiable.")
	}
	api := anaconda.NewTwitterApi(twitterAccessToken, twitterAccessTokenSecret)
	defer api.Close()
	api.EnableThrottling(10*time.Second, 10000000)
	fmt.Println("done")

	fmt.Println("Tweeting Changes...")
	tweetChanges(api, newEmotes, removedEmotes)
	fmt.Println("done")
}

func tweetChanges(api *anaconda.TwitterApi, newEmotes, removedEmotes []Emote) {
	for _, emote := range newEmotes {
		//Get Emote and upload to twitter
		media, err := api.UploadMedia(base64.StdEncoding.EncodeToString(httpGET("https://static-cdn.jtvnw.net/emoticons/v1/" + strconv.Itoa(emote.Id) + "/3.0")))
		checkError(err)
		tweetParams := url.Values{}
		tweetParams.Set("media_ids", media.MediaIDString)
		tweet := "New: " + emote.Code
		api.PostTweet(tweet, tweetParams)
		fmt.Println("Tweeted: \"" + tweet + "\"")
	}

	for _, emote := range removedEmotes {
		tweetParams := url.Values{}
		tweet := "Removed: " + emote.Code
		api.PostTweet(tweet, tweetParams)
		fmt.Println("Tweeted: \"" + tweet + "\"")
	}
}
