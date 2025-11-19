package main

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
