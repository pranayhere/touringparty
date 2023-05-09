package tracing

import (
	"fmt"
	"io"
	"net"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegerconfig "github.com/uber/jaeger-client-go/config"
	jaegerzap "github.com/uber/jaeger-client-go/log/zap"
	"github.com/uber/jaeger-lib/metrics/prometheus"
	"go.uber.org/zap"
)

var (
	// Tracer is used for creating spans for distributed tracing
	tracer opentracing.Tracer
	// Closer holds an instance to the RequestTracing object's Closer.
	closer io.Closer
)

// Config ... struct expected by Init func to initialize jaeger tracing client.
type Config struct {
	LogSpans    bool   `yaml:"logSpan"`     // when set to true, reporter logs all submitted spans
	Host        string `yaml:"host"`        // jaeger-agent UDP binary thrift protocol endpoint
	Port        string `yaml:"port"`        // jaeger-agent UDP binary thrift protocol server port
	ServiceName string `yaml:"serviceName"` // name of this service used by tracer.
	Disabled    bool   `yaml:"disabled"`    // to mock tracer
}

func Init(env string, cnf Config, zlog *zap.Logger) error {
	hostPortPath := net.JoinHostPort(cnf.Host, cnf.Port)

	config := &jaegerconfig.Configuration{
		ServiceName: fmt.Sprintf("%s-%s", cnf.ServiceName, env),
		Sampler: &jaegerconfig.SamplerConfig{
			Type:  jaeger.SamplerTypeProbabilistic,
			Param: 0.05,
		},
		Reporter: &jaegerconfig.ReporterConfig{
			LogSpans:           cnf.LogSpans,
			LocalAgentHostPort: hostPortPath,
		},
		Disabled: cnf.Disabled,
	}

	var err error
	tracer, closer, err = config.NewTracer(
		jaegerconfig.Logger(jaegerzap.NewLogger(zlog)),
		jaegerconfig.Metrics(prometheus.New()),
	)
	if err != nil {
		return err
	}

	opentracing.SetGlobalTracer(tracer)

	return nil
}

// Close calls the closer function if initialized
func Close() error {
	if closer == nil {
		return nil
	}

	return closer.Close()
}
