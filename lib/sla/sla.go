package sla

import (
	"github.com/nlopes/slack"
)

// Run is the main entry point for Slack
func Run(name, token, message, channel string) {

	//api := slack.New(token)
	params := slack.PostMessageParameters{}
	//attachment := slack.Attachment{}

	//params.Attachments = []slack.Attachment{attachment}
	params.Username = name

	//channelID, timestamp, err := api.PostMessage(channel, message, params)

	// _ = channelID
	// _ = timestamp

	// if err != nil {
	// 	fmt.Printf("%s\n", err)
	// 	return
	// }

	return
}
