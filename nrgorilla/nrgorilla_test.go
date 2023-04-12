package nrgorilla_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/newrelic"

	"github.com/gtforge/newrelicutil/v2"
	"github.com/gtforge/newrelicutil/v2/nrgorilla"
)

func TestInstrumentRoutes(t *testing.T) {
	nrapp, _ := newrelic.NewApplication(
		newrelic.ConfigAppName("app_test"),
		newrelic.ConfigEnabled(false),
		newrelic.ConfigInfoLogger(os.Stdout),
	)

	h := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if have := newrelicutil.Transaction(r.Context()); have == nil {
			t.Errorf("want txn, have %+v", have)
		}
	})
	r := mux.NewRouter()
	r.Handle("/api/v1/users", h).Methods("GET")
	r.NotFoundHandler = h
	r = nrgorilla.InstrumentRoutes(r, nrapp)

	req, _ := http.NewRequestWithContext(context.Background(), "GET", "/api/v1/users", http.NoBody)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	req, _ = http.NewRequestWithContext(context.Background(), "GET", "/404", http.NoBody)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
}

func TestRouteName(t *testing.T) {
	r := mux.NewRouter()
	h := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {})

	tt := []struct {
		route *mux.Route
		exp   string
	}{
		{
			route: nil,
			exp:   "",
		}, {
			route: r.Handle("/api/", h).Methods("GET").Name("FOO"),
			exp:   "FOO",
		}, {
			route: r.Handle("/api/v1/users", h).Methods("GET"),
			exp:   "GET /api/v1/users",
		}, {
			route: r.Methods("GET"),
			exp:   "",
		},
	}

	for _, tc := range tt {
		if want, have := tc.exp, nrgorilla.RouteName(tc.route); want != have {
			t.Errorf("want %+v, have %+v", want, have)
		}
	}
}
