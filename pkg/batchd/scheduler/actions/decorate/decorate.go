/*
Copyright 2018 The Kubernetes Authors.

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

package decorate

import (
	"k8s.io/apimachinery/pkg/labels"

	arbapi "github.com/kubernetes-incubator/kube-arbitrator/pkg/batchd/api"
	"github.com/kubernetes-incubator/kube-arbitrator/pkg/batchd/scheduler/framework"
)

type decorateAction struct {
	ssn *framework.Session
}

func New() *decorateAction {
	return &decorateAction{}
}

func (alloc *decorateAction) Name() string {
	return "decorate"
}

func (alloc *decorateAction) Initialize() {}

func (alloc *decorateAction) Execute(ssn *framework.Session) {
	// fetch the nodes that match PodSet NodeSelector and NodeAffinity
	// and store it for following DRF assignment
	jobs := ssn.Jobs
	nodes := ssn.Nodes

	for _, job := range jobs {
		job.Candidates = fetchMatchNodeForPodSet(job, nodes)
	}
}

func (alloc *decorateAction) UnInitialize() {}

func fetchMatchNodeForPodSet(job *arbapi.JobInfo, nodes []*arbapi.NodeInfo) []*arbapi.NodeInfo {
	var matchNodes []*arbapi.NodeInfo

	if len(job.NodeSelector) == 0 {
		return nil
	}

	selector := labels.SelectorFromSet(labels.Set(job.NodeSelector))

	for _, node := range nodes {
		if selector.Matches(labels.Set(node.Node.Labels)) {
			matchNodes = append(matchNodes, node)
		}
	}

	return matchNodes
}
