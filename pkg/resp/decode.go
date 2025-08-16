package resp

import (
	"fmt"
	"strconv"
	"strings"
)

// Decode decodes array of strings from resp protocol.
func Decode(msg []byte) ([]string, error) {
	args := strings.Split(string(msg), "\r\n")

	argsCount, err := strconv.Atoi(args[0][1:])
	if err != nil {
		return nil, fmt.Errorf("invalid number of args: %v", err)
	}

	cmd := make([]string, 0, argsCount)
	for i := 2; i < len(args); i += 2 {
		cmd = append(cmd, args[i])
	}

	return cmd, nil
}
