package main

import (
	"flag"
	"net/http"

	"omniscient"

	log "github.com/Sirupsen/logrus"
	"github.com/getsentry/raven-go"
	"github.com/kouhin/envflag"
)

func main() {
	var (
		redisAddr = flag.String("omniscient-redis-addr", "localhost:6379", "redis address")
		httpAddr  = flag.String("omniscient-http-addr", ":8080", "http server address")
		sentryURL = flag.String("omniscient-sentry-url", "", "sentry url")
	)
	envflag.Parse()

	if len(*sentryURL) > 0 {
		log.Info("configuring sentry")
		raven.SetDSN(*sentryURL)
	}

	rc, err := omniscient.NewRedisClient(*redisAddr)
	if err != nil {
		log.Fatalf("unable to create redis client: %v", err)
	}

	redisPingCheck := func() bool {
		_, err := rc.Ping()
		if err != nil {
			return false
		}
		return true
	}

	nr, err := omniscient.NewRedisNoteRepository(omniscient.RedisClientOption(rc))
	if err != nil {
		log.Fatalf("unable to create note repository: %v", err)
	}

	health, err := omniscient.NewHealth(
		omniscient.HealthCheckOption(redisPingCheck))

	app, err := omniscient.NewApp(
		omniscient.AppNoteRepository(nr),
		omniscient.AppHealth(health))
	if err != nil {
		log.Fatalf("unable to create app: %v", err)
	}

	http.Handle("/", app.Mux)
	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
