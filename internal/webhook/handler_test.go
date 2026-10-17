package webhook_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/project0/external-dns-libdns-webhook/internal/externaldns"
	"github.com/project0/external-dns-libdns-webhook/internal/webhook"
)

const mediaType = "application/external.dns.webhook+json;version=1"

type fakeProvider struct {
	err     error
	changes *externaldns.Changes
}

func (f *fakeProvider) Records(_ context.Context) (externaldns.Endpoints, error) {
	if f.err != nil {
		return nil, f.err
	}

	return externaldns.Endpoints{externaldns.NewEndpointWithTTL("www.example.com", "A", 300, "192.0.2.1")}, nil
}

func (f *fakeProvider) ApplyChanges(_ context.Context, changes *externaldns.Changes) error {
	f.changes = changes

	return f.err
}

func (f *fakeProvider) AdjustEndpoints(_ context.Context, endpoints externaldns.Endpoints) (externaldns.Endpoints, error) {
	return endpoints, f.err
}

func (f *fakeProvider) GetDomainFilter(_ context.Context) externaldns.DomainFilter {
	return externaldns.DomainFilter{Include: []string{"example.com"}}
}

type response struct {
	status      int
	contentType string
	body        string
}

func do(t *testing.T, provider externaldns.Provider, method, path, body string) response {
	t.Helper()

	handler := webhook.New(provider, webhook.Config{Listen: "localhost:0", ReadTimeout: time.Second, WriteTimeout: time.Second}).Handler()
	req := httptest.NewRequestWithContext(t.Context(), method, path, strings.NewReader(body))
	req.Header.Set("Accept", mediaType)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	out, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	return response{status: res.StatusCode, contentType: res.Header.Get("Content-Type"), body: string(out)}
}

func TestNegotiate(t *testing.T) {
	t.Parallel()

	res := do(t, &fakeProvider{}, http.MethodGet, "/", "")
	if res.status != http.StatusOK {
		t.Fatalf("status %d", res.status)
	}

	if ct := res.contentType; ct != mediaType {
		t.Errorf("content type %q", ct)
	}

	if want := `{"include":["example.com"]}`; res.body != want {
		t.Errorf("got %s, want %s", res.body, want)
	}
}

func TestUnknownPath(t *testing.T) {
	t.Parallel()

	res := do(t, &fakeProvider{}, http.MethodGet, "/unknown", "")
	if res.status != http.StatusNotFound {
		t.Errorf("status %d", res.status)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	t.Parallel()

	for _, path := range []string{"/records", "/adjustendpoints"} {
		res := do(t, &fakeProvider{}, http.MethodPut, path, "")
		if res.status != http.StatusMethodNotAllowed {
			t.Errorf("%s: status %d", path, res.status)
		}
	}
}

func TestGetRecords(t *testing.T) {
	t.Parallel()

	res := do(t, &fakeProvider{}, http.MethodGet, "/records", "")
	if res.status != http.StatusOK || res.contentType != mediaType {
		t.Fatalf("status %d content type %q", res.status, res.contentType)
	}

	if !strings.Contains(res.body, `"dnsName":"www.example.com"`) {
		t.Errorf("unexpected body %s", res.body)
	}

	res = do(t, &fakeProvider{err: errors.New("boom")}, http.MethodGet, "/records", "")
	if res.status != http.StatusInternalServerError {
		t.Errorf("error status %d", res.status)
	}
}

func TestApplyChanges(t *testing.T) {
	t.Parallel()

	provider := &fakeProvider{}

	res := do(t, provider, http.MethodPost, "/records",
		`{"create":[{"dnsName":"www.example.com","recordType":"A","targets":["192.0.2.1"]}]}`)
	if res.status != http.StatusNoContent {
		t.Fatalf("status %d", res.status)
	}

	if provider.changes == nil || len(provider.changes.Create) != 1 {
		t.Errorf("changes not passed: %+v", provider.changes)
	}

	res = do(t, &fakeProvider{}, http.MethodPost, "/records", "not json")
	if res.status != http.StatusBadRequest {
		t.Errorf("bad request status %d", res.status)
	}

	res = do(t, &fakeProvider{err: errors.New("boom")}, http.MethodPost, "/records", "{}")
	if res.status != http.StatusInternalServerError {
		t.Errorf("error status %d", res.status)
	}
}

func TestAdjustEndpoints(t *testing.T) {
	t.Parallel()

	res := do(t, &fakeProvider{}, http.MethodPost, "/adjustendpoints",
		`[{"dnsName":"www.example.com","recordType":"A","targets":["192.0.2.1"]}]`)
	if res.status != http.StatusOK || res.contentType != mediaType {
		t.Fatalf("status %d content type %q", res.status, res.contentType)
	}

	if !strings.HasPrefix(res.body, `[{"dnsName":"www.example.com"`) {
		t.Errorf("unexpected body %s", res.body)
	}

	res = do(t, &fakeProvider{err: errors.New("boom")}, http.MethodPost, "/adjustendpoints", "[]")
	if res.status != http.StatusInternalServerError {
		t.Errorf("error status %d", res.status)
	}

	if strings.HasPrefix(res.body, "[") {
		t.Errorf("error response must not contain endpoints: %s", res.body)
	}
}
