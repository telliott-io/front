package projects

import "context"

type Loader interface {
	GetProjects(ctx context.Context) ([]Project, error)
}
