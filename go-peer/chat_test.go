package main

import "testing"

func TestCommandHandler(t *testing.T) {
	commandHandler := CommandHandler{}

	tests := []struct {
		input string
		want  string
	}{
		{
			"!ban @user123",
			"ban user123",
		},
		{
			"!poke @123",
			"poke 123",
		},
		{
			"!whisper @123 hi",
			"whisper 123 hi",
		},
		{
			"!announce hey guys",
			"announce hey guys",
		},
		{
			"!ls",
			"ls",
		},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result, _ := commandHandler.Handle(tc.input)

			if result.Value() != tc.want {
				t.Errorf("got %q, want %q", result.Value(), tc.want)
			}
		})
	}
}

func TestMentionHandler(t *testing.T) {
	mentionHandler := MentionHandler{}

	tests := []struct {
		input string
		want  string
	}{
		{
			"@user hi",
			"user hi",
		},
		{
			"@username",
			"username",
		},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result, _ := mentionHandler.Handle(tc.input)

			if result.Value() != tc.want {
				t.Errorf("want %q, got %q", tc.want, tc.input)
			}
		})
	}
}
