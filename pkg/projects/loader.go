package projects

type Loader interface {
	GetProjects() ([]Project, error)
}
