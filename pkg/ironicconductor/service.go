package ironicconductor

import (
	"context"

	routev1 "github.com/openshift/api/route/v1"
	ironicv1 "github.com/openstack-k8s-operators/ironic-operator/api/v1beta1"
	ironic "github.com/openstack-k8s-operators/ironic-operator/pkg/ironic"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s_labels "k8s.io/apimachinery/pkg/labels"

	"github.com/openstack-k8s-operators/lib-common/modules/common"
	"github.com/openstack-k8s-operators/lib-common/modules/common/helper"
)

// ProvisionServices - Query current conductor provision services
func ProvisionServices(
	ctx context.Context,
	namespace string,
	helper *helper.Helper,
	serviceLabels map[string]string,
) (*corev1.ServiceList, error) {
	selector := k8s_labels.Set(serviceLabels).String()
	return helper.GetKClient().CoreV1().Services(namespace).List(ctx, metav1.ListOptions{LabelSelector: selector})
}

// ProvisionService - Service for conductor pod services exposed for provisioning
func ProvisionService(
	serviceName string,
	instance *ironicv1.IronicConductor,
	serviceLabels map[string]string,
	externalIPs []string,
) *corev1.Service {
	podSelector := map[string]string{
		common.AppSelector:       ironic.ServiceName,
		ironic.ComponentSelector: ironic.ConductorComponent,
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: instance.Namespace,
			Labels:    serviceLabels,
		},
		Spec: corev1.ServiceSpec{
			Selector: podSelector,
			Ports: []corev1.ServicePort{
				{
					Name:     ironic.HttpbootComponent,
					Port:     8088,
					Protocol: corev1.ProtocolTCP,
				},
				{
					Name:     ironic.DhcpComponent,
					Port:     67,
					Protocol: corev1.ProtocolUDP,
				},
			},
			// ExternalIPs: externalIPs,
		},
	}
}

// InternalService - Service to expose conductor JSON-RPC to API pods
func InternalService(
	serviceName string,
	instance *ironicv1.IronicConductor,
	serviceLabels map[string]string,
) *corev1.Service {
	podSelector := map[string]string{
		common.AppSelector:       ironic.ServiceName,
		ironic.ComponentSelector: ironic.ConductorComponent,
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName + "-internal",
			Namespace: instance.Namespace,
			Labels:    serviceLabels,
		},
		Spec: corev1.ServiceSpec{
			Selector: podSelector,
			Ports: []corev1.ServicePort{
				{
					Name:     ironic.JSONRPCComponent,
					Port:     8089,
					Protocol: corev1.ProtocolTCP,
				},
			},
		},
	}
}

// ProvisionRoute - Route for conductor pod services exposed for provisioning
func ProvisionRoute(
	serviceName string,
	instance *ironicv1.IronicConductor,
	serviceLabels map[string]string,
) *routev1.Route {

	return &routev1.Route{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: instance.Namespace,
			Labels:    serviceLabels,
		},
		Spec: routev1.RouteSpec{
			To: routev1.RouteTargetReference{
				Kind: "Service",
				Name: serviceName,
			},
		},
	}
}
