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

func TestSegment(t *testing.T) {
	ctx := context.Background()

	want := newrelic.Segment{}
	if have := newrelicutil.Segment(ctx); want != have {
		t.Errorf("want %+v, have %+v", want, have)
	}
}

func TestWithSegment(t *testing.T) {
	nrapp := newApp()
	txn := nrapp.StartTransaction("foo", nil, nil)
	sgm := newrelic.StartSegment(txn, "bar")
	ctx := context.Background()

	ctx = newrelicutil.WithSegment(ctx, sgm)
	if want, have := sgm, newrelicutil.Segment(ctx); want != have {
		t.Errorf("want %+v, have %+v", want, have)
	}
}

func TestExternalSegment(t *testing.T) {
	ctx := context.Background()

	want := newrelic.ExternalSegment{}
	if have := newrelicutil.ExternalSegment(ctx); want != have {
		t.Errorf("want %+v, have %+v", want, have)
	}
}

func TestWithExternalSegment(t *testing.T) {
	nrapp := newApp()
	txn := nrapp.StartTransaction("foo", nil, nil)
	sgm := newrelic.StartExternalSegment(txn, nil)
	ctx := context.Background()

	ctx = newrelicutil.WithExternalSegment(ctx, sgm)
	if want, have := sgm, newrelicutil.ExternalSegment(ctx); want != have {
		t.Errorf("want %+v, have %+v", want, have)
	}
}

func TestDatastoreSegment(t *testing.T) {
	ctx := context.Background()

	want := newrelic.DatastoreSegment{}
	if have := newrelicutil.DatastoreSegment(ctx); want.StartTime != have.StartTime {
		t.Errorf("want %#v, have %#v", want, have)
	}
}

func TestWithDatastoreSegment(t *testing.T) {
	nrapp := newApp()
	txn := nrapp.StartTransaction("foo", nil, nil)
	sgm := newrelic.DatastoreSegment{
		StartTime:  newrelic.StartSegmentNow(txn),
		Product:    newrelic.DatastoreMySQL,
		Collection: "my_table",
		Operation:  "SELECT",
	}
	ctx := context.Background()

	ctx = newrelicutil.WithDatastoreSegment(ctx, sgm)
	if want, have := sgm, newrelicutil.DatastoreSegment(ctx); want.StartTime != have.StartTime {
		t.Errorf("want %#v, have %#v", want, have)
	}
}

func newApp() newrelic.Application {
	config := newrelic.NewConfig("app_test", "")
	config.Enabled = false
	nrapp, _ := newrelic.NewApplication(config)
	return nrapp
}
