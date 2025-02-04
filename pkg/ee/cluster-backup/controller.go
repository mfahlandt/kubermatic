//go:build ee

/*
                  Kubermatic Enterprise Read-Only License
                         Version 1.0 ("KERO-1.0”)
                     Copyright © 2023 Kubermatic GmbH

   1.	You may only view, read and display for studying purposes the source
      code of the software licensed under this license, and, to the extent
      explicitly provided under this license, the binary code.
   2.	Any use of the software which exceeds the foregoing right, including,
      without limitation, its execution, compilation, copying, modification
      and distribution, is expressly prohibited.
   3.	THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND,
      EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
      MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
      IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
      CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
      TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
      SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

   END OF TERMS AND CONDITIONS
*/

package clusterbackup

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	kubermaticv1 "k8c.io/kubermatic/v2/pkg/apis/kubermatic/v1"
	kubermaticv1helper "k8c.io/kubermatic/v2/pkg/apis/kubermatic/v1/helper"
	clusterclient "k8c.io/kubermatic/v2/pkg/cluster/client"
	predicateutil "k8c.io/kubermatic/v2/pkg/controller/util/predicate"
	seedclusterresources "k8c.io/kubermatic/v2/pkg/ee/cluster-backup/resources/seed-cluster"
	userclusterresources "k8c.io/kubermatic/v2/pkg/ee/cluster-backup/resources/user-cluster"
	"k8c.io/kubermatic/v2/pkg/provider"
	"k8c.io/kubermatic/v2/pkg/resources"
	kkpreconciling "k8c.io/kubermatic/v2/pkg/resources/reconciling"
	"k8c.io/kubermatic/v2/pkg/util/workerlabel"
	"k8c.io/kubermatic/v2/pkg/version/kubermatic"
	"k8c.io/reconciler/pkg/reconciling"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	ControllerName = "cluster-backup-controller"
)

// UserClusterClientProvider provides functionality to get a user cluster client.
type UserClusterClientProvider interface {
	GetClient(ctx context.Context, c *kubermaticv1.Cluster, options ...clusterclient.ConfigOption) (ctrlruntimeclient.Client, error)
}

type reconciler struct {
	ctrlruntimeclient.Client

	seedGetter                    provider.SeedGetter
	workerName                    string
	recorder                      record.EventRecorder
	userClusterConnectionProvider UserClusterClientProvider
	log                           *zap.SugaredLogger
	versions                      kubermatic.Versions
}

func Add(mgr manager.Manager, numWorkers int, workerName string, userClusterConnectionProvider UserClusterClientProvider, seedGetter provider.SeedGetter, log *zap.SugaredLogger, versions kubermatic.Versions) error {
	reconciler := &reconciler{
		Client:                        mgr.GetClient(),
		seedGetter:                    seedGetter,
		workerName:                    workerName,
		recorder:                      mgr.GetEventRecorderFor(ControllerName),
		userClusterConnectionProvider: userClusterConnectionProvider,
		log:                           log,
		versions:                      versions,
	}

	c, err := controller.New(ControllerName, mgr, controller.Options{
		Reconciler:              reconciler,
		MaxConcurrentReconciles: numWorkers,
	})
	if err != nil {
		return fmt.Errorf("failed to create controller: %w", err)
	}

	if err := c.Watch(source.Kind(mgr.GetCache(), &kubermaticv1.Cluster{}), &handler.EnqueueRequestForObject{}, workerlabel.Predicates(workerName), predicateutil.Factory(func(o ctrlruntimeclient.Object) bool {
		cluster := o.(*kubermaticv1.Cluster)
		// Only watch clusters that are in a state where they can be reconciled.
		return !cluster.Spec.Pause && cluster.DeletionTimestamp == nil && cluster.Status.NamespaceName != ""
	})); err != nil {
		return fmt.Errorf("failed to create watch: %w", err)
	}
	return nil
}

func (r *reconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	cluster := &kubermaticv1.Cluster{}
	if err := r.Get(ctx, request.NamespacedName, cluster); err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, fmt.Errorf("failed to get cluster: %w", err)
	}

	// Resource is marked for deletion.
	if cluster.DeletionTimestamp != nil {
		return reconcile.Result{}, nil
	}

	if !cluster.Spec.IsClusterBackupEnabled() {
		return reconcile.Result{}, nil
	}

	// Add a wrapping here so we can emit an event on error
	result, err := kubermaticv1helper.ClusterReconcileWrapper(
		ctx,
		r.Client,
		r.workerName,
		cluster,
		r.versions,
		kubermaticv1.ClusterConditionClusterBackupControllerReconcilingSuccess,
		func() (*reconcile.Result, error) {
			return r.reconcile(ctx, cluster)
		},
	)

	if result == nil || err != nil {
		result = &reconcile.Result{}
	}

	if err != nil {
		r.recorder.Event(cluster, corev1.EventTypeWarning, "ReconcilingError", err.Error())
	}

	return *result, err
}

