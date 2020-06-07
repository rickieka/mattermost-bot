package admin

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"git.rickiekarp.net/rickie/mattermost-bot/bot"
	"git.rickiekarp.net/rickie/mattermost-bot/data/model"
)

type admin struct {
}

// NewScheduleList is a command to list all active schedules in a channel or all
func NewAdmin() *admin {
	return &admin{}
}

func (c *admin) IsAllowed(user bot.BotUser) bool {
	if user.PermissionGroup == "admin" {
		return true
	}

	return false
}

func (c *admin) GetName() string {
	return "admin"
}

func (c *admin) Execute(message string, channelID string, user bot.BotUser) bool {

	if strings.HasPrefix(message, c.GetName()) {

		s := strings.Split(message, " ")
		log.Println(s)

		if len(s) > 1 {
			switch s[1] {
			case "listuser":
				var userlist string
				userlist =
					`| ID  | User  | Email  |  Group  | Permissions |
				| :-----------: | :------------ | :-----------: | :---------------: | :---------------: | :-----: |
				`

				for _, user := range bot.AllowedUsers {
					userlist += "| " + user.UserId + "| " + user.Username + " | " + user.Email + " | " + user.PermissionGroup + " | " + strconv.Itoa(user.PermissionLevel) + " |\n"
				}

				bot.SendMsg(userlist, channelID)
				return true
			}
		}
	}
	return false
}

// isValidEmail checks if the given email format is correct
func isValidEmail(emailInput string) bool {
	mailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return mailRegex.MatchString(emailInput)
}

func (c *admin) GetHelp() []model.Help {
	return []model.Help{
		{
			c.GetName(),
			"common admin functions",
			[]string{
				"admin listuser",
			},
		},
	}
}
