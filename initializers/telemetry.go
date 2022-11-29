package initializers

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/antoniomralmeida/k2/lib"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

// var logger = log.New(os.Stderr, "zipkin-example", log.Ldate|log.Ltime|log.Llongfile)

type Telemetry struct {
	ctx    context.Context
	cancel context.CancelFunc
	trace  trace.Tracer
}

var t Telemetry

// initTracer creates a new trace provider instance and registers it as global trace provider.
func (t *Telemetry) initTracer(url string) (func(context.Context) error, error) {
	// Create Zipkin Exporter and install it as a global tracer.
	//
	// For demoing purposes, always sample. In a production application, you should
	// configure the sampler to a trace.ParentBased(trace.TraceIDRatioBased) set at the desired
	// ratio.
	exporter, err := zipkin.New(
		url,
		zipkin.WithLogger(log.Default()),
	)
	if err != nil {
		return nil, err
	}

	batcher := sdktrace.NewBatchSpanProcessor(exporter)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("k2-telemetry"),
		)),
	)
	otel.SetTracerProvider(tp)
	return tp.Shutdown, nil
}

func InitTelemetry() {
	zipkin := os.Getenv("TELEMETRY")
	Log(lib.Ping(zipkin), Fatal)

	t = Telemetry{}
	url := flag.String("zipkin", zipkin, "zipkin url")
	flag.Parse()

	t.ctx, t.cancel = signal.NotifyContext(context.Background(), os.Interrupt)
	defer t.cancel()
	_, err := t.initTracer(*url)
	Log(err, Fatal)

}

func Begin(spanName string, ctx context.Context) (context.Context, trace.Span) {
	if ctx == nil {
		ctx = t.ctx
		t.trace = otel.GetTracerProvider().Tracer("component-" + spanName)
	}
	return t.trace.Start(ctx, spanName)
}
