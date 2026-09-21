# connect-otel

`connect-otel` describes a Connect RPC on a span: which procedure it is, and
which code it answered with. A server's handler span and a client's outgoing
call span are named through the same two functions, so the attributes are the
same wherever the span was started, and whatever derives metrics from spans
sees one shape.

## Installation

```bash
go get github.com/pbrpc/connect-otel
```

## Usage

```go
span := trace.SpanFromContext(ctx)

connectotel.Name(span, spec.Procedure)

if err := next(ctx, spec, stream); err != nil {
	connectotel.Fail(span, connect.CodeOf(err))

	return err
}
```

## Attributes

| Function | Attributes                                                  |
| -------- | ----------------------------------------------------------- |
| `Name`   | `rpc.system` (`connect_rpc`), `rpc.service`, `rpc.method`   |
| `Fail`   | `rpc.connect_rpc.error_code`; status `Error` with the code  |

The keys are the OpenTelemetry semantic conventions for RPC, so a collector's
`spanmetrics` connector configured with `rpc.service`, `rpc.method`, and
`rpc.connect_rpc.error_code` as dimensions counts every named span by
procedure and outcome.
