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

package ironicconductor

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s_labels "k8s.io/apimachinery/pkg/labels"

	ironicv1 "github.com/openstack-k8s-operators/ironic-operator/api/v1beta1"
	ironic "github.com/openstack-k8s-operators/ironic-operator/pkg/ironic"
	common "github.com/openstack-k8s-operators/lib-common/modules/common"
	env "github.com/openstack-k8s-operators/lib-common/modules/common/env"
	"github.com/openstack-k8s-operators/lib-common/modules/common/helper"
)

// ConductorPods - Query current running ironic-conductor pods managed by the statefulset
func ConductorPods(
	ctx context.Context,
	instance *ironicv1.IronicConductor,
	helper *helper.Helper,
	serviceLabels map[string]string,
) (*corev1.PodList, error) {
	podSelectorString := k8s_labels.Set(serviceLabels).String()
	return helper.GetKClient().CoreV1().Pods(instance.Namespace).List(ctx, metav1.ListOptions{LabelSelector: podSelectorString})
}

// HostNetPod func
func HostNetPod(
	instance *ironicv1.IronicConductor,
	conductorPod *corev1.Pod,
) *corev1.Pod {
	runAsUser := int64(0)
	labels := map[string]string{
		common.AppSelector:       ironic.ServiceName,
		ironic.ComponentSelector: ironic.HostNetComponent,
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

	args := []string{"-c"}
	if instance.Spec.Debug.Service {
		args = append(args, common.DebugCommand)
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
		args = append(args, ServiceCommand)

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

	tftpProxyEnvVars := map[string]env.Setter{}

	dnsmasqEnvVars := map[string]env.Setter{}
	dnsmasqEnvVars["KOLLA_CONFIG_FILE"] = env.SetValue(DnsmasqKollaConfig)
	dnsmasqEnvVars["KOLLA_CONFIG_STRATEGY"] = env.SetValue("COPY_ALWAYS")

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Labels: labels,
			Name:   conductorPod.Name + "-" + conductorPod.Spec.NodeName,
		},
		Spec: corev1.PodSpec{
			HostNetwork:        true,
			ServiceAccountName: ironic.ServiceAccount,
			NodeName:           conductorPod.Spec.NodeName,
			Containers: []corev1.Container{
				// {
				// 	Name: "dhcp-relay",
				// 	Command: []string{
				// 		"/bin/bash",
				// 	},
				// 	Args:  args,
				// 	Image: instance.Spec.PxeContainerImage,
				// 	SecurityContext: &corev1.SecurityContext{
				// 		RunAsUser: &runAsUser,
				// 		Capabilities: &corev1.Capabilities{
				// 			Add: []corev1.Capability{
				// 				"NET_ADMIN",
				// 				"NET_RAW",
				// 			},
				// 		},
				// 	},
				// 	Env:            env.MergeEnvs([]corev1.EnvVar{}, dnsmasqEnvVars),
				// 	VolumeMounts:   GetVolumeMounts(),
				// 	Resources:      instance.Spec.Resources,
				// 	ReadinessProbe: dnsmasqReadinessProbe,
				// 	LivenessProbe:  dnsmasqLivenessProbe,
				// 	// StartupProbe:   startupProbe,
				// },
				{
					Name: "tftp-proxy",
					Command: []string{
						"/bin/bash",
					},
					Args:  args,
					Image: instance.Spec.TftpProxyContainerImage,
					SecurityContext: &corev1.SecurityContext{
						RunAsUser: &runAsUser,
					},
					Env:            env.MergeEnvs([]corev1.EnvVar{}, tftpProxyEnvVars),
					VolumeMounts:   GetVolumeMounts(),
					Resources:      instance.Spec.Resources,
					ReadinessProbe: tftpProxyReadinessProbe,
					LivenessProbe:  tftpProxyLivenessProbe,
					// StartupProbe:   startupProbe,
				},
			},
		},
	}
	return pod
}
