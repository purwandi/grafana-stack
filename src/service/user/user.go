package user

import (
	"context"
	"math/rand"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/purwandi/grafana-stack/src/pkg/logger"
	"go.uber.org/zap"
)

// Service ...
type Service struct{}

// FindAll ..
func (s *Service) FindAll(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user-repository")
	time.Sleep(time.Duration(rand.Intn(200)+200) * time.Millisecond)

	span.SetTag("db.table", "user")
	span.SetTag("db.context", "finduser")

	logger.InfoWithContext(ctx, "", zap.String("query", "user-query"))

	span.Finish()
	return nil
}
