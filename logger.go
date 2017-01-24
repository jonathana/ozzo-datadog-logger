package ozDdLogger

import (
	"sync"
	"net/http"
	"github.com/jonathana/ozzo-routing/access"
	"github.com/DataDog/datadog-go/statsd"
	"fmt"
	"strings"
)

type DdLoggerConfig struct {
	TagProtocol	bool
	TagMethod	bool
	LogResponseCode	bool
	BaseUrl		string
}

func copyConfig(cfg DdLoggerConfig) DdLoggerConfig {
	newCfg := DdLoggerConfig{
		TagProtocol: cfg.TagProtocol,
		TagMethod: cfg.TagMethod,
		LogResponseCode: cfg.LogResponseCode,
		BaseUrl: cfg.BaseUrl,
	}
	return newCfg
}

func ConfigureLogger(client *statsd.Client, config DdLoggerConfig) (access.LogWriterFunc, error) {
	var mutex = &sync.Mutex{}
	mutex.Lock()
	_config := copyConfig(config)
	_client := client
	if len(_client.Namespace) == 0 {
		_client.Namespace = "ozzo.routing"
	}
	mutex.Unlock()
	ddLogFunc := func(req *http.Request, res *access.LogResponseWriter, elapsed float64) {
		go logRequest(_client, _config, req, res, elapsed)
	}
	return ddLogFunc, nil
}

func makeRequestTags(cfg DdLoggerConfig, req *http.Request, res *access.LogResponseWriter) []string {
	tags := []string{fmt.Sprintf("path:%s%s", cfg.BaseUrl, req.URL.EscapedPath())}
	if cfg.TagMethod {
		tags = append(tags, fmt.Sprintf("method:%s", strings.ToLower(req.Method)))
	}
	if cfg.TagProtocol {
		tags = append(tags, fmt.Sprintf("protocol:%s", strings.ToLower(req.Proto)))
	}
	if cfg.LogResponseCode {
		tags = append(tags, fmt.Sprintf("response_code:%d", res.Status))
	}

	return tags
}

func logRequest(client *statsd.Client, cfg DdLoggerConfig, req *http.Request, res *access.LogResponseWriter, elapsed float64) {
	tags := makeRequestTags(cfg, req, res)
	if cfg.LogResponseCode {
		client.Incr(fmt.Sprintf("response_code.%s", res.Status), tags, 1)
		client.Incr("response_code.all", tags, 1)
	}
	client.Histogram("response_time", elapsed, tags, 1)
}
