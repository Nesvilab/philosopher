package sla

import (
	"fmt"

	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
)

// Run is the main entry point for Slack
func Run(name, token, message, channel string) {

	api := slack.New(token)

	// attachment := slack.Attachment{
	// 	//Pretext: "some pretext",
	// 	Text: message,
	// 	// Uncomment the following part to send a field too
	// 	/*
	// 		Fields: []slack.AttachmentField{
	// 			slack.AttachmentField{
	// 				Title: "a",
	// 				Value: "no",
	// 			},
	// 		},
	// 	*/
	// }

	channelID, timestamp, err := api.PostMessage(channel, slack.MsgOptionText(message, false), slack.MsgOptionUsername(name))
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	logrus.Info("Message successfully sent to channel", channelID, "at", timestamp)

	return
}
