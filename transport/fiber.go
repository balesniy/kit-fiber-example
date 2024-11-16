package transport

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/trace"

	"kit-fiber-example/health"
	"kit-fiber-example/metrics"
	"kit-fiber-example/middlewares"
	"kit-fiber-example/service"
)

// Service interface defines our business logic
type StringService interface {
	Uppercase(string) (string, error)
	AskClaude(context.Context, string) (string, error)
}

// Fiber transport layer?? or application layer?
type fiberTransport struct {
	Uppercase middlewares.Endpoint[UppercaseRequest, UppercaseResponse]
	AskClaude middlewares.Endpoint[AskClaudeRequest, AskClaudeResponse]
	//services    []Service
	Metrics *metrics.Metrics // todo interface MetricsCollector
	Health  *health.Health   // todo interface HealthChecker
}

func NewFiberTransport(svc StringService, h *health.Health, m *metrics.Metrics, t trace.Tracer) *fiberTransport {
	uppercaseEndpoint := makeUppercaseEndpoint(svc)
	uppercaseEndpoint = middlewares.LoggingMiddleware(uppercaseEndpoint)
	uppercaseEndpoint = middlewares.WithMetrics(m, uppercaseEndpoint)
	uppercaseEndpoint = middlewares.WithTracing(t, uppercaseEndpoint)

	askClaudeEndpoint := makeAskClaudeEndpoint(svc)

	return &fiberTransport{
		Uppercase: uppercaseEndpoint,
		AskClaude: askClaudeEndpoint,
		Metrics:   m,
		Health:    h,
	}
}

// Health check handler
func (t *fiberTransport) HandleHealth(c *fiber.Ctx) error {
	if t.Health.IsHealthy() {
		return c.SendStatus(200)
	}
	return c.SendStatus(503)
}

// Readiness check handler
func (t *fiberTransport) HandleReady(c *fiber.Ctx) error {
	// Add your readiness check logic here
	// For example, check database connections, external services, etc.
	return c.SendStatus(200)
}

// Error handling middleware
func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var e service.ServiceError
	if errors.As(err, &e) {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}

func InitApp(transport *fiberTransport) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
	})

	// Add fiber middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "${time} ${method} ${path} ${status} ${latency}\n",
	}))

	// Setup routes
	app.Post("/uppercase", transport.HandleUppercase)
	app.Post("/ask", transport.HandleAskClaude)
	app.Get("/health", transport.HandleHealth)
	app.Get("/ready", transport.HandleReady)

	prometheusHandler := adaptor.HTTPHandler(promhttp.Handler())
	app.Get("/metrics", prometheusHandler)
	return app
}
