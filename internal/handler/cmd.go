package handler

import (
	"context"
	"errors"
	"fmt"
	"nova/internal/storage"
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
	cmdRPush  = "rpush"
	cmdLPush  = "lpush"
	cmdLRange = "lrange"
	cmdLLen   = "llen"
)

var (
	ErrUnknownCmd        = "Unknown command"
	ErrWrongNumberOfArgs = "Wrong number of arguments for '%s' command"
	ErrSyntax            = "syntax error"
	ErrProtocol          = "Protocol error"
	ErrWrongType         = "WRONGTYPE Operation against a key holding the wrong kind of value"
	ErrInvalidInt        = "Value is not an integer or out of range"
)

var (
	responseMsg = "request completed"

	nullString = "null"
	nullArray  = "[]"
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

	value, err := h.storage.Get(key)
	switch err {
	case storage.ErrKeyNotFound:
		log.Info(responseMsg, zap.String("response", nullString))
		return resp.NullString
	case storage.ErrWrongType:
		log.Info(responseMsg, zap.String("response", ErrWrongType))
		return resp.EncodeError(ErrWrongType)
	default:
		log.Info(responseMsg, zap.String("response", value))
		return resp.EncodeString(value)
	}
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

func (h *Handler) rPushHandler(ctx context.Context, args []string) []byte {
	var response string
	log := l.FromContext(ctx)

	if len(args) < 3 {
		response = fmt.Sprintf(ErrWrongNumberOfArgs, cmdRPush)
		log.Info(responseMsg, zap.String("response", response))
		return resp.EncodeError(response)
	}

	newLength, err := h.storage.RPush(args[1], args[2:])
	if errors.Is(err, storage.ErrWrongType) {
		log.Info(responseMsg, zap.String("response", ErrWrongType))
		return resp.EncodeError(ErrWrongType)
	}

	log.Info(responseMsg, zap.Int("response", newLength))
	return resp.EncodeInt(newLength)
}

func (h *Handler) lPushHandler(ctx context.Context, args []string) []byte {
	var response string
	log := l.FromContext(ctx)

	if len(args) < 3 {
		response = fmt.Sprintf(ErrWrongNumberOfArgs, cmdRPush)
		log.Info(responseMsg, zap.String("response", response))
		return resp.EncodeError(response)
	}

	newLength, err := h.storage.LPush(args[1], args[2:])
	if errors.Is(err, storage.ErrWrongType) {
		log.Info(responseMsg, zap.String("response", ErrWrongType))
		return resp.EncodeError(ErrWrongType)
	}

	log.Info(responseMsg, zap.Int("response", newLength))
	return resp.EncodeInt(newLength)
}

func (h *Handler) lRangeHandler(ctx context.Context, args []string) []byte {
	log := l.FromContext(ctx)

	if len(args) != 4 {
		response := fmt.Sprintf(ErrWrongNumberOfArgs, cmdLRange)
		log.Info(responseMsg, zap.String("response", response))
		return resp.EncodeError(response)
	}

	start, err := strconv.Atoi(args[2])
	if err != nil {
		log.Info(responseMsg, zap.String("response", ErrInvalidInt))
		return resp.EncodeError(ErrInvalidInt)
	}
	stop, err := strconv.Atoi(args[3])
	if err != nil {
		log.Info(responseMsg, zap.String("response", ErrInvalidInt))
		return resp.EncodeError(ErrInvalidInt)
	}

	values, err := h.storage.LRange(args[1], start, stop)
	if errors.Is(err, storage.ErrWrongType) {
		log.Info(responseMsg, zap.String("response", ErrWrongType))
		return resp.EncodeError(ErrWrongType)
	}
	if errors.Is(err, storage.ErrKeyNotFound) {
		log.Info(responseMsg, zap.String("response", nullArray))
		return resp.NullArray
	}

	log.Info(responseMsg, zap.Strings("response", values))
	return resp.EncodeArray(values)
}

func (h *Handler) lLenHandler(ctx context.Context, args []string) []byte {
	log := l.FromContext(ctx)

	if len(args) != 2 {
		response := fmt.Sprintf(ErrWrongNumberOfArgs, cmdLLen)
		log.Info(responseMsg, zap.String("response", response))
		return resp.EncodeError(response)
	}

	length, err := h.storage.ListLen(args[1])
	if errors.Is(err, storage.ErrWrongType) {
		log.Info(responseMsg, zap.String("response", ErrWrongType))
		return resp.EncodeError(ErrWrongType)
	}
	if errors.Is(err, storage.ErrKeyNotFound) {
		log.Info(responseMsg, zap.Int("response", 0))
		return resp.EncodeInt(0)
	}

	log.Info(responseMsg, zap.Int("response", length))
	return resp.EncodeInt(length)
}
