// Package probe provides a slim interface for adding liveness and readiness
// checks to a service. In order to be as broadly useful as possible, the
// Kubernetes standard has been adopted for the default paths, but strict
// adherance to k8s practices is NOT guaranteed.
//
// A Probe provides a method for mounting the liveness and readiness endpoints
// at their default location for any mux/router/etc. that follows the
// http.ServeMux pattern of including a HandleFunc method.
//
// If you prefer having more control over where your checks are mounted, you are
// still free to mount the handlers manually using your preferred mux-like.
//
// Concurrency
//
// Internally, the Probe maintains a simple ready/not ready state to determine
// the appropriate status code for the readiness check. This state is NOT
// protected against concurrent writes, so attempting to use a Probe in a
// concurrent way (this only applies to setting the ready state; reads are safe
// by default) is generally nonsensical. If you need to mount liveness and
// readiness checks for multiple services in the same application, they should
// all have their own Probe.
package probe

import "net/http"

// The following paths are used to register the handlers when RegisterDefaults
// is called. These paths are compatible with the Kubernetes standard since that
// is widely understood even outside of that ecosystem, including in Prometheus.
const (
	DefaultLivezPath  = "/livez"
	DefaultReadyzPath = "/readyz"
)

// Router describes a common interface for use when registering the default
// liveness and readiness endpoints. If your mux of choice does not satisfy this
// interface by default, some adapters are provided to convert to one that does
// or you can choose to mount the handlers manually.
type Router interface {
	HandleFunc(path string, h func(http.ResponseWriter, *http.Request))
}

// Probe provides a simple interface for handling liveness and readiness in a
// consistent way. A Probe is specifically designed for managing
// liveness/readiness for a single service and is not suitable for concurrent
// access.
type Probe struct {
	ready bool
}

// LivezHandler is intended to be used to signify that a service "up", though it
// may not be ready to handle requests. This endpoint always set a 200 status
// code, since a dead service is indicated by the inability to respond to any
// requests.
func (p *Probe) LivezHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// NotReady puts the Probe into a not ready state. This has the add-on effect
// that the readiness endpoint will begin returning a 503 status code. This is
// the same behavior exhibited by a zero-value Probe. This can be useful if your
// service is setup to continue serving in-flight requests in response to a
// shutdown.
func (p *Probe) NotReady() {
	p.ready = false
}

// Ready puts the Probe into a ready state. This has the add-on effect that the
// readiness endpoint will begin returning a 200 status code.
func (p *Probe) Ready() {
	p.ready = true
}

// ReadyzHandler is intended to be used to signify that a service is ready to
// respond to requests. If the Probe is not ready, a 503 (service unavailable)
// code is returned. Once the Probe is ready, a 200 is returned instead.
func (p *Probe) ReadyzHandler(w http.ResponseWriter, r *http.Request) {
	if p.ready {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
}

// RegisterDefaults mounts the liveness and readiness handlers at the default
// paths indicated by the DefaultLivezPath and DefaultReadyzPath constants.
//
// If you prefer mounting these handlers at alternative locations, you can skip
// this method and use HandleFunc directly on your mux/router/etc.
func (p *Probe) RegisterDefaults(m Router) {
	m.HandleFunc(DefaultLivezPath, p.LivezHandler)
	m.HandleFunc(DefaultReadyzPath, p.ReadyzHandler)
}
