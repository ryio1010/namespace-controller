/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package namespace

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/exp/slices"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	namespacev1alpha1 "github.com/ryio1010/namespace-controller/api/namespace/v1alpha1"

	utilslack "github.com/ryio1010/namespace-controller/internal/slack"
)

type NamespaceData struct {
	Name string
}

func (nd *NamespaceData) ToString() string {
	return fmt.Sprintf("Namespace: %s", nd.Name)

}

// NotificationReconciler reconciles a Notification object
type NotificationReconciler struct {
	client.Client
	Scheme      *runtime.Scheme
	SlackClient *utilslack.Client
}

//+kubebuilder:rbac:groups=namespace.ryio1010.github.io,resources=notifications,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=namespace.ryio1010.github.io,resources=notifications/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=namespace.ryio1010.github.io,resources=notifications/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch

func (r *NotificationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	var namespace corev1.Namespace
	if err := r.Get(ctx, req.NamespacedName, &namespace); err != nil {
		if err := client.IgnoreNotFound(err); err != nil {
			// if the namespace is not found(deleted), do nothing
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	var notificationList namespacev1alpha1.NotificationList
	if err := r.List(ctx, &notificationList); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.notify(notificationList, &namespace); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NotificationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	startAt := time.Now()

	return ctrl.NewControllerManagedBy(mgr).
		For(&namespacev1alpha1.Notification{}).
		Watches(
			&corev1.Namespace{},
			&handler.EnqueueRequestForObject{},
			builder.WithPredicates(predicate.Funcs{
				CreateFunc: func(e event.CreateEvent) bool {
					// remove the namespace which is already existed before the controller starts
					return e.Object.GetCreationTimestamp().After(startAt)
				},
				DeleteFunc: func(e event.DeleteEvent) bool {
					return !e.DeleteStateUnknown
				},
			})).
		Complete(r)
}

func (r *NotificationReconciler) notify(notificationList namespacev1alpha1.NotificationList, namespace *corev1.Namespace) error {
	for _, notification := range notificationList.Items {
		// if the prefix of the namespace name is included in ignorePrefixes, do not notify
		if slices.Contains(notification.Spec.IgnorePrefixes, namespace.Name) {
			return nil
		}

		data := createNamespaceData(namespace)
		if err := r.SlackClient.PostMessage(data.ToString(), notification.Spec.Channel); err != nil {
			return err
		}
	}

	return nil
}

func createNamespaceData(namespace *corev1.Namespace) NamespaceData {
	return NamespaceData{
		Name: namespace.Name,
	}
}
