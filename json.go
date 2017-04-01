package main

type EmoteSet struct {
	Emotes []Emote `json:"emoticons"`
}

type Emote struct {
	Id          int    `json:"id"`
	Code        string `json:"code"`
	EmoticonSet int    `json:"emoticon_set"`
}
