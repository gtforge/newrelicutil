package newrelicutil_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/gettaxi/newrelicutil"
	"github.com/newrelic/go-agent"
)

func TestTransaction(t *testing.T) {
	want := newrelic.Transaction(nil)
	have := newrelicutil.Transaction(context.Background())
	assert.Equal(t, want, have)
}

func TestWithTransaction(t *testing.T) {
	txn := newApp(t).StartTransaction("foo", nil, nil)
	ctx := newrelicutil.WithTransaction(context.Background(), txn)
	assert.Equal(t, txn, newrelicutil.Transaction(ctx))
}

func TestWrapHandler(t *testing.T) {
	h := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		assert.NotNil(t, newrelicutil.Transaction(r.Context()))
	})
	wh := newrelicutil.WrapHandler(newApp(t), "foo", h)
	req, _ := http.NewRequest("GET", "/", nil)
	wh.ServeHTTP(httptest.NewRecorder(), req)
}

func TestSegment(t *testing.T) {
	want := &newrelic.Segment{}
	have := newrelicutil.Segment(context.Background())
	assert.Equal(t, want, have)
}

func TestWithSegment(t *testing.T) {
	txn := newApp(t).StartTransaction("foo", nil, nil)
	want := newrelic.StartSegment(txn, "bar")
	ctx := newrelicutil.WithSegment(context.Background(), want)
	have := newrelicutil.Segment(ctx)
	assert.Equal(t, want, have)
}

func TestExternalSegment(t *testing.T) {
	want := &newrelic.ExternalSegment{}
	have := newrelicutil.ExternalSegment(context.Background())
	assert.Equal(t, want, have)
}

func TestWithExternalSegment(t *testing.T) {
	txn := newApp(t).StartTransaction("foo", nil, nil)
	want := newrelic.StartExternalSegment(txn, nil)
	ctx := context.Background()

	ctx = newrelicutil.WithExternalSegment(ctx, want)
	assert.Equal(t, want, newrelicutil.ExternalSegment(ctx))
}

func TestDatastoreSegment(t *testing.T) {
	want := newrelic.DatastoreSegment{}
	have := newrelicutil.DatastoreSegment(context.Background())
	assert.Equal(t, want.StartTime, have.StartTime)
}

func TestWithDatastoreSegment(t *testing.T) {
	txn := newApp(t).StartTransaction("foo", nil, nil)
	sgm := &newrelic.DatastoreSegment{
		StartTime:  newrelic.StartSegmentNow(txn),
		Product:    newrelic.DatastoreMySQL,
		Collection: "my_table",
		Operation:  "SELECT",
	}
	ctx := newrelicutil.WithDatastoreSegment(context.Background(), sgm)
	assert.Equal(t, sgm.StartTime, newrelicutil.DatastoreSegment(ctx).StartTime)
}

func newApp(t *testing.T) newrelic.Application {
	config := newrelic.NewConfig("app_test", "")
	config.Enabled = false
	nrapp, err := newrelic.NewApplication(config)
	require.NoError(t, err)
	return nrapp
}
