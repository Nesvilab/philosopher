package sla

import (
	"fmt"

	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
)

// Run is the main entry point for Slack
func Run(name, token, message, channel, userID string) {

	var channelID string
	var timestamp string
	var e error

	api := slack.New(token)

	if len(userID) > 0 {

		_, _, channelID, e = api.OpenIMChannel(userID)
		if e != nil {
			fmt.Printf("%s\n", e)
			return
		}

		api.PostMessage(channelID, slack.MsgOptionText(message, false), slack.MsgOptionUsername(name))

	} else {

		channelID, timestamp, e = api.PostMessage(channel, slack.MsgOptionText(message, false), slack.MsgOptionUsername(name))
		if e != nil {
			fmt.Printf("%s\n", e)
			return
		}

	}

	logrus.Info("Message successfully sent to channel", channelID, "at", timestamp)

}
