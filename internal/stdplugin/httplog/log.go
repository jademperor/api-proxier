// Package httplog to log http request
package httplog

import (
	"bytes"
	"net/http"
	"time"

	"github.com/jademperor/api-proxier/plugin"
	"github.com/jademperor/common/pkg/logger"
)

var (
	_ plugin.Plugin = &HTTPLogger{}
)

// New func: generate a new HTTPLogger
func New(logger *logger.Entity) *HTTPLogger {
	return &HTTPLogger{
		logger:      logger,
		logResponse: true,
		enabled:     true,
		status:      plugin.Working,
	}
}

// HTTPLogger ...
type HTTPLogger struct {
	logger      *logger.Entity // logger.Entity is writer to log
	logResponse bool           // logResponse to log response into file or not
	enabled     bool
	status      plugin.PlgStatus
}

// Handle ...
func (h *HTTPLogger) Handle(ctx *plugin.Context) {
	defer plugin.Recover("HTTPLogger")

	var (
		rbw *respBodyWriter
	)
	// to log response
	rbw = &respBodyWriter{
		body:           bytes.NewBufferString(""),
		status:         http.StatusOK,
		ResponseWriter: ctx.ResponseWriter(),
	}
	ctx.SetResponseWriter(rbw)

	// timer
	start := time.Now()
	// continue process
	ctx.Next()
	path := ctx.Path
	end := time.Now()
	latency := end.Sub(start)
	clientIP := ctx.Request().RemoteAddr
	fields := make(map[string]interface{})

	fields["requestForm"] = ctx.Form
	if h.logResponse {
		// set response
		fields["responseBody"] = rbw.body.String()
	}

	// log
	h.logger.WithFields(fields).Infof("[Request] %v |%3d| %13v | %15s |%-7s %s",
		end.Format("2006/01/02 - 15:04:05"),
		rbw.status,
		latency,
		clientIP,
		ctx.Method,
		path,
	)
}

// Status ...
func (h *HTTPLogger) Status() plugin.PlgStatus {
	return h.status
}

// Enabled ...
func (h *HTTPLogger) Enabled() bool {
	return h.enabled
}

// Name ...
func (h *HTTPLogger) Name() string {
	return "plugin.httplog"
}

// Enable ...
func (h *HTTPLogger) Enable(enabled bool) {
	h.enabled = enabled
	if !enabled {
		h.status = plugin.Stopped
	} else {
		h.status = plugin.Working
	}
}

// type respBodyWriter to write log ...
type respBodyWriter struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
}

func (w *respBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *respBodyWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
