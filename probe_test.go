package probe_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/haleyrc/probe"
	"github.com/haleyrc/probe/adapters/gorilla"
)

func TestLivezHandler(t *testing.T) {
	var p probe.Probe
	assertOK(t, get("/", p.LivezHandler))
}

func TestReadyzHandler(t *testing.T) {
	var p probe.Probe
	assertUnavailable(t, get("/", p.ReadyzHandler))

	p.Ready()
	assertOK(t, get("/", p.ReadyzHandler))

	p.NotReady()
	assertUnavailable(t, get("/", p.ReadyzHandler))
}

func TestRegisterDefaults(t *testing.T) {
	t.Run("http.ServeMux", func(t *testing.T) {
		var mux = http.NewServeMux()
		var p probe.Probe

		p.RegisterDefaults(mux)

		assertOK(t, get(probe.DefaultLivezPath, mux.ServeHTTP))
		assertUnavailable(t, get(probe.DefaultReadyzPath, mux.ServeHTTP))
	})

	t.Run("gorilla/mux.Router", func(t *testing.T) {
		mux := mux.NewRouter()
		var p probe.Probe

		p.RegisterDefaults(gorilla.Router(mux))

		assertOK(t, get(probe.DefaultLivezPath, mux.ServeHTTP))
		assertUnavailable(t, get(probe.DefaultReadyzPath, mux.ServeHTTP))
	})
}

func get(path string, h http.HandlerFunc) int {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, path, nil)
	h(w, r)
	return w.Result().StatusCode
}

func assertOK(t *testing.T, got int) {
	t.Helper()
	assertStatus(t, http.StatusOK, got)
}

func assertUnavailable(t *testing.T, got int) {
	t.Helper()
	assertStatus(t, http.StatusServiceUnavailable, got)
}

func assertStatus(t *testing.T, want, got int) {
	t.Helper()
	if got != want {
		t.Errorf(
			"Expected status to be \"%d - %s\", but got \"%d - %s\".",
			want, http.StatusText(want),
			got, http.StatusText(got),
		)
	}
}
