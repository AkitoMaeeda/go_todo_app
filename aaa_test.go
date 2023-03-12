package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestRun(t *testing.T) {

	l, err := net.Listen("tcp", "localhost:0")

	if err != nil {
		t.Fatalf("failed to listen port %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return run(ctx, l)
	})

	in := "message"
	url := fmt.Sprintf("http://%s/%s", l.Addr().String(), in)

	t.Logf("try request to %q", url)

	rsp, err := http.Get(url)
	if err != nil {
		log.Printf("failed to get: %+v ", err)
	}

	defer rsp.Body.Close()

	got, err := io.ReadAll(rsp.Body)

	if err != nil {
		log.Printf("failed to read body: %v", err)
	}

	want := fmt.Sprintf("Hello;, %s!", in)

	if string(got) != want {
		log.Printf("want %q, but got %q", want, got)
	}

	cancel()

	/*if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}*/

	/*err := http.ListenAndServe(
		":18080",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "hello, %s!", r.URL.Path[1:])
		}),
	)

	if err != nil {
		fmt.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}*/
}
