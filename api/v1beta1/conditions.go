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

package v1beta1

import (
	condition "github.com/openstack-k8s-operators/lib-common/modules/common/condition"
)

//
// Ironic Condition Types used by API objects.
//
const (
	// IronicAPIReadyCondition Status=True condition which indicates if the IronicAPI is configured and operational
	IronicAPIReadyCondition condition.Type = "IronicAPIReady"

	// IronicConductorReadyCondition Status=True condition which indicates if the IronicConductor is configured and operational
	IronicConductorReadyCondition condition.Type = "IronicConductorReady"

	// IronicHostNetworkPodsReadyCondiction Status=True condition which indicates host network pods are operational
	IronicHostNetworkPodsReadyCondiction condition.Type = "IronicHostNetworkPods"
)

//
// Ironic Reasons used by API objects.
//
const ()

//
// Common Messages used by API objects.
//
const (
	//
	// IronicAPIReady condition messages
	//
	// IronicAPIReadyInitMessage
	IronicAPIReadyInitMessage = "IronicAPI not started"

	// IronicAPIReadyErrorMessage
	IronicAPIReadyErrorMessage = "IronicAPI error occured %s"

	//
	// IronicConductorReady condition messages
	//
	// IronicConductorReadyInitMessage
	IronicConductorReadyInitMessage = "IronicConductor not started"

	// IronicConductorReadyErrorMessage
	IronicConductorReadyErrorMessage = "IronicConductor error occured %s"

	//
	// IronicHostNetworkPodsReady condition messages
	//
	// IronicHostNetworkPodsReadyInitMessage
	IronicHostNetworkPodsReadyInitMessage = "Ironic host network pods not started"

	// IronicHostNetworkPodsReadyErrorMessage
	IronicHostNetworkPodsReadyErrorMessage = "Ironic host network pods error occured %s"
)
