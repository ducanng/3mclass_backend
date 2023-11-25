package handler

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/ducanng/3mclass_backend/pkg/httputil"
)

type Handler interface {
	Register(r chi.Router)
}

type HandlerAggregator interface {
	RegisterAll()
	AddHandler(Handler, chi.Router)
}

type handlerAggregator struct {
	handlerMap map[Handler]chi.Router
}

func NewHandlerAggregator() HandlerAggregator {
	return &handlerAggregator{
		handlerMap: make(map[Handler]chi.Router),
	}
}

func (ha *handlerAggregator) RegisterAll() {
	for handler, router := range ha.handlerMap {
		handler.Register(router)
	}
}

func (ha *handlerAggregator) AddHandler(h Handler, r chi.Router) {
	ha.handlerMap[h] = r
}

func IndexHandler(w http.ResponseWriter, h *http.Request) {
	httputil.WriteJSONMessage(w, 200, "Ok!")
}
