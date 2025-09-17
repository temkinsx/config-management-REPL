package parser

import (
	"strings"
)

type Input struct {
	Cmd  string
	Args []string
}

func ParseLine(line string) (string, []string, error) {
	if line == "" {
		return "", nil, nil
	}

	parts := strings.Fields(line)
	cmd := parts[0]
	args := parts[1:]

	return cmd, args, nil
}
