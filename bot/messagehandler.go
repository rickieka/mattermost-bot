package bot

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"git.rickiekarp.net/rickie/mattermost-bot/clients/mattermost"
	"git.rickiekarp.net/rickie/mattermost-bot/data/config"
	"git.rickiekarp.net/rickie/mattermost-bot/net/http"
	"github.com/mattermost/mattermost-server/v5/model"
)

func HandleWebSocketResponse(event *model.WebSocketEvent, webSocketClient *model.WebSocketClient, users []BotUser) {
	// read post event data
	post := model.PostFromJson(strings.NewReader(event.Data["post"].(string)))

	// ignore events of bot user
	if post.UserId == config.Conf.Bot.ID {
		return
	}

	if post != nil {

		// send "User typing..." message
		webSocketClient.UserTyping(post.ChannelId, "")

		requestUser := GetAllowedUserById(post.UserId)

		// if the requestUser is not known to the bot, get its Mattermost information and add it to our users.csv
		if requestUser == nil {
			// get User from posted event
			user := mattermost.GetUserById(post.UserId)

			// create new user
			requestUser = createDefaultBotUser(user)

			// add it to the users.csv list
			if !AddUser(*requestUser) {
				SendMsg("Could not add new user "+requestUser.Email, config.Conf.Mattermost.Channels.Debugging)
			} else {
				log.Println(fmt.Sprintf("Added new user:  %s (%s)", requestUser.Email, requestUser.UserId))
			}

			sendFirstInteractionMessage(post.ChannelId)
		}

		// go through the list of existing commands
		for i := range botCommands.Commands {
			if botCommands.Commands[i].IsAllowed(*requestUser) {
				if botCommands.Commands[i].Execute(strings.TrimSpace(post.Message), post.ChannelId, *requestUser) {
					return
				}
			}
		}

		switch msg := post.Message; msg {
		case "hello":
			SendMsg("Hello :shibahey:", post.ChannelId)
			return
		}

		// if you see any word matching 'hello' then respond
		if matched, _ := regexp.MatchString(`hello(?:$|\W)`, post.Message); matched {

			stringSlice := strings.Split(post.Message, " ")
			log.Println(stringSlice)

			value, err := http.MakeGetRequestAndReturnStatusCode(stringSlice[1])
			if err != nil {
				log.Println(err)
				SendMsg(err.Error(), post.ChannelId)
				return
			}
			SendMsg(strconv.Itoa(value), post.ChannelId)

			SendMsg(
				"Hello!\n**bold**text\n[Check out!](https://about.mattermost.com/)",
				post.ChannelId)

			return
		}

	}

	// if at this point no appropriate message was returned, we assume the command was incorrect
	SendMsg("I did not understand you! :cry:", post.ChannelId)
}
