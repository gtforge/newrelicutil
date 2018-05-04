package newrelicutil_test

import (
	"net/http"

	"github.com/gettaxi/newrelicutil"
	"github.com/newrelic/go-agent"
)

func ExampleWrapHandler() {
	config := newrelic.NewConfig("app_test", "")
	config.Enabled = false
	nrapp, _ := newrelic.NewApplication(config)

	handler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {})
	wrapped_handler := newrelicutil.WrapHandler(nrapp, "HandlerName", handler)

	mux := http.NewServeMux()
	mux.Handle("/some/path", wrapped_handler)
}
