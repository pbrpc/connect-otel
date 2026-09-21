//revive:disable:package-comments
package connectotel

import (
	"testing"

	"connectrpc.com/connect/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/pbrpc/otel-testing/mocks/tracer"
)

const procedure = "/example.ExampleService/Echo"

// hasAttribute reports whether attrs carries want.
func hasAttribute(attrs []attribute.KeyValue, want attribute.KeyValue) bool {
	for _, attr := range attrs {
		if attr.Key == want.Key && attr.Value == want.Value {
			return true
		}
	}

	return false
}

func TestSplitProcedure(t *testing.T) {
	cases := []struct {
		procedure, service, method string
	}{
		{"/example.ExampleService/Echo", "example.ExampleService", "Echo"},
		{"example.ExampleService/Echo", "example.ExampleService", "Echo"},
		{"/example.ExampleService/", "example.ExampleService", ""},
		{"/example.ExampleService", "example.ExampleService", ""},
		{"", "", ""},
	}

	for _, c := range cases {
		t.Run(c.procedure, func(t *testing.T) {
			service, method := splitProcedure(c.procedure)

			if service != c.service || method != c.method {
				t.Errorf("got (%q, %q), want (%q, %q)", service, method, c.service, c.method)
			}
		})
	}
}

func TestName(t *testing.T) {
	t.Run("names the system, service, and method", func(t *testing.T) {
		tt, ctx := tracer.New(t)
		defer tt.Shutdown(t)

		Name(trace.SpanFromContext(ctx), procedure)
		tt.EndSpan()

		span := tt.GetSpans()[0]

		for _, want := range []attribute.KeyValue{
			semconv.RPCSystemConnectRPC,
			semconv.RPCService("example.ExampleService"),
			semconv.RPCMethod("Echo"),
		} {
			if !hasAttribute(span.Attributes, want) {
				t.Errorf("attributes = %v, want %v", span.Attributes, want)
			}
		}
		if span.Status.Code != codes.Unset {
			t.Errorf("status = %v, want %v", span.Status.Code, codes.Unset)
		}
	})
}

func TestFail(t *testing.T) {
	t.Run("records the code and marks the span failed", func(t *testing.T) {
		tt, ctx := tracer.New(t)
		defer tt.Shutdown(t)

		Fail(trace.SpanFromContext(ctx), connect.CodeNotFound)
		tt.EndSpan()

		span := tt.GetSpans()[0]

		if !hasAttribute(span.Attributes, semconv.RPCConnectRPCErrorCodeKey.String("not_found")) {
			t.Errorf("attributes = %v, want the not_found code", span.Attributes)
		}
		if span.Status.Code != codes.Error {
			t.Errorf("status = %v, want %v", span.Status.Code, codes.Error)
		}
		if span.Status.Description != "not_found" {
			t.Errorf("status description = %q, want %q", span.Status.Description, "not_found")
		}
	})
}