func (r *reconciler) reconcile(ctx context.Context, cluster *kubermaticv1.Cluster) (*reconcile.Result, error) {
	log := r.log.With("cluster", cluster)
	log.Debug("Reconciling")
	cbsl := &kubermaticv1.ClusterBackupStorageLocation{}
	if err := r.Client.Get(ctx, types.NamespacedName{Namespace: resources.KubermaticNamespace, Name: cluster.Spec.BackupConfig.BackupStorageLocation.Name}, cbsl); err != nil {
		log.Debug("ClusterBackupStorageLocation not found")
		return nil, nil
	}

	if err := r.ensureClusterBackupUserClusterResources(ctx, cluster, cbsl); err != nil {
		return nil, fmt.Errorf("failed to ensure cluster backup user-cluster resources: %w", err)
	}

	if err := r.ensureClusterBackupSeedClusterResources(ctx, cluster, cbsl); err != nil {
		return nil, fmt.Errorf("failed to ensure cluster backup Seed cluster resources: %w", err)
	}
	return nil, nil
}

func (r *reconciler) ensureClusterBackupUserClusterResources(ctx context.Context, cluster *kubermaticv1.Cluster, cbsl *kubermaticv1.ClusterBackupStorageLocation) error {
	userClusterClient, err := r.userClusterConnectionProvider.GetClient(ctx, cluster)
	if err != nil {
		return fmt.Errorf("failed to get user cluster client: %w", err)
	}
	nsReconcilers := []reconciling.NamedNamespaceReconcilerFactory{
		userclusterresources.NamespaceReconciler(),
	}
	if err := reconciling.ReconcileNamespaces(ctx, nsReconcilers, "", userClusterClient); err != nil {
		return fmt.Errorf("failed to reconcile Velero Namespace: %w", err)
	}

	saReconcilers := []reconciling.NamedServiceAccountReconcilerFactory{
		userclusterresources.ServiceAccountReconciler(),
	}
	if err := reconciling.ReconcileServiceAccounts(ctx, saReconcilers, resources.ClusterBackupNamespaceName, userClusterClient); err != nil {
		return fmt.Errorf("failed to reconcile Velero Service Account: %w", err)
	}

	clusterRoleBindingReconciler := []reconciling.NamedClusterRoleBindingReconcilerFactory{
		userclusterresources.ClusterRoleBindingReconciler(),
	}

	if err := reconciling.ReconcileClusterRoleBindings(ctx, clusterRoleBindingReconciler, "", userClusterClient); err != nil {
		return fmt.Errorf("failed to reconcile Velero Cluster Role Binding: %w", err)
	}

	clusterBackupCRDs, err := userclusterresources.CRDs()
	if err != nil {
		return fmt.Errorf("failed to load cluster backup CRDs: %w", err)
	}

	creators := []kkpreconciling.NamedCustomResourceDefinitionReconcilerFactory{}
	for i := range clusterBackupCRDs {
		creators = append(creators, userclusterresources.CRDReconciler(clusterBackupCRDs[i]))
	}

	if err = kkpreconciling.ReconcileCustomResourceDefinitions(ctx, creators, "", userClusterClient); err != nil {
		return fmt.Errorf("failed to reconcile Velero CRDs: %w", err)
	}

	bslReconcilers := []kkpreconciling.NamedBackupStorageLocationReconcilerFactory{
		userclusterresources.BSLReconciler(ctx, cluster, cbsl),
	}

	if err := kkpreconciling.ReconcileBackupStorageLocations(ctx, bslReconcilers, resources.ClusterBackupNamespaceName, userClusterClient); err != nil {
		return fmt.Errorf("failed to reconcile Velero BSL: %w", err)
	}

	return nil
}

func (r *reconciler) ensureClusterBackupSeedClusterResources(ctx context.Context, cluster *kubermaticv1.Cluster, cbsl *kubermaticv1.ClusterBackupStorageLocation) error {
	namespace := cluster.Status.NamespaceName

	clusterData := r.newClusterData(ctx, cluster)
	secretReconcilers := []reconciling.NamedSecretReconcilerFactory{
		resources.GetInternalKubeconfigReconciler(namespace, resources.ClusterBackupKubeconfigSecretName, resources.ClusterBackupUsername, nil, clusterData, r.log),
		seedclusterresources.SecretReconciler(ctx, r.Client, cluster, cbsl),
	}

	// Create kubeconfig secret in the user cluster namespace.
	if err := reconciling.ReconcileSecrets(ctx, secretReconcilers, namespace, r.Client); err != nil {
		return fmt.Errorf("failed to reconcile cluster backup kubeconfig secret: %w", err)
	}

	deploymentReconcilers := []reconciling.NamedDeploymentReconcilerFactory{
		seedclusterresources.DeploymentReconciler(clusterData),
	}
	if err := reconciling.ReconcileDeployments(ctx, deploymentReconcilers, namespace, r.Client); err != nil {
		return fmt.Errorf("failed to reconcile the cluster backup deployment: %w", err)
	}
	return nil
}

func (r *reconciler) newClusterData(ctx context.Context, cluster *kubermaticv1.Cluster) *resources.TemplateData {
	return resources.NewTemplateDataBuilder().
		WithContext(ctx).
		WithCluster(cluster).
		WithClient(r.Client).
		Build()
}
