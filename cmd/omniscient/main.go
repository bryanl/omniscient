package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/bryanl/omniscient"
	"github.com/kouhin/envflag"
)

func main() {
	var (
		redisAddr = flag.String("omniscient-redis-addr", "localhost:6379", "redis address")
		httpAddr  = flag.String("omniscient-http-addr", ":8080", "http server address")
	)
	envflag.Parse()

	rc, err := omniscient.NewRedisClient(*redisAddr)
	if err != nil {
		log.Fatalf("unable to create redis client: %v", err)
	}

	nr, err := omniscient.NewRedisNoteRepository(omniscient.RedisClientOption(rc))
	if err != nil {
		log.Fatalf("unable to create note repository: %v", err)
	}

	app, err := omniscient.NewApp(omniscient.AppNoteRepository(nr))
	if err != nil {
		log.Fatalf("unable to create app: %v", err)
	}

	http.Handle("/", app.Mux)
	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
