package location

import (
	"context"
	"math/rand"
	"time"

	"github.com/opentracing/opentracing-go"
)

// Service ...
type Service struct{}

// FindAll ..
func (s *Service) FindAll(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "location-repository")
	time.Sleep(time.Duration(rand.Intn(200)+200) * time.Millisecond)

	span.SetTag("db.table", "location")
	span.SetTag("db.context", "findlocation")

	span.Finish()
	return nil
}
