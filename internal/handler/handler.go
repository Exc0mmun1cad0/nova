package handler

import (
	"fmt"
	"nova/pkg/resp"
	"strings"
	"time"

	"go.uber.org/zap"
)

type Storage interface {
	Set(key, value string, ttl time.Duration)
	Get(key string) (string, bool)
	DeleteMany(keys []string) int
}

type Handler struct {
	log     *zap.Logger
	storage Storage
	dict    map[string]handlerFunc
}

func NewHandler(log *zap.Logger, storage Storage) *Handler {
	h := &Handler{
		log:     log,
		storage: storage,
	}

	dict := map[string]handlerFunc{
		cmdPing: h.pingHandler,
		cmdEcho: h.echoHandler,

		cmdGet:    h.getHandler,
		cmdSet:    h.setHandler,
		cmdDelete: h.deleteHandler,
	}

	h.dict = dict
	return h
}

func (h *Handler) Serve(input []byte) []byte {
	args, err := resp.Decode(input)
	if err != nil {
		return resp.EncodeError(fmt.Sprintf(ErrProtocol, err.Error()))
	}

	cmd := strings.ToLower(args[0])
	handler, ok := h.dict[cmd]
	if !ok {
		return resp.EncodeError(ErrUnknownCmd)
	}

	return handler(args)
}
