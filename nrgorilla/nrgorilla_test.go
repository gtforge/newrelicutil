package nrgorilla_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gettaxi/newrelicutil/v2"
	"github.com/gettaxi/newrelicutil/v2/nrgorilla"
	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func TestInstrumentRoutes(t *testing.T) {
	nrapp, _ := newrelic.NewApplication(
		newrelic.ConfigAppName("app_test"),
		newrelic.ConfigEnabled(false),
		newrelic.ConfigInfoLogger(os.Stdout),
	)

	h := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if have := newrelicutil.Transaction(r.Context()); nil == have {
			t.Errorf("want txn, have %+v", have)
		}
	})
	r := mux.NewRouter()
	r.Handle("/api/v1/users", h).Methods("GET")
	r.NotFoundHandler = h
	r = nrgorilla.InstrumentRoutes(r, nrapp)

	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	req, _ = http.NewRequest("GET", "/404", nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
}

func TestRouteName(t *testing.T) {
	r := mux.NewRouter()
	h := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {})

	var tt = []struct {
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
