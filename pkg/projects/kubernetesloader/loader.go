package kubernetesloader

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/telliott-io/front/pkg/projects"
)

func New() (projects.Loader, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &loader{
		clientset: clientset,
	}, nil
}

type loader struct {
	clientset *kubernetes.Clientset
}

func (l *loader) GetProjects() ([]projects.Project, error) {
	configMaps, err := l.clientset.CoreV1().ConfigMaps("projectlist").List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var out []projects.Project
	for _, m := range configMaps.Items {
		out = append(
			out,
			projects.Project{
				Name:        m.Data["name"],
				Description: m.Data["description"],
			},
		)
	}

	return out, nil
}

func HitK8s() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {
		configMaps, err := clientset.CoreV1().ConfigMaps("projectlist").List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		for _, m := range configMaps.Items {
			fmt.Printf("ConfigMap %q: %+v\n", m.Name, m.Data)
		}
	}
}
