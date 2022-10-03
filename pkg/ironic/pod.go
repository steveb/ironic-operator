/*
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

package ironic

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ironicv1 "github.com/openstack-k8s-operators/ironic-operator/api/v1beta1"
	common "github.com/openstack-k8s-operators/lib-common/modules/common"
)

const (
	// ServiceCommand -
	ServiceCommand = "/usr/sbin/dnsmasq -k --port=0 --log-dhcp --log-debug --log-facility=- --dhcp-relay=$LOCAL_ADDRESS,$SERVER_ADDRESS"
)

// HostNetPod func
func HostNetPod(
	instance *ironicv1.Ironic,
	service *corev1.Service,
) *corev1.Pod {
	runAsUser := int64(0)
	labels := map[string]string{
		common.AppSelector: ServiceName,
		ComponentSelector:  HostNetComponent,
	}

	dnsmasqLivenessProbe := &corev1.Probe{
		// TODO might need tuning
		TimeoutSeconds:      10,
		PeriodSeconds:       30,
		InitialDelaySeconds: 3,
	}
	dnsmasqReadinessProbe := &corev1.Probe{
		// TODO might need tuning
		TimeoutSeconds:      10,
		PeriodSeconds:       30,
		InitialDelaySeconds: 5,
	}
	tftpProxyLivenessProbe := &corev1.Probe{
		// TODO might need tuning
		TimeoutSeconds:      10,
		PeriodSeconds:       30,
		InitialDelaySeconds: 3,
	}
	tftpProxyReadinessProbe := &corev1.Probe{
		// TODO might need tuning
		TimeoutSeconds:      10,
		PeriodSeconds:       30,
		InitialDelaySeconds: 5,
	}

	dnsmasqArgs := []string{"-c"}
	if instance.Spec.Debug.Service {
		dnsmasqArgs = append(dnsmasqArgs, common.DebugCommand)
		dnsmasqLivenessProbe.Exec = &corev1.ExecAction{
			Command: []string{
				"/bin/true",
			},
		}

		dnsmasqReadinessProbe.Exec = &corev1.ExecAction{
			Command: []string{
				"/bin/true",
			},
		}
	} else {
		dnsmasqArgs = append(dnsmasqArgs, ServiceCommand)

		//
		// https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/
		//
		dnsmasqLivenessProbe.Exec = &corev1.ExecAction{
			Command: []string{
				"/bin/true",
			},
		}

		dnsmasqReadinessProbe.Exec = &corev1.ExecAction{
			Command: []string{
				"/bin/true",
			},
		}
		tftpProxyLivenessProbe.Exec = &corev1.ExecAction{
			Command: []string{
				"/bin/true",
			},
		}

		tftpProxyReadinessProbe.Exec = &corev1.ExecAction{
			Command: []string{
				"/bin/true",
			},
		}
		// dnsmasqLivenessProbe.Exec = &corev1.ExecAction{
		// 	Command: []string{
		// 		"sh", "-c", "ss -lun | grep :67 && ss -lun | grep :69",
		// 	},
		// }

		// dnsmasqReadinessProbe.Exec = &corev1.ExecAction{
		// 	Command: []string{
		// 		"sh", "-c", "ss -lun | grep :67 && ss -lun | grep :69",
		// 	},
		// }
	}

	tftpProxyEnvVars := []corev1.EnvVar{
		{
			Name:  "HTTP_URL",
			Value: fmt.Sprintf("http://%s:%s/tftpboot/", service.Spec.ClusterIP, "8088"),
		},
	}

	dnsmasqEnvVars := []corev1.EnvVar{
		{
			Name: "LOCAL_ADDRESS",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "status.podIP",
				},
			},
		},
		{
			Name:  "SERVER_ADDRESS",
			Value: service.Spec.ClusterIP,
		},
	}

	nodeName := service.Labels[NodeName]

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Labels:    labels,
			Name:      service.Name + "-" + nodeName,
			Namespace: service.Namespace,
		},
		Spec: corev1.PodSpec{
			HostNetwork:        true,
			ServiceAccountName: ServiceAccount,
			NodeName:           nodeName,
			Containers: []corev1.Container{
				{
					Name: "dhcp-relay",
					Command: []string{
						"/bin/bash",
					},
					Args:  dnsmasqArgs,
					Image: instance.Spec.IronicConductor.PxeContainerImage,
					SecurityContext: &corev1.SecurityContext{
						RunAsUser: &runAsUser,
						Capabilities: &corev1.Capabilities{
							Add: []corev1.Capability{
								"NET_ADMIN",
								"NET_RAW",
							},
						},
					},
					Env:            dnsmasqEnvVars,
					ReadinessProbe: dnsmasqReadinessProbe,
					LivenessProbe:  dnsmasqLivenessProbe,
					// StartupProbe:   startupProbe,
				},
				{
					Name:  "tftp-proxy",
					Image: instance.Spec.IronicConductor.TftpProxyContainerImage,
					SecurityContext: &corev1.SecurityContext{
						RunAsUser: &runAsUser,
					},
					Env:            tftpProxyEnvVars,
					ReadinessProbe: tftpProxyReadinessProbe,
					LivenessProbe:  tftpProxyLivenessProbe,
					// StartupProbe:   startupProbe,
				},
			},
		},
	}
	return pod
}
