package ozDdLogger

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/DataDog/datadog-go/statsd"
	"fmt"
	"reflect"
)

func TestOzDdLogger(t *testing.T) {
	assert.True(t, true, "How can true be false?")
}

func TestCreateLoggerFunction(t *testing.T) {
	c, err := statsd.New("127.0.0.1:8125")
	assert.NoError(t, err, "Could not create statsd client")
	cfg := DdLoggerConfig{}
	logFunc, err := ConfigureLogger(c, cfg)
	assert.NoError(t, err, "Error getting configured logger function")
	reflect.TypeOf(logFunc)
	fmt.Printf("logFunc is a %+v", reflect.TypeOf(logFunc).Name())
	assert.Equal(t, "LogWriterFunc", reflect.TypeOf(logFunc).Name())
}
