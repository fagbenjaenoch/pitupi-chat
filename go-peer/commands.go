package main

import (
	"fmt"
	"strings"
)

type Message interface {
	Kind() string
	Value() string
}

type CommandMessage struct {
	command string
	param   string
	message string
}

func (c CommandMessage) Kind() string { return "command" }
func (c CommandMessage) Value() string {
	return strings.TrimSpace(fmt.Sprintf("%s %s %s", c.command, c.param, c.message))
}

type MentionMessage struct {
	user string
	text string
}

func (m MentionMessage) Kind() string  { return "mention" }
func (m MentionMessage) Value() string { return fmt.Sprintf("%s %s", m.user, m.text) }

type PlainMessage struct {
	Text string
}

func (p PlainMessage) Kind() string  { return "plain" }
func (p PlainMessage) Value() string { return p.Text }

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
