// Package nrgorilla provides New Relic instruments for gorilla/mux tool.
package nrgorilla

import (
	"strings"

	"github.com/gettaxi/newrelicutil"
	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent"
)

// InstrumentRoutes adds instrumentation to a router.
// This must be used after the routes have been added to the router.
func InstrumentRoutes(r *mux.Router, app newrelic.Application) *mux.Router {
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		h := newrelicutil.WrapHandler(app, RouteName(route), route.GetHandler())
		route.Handler(h)

		return nil
	})
	if r.NotFoundHandler != nil {
		r.NotFoundHandler = newrelicutil.WrapHandler(app, "NotFoundHandler", r.NotFoundHandler)
	}
	return r
}

// RouteName returns the name that would be used as transaction name in New Relic
func RouteName(route *mux.Route) string {
	if nil == route {
		return ""
	}
	if n := route.GetName(); n != "" {
		return n
	}
	if n, _ := route.GetPathTemplate(); n != "" {
		ms, _ := route.GetMethods()
		return strings.Join(ms, "/") + " " + n
	}
	n, _ := route.GetHostTemplate()
	return n
}
