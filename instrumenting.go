package main

import (
	"context"
	"fmt"
	"time"

	"kit-fiber-example/metrics"
	"kit-fiber-example/transport"
)

type MetricsMiddleware[T any] struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           T
}

func NewMetricsMiddleware[T any](
	requestCount metrics.Counter,
	requestLatency metrics.Histogram,
	next T,
) MetricsMiddleware[T] {
	return MetricsMiddleware[T]{
		requestCount:   requestCount,
		requestLatency: requestLatency,
		next:           next,
	}
}

// InstrumentMethod - функция для инструментирования отдельного метода
func (mw MetricsMiddleware[T]) InstrumentMethod(
	methodName string,
	methodCb func() (any, error),
) (any, error) {
	begin := time.Now()

	result, err := methodCb()

	lvs := []string{
		"method", methodName,
		"error", fmt.Sprint(err != nil),
	}

	mw.requestCount.With(lvs...).Add(1)
	mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())

	return result, err
}

type InstrumentedStringService struct {
	MetricsMiddleware[transport.StringService]
}

func NewInstrumentedStringService(
	requestCount metrics.Counter,
	requestLatency metrics.Histogram,
	next transport.StringService,
) transport.StringService {
	return &InstrumentedStringService{
		MetricsMiddleware: NewMetricsMiddleware(requestCount, requestLatency, next),
	}
}

func (s *InstrumentedStringService) Uppercase(str string) (string, error) {
	result, err := s.InstrumentMethod("uppercase", func() (any, error) {
		return s.next.Uppercase(str)
	})
	if err != nil {
		return "", err
	}
	return result.(string), nil
}

func (s *InstrumentedStringService) AskClaude(ctx context.Context, question string) (string, error) {
	result, err := s.InstrumentMethod("askClaude", func() (any, error) {
		return s.next.AskClaude(ctx, question)
	})
	if err != nil {
		return "", err
	}
	return result.(string), nil
}
