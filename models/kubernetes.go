package models

const (
	// LabelKey k8s label key
	LabelKey = "app"
	// WatchAdded k8s watch added
	WatchAdded   = "ADDED"
	// WatchDeleted k8s watch deleted
	WatchDeleted = "DELETED"
)

// K8sRes struct k8s scheduling result
type K8sRes struct {
	Result string
	Detail string
}
