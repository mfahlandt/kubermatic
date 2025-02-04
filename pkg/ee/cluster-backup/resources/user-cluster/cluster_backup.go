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

package userclusterresources

import (
	"context"
	"fmt"

	velerov1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"

	kubermaticv1 "k8c.io/kubermatic/v2/pkg/apis/kubermatic/v1"
	"k8c.io/kubermatic/v2/pkg/resources"
	kkpreconciling "k8c.io/kubermatic/v2/pkg/resources/reconciling"
	"k8c.io/reconciler/pkg/reconciling"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
)

const (
	clusterRoleBindingName = "velero"
	clusterBackupAppName   = "velero"
	defaultBSLName         = "default-cluster-backup-bsl"
)

// NamespaceReconciler creates the namespace for velero related resources on the user cluster.
func NamespaceReconciler() reconciling.NamedNamespaceReconcilerFactory {
	return func() (string, reconciling.NamespaceReconciler) {
		return resources.ClusterBackupNamespaceName, func(ns *corev1.Namespace) (*corev1.Namespace, error) {
			return ns, nil
		}
	}
}

// ServiceAccountReconciler creates the service account for velero on the user cluster.
func ServiceAccountReconciler() reconciling.NamedServiceAccountReconcilerFactory {
	return func() (string, reconciling.ServiceAccountReconciler) {
		return resources.ClusterBackupServiceAccountName, func(sa *corev1.ServiceAccount) (*corev1.ServiceAccount, error) {
			return sa, nil
		}
	}
}

// ClusterRoleBindingReconciler creates the clusterrolebinding for velero on the user cluster.
func ClusterRoleBindingReconciler() reconciling.NamedClusterRoleBindingReconcilerFactory {
	return func() (string, reconciling.ClusterRoleBindingReconciler) {
		return clusterRoleBindingName, func(crb *rbacv1.ClusterRoleBinding) (*rbacv1.ClusterRoleBinding, error) {
			crb.Labels = resources.BaseAppLabels(clusterBackupAppName, nil)
			crb.RoleRef = rbacv1.RoleRef{
				// too wide but probably needed to be able to do backups and restore.
				Name:     "cluster-admin",
				Kind:     "ClusterRole",
				APIGroup: rbacv1.GroupName,
			}
			crb.Subjects = []rbacv1.Subject{
				{
					Kind:     rbacv1.UserKind,
					Name:     resources.ClusterBackupServiceAccountName,
					APIGroup: rbacv1.GroupName,
				},
			}
			return crb, nil
		}
	}
}

// BSLReconciler creates the default BackupStorage location is created for velero.
func BSLReconciler(ctx context.Context, cluster *kubermaticv1.Cluster, cbsl *kubermaticv1.ClusterBackupStorageLocation) kkpreconciling.NamedBackupStorageLocationReconcilerFactory {
	return func() (string, kkpreconciling.BackupStorageLocationReconciler) {
		return defaultBSLName, func(bsl *velerov1.BackupStorageLocation) (*velerov1.BackupStorageLocation, error) {
			projectID, ok := cluster.Labels[kubermaticv1.ProjectIDLabelKey]
			if !ok {
				return nil, fmt.Errorf("cluster ProjectID label is not set")
			}
			bsl.Spec = *cbsl.Spec.DeepCopy()
			// we set this bsl as default and remove the secret reference to make it use the default velero secret.
			bsl.Spec.Default = true
			bsl.Spec.Credential = nil
			// add bucket prefix using projectID/clusterID to avoid collision.
			bsl.Spec.ObjectStorage.Prefix = fmt.Sprintf("%s/%s", projectID, cluster.Name)
			return bsl, nil
		}
	}
}
