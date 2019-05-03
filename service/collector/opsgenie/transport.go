package opsgenie

import (
	"fmt"
	"net/http"
)

type opsgenieTransport struct {
	transport http.RoundTripper
	key       string
}

func (t opsgenieTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("Authorization", fmt.Sprintf("GenieKey %s", t.key))
	return t.transport.RoundTrip(r)
}
