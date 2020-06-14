package kubernetesloader

import (
	"bytes"
	"context"
	"encoding/base64"
	"log"
	"mime"

	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/opentracing/opentracing-go"
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

func (l *loader) GetProjects(ctx context.Context) ([]projects.Project, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "kubernetesloader/get-projects")
	defer span.Finish()

	configMaps, err := l.clientset.CoreV1().ConfigMaps("projectlist").List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var out []projects.Project
	for _, m := range configMaps.Items {
		p := projects.Project{
			Name:        m.Data["name"],
			Slug:        m.Name,
			Description: m.Data["description"],
			URL:         m.Data["url"],
		}
		if imageBytes, hasImage := m.BinaryData["image"]; hasImage {
			p.Image = base64.StdEncoding.EncodeToString(imageBytes)
			_, format, err := image.DecodeConfig(bytes.NewReader(imageBytes))
			if err != nil {
				log.Printf("[%v] Could not decode image format: %v", m.Name, err)
				format = "jpg"
			}
			p.ImageMimeType = mime.TypeByExtension("." + format)
		}
		out = append(
			out,
			p,
		)
	}

	return out, nil
}
