package models

const (
	// LabelKey k8s label key
	LabelKey = "app"
)

// K8sRes struct k8s scheduling result
type K8sRes struct {
	Result string
	Detail string
}
