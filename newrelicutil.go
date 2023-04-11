// Package newrelicutil provides a set of helpers for working with objects of New Relic Go Agent.
package newrelicutil

import (
	"context"
	"net/http"

	"github.com/newrelic/go-agent/v3/newrelic"
)

type key int

const (
	segment key = iota
	externalSegment
	datastoreSegment
)

// Transaction returns the New Relic Transaction object from context.
// Returns nil if the context does not contain an object.
func Transaction(ctx context.Context) *newrelic.Transaction {
	return newrelic.FromContext(ctx)
}

// WithTransaction puts the New Relic Transaction object to the given context
// and returns the new context.
func WithTransaction(ctx context.Context, txn *newrelic.Transaction) context.Context {
	return newrelic.NewContext(ctx, txn)
}

// WrapHandler returns the given http handler that is wrapped to New Relic Transaction.
// Current New Relic Transaction is placed in the context.
func WrapHandler(app *newrelic.Application, name string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		txn := app.StartTransaction(name)
		defer txn.End()

		txn.SetWebRequestHTTP(r)
		w = txn.SetWebResponse(w)

		ctx := newrelic.NewContext(r.Context(), txn)

		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Segment returns the New Relic Segment object from context.
func Segment(ctx context.Context) *newrelic.Segment {
	if sgm, ok := ctx.Value(segment).(*newrelic.Segment); ok {
		return sgm
	}
	return &newrelic.Segment{}
}

// WithSegment puts the New Relic Segment object to the given context
// and returns the new context.
func WithSegment(ctx context.Context, sgm *newrelic.Segment) context.Context {
	return context.WithValue(ctx, segment, sgm)
}

// ExternalSegment returns the New Relic ExternalSegment object from context.
func ExternalSegment(ctx context.Context) *newrelic.ExternalSegment {
	if sgm, ok := ctx.Value(externalSegment).(*newrelic.ExternalSegment); ok {
		return sgm
	}
	return &newrelic.ExternalSegment{}
}

// WithExternalSegment puts the New Relic ExternalSegment object to the given context
// and returns the new context.
func WithExternalSegment(ctx context.Context, sgm *newrelic.ExternalSegment) context.Context {
	return context.WithValue(ctx, externalSegment, sgm)
}

// DatastoreSegment returns the New Relic DatastoreSegment object from context.
func DatastoreSegment(ctx context.Context) *newrelic.DatastoreSegment {
	if sgm, ok := ctx.Value(datastoreSegment).(*newrelic.DatastoreSegment); ok {
		return sgm
	}
	return &newrelic.DatastoreSegment{}
}

// WithDatastoreSegment puts the New Relic ExternalSegment object to the given context
// and returns the new context.
func WithDatastoreSegment(ctx context.Context, sgm *newrelic.DatastoreSegment) context.Context {
	return context.WithValue(ctx, datastoreSegment, sgm)
}
