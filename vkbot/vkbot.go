package vkbot

import (
	"github.com/SevereCloud/vksdk/v2/api"
	"log"
)

var botAPI *api.VK
var userID int

func Initialize(token string, id int) {
	botAPI = api.NewVK(token)
	userID = id
}

func SendNotice(text string) {
	_, err := botAPI.MessagesSend(api.Params{
		"user_id":   userID,
		"random_id": 0,
		"message":   text,
	})
	log.Println(err)
}
