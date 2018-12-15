package main

import (
	"bytes"
	"fmt"
	"regexp"
	"time"
)

type Command struct {
	Regex       *regexp.Regexp
	Description string
	Run         func(...string) (result string, err error)
}

var Commands = []Command{
	Command{
		Regex:       regexp.MustCompile(`open door`),
		Description: "Opens the front door",
		Run: func(args ...string) (string, error) {
			err := door.Off()
			if err != nil {
				return "", err
			}

			time.Sleep(2 * time.Second)
			err = door.On()

			return "Got it!", err
		},
	},
}

var CommandsWithHelp = append(Commands, Command{
	Regex: regexp.MustCompile(`help`),
	Run: func(args ...string) (string, error) {
		var output bytes.Buffer

		fmt.Fprintf(&output, "Hi, I'm winona. I respond to the following commands:\n\n")

		for _, c := range Commands {
			fmt.Fprintf(&output, "*`%v`*\n", c.Regex)
			fmt.Fprintf(&output, "%v\n\n", c.Description)
		}

		return output.String(), nil
	},
})
