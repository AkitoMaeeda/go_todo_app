package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/AkitoMaeeda/go_todo_app/config"
	"golang.org/x/sync/errgroup"
)

func run(ctx context.Context) error {

	ctx, stop := signal.NotifyContext(ctx, os.Interttupt, syscall.SIGTERM)
	defar stop()

	cfg, err := config.New()
	if err != nil {
		return err
	}

	log.Printf("利用するポート :%d\n", cfg.Port)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to kisten port %d: \n%v", cfg.Port, err)
	}
	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with: %v", url)

	s := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(5 * time.Second)
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if err := s.Serve(l); err != nil &&
			err != http.ErrServerClosed {
			log.Printf("failed to close ; %v", err)
			return err
		}
		return nil
	})

	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {

		log.Printf("failed to shutdown: %+v", err)
	}

	return eg.Wait()

}
