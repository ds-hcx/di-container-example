package middleware

import (
	"net/http"

	"github.com/opentracing/opentracing-go"

	"github.com/VirgilSecurity/virgil-services-core-kit/http/response"
	"github.com/VirgilSecurity/virgil-services-core-kit/tracer"
)

//
// WithTracer wraps request with Tracer functionality.
//
func WithTracer(
	t tracer.Tracer,
	req *http.Request,
	callback func(req *http.Request) response.Provider,
) response.Provider {

	var span tracer.Span
	spanContext, err := t.Extract(tracer.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	if err != nil {
		span = t.StartSpan(tracer.GetCallerInfo())
	} else {
		span = t.StartSpan(tracer.GetCallerInfo(), tracer.RPCServerOption(spanContext))
	}

	span.SetTag(tracer.TagHTTPMethod, req.Method)
	span.SetTag(tracer.TagHTTPRoute, req.RequestURI)
	span.SetTag(tracer.TagComponent, tracer.ComponentMiddleware)
	defer span.Finish()

	return callback(req.WithContext(
		tracer.ContextWithSpan(
			req.Context(), span,
		)),
	)
}
