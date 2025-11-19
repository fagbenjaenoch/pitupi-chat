package main

import "strings"

type Message interface {
	Kind() string
}

type CommandMessage struct {
	command string
	param   string
	message string
}

func (c CommandMessage) Kind() string { return "command" }

type MentionMessage struct {
	user string
	text string
}

func (m MentionMessage) Kind() string { return "mention" }

type PlainMessage struct {
	Text string
}

func (p PlainMessage) Kind() string { return "plain" }

type Handler interface {
	Handle(input string) (Message, bool)
}

type CommandHandler struct{}

func (c *CommandHandler) Handle(input string) (Message, bool) {
	if !strings.HasPrefix(input, "!") {
		return nil, false
	}

	trim := strings.TrimPrefix(input, "!")
	parts := strings.SplitN(trim, " ", 2)

	command := parts[0]

	var param string
	hasParam := false
	var message string

	if strings.HasPrefix(parts[1], "@") && len(parts) > 2 {
		param = strings.TrimPrefix(parts[1], "@")
		hasParam = true
		message = parts[2]
	}

	if !hasParam {
		message = parts[1]
	}

	return CommandMessage{
		command: command,
		param:   param,
		message: message,
	}, true
}
