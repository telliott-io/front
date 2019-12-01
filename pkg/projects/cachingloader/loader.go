package cachingloader

import (
	"context"
	"sync"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/telliott-io/front/pkg/projects"
)

func New(l projects.Loader) projects.Loader {
	return &loader{
		inner: l,
	}
}

type loader struct {
	cache *cache
	inner projects.Loader
	mtx   sync.Mutex
}

type cache struct {
	p         []projects.Project
	err       error
	createdAt time.Time
}

func (l *loader) GetProjects(ctx context.Context) ([]projects.Project, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "cachingloader/get-projects")
	defer span.Finish()

	l.mtx.Lock()
	defer l.mtx.Unlock()

	if l.cache == nil || time.Since(l.cache.createdAt) > 5*time.Second {
		span.SetTag("cache-hit", false)
		p, err := l.inner.GetProjects(ctx)
		l.cache = &cache{
			p:         p,
			err:       err,
			createdAt: time.Now(),
		}
	} else {
		span.SetTag("cache-hit", true)
	}
	return l.cache.p, l.cache.err
}
