package telemetry

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

var logger = log.New(os.Stderr, "k2-zipkin", log.Ldate|log.Ltime|log.Llongfile)
var ctx context.Context
var cancel context.CancelFunc

// initTracer creates a new trace provider instance and registers it as global trace provider.
func initTracer(url string) (func(context.Context) error, error) {
	// Create Zipkin Exporter and install it as a global tracer.
	//
	// For demoing purposes, always sample. In a production application, you should
	// configure the sampler to a trace.ParentBased(trace.TraceIDRatioBased) set at the desired
	// ratio.
	exporter, err := zipkin.New(
		url,
		zipkin.WithLogger(logger),
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
	url := flag.String("zipkin", "http://localhost:9411/api/v2/spans", "zipkin url")
	flag.Parse()

	ctx, cancel = signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := initTracer(*url)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: " + err.Error())
		}
	}()
}

func Begin(component string, spanName string) trace.Span {
	tr := otel.GetTracerProvider().Tracer(component)
	_, span := tr.Start(ctx, spanName)
	return span
}
