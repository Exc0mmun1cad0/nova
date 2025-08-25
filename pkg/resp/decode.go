package resp

import (
	"errors"
	"strconv"
	"strings"
)

var (
	errInvalidMultibulkLength = errors.New("invalid multibulk length")
	errInvalidMultibulkFormat = errors.New("invalid multibulk format")
)

// Decode decodes array of strings from resp protocol.
func Decode(msg []byte) ([]string, error) {
	args := strings.Split(string(msg), "\r\n")

	argsCount, err := strconv.Atoi(args[0][1:])
	if err != nil {
		return []string{}, errInvalidMultibulkLength
	}

	cmd := make([]string, 0, argsCount)
	for i := 2; i < len(args); i += 2 {
		// validate arg length
		if args[i-1][0] != '$' {
			return []string{}, errInvalidMultibulkFormat
		}
		argLen, err := strconv.Atoi(args[i-1][1:])
		if err != nil {
			return []string{}, errInvalidMultibulkFormat
		}
		if len(args[i]) != argLen {
			return []string{}, errInvalidMultibulkFormat
		}

		cmd = append(cmd, args[i])
	}

	return cmd, nil
}
