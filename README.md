# ozzo-datadog-logger

Module implementing an ozzo-routing/access CustomLogger that writes its metrics to a Datadog client.

# Inspiration

This module was inspired by [the node connect middleware Datadog logger](https://github.com/AppPress/node-connect-datadog).

# Usage

See the example in examples/ddlogger.go.  The main things to note are:
- Unlike the connect logger, the Golang Datadog statsd client supports a .Tags attribute that can be used to set
default tags on everything.  Use that to set tags that get used for everything you log.
- Similarly, the Golang Datadog statds client supports a .Namespace attribute that prefixes all logged
value names with the value given for .Namespace.

```go
package main

import (
	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/access"
	"github.com/DataDog/datadog-go/statsd"
	"github.com/jonathana/ozzo-datadog-logger"
	"fmt"
)

func main() {
	router := routing.New()
	c, err := statsd.New("127.0.0.1:8125")
	if err != nil {
		log.Fatal(err)
	}
	// prefix every metric with the app name
	c.Namespace = // Insert your app name here
	// set the default tags
	c.Tags = append(c.Tags, // insert tag(s) you want here)
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
		// Any middleware you want before CustomLogger (minimize to get more accurate timing)
		access.CustomLogger(ddLogger),
		// Any other middleware you want following
	)
	
	// Do whatever else you need for your server here, then run it
}

```
