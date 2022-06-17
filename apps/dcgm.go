// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package apps

import (
   "fmt"
	"github.com/GoogleCloudPlatform/ops-agent/confgenerator"
	"github.com/GoogleCloudPlatform/ops-agent/confgenerator/otel"
)

type MetricsReceiverDcgm struct {
	confgenerator.ConfigComponent `yaml:",inline"`
	confgenerator.MetricsReceiverShared `yaml:",inline"`
}

func (r MetricsReceiverDcgm) Type() string {
	return "dcgm"
}

func (r MetricsReceiverDcgm) Pipelines() []otel.Pipeline {
	return []otel.Pipeline{{
		Receiver: otel.Component{
			Type: "dcgm",
			Config: map[string]interface{}{
				"collection_interval": r.CollectionIntervalString(),
			},
		},
		Processors: []otel.Component{
			otel.NormalizeSums(),
			otel.MetricsTransform(
				otel.AddPrefix("workload.googleapis.com"),
			),
		},
	}}
}

func init() {
   fmt.Printf("Registering DCGM receiver\n")
	confgenerator.MetricsReceiverTypes.RegisterType(func() confgenerator.Component { return &MetricsReceiverDcgm{} })
}


