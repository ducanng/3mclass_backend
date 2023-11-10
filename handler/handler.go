package handler

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/ducanng/no-name/pkg/httputil"
)

type Handler interface {
	Register(r *chi.Mux)
}

type HandlerAggregator interface {
	RegisterAll()
	AddHandler(Handler, *chi.Mux)
}

type handlerAggregator struct {
	handlerMap map[Handler]*chi.Mux
}

func NewHandlerAggregator() HandlerAggregator {
	return &handlerAggregator{
		handlerMap: make(map[Handler]*chi.Mux),
	}
}

func (ha *handlerAggregator) RegisterAll() {
	for handler, router := range ha.handlerMap {
		handler.Register(router)
	}
}

func (ha *handlerAggregator) AddHandler(h Handler, r *chi.Mux) {
	ha.handlerMap[h] = r
}

func IndexHandler(w http.ResponseWriter, h *http.Request) {
	httputil.WriteJSONMessage(w, 200, "Ok!")
}
