package bot

import (
	"log"
	"os"
	"os/signal"

	"github.com/mattermost/mattermost-server/v5/model"
	"git.rickiekarp.net/rickie/mattermost-bot/clients/mattermost"
	"git.rickiekarp.net/rickie/mattermost-bot/data/config"
)

var mattermostClient *mattermost.Mattermost
var webSocketClient *model.WebSocketClient
var botCommands Commands

// Start prepares the mattermost api client and tries to listen on a mattermost websocket
func Start() {
	prepareMattermostClient()
	listenOnMattermostWebSocket()
}

func prepareMattermostClient() {
	mattermostClient = mattermost.NewMattermostClient(config.Conf.Bot.AccessToken)

	// Lets test to see if the mattermost server is up and running
	makeSureServerIsRunning()

}

func SetCommands(commandList Commands) {
	botCommands = commandList
}

func makeSureServerIsRunning() {
	if props, resp := mattermostClient.Client.GetOldClientConfig(""); resp.Error != nil {
		log.Println("There was a problem pinging the Mattermost server.  Are you sure it's running?")
		log.Println(resp.Error)
		os.Exit(1)
	} else {
		log.Println("Server detected and is running version " + props["Version"])
	}
}

func listenOnMattermostWebSocket() {
	setupGracefulShutdown()

	// Lets start listening to channels via the websocket
	webSocketClient, err := model.NewWebSocketClient4(config.Conf.Mattermost.WebsocketUrl, mattermostClient.Client.AuthToken)
	if err != nil {
		log.Println("We failed to connect to the web socket")
		log.Fatal(err)
	}

	webSocketClient.Listen()

	log.Println("Bot has started")
	//SendMsg("_Bot has **started** running_", config.Conf.Mattermost.Channels.Debugging)

	go func() {
		for {
			select {
			case resp := <-webSocketClient.EventChannel:

				// only respond to posted events
				if resp.Event == model.WEBSOCKET_EVENT_POSTED {
					HandleWebSocketResponse(resp, webSocketClient, AllowedUsers)
				}

			}
		}
	}()

	// You can block forever with
	select {}
}

func setupGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			if webSocketClient != nil {
				log.Println("Closing websocket client")
				webSocketClient.Close()
			}

			//SendMsg("_Bot has **stopped** running_", config.Conf.Mattermost.Channels.Debugging)

			os.Exit(0)
		}
	}()
}
