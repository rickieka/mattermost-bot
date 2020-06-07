package mattermost

import (
	"log"

	"github.com/mattermost/mattermost-server/v5/model"
	"git.rickiekarp.net/rickie/mattermost-bot/data/config"
)

var mattermostClient *model.Client4

type Mattermost struct {
	Client *model.Client4
}

func NewMattermostClient(accessToken string) *Mattermost {
	mattermostClient = model.NewAPIv4Client(config.Conf.Mattermost.HttpUrl)
	mattermostClient.SetToken(accessToken)

	return &Mattermost{
		mattermostClient,
	}
}

func PrintMattermostError(err *model.AppError) {
	log.Println("\tError Details:")
	log.Println("\t\t" + err.Message)
	log.Println("\t\t" + err.Id)
	log.Println("\t\t" + err.DetailedError)
}
