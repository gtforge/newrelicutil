package newrelicutil_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gettaxi/newrelicutil"
	"github.com/newrelic/go-agent"
)

func TestTransaction(t *testing.T) {
	ctx := context.Background()

	if want, have := newrelic.Transaction(nil), newrelicutil.Transaction(ctx); want != have {
		t.Errorf("want %+v, have %+v", want, have)
	}
}

func TestWithTransaction(t *testing.T) {
	nrapp := newApp()
	txn := nrapp.StartTransaction("foo", nil, nil)
	ctx := context.Background()

	ctx = newrelicutil.WithTransaction(ctx, txn)
	if want, have := txn, newrelicutil.Transaction(ctx); want != have {
		t.Errorf("want %+v, have %+v", want, have)
	}
}

func TestWrapHandler(t *testing.T) {
	nrapp := newApp()
	h := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if have := newrelicutil.Transaction(r.Context()); nil == have {
			t.Errorf("want txn, have %+v", have)
		}
	})
	wh := newrelicutil.WrapHandler(nrapp, "foo", h)
	req, _ := http.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()
	wh.ServeHTTP(resp, req)
}

func newApp() newrelic.Application {
	config := newrelic.NewConfig("app_test", "")
	config.Enabled = false
	nrapp, _ := newrelic.NewApplication(config)
	return nrapp
}
