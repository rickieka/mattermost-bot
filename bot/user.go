package bot

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"git.rickiekarp.net/rickie/mattermost-bot/utils"
)

const (
	// NoPermissionReply - Default reply if permission not granted
	NoPermissionReply = "You are not allowed to use this command. :cry:"
	// PermissionLevelUser - Basic user permission
	PermissionLevelBasic = 0
	// PermissionLevelAdmin - All permissions
	PermissionLevelAll = 1
)

var userList = "conf/users.csv"

// users allowed to execute specific commands
var AllowedUsers []BotUser = ReadAllowedUsers()

// BotUser is a wrapper with the user name and permission level
type BotUser struct {
	UserId          string
	Username        string
	Email           string
	PermissionGroup string
	PermissionLevel int
}

// GetAllowedUser returns the user for the given email
func GetAllowedUserByEmail(email string) *BotUser {
	for _, user := range AllowedUsers {
		if email == user.Email {
			return &user
		}
	}
	return nil
}

func GetAllowedUserById(userId string) *BotUser {
	for _, user := range AllowedUsers {
		if userId == user.UserId {
			return &user
		}
	}
	return nil
}

func createDefaultBotUser(modelUser *model.User) *BotUser {
	user := BotUser{
		UserId:          modelUser.Id,
		Username:        modelUser.Username,
		Email:           modelUser.Email,
		PermissionGroup: "user",
		PermissionLevel: PermissionLevelBasic,
	}

	return &user
}

// AddUser adds a new user to the AllowedUsers array
// Returns true, if user was added, false otherwise
func AddUser(newUser BotUser) bool {

	// if the new user is written to the users.csv file, also add it to the AllowedUsers array
	if writeNewUserToUserlistFile(newUser) {
		AllowedUsers = append(AllowedUsers, newUser)
		return true
	}

	return false
}

// writeNewUserToUserlistFile writes a new user to the users.csv file
// if the file does not exist, it is created first
func writeNewUserToUserlistFile(newUser BotUser) bool {

	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(userList, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return false
	}

	defer f.Close()

	newUserString := newUser.UserId + "," +
		newUser.Username + "," +
		newUser.Email + "," +
		newUser.PermissionGroup + "," +
		strconv.Itoa(newUser.PermissionLevel)

	w := bufio.NewWriter(f)
	_, err = w.WriteString(newUserString + "\n")
	if err != nil {
		log.Println(err)
		return false
	}

	//Use Flush to ensure all buffered operations have been applied to the underlying writer.
	w.Flush()

	return true
}

// GetAllowedUsers reads the userList file from disk and returns all users
// in a BotUser array.
func ReadAllowedUsers() []BotUser {
	var users []BotUser

	// if the users.csv exists, try to load the users from it
	if utils.FileExists(userList) {
		file, err := os.Open(userList)
		if err != nil {
			panic(err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(file)

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.Split(scanner.Text(), ",")

			// set user group
			group := "none"
			if len(line) > 3 {
				group = line[3]
			}

			// set user permission
			right := PermissionLevelBasic
			if len(line) > 4 {
				x, err := strconv.Atoi(line[4])
				if err != nil {
					log.Println("Could not get user rights")
				}
				right = x
			}

			// instantiate user
			botUser := BotUser{
				UserId:          line[0],
				Username:        line[1],
				Email:           line[2],
				PermissionGroup: group,
				PermissionLevel: right,
			}

			users = append(users, botUser)
		}
	}

	return users
}
