> **NOTE:** This repository contains code for my own personal use and I make no
  guarantees regarding suitability, compatibility, etc. If you make use of this
	work anyway, pull requests are welcome but I also don't guarantee I will merge
	or even address them.

The `probe` package provides a slim interface for adding liveness and readiness
checks to a service.

## Usage

```
go get github.com/haleyrc/probe
```

### Defaults

The following examples demonstrate how to use `probe` with the two most common
muxes for me. The process is nearly identical for both, but since the `gorilla`
libraries deviate from the standard library in a somewhat obnoxious way, there
is additional conversion step to make it all work. That said, I've provided an
adapter to make it easier to get going.

Regardless, liveness and readiness checks will be available at `/livez` and
`/readyz` respectively.

**`http.ServeMux`**

```go
import (
	"net/http"

	"github.com/haleyrc/probe"
	"github.com/haleyrc/probe/adapters/gorilla"
)

mux := http.NewServeMux()
var p probe.Probe
p.RegisterDefaults(mux)
```

**`gorilla/mux.Router`**

```go
import (
	"github.com/gorilla/mux"
	"github.com/haleyrc/probe"
	"github.com/haleyrc/probe/adapters/gorilla"
)

mux := mux.NewRouter()
var p probe.Probe
p.RegisterDefaults(gorilla.Router(mux))
```

### Manual

If you prefer to have more control over where your checks are mounted, you can
use the handler functions directly:

```go
var p probe.Probe
mux.HandlerFunc("/myreadyz", p.ReadyzHandler)
```
