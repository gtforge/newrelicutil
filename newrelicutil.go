package newrelicutil

import (
	"context"
	"net/http"

	"github.com/newrelic/go-agent"
)

type key int

const (
	transaction key = iota
	segment
	externalSegment
)

// Transaction returns the New Relic Transaction object from context.
func Transaction(ctx context.Context) newrelic.Transaction {
	if txn, ok := ctx.Value(transaction).(newrelic.Transaction); ok {
		return txn
	}
	return nil
}

// WithTransaction puts the New Relic Transaction object to the given context
// and returns the new context.
func WithTransaction(ctx context.Context, txn newrelic.Transaction) context.Context {
	return context.WithValue(ctx, transaction, txn)
}

// WrapHandler return the given http handler that is wrapped to New Relic Transaction.
// Current New Relic Transaction is placed in the context.
func WrapHandler(app newrelic.Application, name string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		txn := app.StartTransaction(name, w, r)
		defer txn.End()

		ctx := WithTransaction(r.Context(), txn)

		handler.ServeHTTP(txn, r.WithContext(ctx))
	})
}
