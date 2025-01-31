package api

import (
	goctx "context"

	v1 "github.com/flanksource/config-db/api/v1"
	"github.com/google/uuid"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var KubernetesClient *kubernetes.Clientset
var KubernetesRestConfig *rest.Config
var Namespace string

func NewScrapeContext(scraper *v1.ConfigScraper, id *uuid.UUID) *v1.ScrapeContext {
	return &v1.ScrapeContext{
		Context:              goctx.Background(),
		Scraper:              scraper,
		ScraperID:            id,
		Namespace:            Namespace,
		Kubernetes:           KubernetesClient,
		KubernetesRestConfig: KubernetesRestConfig,
	}
}
