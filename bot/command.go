package bot

import "git.rickiekarp.net/rickie/mattermost-bot/data/model"

// Commands is a wrapper of a list of commands
type Commands struct {
	Commands []Command
}

// Command is the main command struct which needs to provide the actual executed action, validation, a name and a help for the user
type Command interface {

	// check if the user is allowed to use this command
	IsAllowed(user BotUser) bool

	// return true in case command did a response
	Execute(message string, channelId string, user BotUser) bool

	// each command has a name
	GetName() string

	// information on how to use the command with examples
	GetHelp() []model.Help
}
