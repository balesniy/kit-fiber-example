package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kit-fiber-example/config"
	"kit-fiber-example/health"
	"kit-fiber-example/metrics"
	"kit-fiber-example/service"
	"kit-fiber-example/tracing"
	"kit-fiber-example/transport"
)

// Options содержит все зависимости приложения
/*type Options struct {
	Config      config.AppConfig
	Logger      Logger
	Metrics     MetricsCollector
	Services    []Service
	HealthCheck HealthChecker
}*/

func main() {
	// Load config
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		panic(err)
	}

	// Initialize tracer
	tp, err := tracing.InitOtel(cfg)
	if err != nil {
		panic(err)
	}
	defer tp.Shutdown(context.Background())

	tracer := tp.Tracer("string-service")

	// Initialize service
	h := &health.Health{}
	metricsSet := metrics.Setup()
	claudeClient := service.NewClaudeClient(cfg)
	svc := service.String{
		ClaudeClient: claudeClient,
	}

	// Set initial health status
	h.SetHealthy()

	// Собираем все зависимости
	/*opts := Options{
		Config:      cfg,
		Logger:      logger,
		Metrics:     metrics,
		HealthCheck: healthCheck,
		Services:    services,
	}*/

	// Create Fiber transport
	tr := transport.NewFiberTransport(&svc, h, metricsSet, tracer)

	// Create Fiber app
	server := transport.InitApp(tr)

	// Graceful shutdown setup
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)

	// Start server
	serverError := make(chan error, 1)
	go func() {
		if err := server.Listen(cfg.Server.Port); err != nil {
			serverError <- err
		}
	}()

	// Wait for interrupt signal or server error
	select {
	case err := <-serverError:
		log.Printf("Server error: %v", err)
	case sig := <-shutdown:
		log.Printf("Start shutdown... Signal: %v", sig)

		// Create shutdown context with timeout
		ctx, cancel := context.WithTimeout(
			context.Background(),
			time.Duration(cfg.Server.ShutdownTimeout)*time.Second,
		)
		defer cancel()

		// Shutdown server
		if err := server.ShutdownWithContext(ctx); err != nil {
			log.Printf("Server forced to shutdown: %v", err)
		}
	}
}
