# toolbox

[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/go-macaron/toolbox/Go?logo=github&style=for-the-badge)](https://github.com/go-macaron/toolbox/actions?query=workflow%3AGo)
[![codecov](https://img.shields.io/codecov/c/github/go-macaron/toolbox/master?logo=codecov&style=for-the-badge)](https://codecov.io/gh/go-macaron/toolbox)
[![GoDoc](https://img.shields.io/badge/GoDoc-Reference-blue?style=for-the-badge&logo=go)](https://pkg.go.dev/github.com/go-macaron/toolbox?tab=doc)
[![Sourcegraph](https://img.shields.io/badge/view%20on-Sourcegraph-brightgreen.svg?style=for-the-badge&logo=sourcegraph)](https://sourcegraph.com/github.com/go-macaron/toolbox)

Middleware toolbox provides health chcek, pprof, profile and statistic services for [Macaron](https://github.com/go-macaron/macaron).

[API Reference](https://gowalker.org/github.com/go-macaron/toolbox)

### Installation

	go get github.com/go-macaron/toolbox
	
## Usage

```go
// main.go
import (
	"gopkg.in/macaron.v1"
	"github.com/go-macaron/toolbox"
)

func main() {
  	m := macaron.Classic()
  	m.Use(toolbox.Toolboxer(m))
	m.Run()
}
```

Open your browser and visit `http://localhost:4000/debug` to see the effects.

## Options

`toolbox.Toolboxer` comes with a variety of configuration options:

```go
type dummyChecker struct {
}

func (dc *dummyChecker) Desc() string {
	return "Dummy checker"
}

func (dc *dummyChecker) Check() error {
	return nil
}

// ...
m.Use(toolbox.Toolboxer(m, toolbox.Options{
	URLPrefix:			"/debug",			// URL prefix for toolbox dashboard
	HealthCheckURL:		"/healthcheck", 	// URL for health check request
	HealthCheckers: []HealthChecker{
		new(dummyChecker),
	},										// Health checkers
	HealthCheckFuncs: []*toolbox.HealthCheckFuncDesc{
		&toolbox.HealthCheckFuncDesc{
			Desc: "Database connection",
			Func: func() error { return "OK" },
		},
	},										// Health check functions
	DisableDebug:		false,				// Turns off all debug functionality when true
	PprofURLPrefix:		"/debug/pprof/", 	// URL prefix of pprof
	ProfileURLPrefix:	"/debug/profile/", 	// URL prefix of profile
	ProfilePath:		"profile",			// Path store profile files
}))
// ...
```

## Route Statistic

Toolbox also comes with a route call statistic functionality:

```go
import (
	"os"
	"time"
	//...
	"github.com/go-macaron/toolbox"
)

func main() {
	//...
	m.Get("/", func(t toolbox.Toolbox) {
		start := time.Now()
		
		// Other operations.
		
		t.AddStatistics("GET", "/", time.Since(start))
	})
	
	m.Get("/dump", func(t toolbox.Toolbox) {
		t.GetMap(os.Stdout)
	})
}
```

Output take from test:

```
+---------------------------------------------------+------------+------------------+------------------+------------------+------------------+------------------+
| Request URL                                       | Method     | Times            | Total Used(s)    | Max Used(μs)     | Min Used(μs)     | Avg Used(μs)     |
+---------------------------------------------------+------------+------------------+------------------+------------------+------------------+------------------+
| /api/user                                         | POST       |                2 |         0.000122 |       120.000000 |         2.000000 |        61.000000 |
| /api/user                                         | GET        |                1 |         0.000013 |        13.000000 |        13.000000 |        13.000000 |
| /api/user                                         | DELETE     |                1 |         0.000001 |         1.400000 |         1.400000 |         1.400000 |
| /api/admin                                        | POST       |                1 |         0.000014 |        14.000000 |        14.000000 |        14.000000 |
| /api/user/unknwon                                 | POST       |                1 |         0.000012 |        12.000000 |        12.000000 |        12.000000 |
+---------------------------------------------------+------------+------------------+------------------+------------------+------------------+------------------+
```

## License

This project is under Apache v2 License. See the [LICENSE](LICENSE) file for the full license text.