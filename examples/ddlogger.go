package main

import (
	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/fault"
	"github.com/jonathana/ozzo-routing/access"
	"log"
	"net/http"
	"github.com/DataDog/datadog-go/statsd"
	"fmt"
	"github.com/jonathana/ozzo-datadog-logger"
)

func main() {
	router := routing.New()

	c, err := statsd.New("127.0.0.1:8125")
	if err != nil {
		log.Fatal(err)
	}
	// prefix every metric with the app name
	c.Namespace = "ozddLogger."
	// set the default tags
	c.Tags = append(c.Tags, "us-east-1a")
	cfg := ozDdLogger.DdLoggerConfig{
		TagProtocol: true,
		TagMethod: true,
		LogResponseCode: true,
		BaseUrl: "foo/",
	}
	ddLogger, err := ozDdLogger.ConfigureLogger(c, cfg)
	if err != nil {
		fmt.Errorf("Got error creating Datadog logger: %+v", err)
	}
	router.Use(
		// all these handlers are shared by every route
		access.Logger(log.Printf),
		access.CustomLogger(ddLogger),
		fault.Recovery(log.Printf),
	)

	// serve index file
	router.Get("/", func(c *routing.Context) error {
		return c.Write("<html><head><title>Hello, world</title></head><body><h1>Hello, world</h1></body></html>")
	})

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
