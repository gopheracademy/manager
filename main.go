package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gopheracademy/manager/pkg/log"
	"github.com/gopheracademy/manager/pkg/tracing"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	jexpvar "github.com/uber/jaeger-lib/metrics/expvar"
	jprom "github.com/uber/jaeger-lib/metrics/prometheus"

	"github.com/pacedotdev/oto/otohttp"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	metricsBackend string = "prometheus"
	logger         *zap.Logger
	metricsFactory metrics.Factory

	basepath string
	jaegerUI string
)

// spaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type spaHandler struct {
	staticPath string
	indexPath  string
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func main() {
	onInitialize()
	zapLogger := logger.With(zap.String("service", "showrunner"))
	logg := log.NewFactory(zapLogger)
	mytracer := tracing.Init("showrunner", metricsFactory, logg)
	tracedRouter := tracing.NewServeMux(mytracer)
	tracedRouter.Mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})
	conferenceService := newconferenceService(mytracer, metricsFactory, logg)
	server := otohttp.NewServer()

	RegisterConferenceService(metricsFactory.Namespace(metrics.NSOptions{Name: "conference.service"}), mytracer,
		logg, server, conferenceService)

	tracedRouter.Handle("/oto/", server)
	spa := spaHandler{staticPath: "./www/public", indexPath: "index.html"}

	tracedRouter.Mux.Handle("/metrics", promhttp.Handler()) // Prometheus
	tracedRouter.Mux.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler: tracedRouter,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	zapLogger.Fatal(srv.ListenAndServe().Error())

}

// onInitialize is called before the command is executed.
func onInitialize() {
	rand.Seed(int64(time.Now().Nanosecond()))
	logger, _ = zap.NewDevelopment(
		zap.AddStacktrace(zapcore.FatalLevel),
		zap.AddCallerSkip(1),
	)

	switch metricsBackend {
	case "expvar":
		metricsFactory = jexpvar.NewFactory(10) // 10 buckets for histograms
		logger.Info("Using expvar as metrics backend")
	case "prometheus":
		metricsFactory = jprom.New().Namespace(metrics.NSOptions{Name: "showrunner", Tags: nil})
		logger.Info("Using Prometheus as metrics backend")
	default:
		logger.Fatal("unsupported metrics backend " + metricsBackend)
	}
}
