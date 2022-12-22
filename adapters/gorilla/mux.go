// Package gorilla provides an adapter for using a gorilla/mux.Router to a
// probe.Router for use when registering the default liveness and readiness
// handlers.
package gorilla

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/haleyrc/probe"
)

// Router converts a gorilla/mux.Router to an interface that can be used for
// registering the default liveness and readiness handlers.
func Router(gr *mux.Router) probe.Router {
	return &router{gr: gr}
}

type router struct {
	gr *mux.Router
}

func (r *router) HandleFunc(path string, h func(http.ResponseWriter, *http.Request)) {
	r.gr.HandleFunc(path, h)
}
