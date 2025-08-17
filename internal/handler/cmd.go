package handler

import (
	"context"
	"fmt"
	l "nova/pkg/logger"
	"nova/pkg/resp"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

type handlerFunc func(context.Context, []string) []byte

var (
	cmdInfo   = "info"
	cmdPing   = "ping"
	cmdEcho   = "echo"
	cmdGet    = "get"
	cmdSet    = "set"
	cmdDelete = "del"
)

var (
	ErrUnknownCmd        = "Unknown command"
	ErrWrongNumberOfArgs = "Wrong number of arguments for '%s' command"
	ErrSyntax            = "syntax error"
	ErrProtocol          = "Protocol error"
)

var (
	responseMsg = "request completed"

	nullString = "null"
)

func (h *Handler) pingHandler(ctx context.Context, args []string) []byte {
	response := "PONG"

	log := l.FromContext(ctx)
	log.Info(responseMsg, zap.String("response", response))
	return resp.EncodeSimpleString(response)
}

func (h *Handler) echoHandler(ctx context.Context, args []string) []byte {
	var response string
	log := l.FromContext(ctx)

	if len(args) != 2 {
		response = fmt.Sprintf(ErrWrongNumberOfArgs, cmdEcho)
		log.Info(responseMsg, zap.String("response", response))
		return resp.EncodeError(response)
	}

	log.Info(responseMsg, zap.String("response", response))
	return resp.EncodeString(args[1])
}

func (h *Handler) getHandler(ctx context.Context, args []string) []byte {
	var response string
	log := l.FromContext(ctx)

	if len(args) != 2 {
		response = fmt.Sprintf(ErrWrongNumberOfArgs, cmdGet)
		log.Info(responseMsg, zap.String("response", response))
		return resp.EncodeError(response)
	}
	key := args[1]
	value, ok := h.storage.Get(key)
	if ok {
		log.Info(responseMsg, zap.String("response", value))
		return resp.EncodeString(value)
	}

	log.Info(responseMsg, zap.String("response", nullString))
	return resp.NullString
}

func (h *Handler) setHandler(ctx context.Context, args []string) []byte {
	log := l.FromContext(ctx)

	switch {
	case len(args) == 3:
		key, value := args[1], args[2]
		h.storage.Set(key, value, 0)

		log.Info(responseMsg, zap.String("response", "OK"))
		return resp.EncodeSimpleString("OK")

	case len(args) < 3:
		response := fmt.Sprintf(ErrWrongNumberOfArgs, cmdGet)
		log.Info(responseMsg, zap.String("response", response))
		return resp.EncodeError(response)

	case len(args) > 3:
		if strings.ToLower(args[3]) == "px" {
			if len(args) > 5 {
				log.Info(responseMsg, zap.String("response", ErrSyntax))
				return resp.EncodeError(ErrSyntax)
			}

			ttl, err := strconv.Atoi(args[4])
			if err != nil {
				// TODO: should be completely another error. I'll definitely fix it
				log.Info(responseMsg, zap.String("response", ErrSyntax))
				return resp.EncodeError(ErrSyntax)
			}

			key, value := args[1], args[2]
			h.storage.Set(key, value, time.Duration(ttl)*time.Millisecond)
			log.Info(responseMsg, zap.String("response", "OK"))
			return resp.EncodeSimpleString("OK")
		}

		log.Info(responseMsg, zap.String("response", ErrSyntax))
		return resp.EncodeError(ErrSyntax)
	}

	return resp.EncodeSimpleString("OK")
}

func (h *Handler) deleteHandler(ctx context.Context, args []string) []byte {
	var response string
	log := l.FromContext(ctx)

	if len(args) < 2 {
		response = fmt.Sprintf(ErrWrongNumberOfArgs, cmdDelete)
		log.Info(responseMsg, zap.String("response", response))
		return resp.EncodeError(response)
	}

	count := h.storage.DeleteMany(args[1:])
	log.Info(responseMsg, zap.Int("response", count))
	return resp.EncodeInt(count)
}
