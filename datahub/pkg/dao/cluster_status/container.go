package clusterstatus

import (
	datahub_api "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
)

// ContainerOperation provides container measurement operations
type ContainerOperation interface {
	AddPods([]*datahub_api.Pod) error
	DeletePods([]*datahub_api.Pod) error
	ListAlamedaPods(string, string) ([]*datahub_api.Pod, error)
}
