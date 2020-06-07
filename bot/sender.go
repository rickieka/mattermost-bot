package bot

import (
	"fmt"
	"log"

	"git.rickiekarp.net/rickie/mattermost-bot/data/config"
	"github.com/mattermost/mattermost-server/v5/model"
)

func SendMsg(msg string, channelId string) {
	SendMsgWithProps(msg, channelId, nil)
}

func SendMsgWithProps(msg string, channelId string, props map[string]interface{}) {
	post := &model.Post{}
	post.ChannelId = channelId
	post.Message = msg

	if props != nil {
		post.Props = props
	}

	log.Println(fmt.Sprintf("SEND_MESSAGE: %s (%s)", post.Message, post.ChannelId))

	//log.Println(post.ToUnsanitizedJson())

	if _, resp := mattermostClient.Client.CreatePost(post); resp.Error != nil {
		log.Println("We failed to send a message to the logging channel")
		log.Println(resp.Error)
	}
}

// sendFirstInteractionMessage sends a welcome message to new users
func sendFirstInteractionMessage(channelId string) {
	welcomeMessage := "Hello there :wave:\n"
	welcomeMessage += "It looks it's your first time here, so let me introduce myself.\n"
	welcomeMessage += "I'm " + config.Conf.Bot.Name + ", your digital assistant and I am here to help you.\n"
	welcomeMessage += "Want me to do something for you? Sure! Just let me know what you need.\n"
	welcomeMessage += "To get started, type **help** for an overview of how I can help you.\n"
	SendMsg(welcomeMessage, channelId)
}
