package mattermost

import (
	"log"
	"os"

	"github.com/mattermost/mattermost-server/v5/model"
)

func GetUserById(userID string) *model.User {
	if user, resp := mattermostClient.GetUser(userID, ""); resp.Error != nil {
		log.Println("We failed to get the user")
		PrintMattermostError(resp.Error)
		return nil
	} else {
		return user
	}
}

func GetUserByEmail(userEmail string) *model.User {
	if user, resp := mattermostClient.GetUserByEmail(userEmail, ""); resp.Error != nil {
		log.Println("We failed to get the user")
		PrintMattermostError(resp.Error)
		return nil
	} else {
		return user
	}
}

func GetChannelByName(channelName string, botTeamId string) {
	if channel, resp := mattermostClient.GetChannelByName(channelName, botTeamId, ""); resp.Error != nil {
		log.Println("We failed to get the channel")
		PrintMattermostError(resp.Error)
	} else {
		log.Println(channel)
	}
}

func FindBotTeam(teamName string) {
	if team, resp := mattermostClient.GetTeamByName(teamName, ""); resp.Error != nil {
		log.Println("We failed to get the initial load")
		log.Println("or we do not appear to be a member of the team '" + teamName + "'")
		PrintMattermostError(resp.Error)
		os.Exit(1)
	} else {
		log.Println(team)
	}
}
