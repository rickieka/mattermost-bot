package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"git.rickiekarp.net/rickie/mattermost-bot/bot"
	"git.rickiekarp.net/rickie/mattermost-bot/command"
	"git.rickiekarp.net/rickie/mattermost-bot/command/admin"
	"git.rickiekarp.net/rickie/mattermost-bot/data/config"
	"git.rickiekarp.net/rickie/mattermost-bot/logging"
)

var (
	// Version set during go build using ldflags
	Version = "development"
)

func init() {
	argsWithoutProg := os.Args[1:]

	for i := 0; i < len(argsWithoutProg); i++ {
		switch argsWithoutProg[i] {
		case "version":
			fmt.Println(Version)
			os.Exit(0)
		}
	}
}

// Documentation for the Go driver can be found
// at https://godoc.org/github.com/mattermost/platform/model#Client
func main() {

	config.InitConfig("conf/config.yaml")

	logBool := flag.Bool("logging", config.Conf.Logging.Enabled, "logging to file enabled?")

	flag.Parse()

	config.Conf.Logging.Enabled = *logBool

	args := flag.Args()
	for i := 0; i < len(args); i++ {
		log.Println(args[i])
		switch args[i] {
		case "version":
			log.Println(Version)
			os.Exit(0)
		}
	}

	logging.ConfigureLogger()

	commands := getCommands()
	bot.SetCommands(commands)

	bot.Start()
}

// getCommands returns a list of all commands
func getCommands() bot.Commands {
	var commands bot.Commands

	commands = bot.Commands{
		Commands: []bot.Command{
			admin.NewAdmin(),
			command.NewHelp(&commands),
		},
	}

	return commands
}
