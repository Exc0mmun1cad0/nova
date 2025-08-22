package handler

import (
	"context"
	"fmt"
	l "nova/pkg/logger"
	"nova/pkg/resp"
	"strings"
	"time"

	"go.uber.org/zap"
)

type Storage interface {
	Set(key, value string, ttl time.Duration)
	Get(key string) (string, error)
	DeleteMany(keys []string) int

	RPush(key string, values []string) (int, error)
	LPush(key string, values []string) (int, error)
	LRange(key string, start, stop int) ([]string, error)
}

type Handler struct {
	storage Storage
	dict    map[string]handlerFunc
}

func NewHandler(storage Storage) *Handler {
	h := &Handler{
		storage: storage,
	}

	dict := map[string]handlerFunc{
		cmdPing: h.pingHandler,
		cmdEcho: h.echoHandler,

		cmdGet:    h.getHandler,
		cmdSet:    h.setHandler,
		cmdDelete: h.deleteHandler,

		cmdRPush:  h.rPushHandler,
		cmdLPush:  h.lPushHandler,
		cmdLRange: h.lRangeHandler,
	}

	h.dict = dict
	return h
}

func (h *Handler) Serve(ctx context.Context, input []byte) []byte {
	log := l.FromContext(ctx)

	args, err := resp.Decode(input)
	if err != nil {
		return resp.EncodeError(fmt.Sprintf(ErrProtocol, err.Error()))
	}
	log.Info("decoded request", zap.Strings("args", args))

	cmd := strings.ToLower(args[0])
	handler, ok := h.dict[cmd]
	if !ok {
		return resp.EncodeError(ErrUnknownCmd)
	}

	return handler(ctx, args)
}
