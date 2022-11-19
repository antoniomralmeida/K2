package telemetry

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"time"

	"github.com/antoniomralmeida/k2/initializers"
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

func Init() {
	t = Telemetry{}
	url := flag.String("zipkin", "http://localhost:9411/api/v2/spans", "zipkin url")
	flag.Parse()

	t.ctx, t.cancel = signal.NotifyContext(context.Background(), os.Interrupt)
	defer t.cancel()

	shutdown, err := t.initTracer(*url)
	if err != nil {
		initializers.Log(err, initializers.Fatal)
	}
	defer func() {
		if err := shutdown(t.ctx); err != nil {
			initializers.Log("failed to shutdown TracerProvider: "+err.Error(), initializers.Fatal)
		}
	}()
}

func Begin(ct context.Context, spanName string) (context.Context, trace.Span) {
	if ct == context.TODO() {
		ct = t.ctx
	}
	tr := otel.GetTracerProvider().Tracer("component-" + spanName)
	return tr.Start(ct, spanName)
}

func End(span trace.Span) {
	time.Sleep(1 * time.Millisecond)
	span.End()
}

// clear server
/*
	tr := otel.GetTracerProvider().Tracer("component-main")
	ctx, span := tr.Start(ctx, "foo", trace.WithSpanKind(trace.SpanKindServer))
	<-time.After(6 * time.Millisecond)
	bar(ctx)
	<-time.After(6 * time.Millisecond)
	span.End()
}

func bar(ctx context.Context) {
	tr := otel.GetTracerProvider().Tracer("component-bar")
	_, span := tr.Start(ctx, "bar")
	<-time.After(6 * time.Millisecond)
	span.End()
}
*/
