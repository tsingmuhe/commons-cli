package cli

import (
	"reflect"
	"testing"
)

func Test_tokenize(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []token
	}{
		{
			name:     "empty input",
			input:    []string{},
			expected: nil,
		},
		{
			name:  "command only",
			input: []string{"git"},
			expected: []token{
				{typ: tokenCommand, val: []string{"git"}},
			},
		},
		{
			name:  "command with subcommand",
			input: []string{"git", "commit"},
			expected: []token{
				{typ: tokenCommand, val: []string{"git"}},
				{typ: tokenArgument, val: []string{"commit"}},
			},
		},
		{
			name:  "short options",
			input: []string{"ls", "-a", "-l", "-h"},
			expected: []token{
				{typ: tokenCommand, val: []string{"ls"}},
				{typ: tokenShort, val: []string{"-a"}},
				{typ: tokenShort, val: []string{"-l"}},
				{typ: tokenShort, val: []string{"-h"}},
			},
		},
		{
			name:  "multiple short options",
			input: []string{"git", "-abc"},
			expected: []token{
				{typ: tokenCommand, val: []string{"git"}},
				{typ: tokenShort, val: []string{"-a"}},
				{typ: tokenShort, val: []string{"-b"}},
				{typ: tokenShort, val: []string{"-c"}},
			},
		},
		{
			name:  "long options",
			input: []string{"docker", "--verbose", "--all"},
			expected: []token{
				{typ: tokenCommand, val: []string{"docker"}},
				{typ: tokenLong, val: []string{"--verbose"}},
				{typ: tokenLong, val: []string{"--all"}},
			},
		},
		{
			name:  "long option with value",
			input: []string{"cmd", "--name=value"},
			expected: []token{
				{typ: tokenCommand, val: []string{"cmd"}},
				{typ: tokenLong, val: []string{"--name"}},
				{typ: tokenArgument, val: []string{"value"}},
			},
		},
		{
			name:  "mixed options and arguments",
			input: []string{"grep", "-i", "pattern", "--color", "file.txt"},
			expected: []token{
				{typ: tokenCommand, val: []string{"grep"}},
				{typ: tokenShort, val: []string{"-i"}},
				{typ: tokenArgument, val: []string{"pattern"}},
				{typ: tokenLong, val: []string{"--color"}},
				{typ: tokenArgument, val: []string{"file.txt"}},
			},
		},
		{
			name:  "separator",
			input: []string{"rm", "--", "-f", "--filename", "file.txt"},
			expected: []token{
				{typ: tokenCommand, val: []string{"rm"}},
				{typ: tokenSeparator, val: []string{"--"}},
				{typ: tokenArgument, val: []string{"-f"}},
				{typ: tokenArgument, val: []string{"--filename"}},
				{typ: tokenArgument, val: []string{"file.txt"}},
			},
		},
		{
			name:  "complex command",
			input: []string{"git", "commit", "-rf", "message", "--no-verify", "file1.txt", "--", "-f", "--option", "img.png"},
			expected: []token{
				{typ: tokenCommand, val: []string{"git"}},
				{typ: tokenArgument, val: []string{"commit"}},
				{typ: tokenShort, val: []string{"-r"}},
				{typ: tokenShort, val: []string{"-f"}},
				{typ: tokenArgument, val: []string{"message"}},
				{typ: tokenLong, val: []string{"--no-verify"}},
				{typ: tokenArgument, val: []string{"file1.txt"}},
				{typ: tokenSeparator, val: []string{"--"}},
				{typ: tokenArgument, val: []string{"-f"}},
				{typ: tokenArgument, val: []string{"--option"}},
				{typ: tokenArgument, val: []string{"img.png"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tokenize(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("tokenize() = %v, want %v", got, tt.expected)
			}
		})
	}
}
