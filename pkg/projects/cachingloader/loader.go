package cachingloader

import (
	"context"
	"sync"
	"time"

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
	l.mtx.Lock()
	defer l.mtx.Unlock()

	if l.cache == nil || l.cache.createdAt.Sub(time.Now()) > time.Second {
		p, err := l.inner.GetProjects(ctx)
		l.cache = &cache{
			p:         p,
			err:       err,
			createdAt: time.Now(),
		}
	}
	return l.cache.p, l.cache.err
}
