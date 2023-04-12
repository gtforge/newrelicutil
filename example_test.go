package newrelicutil_test

import (
	"net/http"
	"os"

	"github.com/newrelic/go-agent/v3/newrelic"

	"github.com/gtforge/newrelicutil/v2"
)

func ExampleWrapHandler() {
	nrapp, _ := newrelic.NewApplication(
		newrelic.ConfigAppName("app_test"),
		newrelic.ConfigEnabled(false),
		newrelic.ConfigInfoLogger(os.Stdout),
	)

	handler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {})
	wrappedHandler := newrelicutil.WrapHandler(nrapp, "HandlerName", handler)

	mux := http.NewServeMux()
	mux.Handle("/some/path", wrappedHandler)
}
