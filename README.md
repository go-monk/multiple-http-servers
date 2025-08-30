In this post, we’ll build a command-line Go program that spins up multiple HTTP servers concurrently, each listening on its own port and serving a simple response. Usually you need just one server so it may seem like an odd thing to do at first glance. However it’s a practical way to play with and understand concurrency and the `net/http` standard library.

We’ll write a Go program named `multiserv` that:

* Accepts a number `n` from the command line
* Starts `n` HTTP servers, each listening on localhost port above 1023
* Assigns each server its own `ServeMux` and a handler that responds with its server number

We create separate `ServeMux` instances for each server:

```go
func newMuxers(n int) []*http.ServeMux {
	muxers := make([]*http.ServeMux, n)
	for i := range n {
		muxers[i] = http.NewServeMux()
	}
	return muxers
}
```

Each handler is just an `int`, but by defining a method on it, it satisfies the `http.Handler` interface:

```go
type handler int

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello from handler %d", h)
}
```

This is a neat trick in Go — using a primitive type (like `int`) and giving it behavior via methods.

Run the program and test each server:

```bash
go run multiserv.go 3
```

Then in another terminal:

```bash
$ curl localhost:1024
hello from handler 0
$ curl localhost:1025
hello from handler 1
$ curl localhost:1026
hello from handler 2
```

Running multiple HTTP servers can be useful in:

* Load testing: Simulate multiple nodes with lightweight stubs.
* Multi-tenancy demos: Each port could represent a different tenant.
* Local sharding: Simulate partitioned services.
* Teaching concurrency: Excellent for understanding goroutines and handler isolation.

This tiny `multiserv` tool may not seem like much, but it encapsulates several core ideas in Go:

* Concurrency with goroutines
* Safe closure handling
* HTTP routing with `ServeMux`
* Interface satisfaction with custom types

It’s a solid building block for more advanced server architectures or just a neat utility to throw in your developer toolbox.
