package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
)

const (
	SAMPLE_NAME = "Mattermost Bot Sample"

	MATTERMOST_PROTOCOL_HTTP = "http://"
	MATTERMOST_PROTOCOL_WS   = "ws://"
	MATTERMOST_HOST          = "localhost:8065"

	BOT_ID         = "j69dqnrkqfdrzyh1tmyc697whh"
	BOT_AUTH_TOKEN = "y36o13cnqbdjdj6mwgsrnf7pzh"
)

var client *model.Client4
var webSocketClient *model.WebSocketClient

var debuggingChannel *model.Channel

// Documentation for the Go driver can be found
// at https://godoc.org/github.com/mattermost/platform/model#Client
func main() {
	println(SAMPLE_NAME)

	SetupGracefulShutdown()

	client = model.NewAPIv4Client(MATTERMOST_PROTOCOL_HTTP + MATTERMOST_HOST)

	// Lets test to see if the mattermost server is up and running
	MakeSureServerIsRunning()

	// This is an important step.  Lets make sure we use the team
	// for all future web service requests that require a team.
	client.SetToken(BOT_AUTH_TOKEN)

	// Lets start listening to some channels via the websocket!
	webSocketClient, err := model.NewWebSocketClient4(MATTERMOST_PROTOCOL_WS+MATTERMOST_HOST, client.AuthToken)
	if err != nil {
		println("We failed to connect to the web socket")
		PrintError(err)
	}

	webSocketClient.Listen()

	go func() {
		for {
			select {
			case resp := <-webSocketClient.EventChannel:
				HandleWebSocketResponse(resp)
			}
		}
	}()

	// You can block forever with
	select {}
}

func MakeSureServerIsRunning() {
	if props, resp := client.GetOldClientConfig(""); resp.Error != nil {
		println("There was a problem pinging the Mattermost server.  Are you sure it's running?")
		PrintError(resp.Error)
		os.Exit(1)
	} else {
		println("Server detected and is running version " + props["Version"])
	}
}

func SendMsgToDebuggingChannel(msg string, replyToId string, userId string, channelId string) {
	post := &model.Post{}

	post.UserId = userId
	post.ChannelId = channelId
	post.Message = msg

	fmt.Println(post.ToUnsanitizedJson())

	if _, resp := client.CreatePost(post); resp.Error != nil {
		println("We failed to send a message to the logging channel")
		PrintError(resp.Error)
	}
}

func HandleWebSocketResponse(event *model.WebSocketEvent) {
	// Lets only reponded to messaged posted events
	if event.Event != model.WEBSOCKET_EVENT_POSTED {
		return
	}

	post := model.PostFromJson(strings.NewReader(event.Data["post"].(string)))

	// ignore my events
	if post.UserId == BOT_ID {
		return
	}

	println("responding to debugging channel msg")

	if post != nil {

		// if you see any word matching 'alive' then respond
		if matched, _ := regexp.MatchString(`(?:^|\W)alive(?:$|\W)`, post.Message); matched {
			SendMsgToDebuggingChannel("Yes I'm running", post.Id, post.UserId, post.ChannelId)
			return
		}

		// if you see any word matching 'up' then respond
		if matched, _ := regexp.MatchString(`(?:^|\W)up(?:$|\W)`, post.Message); matched {
			SendMsgToDebuggingChannel("Yes I'm running", post.Id, post.UserId, post.ChannelId)
			return
		}

		// if you see any word matching 'running' then respond
		if matched, _ := regexp.MatchString(`(?:^|\W)running(?:$|\W)`, post.Message); matched {
			SendMsgToDebuggingChannel("Yes I'm running", post.Id, post.UserId, post.ChannelId)
			return
		}

		// if you see any word matching 'hello' then respond
		if matched, _ := regexp.MatchString(`(?:^|\W)Hi(?:$|\W)`, post.Message); matched {
			SendMsgToDebuggingChannel("Hi", post.Id, post.UserId, post.ChannelId)
			return
		}
	}

	SendMsgToDebuggingChannel(
		"Hello!\n**bold**sdjfh\n[Check out!](https://about.mattermost.com/)",
		post.Id,
		post.UserId,
		post.ChannelId)

	SendMsgToDebuggingChannel(
		`| Left-Aligned  | Center Aligned  | Right Aligned |
		| :------------ |:---------------:| -----:|
		| Left column 1 | this text       |  $100 |
		| Left column 2 | is              |   $10 |
		| Left column 3 | centered        |    $1 |`,
		post.Id,
		post.UserId,
		post.ChannelId)

}

func PrintError(err *model.AppError) {
	println("\tError Details:")
	println("\t\t" + err.Message)
	println("\t\t" + err.Id)
	println("\t\t" + err.DetailedError)
}

func SetupGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			if webSocketClient != nil {
				log.Println("Closing websocket client")
				webSocketClient.Close()
			}

			os.Exit(0)
		}
	}()
}
