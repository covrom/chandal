package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/covrom/chandal/api"
	"github.com/covrom/chandal/dal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	db, err := dal.MigrateDB("postgres://postgres:123@localhost:5432/db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	wg := &sync.WaitGroup{}

	dal.GoGetUsers(ctx, db, 10, 100, wg)
	dal.GoCreateUser(ctx, db, 5, 10, wg)

	a := api.NewApi()

	server := &http.Server{
		Addr:        ":8000",
		Handler:     a,
		BaseContext: func(l net.Listener) context.Context { return ctx },
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		server.ListenAndServe()
	}()
	go func() {
		defer wg.Done()
		<-ctx.Done()
		server.Shutdown(ctx)
	}()
	wg.Wait()
}
