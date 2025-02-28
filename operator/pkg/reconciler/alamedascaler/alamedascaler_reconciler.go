package alamedascaler

import (
	"fmt"

	autoscaling_v1alpha1 "github.com/containers-ai/alameda/operator/pkg/apis/autoscaling/v1alpha1"
	utils "github.com/containers-ai/alameda/operator/pkg/utils"
	logUtil "github.com/containers-ai/alameda/operator/pkg/utils/log"
	utilsresource "github.com/containers-ai/alameda/operator/pkg/utils/resources"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	alamedascalerReconcilerScope = logUtil.RegisterScope("alamedascaler_reconciler", "alamedascaler_reconciler", 0)
)

// Reconciler reconciles AlamedaScaler object
type Reconciler struct {
	client        client.Client
	alamedascaler *autoscaling_v1alpha1.AlamedaScaler
}

// NewReconciler creates Reconciler object
func NewReconciler(client client.Client, alamedascaler *autoscaling_v1alpha1.AlamedaScaler) *Reconciler {
	return &Reconciler{
		client:        client,
		alamedascaler: alamedascaler,
	}
}

// HasAlamedaDeployment checks the AlamedaScaler has the deployment or not
func (reconciler *Reconciler) HasAlamedaDeployment(deploymentNS, deploymentName string) bool {
	key := utils.GetNamespacedNameKey(deploymentNS, deploymentName)
	_, ok := reconciler.alamedascaler.Status.AlamedaController.Deployments[autoscaling_v1alpha1.NamespacedName(key)]
	return ok
}

// HasAlamedaPod checks the AlamedaScaler has the AlamedaPod or not
func (reconciler *Reconciler) HasAlamedaPod(podNS, podName string) bool {
	for _, deployment := range reconciler.alamedascaler.Status.AlamedaController.Deployments {
		deploymentNS := deployment.Namespace
		for _, pod := range deployment.Pods {
			if deploymentNS == podNS && pod.Name == podName {
				return true
			}
		}
	}
	return false
}

// RemoveAlamedaDeployment removes deployment from alamedaController of AlamedaScaler
func (reconciler *Reconciler) RemoveAlamedaDeployment(deploymentNS, deploymentName string) *autoscaling_v1alpha1.AlamedaScaler {
	key := utils.GetNamespacedNameKey(deploymentNS, deploymentName)

	if _, ok := reconciler.alamedascaler.Status.AlamedaController.Deployments[autoscaling_v1alpha1.NamespacedName(key)]; ok {
		delete(reconciler.alamedascaler.Status.AlamedaController.Deployments, autoscaling_v1alpha1.NamespacedName(key))
		return reconciler.alamedascaler
	}
	return reconciler.alamedascaler
}

// InitAlamedaController try to initialize alamedaController field of AlamedaScaler
func (reconciler *Reconciler) InitAlamedaController() (alamedascaler *autoscaling_v1alpha1.AlamedaScaler, needUpdated bool) {
	if reconciler.alamedascaler.Status.AlamedaController.Deployments == nil {
		reconciler.alamedascaler.Status.AlamedaController.Deployments = map[autoscaling_v1alpha1.NamespacedName]autoscaling_v1alpha1.AlamedaDeployment{}
		return reconciler.alamedascaler, true
	}
	return reconciler.alamedascaler, false
}

// UpdateStatusByDeployment updates status by deployment
func (reconciler *Reconciler) UpdateStatusByDeployment(deployment *appsv1.Deployment) *autoscaling_v1alpha1.AlamedaScaler {
	alamedaScalerNS := reconciler.alamedascaler.GetNamespace()
	alamedaScalerName := reconciler.alamedascaler.GetName()

	listResources := utilsresource.NewListResources(reconciler.client)
	alamedaDeploymentNS := deployment.GetNamespace()
	alamedaDeploymentName := deployment.GetName()
	alamedaDeploymentUID := deployment.GetUID()
	alamedaPodsMap := map[autoscaling_v1alpha1.NamespacedName]autoscaling_v1alpha1.AlamedaPod{}
	alamedaDeploymentsMap := reconciler.alamedascaler.Status.AlamedaController.Deployments
	if alamedaPods, err := listResources.ListPodsByDeployment(alamedaDeploymentNS, alamedaDeploymentName); err == nil && len(alamedaPods) > 0 {
		for _, alamedaPod := range alamedaPods {
			alamedaPodName := alamedaPod.GetName()
			alamedaPodUID := alamedaPod.GetUID()
			alamedascalerReconcilerScope.Infof(fmt.Sprintf("Pod (%s/%s) belongs to AlamedaScaler (%s/%s).", alamedaDeploymentNS, alamedaPodName, alamedaScalerNS, alamedaScalerName))
			alamedaContainers := []autoscaling_v1alpha1.AlamedaContainer{}
			for _, alamedaContainer := range alamedaPod.Spec.Containers {
				alamedaContainers = append(alamedaContainers, autoscaling_v1alpha1.AlamedaContainer{
					Name: alamedaContainer.Name,
				})
			}
			alamedaPodsMap[autoscaling_v1alpha1.NamespacedName(utils.GetNamespacedNameKey(alamedaPod.GetNamespace(), alamedaPodName))] = autoscaling_v1alpha1.AlamedaPod{
				Name:       alamedaPodName,
				UID:        string(alamedaPodUID),
				Containers: alamedaContainers,
			}
		}
	}

	alamedaDeploymentsMap[autoscaling_v1alpha1.NamespacedName(utils.GetNamespacedNameKey(deployment.GetNamespace(), deployment.GetName()))] = autoscaling_v1alpha1.AlamedaDeployment{
		Namespace: alamedaDeploymentNS,
		Name:      alamedaDeploymentName,
		UID:       string(alamedaDeploymentUID),
		Pods:      alamedaPodsMap,
	}
	reconciler.alamedascaler.Status.AlamedaController.Deployments = alamedaDeploymentsMap
	return reconciler.alamedascaler
}
