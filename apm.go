package apm

import (
	"os"
	"time"

	"cloud.google.com/go/profiler"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

func Init(service, version string) {
	initLogging()
	initTracing()
	initProfiling(service, version)
}

func initLogging() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
		TimestampFormat: time.RFC3339Nano,
	})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
}

func initStats(exporter *stackdriver.Exporter) {
	view.SetReportingPeriod(60 * time.Second)
	view.RegisterExporter(exporter)
}

func initTracing() {
	exporter, err := stackdriver.NewExporter(stackdriver.Options{})
	if err != nil {
		logrus.Warnf("failed to initialize Stackdriver exporter: %+v", err)
	} else {
		trace.RegisterExporter(exporter)
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
		logrus.Info("registered Stackdriver tracing")

		initStats(exporter)
	}
}

func initProfiling(service, version string) {
	if err := profiler.Start(profiler.Config{
		Service:        service,
		ServiceVersion: version,
	}); err != nil {
		logrus.Warnf("failed to start profiler: %+v", err)
	}
}
