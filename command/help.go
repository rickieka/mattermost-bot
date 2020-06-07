package command

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	mattermostModel "github.com/mattermost/mattermost-server/v5/model"
	"git.rickiekarp.net/rickie/mattermost-bot/bot"
	"git.rickiekarp.net/rickie/mattermost-bot/data/model"
)

type help struct {
	commands *bot.Commands
}

// NewHelp is a command to provide helpful information for various commands
func NewHelp(commands *bot.Commands) *help {
	return &help{
		commands,
	}
}

func (c *help) IsAllowed(user bot.BotUser) bool {
	// allow to everyone
	return true
}

func (c *help) Execute(message string, channelID string, user bot.BotUser) bool {
	if matched, _ := regexp.MatchString(c.GetName()+`(?:$|\W)`, message); matched {

		var re = regexp.MustCompile("^(?i:help)(.*)")
		match := re.FindAllStringSubmatch(message, 1)
		if len(match) == 0 {
			return false
		}
		message = strings.Trim(match[0][1], " ")

		// build the help for the current user
		help, names := c.buildHelpTree(user)

		if message == "" {
			text := fmt.Sprintf("Hi %s!\n", user.Username)
			text += "The following commands are currently available:\n "
			for _, name := range names {
				text += fmt.Sprintf("- **%s**", name)
				if len(help[name].Description) > 0 {
					text += fmt.Sprintf(" _(%s)_", help[name].Description)
				}
				text += "\n"
			}

			text += fmt.Sprintf("With **help <command>** I can provide you with more details!")

			bot.SendMsg(text, channelID)
			return true

		} else {
			commandHelp, ok := help[message]
			if !ok {
				bot.SendMsg("Invalid command!", channelID)
				return true
			}

			var text string
			if len(commandHelp.Description) > 0 {
				text += commandHelp.Description + "\n"
			}
			var exampleUsages string
			for _, example := range commandHelp.Examples {
				exampleUsages += fmt.Sprintf(" - %s\n", example)
			}

			props := map[string]interface{}{
				"attachments": []mattermostModel.SlackAttachment{
					{
						Color: "#ACCAEB",
						Title: message,
						Text:  text,
						Fields: []*mattermostModel.SlackAttachmentField{
							{
								Title: "Example usage",
								Value: exampleUsages,
								Short: true,
							},
						},
					},
				},
			}

			bot.SendMsgWithProps("Here you go:", channelID, props)
			return true
		}
	}

	return false
}

func (c *help) GetName() string {
	return "help"
}

func (c *help) buildHelpTree(user bot.BotUser) (map[string]model.Help, []string) {
	var names []string
	help := map[string]model.Help{}

	for _, commands := range c.commands.Commands {
		// continue loop if current user has no permission for the command
		if !commands.IsAllowed(user) {
			continue
		}

		for _, commandHelp := range commands.GetHelp() {
			if _, ok := help[commandHelp.Command]; ok {
				// main command already defined
				continue
			}

			help[commandHelp.Command] = commandHelp
			names = append(names, commandHelp.Command)
		}
	}

	sort.Strings(names)

	return help, names
}

func (c *help) GetHelp() []model.Help {
	return []model.Help{
		{
			"help",
			"displays all available commands",
			[]string{
				"help",
			},
		},
	}
}
