// Package connectotel describes a Connect RPC on a span: which procedure it
// is, and which code it answered with. Every span standing for a Connect
// call, whether a server's handler or a client's outgoing call, is named
// through it, so the attributes are the same wherever the span was started.
package connectotel

import (
	"strings"

	"connectrpc.com/connect/v2"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

// splitProcedure breaks "/package.Service/Method" into its service and method.
func splitProcedure(procedure string) (service, method string) {
	service, method, _ = strings.Cut(strings.TrimPrefix(procedure, "/"), "/")

	return service, method
}

// Name marks span as the Connect RPC procedure names: the system, and the
// service and method split from "/package.Service/Method".
func Name(span trace.Span, procedure string) {
	service, method := splitProcedure(procedure)

	span.SetAttributes(
		semconv.RPCSystemConnectRPC,
		semconv.RPCService(service),
		semconv.RPCMethod(method),
	)
}

// Fail records code as the Connect code the RPC answered with, and marks
// span as failed by it.
func Fail(span trace.Span, code connect.Code) {
	span.SetAttributes(semconv.RPCConnectRPCErrorCodeKey.String(code.String()))
	span.SetStatus(codes.Error, code.String())
}
