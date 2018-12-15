package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/nlopes/slack"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/chip"
)

type cmd struct {
	Action string `json:"action"`
}

var door *gpio.RelayDriver

func init() {
	chipAdaptor := chip.NewAdaptor()
	door = gpio.NewRelayDriver(chipAdaptor, "XIO-P1")
	door.On()
}

func main() {
	token := os.Getenv("SLACK_TOKEN")
	api := slack.New(token)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			info := rtm.GetInfo()
			prefix := fmt.Sprintf("<@%s>", info.User.ID)
			if ev.User != info.User.ID && strings.HasPrefix(ev.Text, prefix) {

				var matched bool
				for _, cmd := range CommandsWithHelp {
					if cmd.Regex.Match([]byte(ev.Text)) {
						matched = true
						var msg string
						var err error

						args := cmd.Regex.FindStringSubmatch(ev.Text)
						msg, err = cmd.Run(args...)

						if err != nil {
							rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("ERROR: %s", err), ev.Channel))
						}

						rtm.SendMessage(rtm.NewOutgoingMessage(msg, ev.Channel))
					}
				}

				if !matched {
					rtm.SendMessage(rtm.NewOutgoingMessage("I'm not sure what that means-- ask me for help to see what I can respond to.", ev.Channel))
				}
			}
		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:
			// Ignore other events..
		}
	}
}
