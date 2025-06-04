package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
	"time"

	"github.com/rixtrayker/go-loadbalancer/internal/logging"
	"github.com/rixtrayker/go-loadbalancer/internal/monitoring"
	"github.com/rixtrayker/go-loadbalancer/internal/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)
