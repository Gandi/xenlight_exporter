# Xenlight Exporter
Prometheus exporter for [xen](https://xenproject.org/) using `libxl` go bindings.

## Installation
You can build the latest version using Go v1.11+ via `go get`:
```
go get -u github.com/Gandi/xenlight_exporter
```

You need `xen` headers as well as `yajl` headers to be able to compile xenlight go
bindings

## Usage

```
usage: xenlight_exporter [<flags>]

Flags:
  -h, --help                Show context-sensitive help (also try --help-long and --help-man).
      --collector.domain.show-vcpus-details
                            Enable the collection of per-vcpu time
      --collector.domain    Enable the domain collector (default: enabled).
      --collector.physical  Enable the physical collector (default: enabled).
      --collector.version   Enable the version collector (default: enabled).
      --web.listen-address=":9603"
                            Address on which to expose metrics and web interface.
      --web.telemetry-path="/metrics"
                            Path under which to expose metrics.
      --log.level="info"    Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal]
      --log.format="logger:stderr"
                            Set the log target and format. Example: "logger:syslog?appname=bob&local=7" or "logger:stdout?json=true"
      --version             Show application version.
```

## Notes about golang bindings of `libxl`

The Go bindings are expected to be imported (according to Xen makefiles) using
the following import path: `golang.xenproject.org/xenlight`. However the domain
`golang.xenproject.org` doesn't exists (thus not allowing the use of Go modules)
so I chose to import the bindings from their Github mirror of the repository.

Would the situation evolve and the Xen project provide an usable import path, I
will reconsider this choice and switch to its official import path.
