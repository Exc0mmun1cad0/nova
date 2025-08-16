package handler

import (
	"fmt"
	"nova/pkg/resp"
	"strconv"
	"strings"
	"time"
)

type handlerFunc func([]string) []byte

var (
	cmdInfo = "info"
	cmdPing = "ping"
	cmdEcho = "echo"
	cmdGet  = "get"
	cmdSet  = "set"
)

var (
	ErrUnknownCmd        = "Unknown command"
	ErrWrongNumberOfArgs = "Wrong number of arguments for '%s' command"
	ErrSyntax            = "syntax error"
)

func (h *Handler) pingHandler(args []string) []byte {
	return resp.EncodeSimpleString("PONG")
}

func (h *Handler) echoHandler(args []string) []byte {
	if len(args) != 2 {
		return resp.EncodeError(fmt.Sprintf(ErrWrongNumberOfArgs, cmdEcho))
	}
	return resp.EncodeString(args[1])
}

func (h *Handler) getHandler(args []string) []byte {
	if len(args) != 2 {
		return resp.EncodeError(fmt.Sprintf(ErrWrongNumberOfArgs, cmdGet))
	}
	key := args[1]
	value, ok := h.storage.Get(key)
	if ok {
		return resp.EncodeString(value)
	}
	return resp.NullString
}

func (h *Handler) setHandler(args []string) []byte {
	if len(args) < 3 {
		return resp.EncodeError(fmt.Sprintf(ErrWrongNumberOfArgs, cmdGet))
	} else if len(args) > 3 {
		switch strings.ToLower(args[3]) {
		case "px":
			if len(args) > 5 {
				return resp.EncodeError(ErrSyntax)
			}
			ttl, err := strconv.Atoi(args[4])
			if err != nil {
				// TODO: should be completely another error. I'll definitely fix it
				return resp.EncodeError(ErrSyntax)
			}
			key, value := args[1], args[2]
			h.storage.Set(key, value, time.Duration(ttl)*time.Millisecond)
		default:
			return resp.EncodeError(ErrSyntax)
		}
	} else {
		key, value := args[1], args[2]
		h.storage.Set(key, value, 0)
	}

	return resp.EncodeSimpleString("OK")
}
