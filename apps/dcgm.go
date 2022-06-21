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
	"github.com/GoogleCloudPlatform/ops-agent/confgenerator"
	"github.com/GoogleCloudPlatform/ops-agent/confgenerator/otel"
)

type MetricsReceiverDcgm struct {
	confgenerator.ConfigComponent `yaml:",inline"`
	confgenerator.MetricsReceiverShared `yaml:",inline"`

	Endpoint string `yaml:"endpoint" validate:"omitempty"`
   ProfilingMetrics bool `yaml:"profiling_metrics"`
}

const defaultDcgmEndpoint = "localhost:5555"

func (r MetricsReceiverDcgm) Type() string {
	return "dcgm"
}

func (r MetricsReceiverDcgm) Pipelines() []otel.Pipeline {
   if r.Endpoint == "" {
      r.Endpoint = defaultDcgmEndpoint
   }

   metrics := map[string]interface{}{
     "dcgm.gpu.utilization": map[string]bool {
        "enabled": true,
     },
     "dcgm.gpu.memory.bytes_used": map[string]bool {
        "enabled": true,
     },
   }

   if r.ProfilingMetrics {
      metrics = map[string]interface{}{
         "dcgm.gpu.profiling.sm_utilization": map[string]bool {
            "enabled": true,
         },
         "dcgm.gpu.profiling.sm_occupancy": map[string]bool {
            "enabled": true,
         },
         "dcgm.gpu.profiling.tensor_utilization": map[string]bool {
            "enabled": true,
         },
         "dcgm.gpu.profiling.fp64_utilization": map[string]bool {
            "enabled": true,
         },
         "dcgm.gpu.profiling.fp32_utilization": map[string]bool {
            "enabled": true,
         },
         "dcgm.gpu.profiling.fp16_utilization": map[string]bool {
            "enabled": true,
         },
         "dcgm.gpu.profiling.dram_utilization": map[string]bool {
            "enabled": true,
         },
         "dcgm.gpu.profiling.pcie_sent_bytes": map[string]bool {
            "enabled": true,
         },
         "dcgm.gpu.profiling.pcie_received_bytes": map[string]bool {
            "enabled": true,
         },
         "dcgm.gpu.profiling.nvlink_sent_bytes": map[string]bool {
            "enabled": true,
         },
         "dcgm.gpu.profiling.nvlink_received_bytes": map[string]bool {
            "enabled": true,
         },
      }
   }

	return []otel.Pipeline{{
		Receiver: otel.Component{
			Type: "dcgm",
			Config: map[string]interface{}{
				"collection_interval": r.CollectionIntervalString(),
            "endpoint": r.Endpoint,
            "metrics": metrics,
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
	confgenerator.MetricsReceiverTypes.RegisterType(func() confgenerator.Component { return &MetricsReceiverDcgm{} })
}


