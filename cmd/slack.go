package cmd

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/nlopes/slack"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/spf13/cobra"
)

// if you are here that's probably because you are curious on how Philosopher deals with the Slack user data.
// If you inspect the following code you will see that the organization is quite different from the other classes,
// all the information collected from the command-line execution is NOT internalized, and it's not added to the Meta
// data structure, which gathers the users input, so your private data is not stored and shared!

var name string
var token string
var message string
var channel string

// slackCmd represents the slack command
var slackCmd = &cobra.Command{
	Use:   "slack",
	Short: "Push notifications to Slack",
	Run: func(cmd *cobra.Command, args []string) {

		if len(token) < 1 {
			logrus.Fatal("You need to specify yout toke in order to push a notification.")
		}

		api := slack.New(token)
		params := slack.PostMessageParameters{}
		attachment := slack.Attachment{
		//Pretext: "some pretext",
		//Text:    "some text",
		// Uncomment the following part to send a field too
		/*
			Fields: []slack.AttachmentField{
				slack.AttachmentField{
					Title: "a",
					Value: "no",
				},
			},
		*/
		}

		params.Attachments = []slack.Attachment{attachment}
		params.Username = name

		channelID, timestamp, err := api.PostMessage(channel, message, params)

		_ = channelID
		_ = timestamp

		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		//logrus.Info("Message successfully sent to channel %s at %s", channelID, timestamp)

		return
	},
}

func init() {
	if len(os.Args) > 1 && os.Args[1] == "slack" {

		m.Restore(sys.Meta())

		slackCmd.Flags().StringVarP(&name, "username", "", "Philosopher", "specify the name of the bot")
		slackCmd.Flags().StringVarP(&token, "token", "", "", "specify the Slack API token")
		slackCmd.Flags().StringVarP(&message, "message", "", "", "specify the text of the message to send")
		slackCmd.Flags().StringVarP(&channel, "channel", "", "", "specify the channel name")
	}

	RootCmd.AddCommand(slackCmd)
}
