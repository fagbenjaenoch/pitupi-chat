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
	var result string

	result = result + c.command

	if len(c.param) > 0 {
		result = result + " " + c.param
	}

	if len(c.message) > 0 {
		result = result + " " + c.message
	}

	return result
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
	parts := strings.SplitN(trim, " ", 3)

	command := parts[0]

	if len(parts) < 2 {
		return CommandMessage{
			command: command,
			param:   "",
			message: "",
		}, true
	}

	var param string
	hasParam := false
	var message string

	if after, ok := strings.CutPrefix(parts[1], "@"); ok {
		param = strings.TrimSpace(after)
		hasParam = true
	}

	if len(parts) == 2 && !hasParam {
		param = parts[1]
	}

	if len(parts) > 2 {
		message = parts[2]
	}

	if !hasParam {
		message = strings.TrimSpace(parts[1])
	}

	fmt.Printf("command:%q, param: %q, message:%q\n", command, param, message)

	return CommandMessage{
		command: command,
		param:   param,
		message: message,
	}, true
}
