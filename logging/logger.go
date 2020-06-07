package logging

import (
	"log"
	"os"

	"git.rickiekarp.net/rickie/mattermost-bot/data/config"
)

// ConfigureLogger sets log format and output options
func ConfigureLogger() {
	log.SetFlags(log.Ltime | log.Lshortfile)

	if config.Conf.Logging.Enabled {
		outfile, _ := checkFile(config.Conf.Logging.Logfile)
		log.SetOutput(outfile)
	}
}

func checkFile(logFile string) (*os.File, error) {
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Error opening file:", err)
		f, err = getFallbackLog()
	}
	return f, err
}

func getFallbackLog() (*os.File, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("error when getting home directory")
	}

	f, err := os.OpenFile(home+"/fallback.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	log.Println("Writing to fallback log: " + home + "/fallback.log")
	return f, err
}
