# Xenlight Exporter
Prometheus exporter for [xen](https://xenproject.org/) using `libxl` go bindings.

## Installation
You can build the latest version using Go v1.11+ via `go get`:
```
go get -u github.com/Gandi/xenlight_exporter
```

You need `xen` headers as well as `yajl` headers to be able to compile xenlight go
bindings
